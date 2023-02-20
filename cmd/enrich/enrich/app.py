"""
MediaWatch Enrich Service v2.0
Python implementation of the gRPC Enrich server.
"""

import os
import io
import logging
import grpc
import time
import glob
import json
import spacy
import html
import nltk

from concurrent import futures
from collections import namedtuple

from dotenv import load_dotenv
from config.config import AppConfig

from mediawatch.enrich.v2 import enrich_pb2_grpc
from mediawatch.enrich.v2 import enrich_pb2

from google.rpc import code_pb2
from google.rpc import status_pb2
from grpc_status import rpc_status

from transformers import AutoTokenizer, AutoModelForSequenceClassification, pipeline
from nlp.nlp import (
    extract_stopwords,
    extract_keywords,
    extract_entities,
    summarize_doc,
    extract_topics,
    extract_quotes,
    extract_claims
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
    - tokenizer: path or name of the tokenizer (via @huggingfaces)
    - classifier: path or name of the classifier (via @huggingfaces). for the moment
    we only support the { nlp: topics } field. in the future we should
    add more classifiers on each model to support multiple classification
    tasks.
    """

    def __init__(self, file):
        """
        initialize model object
        """

        # create an empty namedtuple with default values
        p = namedtuple(
            "model",
            ["lang", "tokenizer", "classifier", "spacy", "stopwords"],
            defaults=(None,)
            * len(["lang", "tokenizer", "classifier", "spacy", "stopwords"]),
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
            self.get_transformers()
            if (self.model.tokenizer != None and self.model.classifier != None)
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
        if os.environ["ENV"] != "development":
            spacy.cli.download(self.model.spacy)
        m = spacy.load(self.model.spacy)
        m.Defaults.stop_words.update(self.model.stopwords)
        logging.debug("Spacy model {} loaded".format(self.model.spacy))
        return m

    def get_transformers(self):
        """
        load tokenizer, text-classification and return the pipeline
        """

        # tokenizer
        tokenizer = AutoTokenizer.from_pretrained(
            self.model.tokenizer,
            use_fast=True,
            truncate=True,
            max_length=512,
            use_auth_token=os.environ["HUGGINGFACES_API_TOKEN"],
        )
        # text-calssification model
        classifier = AutoModelForSequenceClassification.from_pretrained(
            self.model.classifier,
            use_auth_token=os.environ["HUGGINGFACES_API_TOKEN"],
        )

        return (
            pipeline(
                "text-classification",
                model=classifier,
                tokenizer=tokenizer,
                top_k=4,
                device=-1,
            )
            if (tokenizer != None and classifier != None)
            else None
        )


class EnrichService(enrich_pb2_grpc.EnrichServiceServicer):
    """
    EnrichService stub implementation of the gRPC EnrichService servicer.
    """

    def __init__(self, models=None):
        """
        initialize gRPC stub
        """
        # add the models inside the stub
        self.models = models


    def StopWords(self, data, context):
        """
        StopWords gRPC endpoint
        """
        if self.models == None:
            # if there are no models return an error
            return context.abort_with_status(
                rpc_status.to_status(
                    status_pb2.Status(
                        code=code_pb2.INTERNAL,
                        message="unable to extract stopwords from document, no models defined",
                    )
                )
            )

        if data.body == "" or len(data.body) < 24:
            return context.abort_with_status(
                rpc_status.to_status(
                    status_pb2.Status(
                        code=code_pb2.INTERNAL,
                        message="document body empty or too short",
                    )
                )
            )

        logging.info("Retrieve StopWords from document with language model: {}".format(data.lang))

        # select language specific model (Model). language is defined
        # by the incoming request (field: lang).
        m = next(
            model for model in self.models if model.model.lang == data.lang.lower()
        )

        output = {
            "nlp": {
                "stopwords": [],
            }
        }

        # escape text
        body = html.unescape(data.body)

        try:
            # get the stopwords
            stopwords = extract_stopwords(body, list(m.spacy.Defaults.stop_words))
            output["nlp"]["stopwords"] = stopwords
            logging.debug("Stopwords: {}".format(stopwords))
        except Exception as err:
            logging.error("Error while getting stopwords: {}".format(str(err)))
            return context.abort_with_status(
                rpc_status.to_status(
                    status_pb2.Status(
                        code=code_pb2.INTERNAL,
                        message=str(err),
                    )
                )
            )

        logging.info(
            "Stopwords ({})".format(
                len(output["nlp"]["stopwords"]),
            )
        )
        return enrich_pb2.EnrichResponse(code=200, status="success", data=output)


    def Keywords(self, data, context):
        """
        Keywords gRPC endpoint
        """
        if self.models == None:
            # if there are no models return an error
            return context.abort_with_status(
                rpc_status.to_status(
                    status_pb2.Status(
                        code=code_pb2.INTERNAL,
                        message="unable to extract keywords from document, no models defined",
                    )
                )
            )

        if data.body == "" or len(data.body) < 24:
            return context.abort_with_status(
                rpc_status.to_status(
                    status_pb2.Status(
                        code=code_pb2.INTERNAL,
                        message="document body empty or too short",
                    )
                )
            )

        logging.info("Retrieve keywords from document with language model: {}".format(data.lang))

        # select language specific model (Model). language is defined
        # by the incoming request (field: lang).
        m = next(
            model for model in self.models if model.model.lang == data.lang.lower()
        )

        output = {
            "nlp": {
                "keywords": [],
            }
        }

        # escape text
        body = html.unescape(data.body)

        # parse the document usign a spacy model
        doc = m.spacy(body)
        logging.debug("Doc: {}".format(doc))

        try:
            # get the keywords
            keywords = extract_keywords(doc)
            output["nlp"]["keywords"] = keywords
            logging.debug("Keywords: {}".format(keywords))
        except Exception as err:
            logging.error("Error while getting keywords: {}".format(str(err)))
            return context.abort_with_status(
                rpc_status.to_status(
                    status_pb2.Status(
                        code=code_pb2.INTERNAL,
                        message=str(err),
                    )
                )
            )

        logging.info(
            "Keywords ({})".format(
                len(output["nlp"]["keywords"]),
            )
        )
        return enrich_pb2.EnrichResponse(code=200, status="success", data=output)


    def Entities(self, data, context):
        """
        Entities gRPC endpoint
        """
        if self.models == None:
            # if there are no models return an error
            return context.abort_with_status(
                rpc_status.to_status(
                    status_pb2.Status(
                        code=code_pb2.INTERNAL,
                        message="unable to extrct entities from document, no models defined",
                    )
                )
            )

        if data.body == "" or len(data.body) < 24:
            return context.abort_with_status(
                rpc_status.to_status(
                    status_pb2.Status(
                        code=code_pb2.INTERNAL,
                        message="document body empty or too short",
                    )
                )
            )

        logging.info("Enrirch document with language model: {}".format(data.lang))

        # select language specific model (Model). language is defined
        # by the incoming request (field: lang).
        m = next(
            model for model in self.models if model.model.lang == data.lang.lower()
        )

        output = {
            "nlp": {
                "entities": [],
            }
        }

        # escape text
        body = html.unescape(data.body)

        # parse the document usign a spacy model
        doc = m.spacy(body)
        logging.debug("Doc: {}".format(doc))

        try:
            # get the extracted entites
            entities = extract_entities(doc)
            output["nlp"]["entities"] = entities
            logging.debug("Entities: {}".format(entities))
        except Exception as err:
            logging.error("Error while getting entities: {}".format(str(err)))
            return context.abort_with_status(
                rpc_status.to_status(
                    status_pb2.Status(
                        code=code_pb2.INTERNAL,
                        message=str(err),
                    )
                )
            )

        logging.info(
            "Entities ({})".format(
                len(output["nlp"]["entities"]),
            )
        )
        return enrich_pb2.EnrichResponse(code=200, status="success", data=output)


    def Summary(self, data, context):
        """
        Summary gRPC endpoint
        """
        if self.models == None:
            # if there are no models return an error
            return context.abort_with_status(
                rpc_status.to_status(
                    status_pb2.Status(
                        code=code_pb2.INTERNAL,
                        message="unable to extract summary from document, no models defined",
                    )
                )
            )

        if data.body == "" or len(data.body) < 24:
            return context.abort_with_status(
                rpc_status.to_status(
                    status_pb2.Status(
                        code=code_pb2.INTERNAL,
                        message="document body empty or too short",
                    )
                )
            )

        logging.info("Retrieve summary from document with language model: {}".format(data.lang))

        # select language specific model (Model). language is defined
        # by the incoming request (field: lang).
        m = next(
            model for model in self.models if model.model.lang == data.lang.lower()
        )

        output = {
            "nlp": {
                "summary": "",
            }
        }

        # escape text
        body = html.unescape(data.body)

        # parse the document usign a spacy model
        doc = m.spacy(body)
        logging.debug("Doc: {}".format(doc))

        try:
            # generate a summary of the text
            summary = summarize_doc(doc, 3)
            output["nlp"]["summary"] = summary
            logging.debug("Summary: {}".format(summary))
        except Exception as err:
            logging.error("Error while getting summary: {}".format(str(err)))
            return context.abort_with_status(
                rpc_status.to_status(
                    status_pb2.Status(
                        code=code_pb2.INTERNAL,
                        message=str(err),
                    )
                )
            )

        return enrich_pb2.EnrichResponse(code=200, status="success", data=output)


    def Topics(self, data, context):
        """
        Topics gRPC endpoint
        """
        if self.models == None:
            # if there are no models return an error
            return context.abort_with_status(
                rpc_status.to_status(
                    status_pb2.Status(
                        code=code_pb2.INTERNAL,
                        message="unable to extract topics from document, no models defined",
                    )
                )
            )

        if data.body == "" or len(data.body) < 24:
            return context.abort_with_status(
                rpc_status.to_status(
                    status_pb2.Status(
                        code=code_pb2.INTERNAL,
                        message="document body empty or too short",
                    )
                )
            )

        logging.info("Retrieve topics from document with language model: {}".format(data.lang))

        # select language specific model (Model). language is defined
        # by the incoming request (field: lang).
        m = next(
            model for model in self.models if model.model.lang == data.lang.lower()
        )

        output = {
            "nlp": {
                "topics": [],
            }
        }

        # escape text
        body = html.unescape(data.body)

        try:
            # classify the text usign a pretrained classifier
            topics = (
                extract_topics(body, m.topic_classification_pipeline)
                if m.topic_classification_pipeline != None
                else []
            )
            output["nlp"]["topics"] = topics
            logging.debug("Topics: {}".format(topics))
        except Exception as err:
            logging.error("Error while getting topics: {}".format(str(err)))
            return context.abort_with_status(
                rpc_status.to_status(
                    status_pb2.Status(
                        code=code_pb2.INTERNAL,
                        message=str(err),
                    )
                )
            )

        logging.info(
            "Topics ({}) {}".format(
                len(output["nlp"]["topics"]),
                ", ".join(topics)
            )
        )
        return enrich_pb2.EnrichResponse(code=200, status="success", data=output)


    def Quotes(self, data, context):
        """
        Quotes gRPC endpoint
        """

        if data.body == "" or len(data.body) < 24:
            return context.abort_with_status(
                rpc_status.to_status(
                    status_pb2.Status(
                        code=code_pb2.INTERNAL,
                        message="document body empty or too short",
                    )
                )
            )

        logging.info("Retrieve Quotes from document with language model: {}".format(data.lang))

        output = {
            "nlp": {
                "quotes": [],
            }
        }

        # escape text
        body = html.unescape(data.body)

        try:
            # get the quotes
            quotes = extract_quotes(body)
            output["nlp"]["quotes"] = quotes
            logging.debug("Quotes: {}".format(quotes))
        except Exception as err:
            logging.error("Error while getting quotes: {}".format(str(err)))
            return context.abort_with_status(
                rpc_status.to_status(
                    status_pb2.Status(
                        code=code_pb2.INTERNAL,
                        message=str(err),
                    )
                )
            )

        logging.info(
            "Quotes ({})".format(
                len(output["nlp"]["quotes"]),
            )
        )
        return enrich_pb2.EnrichResponse(code=200, status="success", data=output)


    def Claims(self, data, context):
        """
        Claims gRPC endpoint
        """
        if self.models == None:
            # if there are no models return an error
            return context.abort_with_status(
                rpc_status.to_status(
                    status_pb2.Status(
                        code=code_pb2.INTERNAL,
                        message="unable to extract cclaims from document, no models defined",
                    )
                )
            )

        if data.body == "" or len(data.body) < 24:
            return context.abort_with_status(
                rpc_status.to_status(
                    status_pb2.Status(
                        code=code_pb2.INTERNAL,
                        message="document body empty or too short",
                    )
                )
            )

        logging.info("Retrive claims from document with language model: {}".format(data.lang))

        # select language specific model (Model). language is defined
        # by the incoming request (field: lang).
        m = next(
            model for model in self.models if model.model.lang == data.lang.lower()
        )

        output = {
            "nlp": {
                "claims": [],
            }
        }

        # escape text
        body = html.unescape(data.body)

        # parse the document usign a spacy model
        doc = m.spacy(body)
        logging.debug("Doc: {}".format(doc))

        try:
            # get the claims (top 50%)
            claims = extract_claims(doc, list(m.spacy.Defaults.stop_words), 0.5)
            output["nlp"]["claims"] = claims
            logging.debug("Claims: {}".format(claims))
        except Exception as err:
            logging.error("Error while getting claims: {}".format(str(err)))
            return context.abort_with_status(
                rpc_status.to_status(
                    status_pb2.Status(
                        code=code_pb2.INTERNAL,
                        message=str(err),
                    )
                )
            )

        logging.info(
            "Claims ({})".format(
                len(output["nlp"]["claims"]),
            )
        )
        return enrich_pb2.EnrichResponse(code=200, status="success", data=output)


    def NLP(self, data, context):
        """
        NLP gRPC endpoint
        """

        if self.models == None:
            # if there are no models return an error
            return context.abort_with_status(
                rpc_status.to_status(
                    status_pb2.Status(
                        code=code_pb2.INTERNAL,
                        message="unable to enrich document, no models defined",
                    )
                )
            )

        if data.body == "" or len(data.body) < 24:
            return context.abort_with_status(
                rpc_status.to_status(
                    status_pb2.Status(
                        code=code_pb2.INTERNAL,
                        message="document body empty or too short",
                    )
                )
            )

        logging.info("Enrirch document with language model: {}".format(data.lang))

        # select language specific model (Model). language is defined
        # by the incoming request (field: lang).
        m = next(
            model for model in self.models if model.model.lang == data.lang.lower()
        )

        output = {
            "nlp": {
                "stopwords": [],
                "keywords": [],
                "entities": [],
                "summary": "",
                "topics": [],
                "claims": [],
                "quotes": [],
            }
        }

        # escape text
        body = html.unescape(data.body)

        try:
            # get the stopwords
            stopwords = extract_stopwords(body, list(m.spacy.Defaults.stop_words))
            output["nlp"]["stopwords"] = stopwords
            logging.debug("Stopwords: {}".format(stopwords))
        except Exception as err:
            logging.error("Error while getting stopwords: {}".format(str(err)))

        try:
            # classify the text usign a pretrained classifier
            topics = (
                extract_topics(body, m.topic_classification_pipeline)
                if m.topic_classification_pipeline != None
                else []
            )
            output["nlp"]["topics"] = topics
            logging.debug("Topics: {}".format(topics))
        except Exception as err:
            logging.error("Error while getting topics: {}".format(str(err)))

        try:
            # get the quotes
            quotes = extract_quotes(body)
            output["nlp"]["quotes"] = quotes
            logging.debug("Quotes: {}".format(quotes))
        except Exception as err:
            logging.error("Error while getting quotes: {}".format(str(err)))

        # parse the document usign a spacy model
        doc = m.spacy(body)
        logging.debug("Doc: {}".format(doc))

        try:
            # get the keywords
            keywords = extract_keywords(doc)
            output["nlp"]["keywords"] = keywords
            logging.debug("Keywords: {}".format(keywords))
        except Exception as err:
            logging.error("Error while getting keywords: {}".format(str(err)))

        try:
            # get the extracted entites
            entities = extract_entities(doc)
            output["nlp"]["entities"] = entities
            logging.debug("Entities: {}".format(entities))
        except Exception as err:
            logging.error("Error while getting entities: {}".format(str(err)))

        try:
            # generate a summary of the text
            summary = summarize_doc(doc, 3)
            output["nlp"]["summary"] = summary
            logging.debug("Summary: {}".format(summary))
        except Exception as err:
            logging.error("Error while getting summary: {}".format(str(err)))

        try:
            # get the claims (top 50%)
            claims = extract_claims(doc, list(m.spacy.Defaults.stop_words), 0.5)
            output["nlp"]["claims"] = claims
            logging.debug("Claims: {}".format(claims))
        except Exception as err:
            logging.error("Error while getting claims: {}".format(str(err)))

        if len(output["nlp"]["keywords"]) < 3 or len(stopwords) < 8:
            return context.abort_with_status(
                rpc_status.to_status(
                    status_pb2.Status(
                        code=code_pb2.INTERNAL,
                        message="document body too short",
                    )
                )
            )

        logging.info(
            "Keywords ({}), Stopwords ({}), Entities ({}), Topics ({}), Quotes ({}), Claims ({})".format(
                len(output["nlp"]["keywords"]),
                len(output["nlp"]["stopwords"]),
                len(output["nlp"]["entities"]),
                len(output["nlp"]["topics"]),
                len(output["nlp"]["quotes"]),
                len(output["nlp"]["claims"]),
            )
        )
        return enrich_pb2.EnrichResponse(code=200, status="success", data=output)


def main():
    """
    Start the gRPC service
    """

    # take environment variables from .env.
    load_dotenv(".env")

    # load conf from
    AppConfig(os.environ)

    # set log format and level
    logging.basicConfig(level=os.environ["LOG_LEVEL"], format=os.environ["LOG_FORMAT"])

    if os.environ["ENV"] != "development":
        nltk.download("punkt")

    """
    Load Models
    """
    for filename in glob.glob("models/*.json"):
        logging.info("Load model configuration {}".format(filename))

    models = [Model(filename) for filename in glob.glob("./models/*.json")]

    """
    Start GRPC Server
    """
    logging.info("Starting gRPC server")

    server = grpc.server(
        futures.ThreadPoolExecutor(max_workers=int(os.environ["MAX_WORKERS"]))
    )

    server.add_insecure_port("%s:%s" % (os.environ["HOST"], os.environ["PORT"]))
    enrich_pb2_grpc.add_EnrichServiceServicer_to_server(EnrichService(models), server)
    server.start()

    logging.info(
        "GRPC Server Listening on %s:%s", os.environ["HOST"], os.environ["PORT"]
    )

    try:
        while True:
            time.sleep(60 * 60)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == "__main__":
    main()
