"""
MediaWatch Enrich Service v2.0
Python implementation of the gRPC Enrich server.
"""

import os
import logging
import asyncio

import nltk
import uvloop

from dotenv import load_dotenv
from config.config import AppConfig

from ai.model import AIModel

from server.server import GRPCServer
from services.enrich import EnrichService
from mediawatch.enrich.v2 import enrich_pb2_grpc


async def main():
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

    # Load Models
    models = []
    for lang in env.SUPPORTED_LANGUAGES:
        if not os.path.exists(f"models/{lang}.json"):
            env.SUPPORTED_LANGUAGES.remove(lang)
            logging.warning("No model configuration found for lang: %s, omitting", lang)
        else:
            logging.info("Load model configuration for lang: %s", lang)
            models.append(AIModel(f"models/{lang}.json"))

    logging.info("Loaded %d models", len(models))

    # Start GRPC Server
    logging.info("Starting gRPC server")

    server = GRPCServer(env.HOST, env.PORT, env.MAX_WORKERS)
    server.register_service(
        enrich_pb2_grpc.add_EnrichServiceServicer_to_server, EnrichService, *models
    )

    await server.serve()


if __name__ == "__main__":
    uvloop.install()
    asyncio.run(main())
