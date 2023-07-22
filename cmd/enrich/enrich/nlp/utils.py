import re
from unicodedata import normalize, category
from nltk.tokenize import RegexpTokenizer
from nltk.stem.porter import PorterStemmer


def normalize_nfd(s):
    """
        Remove accents from string
    """
    return "".join(c for c in normalize("NFD", s) if category(c) != "Mn")


def unique(l):
    """
        Keep only unique elements in list
    """
    unique_list = []
    for x in l:
        if x not in unique_list:
            unique_list.append(x)

    return unique_list


def normalize_text(s):
    """
    Normalize text - Normalize input string
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


def normalize_keyword(s, strict=False):
    """
    Normalize keyword - Remove accents from input string
    """
    if s is None or len(s) == 0:
        return ""
    text = "".join(s)
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


def tokenize_to_max_length(text, max_length:int = 512):
    tokens = text.split(" ")
    chunks = [tokens[i:i + max_length] for i in range(0, len(tokens), max_length)]
    chunks = [" ".join([c for c in chunk]) for chunk in chunks]
    return chunks


def unique_entities(l):
    """
        Keep only unique elements in list of entities
    """
    unique_list = []
    unique_items = []
    for x in l:
        if x["text"] not in unique_list:
            unique_list.append(x["text"])
            unique_items.append(x)

    return unique_items
