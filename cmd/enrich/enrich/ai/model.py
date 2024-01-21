"""AIModel module loads pretrained models
for a specific language and multiple tasks."""
from __future__ import annotations

import os
import io
import json
import logging

from typing import Any

from collections import namedtuple
from config.config import AppConfig

import spacy
import pytextrank # noqa # pylint: disable=unused-import

from transformers import (
    pipeline,
    AutoTokenizer,
    AutoModelForSequenceClassification,
    PreTrainedTokenizer,
    PreTrainedTokenizerFast,
)

from transformers.pipelines import Pipeline


def read_configuration(file: str) -> Any:
    """Read model from file"""
    with io.open(os.path.join(file), encoding="utf-8") as f:
        return json.loads(
            f.read(), object_hook=lambda d: namedtuple("field", d.keys())(*d.values())
        )


class AIModel:
    """AIModel class loads pretrained models for a specific language."""

    def __init__(self, file: str):
        self.conf = read_configuration(file)
        self.lang = self.conf.lang

        self.tokenizer = (
            self.load_transformers_tokenizer(self.conf.tokenizer.path)
            if getattr(self.conf, "tokenizer", None)
            else None
        )

        self.topic_classification_pipeline = (
            self.load_transformers_pipeline(
                self.conf.topics.path,
                "text-classification",
            )
            if getattr(self.conf, "topics", None)
            else None
        )

        self.ner_classification_pipeline = (
            self.load_transformers_pipeline(
                self.conf.ner.path, "ner", aggregation_strategy="first"
            )
            if getattr(self.conf, "ner", None)
            else None
        )

        self.spacy = (
            self.load_spacy_model(self.conf.spacy.path)
            if getattr(self.conf, "spacy", None)
            else None
        )

        logging.info("Loaded model for lang: %s", self.conf.lang)

    def load_transformers_pipeline(
        self, path: str, task: str = "text-classification", **kwargs
    ) -> Pipeline:
        """Load transformers model"""
        return pipeline(
            task=task,
            model=AutoModelForSequenceClassification.from_pretrained(
                path,
                # token=True if AppConfig(os.environ).ENV == "development" else AppConfig(os.environ).HUGGING_FACE_HUB_TOKEN,
            ),
            tokenizer=self.tokenizer,
            device=AppConfig(os.environ).DEVICE,
            **kwargs,
        )

    def load_transformers_tokenizer(
        self, path: str
    ) -> (PreTrainedTokenizer | PreTrainedTokenizerFast):
        """Load transformers tokenizer"""
        return AutoTokenizer.from_pretrained(
            path,
            use_fast=True,
            truncate=True,
            max_length=512,
            # token=True if AppConfig(os.environ).ENV == "development" else AppConfig(os.environ).HUGGING_FACE_HUB_TOKEN,
        )

    def load_spacy_model(self, path: str) -> spacy.language.Language:
        """Load spacy model"""
        if AppConfig(os.environ).ENV != "development":
            spacy.cli.download(path)

        model = spacy.load(path)
        model.Defaults.stop_words.update(self.conf.stopwords)
        model.add_pipe("textrank", config={"token_lookback": 10})
        return model
