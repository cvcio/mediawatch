import re
from unicodedata import normalize, category
from nltk.tokenize import RegexpTokenizer
from nltk.stem.porter import PorterStemmer


"""
    Remove accents from string
"""


def normalize_nfd(s):
    return "".join(c for c in normalize("NFD", s) if category(c) != "Mn")


"""
    Keep only unique elements in list
"""


def unique(l):
    unique_list = []
    for x in l:
        if x not in unique_list:
            unique_list.append(x)

    return unique_list


"""
    Normalize text
"""


def normalize_text(s):
    """
    Normalize input string
    """
    if s is None or len(s) == 0:
        return ""
    # unicode
    # text = unicode(s, 'utf-8')
    # trim text
    text = s.strip()
    # remove html forgotten tags
    text = re.sub(re.compile("<.*?>"), "", text)
    # remove double space
    text = re.sub(" +", " ", text)

    return text


"""
    Normalize keyword
"""


def normalize_keyword(s, strict=False):
    """
    Remove accents from input string
    """
    if s is None or len(s) == 0:
        return ""
    text = "".join(s)
    text = re.sub("[!@#$<>|]", "", text)
    text = normalize_text(text)
    if False is strict:
        text = normalize_nfd(text)

    return text.upper()


"""
    Prepare text for keyword extraction using gensim
    Tokenize and Stem the text to token list
"""


def prepare_text(text, stopwords):
    tokenizer = RegexpTokenizer(r"\w+")
    p_stemmer = PorterStemmer()

    text = text.lower()
    text = normalize_nfd(text)

    tokens = tokenizer.tokenize(text)
    tokens = [token for token in tokens if len(token) > 4]
    tokens = [token for token in tokens if token not in stopwords]
    tokens = [p_stemmer.stem(i) for i in tokens]

    return tokens
