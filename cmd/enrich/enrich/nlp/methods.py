"""nlp.nlp module"""

import re
import logging

from typing import Union, List

from collections import Counter
from heapq import nlargest
from string import punctuation
from gensim.summarization import keywords

from nlp.utils import (
    normalize_nfd,
    normalize_keyword,
    unique,
    unique_entities,
    normalize_text,
    tokenize_to_max_length,
    has_numbers,
    remove_punctuation,
)

from nltk.cluster.util import cosine_distance
from nltk import word_tokenize

import numpy as np
import networkx as nx

from transformers import PreTrainedTokenizer, PreTrainedTokenizerFast


async def extract_stopwords(text: str, stopwords: list[str]) -> list[str]:
    """Extract stopwords from text

    Args:
        text (str): The text to extract stopwords from
        stopwords (list[str]): The list of stopwords

    Returns:
        list[str]: A list of stopwords
    """

    tokens = word_tokenize(normalize_nfd(text.lower()))
    return [token for token in tokens if token in stopwords]


async def extract_keywords(doc) -> list[str]:
    """Extract keywords from text using spacy

    Args:
        doc (spacy.tokens.doc.Doc): Spacy Doc object

    Returns:
        list[str]: A list of keywords
    """

    pos_tags = ["PROPN", "ADJ", "NOUN", "VERB"]
    tokens = []
    for token in doc:
        if token.is_stop or token.is_punct or len(token.text) < 4:
            continue
        if token.pos_ in pos_tags:
            tokens.append(normalize_keyword(token.text))

    text = " ".join(tokens)
    keys = []
    try:
        keys = keywords(text, words=6, scores=False, lemmatize=False, deacc=True).split(
            "\n"
        )
        return " ".join(keys).upper().split()
    except (ValueError, TypeError):
        return [x[0] for x in Counter(tokens).most_common(6)]


async def extract_entities(doc) -> list[dict]:
    """Extract entities from text using spacy

    Args:
        doc (spacy.tokens.doc.Doc): Spacy Doc object

    Returns:
        list[dict]: A list of entities
    """

    entities = [
        {"text": normalize_keyword(w.text), "type": w.label_, "index": [w.start, w.end]}
        for w in doc.ents
        if w.label_ in ["GPE", "ORG", "PRESON"]
        and len(remove_punctuation(w.text)) >= 2
        and not has_numbers(w.text)
    ]
    return unique(entities)


def get_freq_word(doc) -> Counter:
    """Get frequency of words in text

    Args:
        doc (spacy.tokens.doc.Doc): Spacy Doc object

    Returns:
        Counter: A Counter object
    """

    # Filter Tokens, remove stopwords, punctuation
    # && keep only specific pos_tags

    keys = []
    pos_tags = ["PROPN", "ADJ", "NOUN", "VERB"]
    for token in doc:
        if token.is_stop or token.is_punct or len(token.text) < 4:
            continue
        if token.pos_ in pos_tags:
            keys.append(token.text)

    # Calculate the frequency of each token
    freq_word = Counter(keys)

    # Normalize frequency
    max_freq = Counter(keys).most_common(1)[0][1]
    for word in freq_word.keys():
        freq_word[word] = freq_word[word] / max_freq

    return freq_word


async def summarize_doc(doc, limit: int = 5) -> str:
    """Simple document extractive summarization

    Args:
        doc (spacy.tokens.doc.Doc): Spacy Doc object
        limit (int, optional): Number of sentenses. Defaults to 5.

    Returns:
        str: The summarized text
    """
    try:
        freq_word = get_freq_word(doc)
    except Exception as err:
        logging.error("Document summarization error: %s", err)
        return ""

    # Calculate sentences weight
    sent_strength = {}
    for sent in doc.sents:
        for word in sent:
            if word.text in freq_word:
                if sent in sent_strength:
                    sent_strength[sent] += freq_word[word.text]
                else:
                    sent_strength[sent] = freq_word[word.text]
    # Summarize sentences
    summarized_sentences = nlargest(
        int(len(list(doc.sents)) / limit), sent_strength, key=sent_strength.get
    )
    final_sentences = [w.text for w in summarized_sentences]
    return format_sentences(final_sentences)


async def extract_claims(doc, stopwords: list[str], per: float = 0.25) -> list[dict]:
    """Extract claims from text

    Args:
        doc (spacy.tokens.doc.Doc): Spacy Doc object
        stopwords (list[str]): A list of stopwords
        per (float, optional): Percentage. Defaults to 0.25.

    Returns:
        list[dict]: _description_
    """
    word_frequencies = {}
    for word in doc:
        if word.text.lower() not in stopwords:
            if word.text.lower() not in punctuation:
                if word.text not in word_frequencies:
                    word_frequencies[word.text] = 1
                else:
                    word_frequencies[word.text] += 1
    max_frequency = max(word_frequencies.values())
    for word in word_frequencies:
        word_frequencies[word] = word_frequencies[word] / max_frequency
    sentence_tokens = list(doc.sents)

    sentence_scores = {}
    for sent in sentence_tokens:
        for word in sent:
            if word.text.lower() in word_frequencies:
                if sent not in sentence_scores:
                    sentence_scores[sent] = word_frequencies[word.text.lower()]
                else:
                    sentence_scores[sent] += word_frequencies[word.text.lower()]

    select_length = int(len(sentence_tokens) * per)
    summary = nlargest(select_length, sentence_scores, key=sentence_scores.get)
    if len(summary) == 0:
        return []

    claims = []
    for word in summary:
        text = word.text.strip()
        text = text[0].upper() + text[1:] if len(text) > 1 else text
        claim = {
            "type": "claim",
            "index": [word.start, word.end],
            "score": sentence_scores[word],
            "text": text,
        }
        claims.append(claim)

    return claims


async def extract_topics(
    body: Union[str, List[str]],
    pipeline,
    tokenizer: Union[PreTrainedTokenizer, PreTrainedTokenizerFast] = None,
) -> list[dict]:
    body = tokenize_to_max_length(body, 512, tokenizer)
    topics = []

    try:
        topics = [
            pipeline(text, return_all_scores=True, top_k=4, max_length=512, truncation=True)
            for text in body
        ]
        topics = [topic for sublist in topics for topic in sublist]
    except Exception as err:
        logging.error("Topic extraction error: %s", err)

    topics = (
        [
            {"text": x["label"], "type": "topic", "score": x["score"]}
            for x in topics
            if x["score"] > 0.25
        ]
        if len(topics) > 0
        else []
    )
    topics = sorted(topics, key=lambda x: x["score"], reverse=True)
    return unique_entities(topics)


async def extract_named_entities(
    body: Union[str, List[str]],
    pipeline,
    tokenizer: Union[PreTrainedTokenizer, PreTrainedTokenizerFast] = None,
) -> list[dict]:
    body = tokenize_to_max_length(body, 512, tokenizer)
    named_entities = []

    try:
        named_entities = [pipeline(text, aggregation_strategy="first") for text in body]
        named_entities = [entity for sublist in named_entities for entity in sublist]
    except:
        pass
    named_entities = (
        [
            {
                "text": x["word"],
                "type": x["entity_group"],
                "score": x["score"],
                "index": [x["start"], x["end"]],
            }
            for x in named_entities
            if x["score"] > 0.2
        ]
        if len(named_entities) > 0
        else []
    )
    return named_entities


async def extractive_summarization(doc, stopwords, top_n=3):
    sentences = extract_sentences(doc)
    sentences_to_summarize = []
    sentence_similarity_martix = build_similarity_matrix(sentences, stopwords)
    sentence_similarity_graph = nx.from_numpy_array(sentence_similarity_martix)
    scores = nx.pagerank_numpy(sentence_similarity_graph)
    ranked_sentence = sorted(
        ((scores[i], s) for i, s in enumerate(sentences)), reverse=True
    )
    for i in range(int(len(sentences) / top_n)):
        sentences_to_summarize.append(ranked_sentence[i][1])
    return format_sentences(sentences_to_summarize)


def format_sentences(sentences_to_summarize: list[str]) -> str:
    """Format sentences to summarize

    Args:
        sentences_to_summarize (list[str]): A list of sentences

    Returns:
        str: A string of sentences
    """
    sentences = [
        sentence[0].upper() + sentence[1:] for sentence in sentences_to_summarize
    ]
    sentences = [sentence.strip() for sentence in sentences]
    return " ".join(sentences)


def extract_sentences(doc):
    sentences = [sentence.text for sentence in doc.sents]
    return list(set(sentences))


def build_similarity_matrix(sentences, stopwords) -> np.ndarray:
    similarity_matrix = np.zeros((len(sentences), len(sentences)))
    for i, _ in enumerate(sentences):
        for j, _ in enumerate(sentences):
            if i == j:  # ignore if both are same sentences
                continue
            similarity_matrix[i][j] = sentence_similarity(
                sentences[i], sentences[j], stopwords
            )
    return similarity_matrix


def sentence_similarity(sent1, sent2, stopwords):
    sent1 = [w.lower() for w in sent1]
    sent2 = [w.lower() for w in sent2]

    all_words = list(set(sent1 + sent2))

    vector1 = [0] * len(all_words)
    vector2 = [0] * len(all_words)

    # build the vector for the first sentence
    for w in sent1:
        if w in stopwords:
            continue
        vector1[all_words.index(w)] += 1

    # build the vector for the second sentence
    for w in sent2:
        if w in stopwords:
            continue
        vector2[all_words.index(w)] += 1

    return 1 - cosine_distance(vector1, vector2)


async def extract_quotes(body):
    expressions = [
        r"([\"'«‹])((.*?))(\1|[\»›])",
        r"([“])((.*?))([”])",
        r"([«])((.*?))([»])",
        r"([‹])((.*?))([›])",
    ]
    quotes = []
    for expression in expressions:
        regex = re.finditer(expression, body)
        for group in regex:
            if group:
                size = len(group.group().split())
                if size > 2 and size <= 48:
                    quotes.append(
                        {
                            "text": normalize_text(group.group()),
                            "type": "quote",
                            "index": [group.start(), group.end()],
                        }
                    )
    return quotes


async def summarize(doc, limit: int = 5) -> str:
    textrank = doc._.textrank
    sentences = [
        sentence.text
        for sentence in textrank.summary(
            limit_sentences=limit, limit_phrases=limit, preserve_order=True
        )
    ]
    return format_sentences(list(set(sentences)))
