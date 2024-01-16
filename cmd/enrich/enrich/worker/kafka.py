import json
import logging

from aiokafka import AIOKafkaConsumer, AIOKafkaProducer

from mediawatch.enrich.v2 import enrich_pb2
from config.config import AppConfig
from google.protobuf.json_format import MessageToDict


class Worker:
    def __init__(self, config: AppConfig):
        self.config: AppConfig = config
        self.consuming: bool = False
        self.consumer: AIOKafkaConsumer = None
        self.producer: AIOKafkaProducer = None

    async def connect(self):
        await self.init_consumer()
        await self.init_producer()

    async def init_consumer(self):
        self.consumer = AIOKafkaConsumer(
            self.config.KAFKA_CONSUMER_TOPIC,
            bootstrap_servers=self.config.KAFKA_BOOTSTRAP_SERVERS,
            group_id=self.config.KAFKA_CONSUMER_GROUP_ID,
            auto_offset_reset=self.config.KAFKA_AUTO_OFFSET_RESET,
            enable_auto_commit=self.config.KAFKA_ENABLE_AUTO_COMMIT,
            value_deserializer=lambda m: json.loads(m.decode("utf-8")),
        )
        await self.consumer.start()

    async def init_producer(self):
        self.producer =  AIOKafkaProducer(
            bootstrap_servers=self.config.KAFKA_BOOTSTRAP_SERVERS,
            value_serializer=lambda m: json.dumps(m, ensure_ascii=False).encode("utf-8"),
        )
        await self.producer.start()

    async def run_consumer(self, callback, *args):
        self.consuming = True
        try:
            async for msg in self.consumer:
                await callback(msg, *args)
        except KeyboardInterrupt:
            logging.info("Stopping Kafka Consumer")
        finally:
            await self.stop()

    async def process(self, msg, method):
        try:
            nlp = await method(enrich_pb2.EnrichRequest(
                body = msg.value["content"]["body"],
                lang = msg.value["lang"].lower(),
            ), None)
            article = msg.value
            article["nlp"] = MessageToDict(nlp.data.nlp)
            await self.producer.send(self.config.KAFKA_PRODUCER_TOPIC, article)
        except Exception as e:
            logging.error("Error processing message: %s", e, stack_info=True)

    async def stop(self):
        if self.consuming:
            logging.info("Stopping Kafka Consumer")
            await self.consumer.stop()
            self.consuming = False
