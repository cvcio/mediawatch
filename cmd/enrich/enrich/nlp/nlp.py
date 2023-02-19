import re
import logging

from nlp.utils import normalize_nfd, normalize_keyword, unique, normalize_text
from nltk.cluster.util import cosine_distance
from nltk import word_tokenize

from gensim.summarization import keywords
from collections import Counter
from heapq import nlargest

from string import punctuation

import numpy as np
import networkx as nx


def extract_stopwords(text, stopwords):
    tokens = word_tokenize(normalize_nfd(text.lower()))
    return [token for token in tokens if token in stopwords]


def extract_keywords(doc):
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
    except (ValueError, TypeError, Exception):
        return [x[0] for x in Counter(tokens).most_common(6)]


def extract_entities(doc):
    entities = [
        {"text": normalize_keyword(w.text), "type": w.label_, "index": [w.start, w.end]}
        for w in doc.ents
        if w.label_ in ["GPE", "ORG", "PRESON"]
    ]
    return unique(entities)


def get_freq_word(doc):
    # Filter Tokens, remove stopwords, punctuation
    # && keep only specific pos_tags
    keywords = []
    pos_tags = ["PROPN", "ADJ", "NOUN", "VERB"]
    for token in doc:
        if token.is_stop or token.is_punct or len(token.text) < 4:
            continue
        if token.pos_ in pos_tags:
            keywords.append(token.text)

    # Calculate the frequency of each token
    freq_word = Counter(keywords)

    # Normalize frequency
    max_freq = Counter(keywords).most_common(1)[0][1]
    for word in freq_word.keys():
        freq_word[word] = freq_word[word] / max_freq

    return freq_word


def summarize_doc(doc, limit):
    try:
        freq_word = get_freq_word(doc)
    except Exception as err:
        logging.error("Document summarization error: {}".format(err))
        return ""

    # Calculate sentences weight
    sent_strength = {}
    for sent in doc.sents:
        for word in sent:
            if word.text in freq_word.keys():
                if sent in sent_strength.keys():
                    sent_strength[sent] += freq_word[word.text]
                else:
                    sent_strength[sent] = freq_word[word.text]
    # Summarize sentences
    summarized_sentences = nlargest(
        int(len(list(doc.sents)) / limit), sent_strength, key=sent_strength.get
    )
    final_sentences = [w.text for w in summarized_sentences]
    return format_sentences(final_sentences)


def extract_claims(doc, stopwords, per=0.25):
    word_frequencies = {}
    for word in doc:
        if word.text.lower() not in list(stopwords):
            if word.text.lower() not in punctuation:
                if word.text not in word_frequencies.keys():
                    word_frequencies[word.text] = 1
                else:
                    word_frequencies[word.text] += 1
    max_frequency = max(word_frequencies.values())
    for word in word_frequencies.keys():
        word_frequencies[word] = word_frequencies[word] / max_frequency
    sentence_tokens = [sent for sent in doc.sents]

    sentence_scores = {}
    for sent in sentence_tokens:
        for word in sent:
            if word.text.lower() in word_frequencies.keys():
                if sent not in sentence_scores.keys():
                    sentence_scores[sent] = word_frequencies[word.text.lower()]
                else:
                    sentence_scores[sent] += word_frequencies[word.text.lower()]

    select_length = int(len(sentence_tokens) * per)
    summary = nlargest(select_length, sentence_scores, key=sentence_scores.get)
    if len(summary) == 0:
        return []

    claims = [{"text": word.text, "type": "claim", "index": [word.start, word.end], "score": sentence_scores[word]} for word in summary]
    claims = [{"text": c["text"].strip(), "type": c["type"], "index": c["index"], "score": c["score"]} for c in claims]
    claims = [{"text": c["text"][0].upper() + c["text"][1:], "type": c["type"], "index": c["index"], "score": c["score"]} for c in claims]

    return claims


def extract_topics(body, pipeline):
    topics = []
    try:
        topics = pipeline(
            body,
            padding=True,
            truncation=True,
            max_length=512,
            return_all_scores=True,
        )[0]
        logging.info(topics)
    except:
        pass
    topics = [
        {"text": x["label"], "type": "topic", "score": x["score"]}
        for x in topics if x["score"] > 0.2
    ] if len(topics) > 0 else []
    return topics


def extractive_summarization(doc, stopwords, top_n=3):
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


def format_sentences(sentences_to_summarize):
    sentences = [
        sentence[0].upper() + sentence[1:] for sentence in sentences_to_summarize
    ]
    sentences = [sentence.strip() for sentence in sentences]
    return " ".join(sentences)


def extract_sentences(doc):
    sentences = [sentence.text for sentence in doc.sents]
    return list(set(sentences))


def build_similarity_matrix(sentences, stopwords):
    similarity_matrix = np.zeros((len(sentences), len(sentences)))
    for idx1 in range(len(sentences)):
        for idx2 in range(len(sentences)):
            if idx1 == idx2:  # ignore if both are same sentences
                continue
            similarity_matrix[idx1][idx2] = sentence_similarity(
                sentences[idx1], sentences[idx2], stopwords
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


def extract_quotes(body):
    # ([\"'“«‹])(.*?)(\1|[\”»›])
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
                    quotes.append({
                        "text": normalize_text(group.group()),
                        "type": "quote",
                        "index": [group.start(), group.end()]
                    })
    return quotes
