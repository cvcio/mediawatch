import os
import io
import json
import logging

import spacy

from config.config import AppConfig
from collections import namedtuple
from transformers import (
    AutoTokenizer,
    AutoModelForSequenceClassification,
    AutoModelForTokenClassification,
    pipeline,
)


class Model:
    """
    Model class loads pretrained models for a specific language.
    Options are defined inside each language configuration file
    located under the models folder.

    Available options are:

    - lang: iso 639-1 language code (ex. el)
    - spacy:spacy's base language model name (ex. en_core_web_sm). we only use
    spacy's document parser (nlp) to extract entities and tokens.
    - stopwods: a list of additional stopwords
    - tokenizer: path or name of the tokenizer (via @huggingface)
    - topics: path or name of the topic classifier (via @huggingface), for the moment
    we only support the { nlp: topics } field. in the future we should
    add more classifiers on each model to support multiple classification
    tasks.
    - ner: path or name of the ner classifier (via @huggingface).
    """

    def __init__(self, file):
        """
        initialize model object
        """

        # create an empty namedtuple with default values
        p = namedtuple(
            "model",
            ["lang", "tokenizer", "topics", "ner", "spacy", "stopwords"],
            defaults=(None,)
            * len(["lang", "tokenizer", "topics", "ner", "spacy", "stopwords"]),
        )
        # read the configuration file
        self.model = p(**self.read_file(file))
        logging.debug("Config for model {} loaded ({})".format(self.model.lang, file))
        # load spacy model (if available)
        self.spacy = (
            self.get_spacy() if self.model.spacy and self.model.spacy != "" else None
        )
        # load the topic classifier (if available)
        self.topic_classification_pipeline = (
            self.get_topics_classifier() if (self.model.topics != None) else None
        )
        # load the ner classifier (if available)
        self.ner_classification_pipeline = (
            self.get_ner_classifier()
            if (self.model.tokenizer != None and self.model.ner != None)
            else None
        )

    def read_file(self, file):
        """
        read model configuration file and return a dict
        """
        with io.open(os.path.join(file), encoding="utf-8") as f:
            return json.loads(
                f.read(),
            )

    def get_spacy(self):
        """
        load spacy model
        """
        if AppConfig(os.environ).ENV != "development":
            spacy.cli.download(self.model.spacy)

        m = spacy.load(self.model.spacy)
        m.Defaults.stop_words.update(self.model.stopwords)
        logging.debug("Spacy model {} loaded".format(self.model.spacy))
        return m

    def get_topics_classifier(self):
        """
        load tokenizer, text-classification and return the pipeline
        """

        # tokenizer
        tokenizer = AutoTokenizer.from_pretrained(
            self.model.tokenizer,
            use_fast=True,
            truncate=True,
            max_length=512,
            use_auth_token=True,
        )
        # text-calssification model
        model = AutoModelForSequenceClassification.from_pretrained(
            self.model.topics,
            use_auth_token=True,
        )

        return (
            pipeline(
                "text-classification",
                model=model,
                tokenizer=self.model.topics,
                top_k=4,
                device=-1,
            )
            if (tokenizer != None and model != None)
            else None
        )

    def get_ner_classifier(self):
        """
        load tokenizer, text-classification and return the pipeline
        """

        # tokenizer
        tokenizer = AutoTokenizer.from_pretrained(
            self.model.tokenizer,
            use_fast=True,
            truncate=True,
            max_length=512,
            use_auth_token=True,
        )
        # ner model
        model = AutoModelForTokenClassification.from_pretrained(
            self.model.ner,
            use_auth_token=True,
        )

        return (
            pipeline(
                "ner",
                model=model,
                tokenizer=tokenizer,
                aggregation_strategy="first",
                device=-1,
            )
            if (tokenizer != None and model != None)
            else None
        )
