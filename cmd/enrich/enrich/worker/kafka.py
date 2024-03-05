import json
import logging
import socket

from asyncio import AbstractEventLoop
from aiokafka import AIOKafkaConsumer, AIOKafkaProducer

from mediawatch.enrich.v2 import enrich_pb2
from config.config import AppConfig

from google.protobuf.json_format import MessageToDict
from grpc_interceptor.exceptions import NotFound, Internal


class Worker:
    def __init__(self, config: AppConfig, loop: AbstractEventLoop):
        self.config: AppConfig = config
        self.consuming: bool = False

        # TODO: Investigate the following error:
        # Commit cannot be completed since the group has already rebalanced and assigned the partitions to another member.
        # This means that the time between subsequent calls to poll() was longer than the configured max_poll_interval_ms,
        # which typically implies that the poll loop is spending too much time message processing.
        # You can address this either by increasing the rebalance timeout with max_poll_interval_ms,
        # or by reducing the maximum size of batches returned in poll() with max_poll_records.
        self.consumer: AIOKafkaConsumer = AIOKafkaConsumer(
            self.config.KAFKA_CONSUMER_TOPIC,
            bootstrap_servers=self.config.KAFKA_BOOTSTRAP_SERVERS,
            client_id=socket.getfqdn(),
            group_id=self.config.KAFKA_CONSUMER_GROUP_ID,
            fetch_max_bytes=1024 * 1024 * 2,
            enable_auto_commit=self.config.KAFKA_ENABLE_AUTO_COMMIT,
            value_deserializer=lambda m: json.loads(m.decode("utf-8")),
            max_poll_interval_ms=45000,
            max_poll_records=4,
            session_timeout_ms=30000,
            loop=loop,
        )

        self.producer: AIOKafkaProducer = AIOKafkaProducer(
            client_id=socket.getfqdn(),
            bootstrap_servers=self.config.KAFKA_BOOTSTRAP_SERVERS,
            value_serializer=lambda m: json.dumps(m, ensure_ascii=False).encode(
                "utf-8"
            ),
            loop=loop,
        )
    async def connect(self):
        await self.init_producer()

    async def init_producer(self):
        await self.producer.start()

    async def run_consumer(self, callback, *args):
        await self.consumer.start()
        self.consuming = True
        try:
            async for msg in self.consumer:
                try:
                    await callback(msg, *args)
                except Exception as e:
                    logging.error("Error processing message (run_consumer): %s", e)
        except KeyboardInterrupt:
            logging.info("Stopping Kafka Consumer")
        finally:
            await self.stop()

    async def process(self, msg, method):
        try:
            article = msg.value
            req = enrich_pb2.EnrichRequest(
                body=article["content"]["body"], lang=article["lang"].lower()
            )
            logging.debug("Processing message: %s", req)
            nlp = await method(req, None)
            article["nlp"] = MessageToDict(nlp.data.nlp)
            await self.producer.send(self.config.KAFKA_PRODUCER_TOPIC, article)
        except (NotFound, Internal) as e:
            logging.error("Error processing message (process): %s", e)
        except (KeyError, ValueError):
            logging.error("Malformed message: %s", msg.value)

    async def stop(self):
        if self.consuming:
            logging.info("Stopping Kafka Consumer")
            await self.consumer.stop()
            self.consuming = False
