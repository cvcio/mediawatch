"""String utils for NLP"""
from __future__ import annotations

import re
import string

from unicodedata import normalize, category
from nltk.tokenize import RegexpTokenizer
from nltk.stem.porter import PorterStemmer
from transformers import PreTrainedTokenizer, PreTrainedTokenizerFast


def normalize_nfd(value):
    """
    Remove accents from string
    """
    return "".join(char for char in normalize("NFD", value) if category(char) != "Mn")


def unique(items):
    """
    Keep only unique elements in list
    """
    unique_list = []
    for item in items:
        if item not in unique_list:
            unique_list.append(item)

    return unique_list


def normalize_text(value):
    """
    Normalize text - Normalize input string
    """
    if value is None or len(value) == 0:
        return ""
    # unicode
    # text = unicode(s, 'utf-8')
    # trim text
    text = value.strip()
    # remove html forgotten tags
    text = re.sub(re.compile("<.*?>"), "", text)
    # remove double space
    text = re.sub(" +", " ", text)

    return text


def normalize_keyword(value, strict=False):
    """
    Normalize keyword - Remove accents from input string
    """
    if value is None or len(value) == 0:
        return ""
    text = "".join(value)
    text = re.sub("[!@#$<>|]", "", text)
    text = normalize_text(text)
    if False is strict:
        text = normalize_nfd(text)

    return text.upper()


def prepare_text(text, stopwords):
    """
    Prepare text for keyword extraction using gensim
    Tokenize and Stem the text to token list
    """
    tokenizer = RegexpTokenizer(r"\w+")
    p_stemmer = PorterStemmer()

    text = text.lower()
    text = normalize_nfd(text)

    tokens = tokenizer.tokenize(text)
    tokens = [token for token in tokens if len(token) > 4]
    tokens = [token for token in tokens if token not in stopwords]
    tokens = [p_stemmer.stem(i) for i in tokens]

    return tokens


def tokenize_to_max_length(
    text,
    max_length: int = 512,
    tokenizer: (PreTrainedTokenizer | PreTrainedTokenizerFast) = None,
):
    """Tokenize text to max length

    Args:
        text (str): Input test
        max_length (int, optional): Max length. Defaults to 512
        tokenizer: Tokenizer. Defaults to None

    Returns:
        list[str]: A list of text chunks of max length
    """
    if tokenizer is None:
        tokens = text.split(" ")
        chunks = [tokens[i : i + max_length] for i in range(0, len(tokens), max_length)]
        chunks = [" ".join(list(chunk)) for chunk in chunks]
        return chunks

    output = tokenizer.encode(text, add_special_tokens=False)
    encoded_chunks = [
        output[i : i + max_length] for i in range(0, len(output), max_length)
    ]
    decoded_chunks = [
        tokenizer.decode(c, skip_special_tokens=True) for c in encoded_chunks
    ]
    return decoded_chunks


def unique_entities(items):
    """
    Keep only unique elements in list of entities
    """
    unique_list = []
    unique_items = []
    for item in items:
        if item["text"] not in unique_list:
            unique_list.append(item["text"])
            unique_items.append(item)

    return unique_items


def has_numbers(value):
    return bool(re.search(r"\d", value))


def remove_punctuation(value):
    translator = str.maketrans("", "", string.punctuation)
    return value.translate(translator)


def in_black_list(value, black_list):
    return value in black_list
