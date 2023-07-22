"""
MediaWatch Enrich Service v2.0
Python implementation of the gRPC Enrich server.
"""

import os
import logging
import grpc
import time
import glob
import html
import nltk

from concurrent import futures

from dotenv import load_dotenv
from config.config import AppConfig
from nlp.model import Model

from mediawatch.enrich.v2 import enrich_pb2_grpc
from mediawatch.enrich.v2 import enrich_pb2

from google.rpc import code_pb2
from google.rpc import status_pb2
from grpc_status import rpc_status

from nlp.nlp import (
    extract_stopwords,
    extract_keywords,
    extract_entities,
    summarize_doc,
    extract_topics,
    extract_quotes,
    extract_claims,
    extract_named_entities,
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

        logging.info(
            "Retrieve StopWords from document with language model: {}".format(data.lang)
        )

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

        logging.info(
            "Retrieve keywords from document with language model: {}".format(data.lang)
        )

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

        logging.info(
            "Retrieve summary from document with language model: {}".format(data.lang)
        )

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

        logging.info(
            "Retrieve topics from document with language model: {}".format(data.lang)
        )

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
            "Topics ({}) {}".format(len(output["nlp"]["topics"]), ", ".join(topics))
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

        logging.info(
            "Retrieve Quotes from document with language model: {}".format(data.lang)
        )

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

        logging.info(
            "Retrive claims from document with language model: {}".format(data.lang)
        )

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

        logging.debug("Enrirch document with language model: {}".format(data.lang))

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
            # classify the text usign a pretrained classifier
            named_entities = (
                extract_named_entities(body, m.ner_classification_pipeline)
                if m.ner_classification_pipeline != None
                else []
            )
            # output["nlp"]["topics"] = topics
            logging.info("Named Entities: {}".format(named_entities))
        except Exception as err:
            logging.error("Error while getting named entities: {}".format(str(err)))

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
            logging.info("Entities: {}".format(entities))
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
    env = AppConfig(os.environ)

    # set log format and level
    logging.basicConfig(level=env.LOG_LEVEL, format=env.LOG_FORMAT)

    if env.ENV != "development":
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

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=int(env.MAX_WORKERS)))

    server.add_insecure_port("%s:%s" % (env.HOST, env.PORT))
    enrich_pb2_grpc.add_EnrichServiceServicer_to_server(EnrichService(models), server)
    server.start()

    logging.info("GRPC Server Listening on %s:%s", env.HOST, env.PORT)

    try:
        while True:
            time.sleep(60 * 60)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == "__main__":
    main()
