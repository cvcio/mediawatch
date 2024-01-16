"""
MediaWatch Enrich Service v2.0
Python implementation of the gRPC Enrich server.
"""

import os
import logging
import asyncio
import signal

import nltk
import uvloop

from dotenv import load_dotenv
from config.config import AppConfig

from ai.model import AIModel

from server.server import GRPCServer
from services.enrich import EnrichService
from mediawatch.enrich.v2 import enrich_pb2_grpc, enrich_pb2
from worker.kafka import Worker
import inspect


async def process(msg, method):
    try:
        a = enrich_pb2.EnrichRequest(
            body=msg.value["content"]["body"], lang=msg.value["lang"].lower()
        )
        response = await method(a, None)
    except Exception as e:
        logging.error(e, stack_info=True)
        pass


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

    enrich_service = EnrichService(models)

    worker = Worker(env)
    await worker.connect()

    server = GRPCServer(env.HOST, env.PORT, env.MAX_WORKERS)
    server.register_service_method(
        enrich_pb2_grpc.add_EnrichServiceServicer_to_server, enrich_service
    )

    def on_signal_exit():
        logging.info("Received exit signal")
        asyncio.create_task(server.stop())
        asyncio.create_task(worker.stop())

    loop = asyncio.get_running_loop()

    loop.add_signal_handler(signal.SIGTERM, on_signal_exit)
    loop.add_signal_handler(signal.SIGINT, on_signal_exit)

    # loop.create_task(worker.run_consumer(process, enrich_service.NLP))
    # Start GRPC Server
    # logging.info("Starting gRPC server")
    # await server.serve()

    await asyncio.gather(
        worker.run_consumer(worker.process, enrich_service.NLP), server.serve()
    )


if __name__ == "__main__":
    uvloop.install()
    asyncio.run(main())
