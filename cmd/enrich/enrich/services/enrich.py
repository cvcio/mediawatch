"""Enrich Service Module"""

import logging
import html

from ai.model import AIModel
from nlp.methods import (
    extract_stopwords,
    extract_keywords,
    extract_entities,
    extract_topics,
    extract_quotes,
    extract_claims,
    summarize,
)

from mediawatch.enrich.v2.enrich_pb2 import EnrichRequest, EnrichResponse
from mediawatch.enrich.v2.enrich_pb2_grpc import EnrichServiceServicer

from grpc_interceptor.exceptions import NotFound, Internal


class NLPException(Exception):
    """Custom exception for NLP errors"""


class EnrichService(EnrichServiceServicer):
    """
    EnrichService stub implementation of the gRPC EnrichService servicer.
    """

    def __init__(self, models: list[AIModel]) -> None:
        """
        initialize gRPC stub
        """
        # add the models inside the stub
        self.models = models

    def _get_model_by_lang(self, lang: str) -> AIModel:
        return next(
            (model for model in self.models if model.lang == lang.lower()), None
        )

    # pylint: disable-next=invalid-overridden-method
    async def StopWords(self, request: EnrichRequest, context):
        """
        StopWords gRPC endpoint
        """
        if self.models is None:
            # if there are no models return an error
            raise NotFound("unable to enrich document, no models defined")

        if request.body == "" or len(request.body) < 24:
            raise Internal("document body empty or too short")

        logging.debug(
            "Retrieve StopWords from document with language model: %s", request.lang
        )

        # select language specific model (Model). language is defined
        # by the incoming request (field: lang).
        base_model = self._get_model_by_lang(request.lang)

        output = {
            "nlp": {
                "stopwords": [],
            }
        }

        # escape text
        body = html.unescape(request.body)

        try:
            # get the stopwords
            stopwords = await extract_stopwords(
                body, list(base_model.spacy.Defaults.stop_words)
            )
            output["nlp"]["stopwords"] = stopwords
            logging.debug("Stopwords: %s", stopwords)
        except NLPException as err:
            logging.error("Error while getting stopwords: %s", err)
            raise Internal(f"An error occurred: {type(err).__name__} - {err}") from err

        logging.info("Stopwords (%d)", len(output["nlp"]["stopwords"]))
        return EnrichResponse(code=200, status="success", data=output)

    # pylint: disable-next=invalid-overridden-method
    async def Keywords(self, request: EnrichRequest, context):
        """
        Keywords gRPC endpoint
        """
        if self.models is not None:
            # if there are no models return an error
            raise NotFound("unable to enrich document, no models defined")

        if request.body == "" or len(request.body) < 24:
            raise Internal("document body empty or too short")

        logging.info(
            "Retrieve keywords from document with language model: %s", request.lang
        )

        # select language specific model (Model). language is defined
        # by the incoming request (field: lang).
        base_model = self._get_model_by_lang(request.lang)

        output = {
            "nlp": {
                "keywords": [],
            }
        }

        # escape text
        body = html.unescape(request.body)

        # parse the document usign a spacy model
        doc = base_model.spacy(body)
        logging.debug("Doc: %s", doc)

        try:
            # get the keywords
            keywords = await extract_keywords(doc)
            output["nlp"]["keywords"] = keywords
            logging.debug("Keywords: %s", keywords)
        except NLPException as err:
            logging.error("Error while getting keywords: %s", err)
            raise Internal(f"An error occurred: {type(err).__name__} - {err}") from err

        logging.info("Keywords (%d)", len(output["nlp"]["keywords"]))
        return EnrichResponse(code=200, status="success", data=output)

    # pylint: disable-next=invalid-overridden-method
    async def Entities(self, request: EnrichRequest, context):
        """
        Entities gRPC endpoint
        """
        if self.models is None:
            # if there are no models return an error
            raise NotFound("unable to enrich document, no models defined")

        if request.body == "" or len(request.body) < 24:
            raise Internal("document body empty or too short")

        logging.info("Enrich document with language model: %s", request.lang)

        # select language specific model (Model). language is defined
        # by the incoming request (field: lang).
        base_model = self._get_model_by_lang(request.lang)

        output = {
            "nlp": {
                "entities": [],
            }
        }

        # escape text
        body = html.unescape(request.body)

        # parse the document usign a spacy model
        doc = base_model.spacy(body)
        logging.debug("Doc: %s", doc)

        try:
            # get the extracted entites
            entities = await extract_entities(doc)
            output["nlp"]["entities"] = entities
            logging.debug("Entities: %s", entities)
        except NLPException as err:
            logging.error("Error while getting entities: %s", err)
            raise Internal(f"An error occurred: {type(err).__name__} - {err}") from err

        logging.info("Entities (%d)", len(output["nlp"]["entities"]))
        return EnrichResponse(code=200, status="success", data=output)

    # pylint: disable-next=invalid-overridden-method
    async def Summary(self, request: EnrichRequest, context):
        """
        Summary gRPC endpoint
        """
        if self.models is not None:
            # if there are no models return an error
            raise NotFound("unable to enrich document, no models defined")

        if request.body == "" or len(request.body) < 24:
            raise Internal("document body empty or too short")

        logging.info(
            "Retrieve summary from document with language model: %s", request.lang
        )

        # select language specific model (Model). language is defined
        # by the incoming request (field: lang).
        base_model = self._get_model_by_lang(request.lang)

        output = {
            "nlp": {
                "summary": "",
            }
        }

        # escape text
        body = html.unescape(request.body)

        # parse the document usign a spacy model
        doc = base_model.spacy(body)
        logging.debug("Doc: %s", doc)

        try:
            # generate a summary of the text
            summary = await summarize(doc, 3)
            output["nlp"]["summary"] = summary
            logging.debug("Summary: %s", summary)
        except NLPException as err:
            logging.error("Error while getting summary: %s", err)
            raise Internal(f"An error occurred: {type(err).__name__} - {err}") from err

        return EnrichResponse(code=200, status="success", data=output)

    # pylint: disable-next=invalid-overridden-method
    async def Topics(self, request: EnrichRequest, context):
        """
        Topics gRPC endpoint
        """
        if self.models is None:
            # if there are no models return an error
            raise NotFound("unable to enrich document, no models defined")

        if request.body == "" or len(request.body) < 24:
            raise Internal("document body empty or too short")

        logging.info(
            "Retrieve topics from document with language model: %s", request.lang
        )

        # select language specific model (Model). language is defined
        # by the incoming request (field: lang).
        base_model = self._get_model_by_lang(request.lang)

        output = {
            "nlp": {
                "topics": [],
            }
        }

        # escape text
        body = html.unescape(request.body)

        try:
            # classify the text usign a pretrained classifier
            topics = (
                await extract_topics(body, base_model.topic_classification_pipeline)
                if base_model.topic_classification_pipeline is not None
                else []
            )
            output["nlp"]["topics"] = topics
            logging.debug("Topics: %s", topics)
        except NLPException as err:
            logging.error("Error while getting topics: %s", err)
            raise Internal(f"An error occurred: {type(err).__name__} - {err}") from err

        logging.info("Topics (%s) %s", len(output["nlp"]["topics"]), ", ".join(topics))
        return EnrichResponse(code=200, status="success", data=output)

    # pylint: disable-next=invalid-overridden-method
    async def Quotes(self, request: EnrichRequest, context):
        """
        Quotes gRPC endpoint
        """

        if request.body == "" or len(request.body) < 24:
            raise Internal("document body empty or too short")

        logging.info(
            "Retrieve Quotes from document with language model: %s", request.lang
        )

        output = {
            "nlp": {
                "quotes": [],
            }
        }

        # escape text
        body = html.unescape(request.body)

        try:
            # get the quotes
            quotes = await extract_quotes(body)
            output["nlp"]["quotes"] = quotes
            logging.debug("Quotes: %s", quotes)
        except NLPException as err:
            logging.error("Error while getting quotes: %s", err)
            raise Internal(f"An error occurred: {type(err).__name__} - {err}") from err

        logging.info("Quotes (%s)", len(output["nlp"]["quotes"]))
        return EnrichResponse(code=200, status="success", data=output)

    # pylint: disable-next=invalid-overridden-method
    async def Claims(self, request: EnrichRequest, context):
        """
        Claims gRPC endpoint
        """
        if self.models is None:
            # if there are no models return an error
            raise NotFound("unable to enrich document, no models defined")

        if request.body == "" or len(request.body) < 24:
            raise Internal("document body empty or too short")

        logging.info(
            "Retrive claims from document with language model: %s", request.lang
        )

        # select language specific model (Model). language is defined
        # by the incoming request (field: lang).
        base_model = self._get_model_by_lang(request.lang)

        output = {
            "nlp": {
                "claims": [],
            }
        }

        # escape text
        body = html.unescape(request.body)

        # parse the document usign a spacy model
        doc = base_model.spacy(body)
        logging.debug("Doc: %s", doc)

        try:
            # get the claims (top 50%)
            claims = await extract_claims(
                doc, list(base_model.spacy.Defaults.stop_words), 0.5
            )
            output["nlp"]["claims"] = claims
            logging.debug("Claims: %s", claims)
        except NLPException as err:
            logging.error("Error while getting claims: %s", err)
            raise Internal(f"An error occurred: {type(err).__name__} - {err}") from err

        logging.info("Claims (%d)", len(output["nlp"]["claims"]))
        return EnrichResponse(code=200, status="success", data=output)

    # pylint: disable-next=invalid-overridden-method
    async def NLP(self, request: EnrichRequest, context):
        """
        NLP gRPC endpoint
        """

        if self.models is None:
            # if there are no models return an error
            raise NotFound("unable to enrich document, no models defined")

        if request.body == "" or len(request.body) < 24:
            raise Internal("document body empty or too short")

        logging.debug("Enrich document with language model: %s", request.lang)

        # select language specific model (Model). language is defined
        # by the incoming request (field: lang).
        base_model = self._get_model_by_lang(request.lang)

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
        body = html.unescape(request.body)
        stopwords = []

        try:
            # get the stopwords
            stopwords = await extract_stopwords(
                body, list(base_model.spacy.Defaults.stop_words)
            )
            output["nlp"]["stopwords"] = stopwords
            logging.debug("Stopwords: %s", stopwords)
        except NLPException as err:
            logging.error("Error while getting stopwords: %s", err)

        try:
            # classify the text usign a pretrained classifier
            topics = (
                await extract_topics(body, base_model.topic_classification_pipeline)
                if base_model.topic_classification_pipeline is not None
                else []
            )
            output["nlp"]["topics"] = topics
            logging.debug("Topics: %s", topics)
        except NLPException as err:
            logging.error("Error while getting topics: %s", err)

        try:
            # get the quotes
            quotes = await extract_quotes(body)
            output["nlp"]["quotes"] = quotes
            logging.debug("Quotes: %s", quotes)
        except NLPException as err:
            logging.error("Error while getting quotes: %s", err)

        # parse the document usign a spacy model
        doc = base_model.spacy(body)
        logging.debug("Doc: %s", doc)

        try:
            # get the keywords
            keywords = await extract_keywords(doc)
            output["nlp"]["keywords"] = keywords
            logging.debug("Keywords: %s", keywords)
        except NLPException as err:
            logging.error("Error while getting keywords: %s", err)

        try:
            # get the extracted entites
            entities = await extract_entities(doc)
            output["nlp"]["entities"] = entities
            logging.debug("Entities: %s", entities)
        except NLPException as err:
            logging.error("Error while getting entities: %s", err)

        try:
            # generate a summary of the text
            summary = await summarize(doc, 3)
            output["nlp"]["summary"] = summary
            logging.debug("Summary: %s", summary)
        except NLPException as err:
            logging.error("Error while getting summary: %s", err)

        try:
            # get the claims (top 50%)
            claims = await extract_claims(
                doc, list(base_model.spacy.Defaults.stop_words), 0.5
            )
            output["nlp"]["claims"] = claims
            logging.debug("Claims: %s", claims)
        except NLPException as err:
            logging.error("Error while getting claims: %s", err)

        if len(output["nlp"]["keywords"]) < 3 or len(stopwords) < 8:
            raise Internal("document body empty or too short")

        logging.info(
            "Entities (%s), Topics (%s), Quotes (%s), Claims (%s)",
            len(output["nlp"]["entities"]),
            len(output["nlp"]["topics"]),
            len(output["nlp"]["quotes"]),
            len(output["nlp"]["claims"]),
        )

        return EnrichResponse(code=200, status="success", data=output)
