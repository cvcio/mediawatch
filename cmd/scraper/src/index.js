require('dotenv').config();

const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const { MongoClient } = require('mongodb');
const { Kafka } = require('kafkajs');
const moment = require('moment');

const { Worker } = require('./worker');
const services = require('./services');
const logger = require('./logger');
const { mergeArticle } = require('./utils/strings');

let server;

// Connect to Kafka
const kafka = new Kafka({ brokers: process.env.KAFKA_BROKERS.split(',') });
const consumer = kafka.consumer({ groupId: 'mw-scraper' });
const producer = kafka.producer();
const worker = new Worker(consumer, producer);

const shutdown = (err) => {
	if (err) {
		logger.error(`[SVC-SCRAPER] gRPC server error: ${err.message}`);
	}

	if (server) {
		server.tryShutdown(() => {
			logger.info('[SVC-SCRAPER] gRPC server closed');
			process.exit(err ? 1 : 0);
		});
	} else {
		process.exit(err ? 1 : 0);
	}
};

const main = async () => {
	logger.info('[SVC-SCRAPER] Starting gRPC server');

	// Connect to mongo and retrieve passages
	const mongo = new MongoClient(process.env.MONGODB_URL);
	mongo.on('serverClosed', (event) => logger.info(`[SVC-SCRAPER] MongoDB connection closed: ${event.address}`));
	await mongo.connect();
	const db = mongo.db(process.env.MONGODB_DB);
	const collection = db.collection('passages');
	const passages = await collection.find({}).limit(5000).toArray();
	mongo.close();

	// Server Constructor
	server = new grpc.Server();

	// Load Scraper protobuf
	const packageDefinition = protoLoader
		.loadSync(`${process.env.PROTO_PATH}/scrape.proto`, {
			keepCase: true,
			longs: String,
			enums: String,
			defaults: true,
			oneofs: true
		});

	const scrapeProto = grpc.loadPackageDefinition(packageDefinition);
	const service = new services.ScrapeService(passages);

	// Add Services (Endpoints)
	server.addService(scrapeProto.mediawatch.scrape.v2.ScrapeService.service, service);
	server.bindAsync(process.env.SERVER_ADDRESS, grpc.ServerCredentials.createInsecure(), (err) => {
		if (err !== null) {
			shutdown(err);
		}

		server.start();
		logger.info(`[SVC-SCRAPER] gRPC server started at: ${process.env.SERVER_ADDRESS}`);
	});

	await worker.initialize();

	worker.consume('scrape', (message) => {
		const producedAt = moment.unix(message.timestamp / 1000);
		const start = moment();
		const request = JSON.parse(message.value.toString());
		service.Scrape({ request }, async (err, response) => {
			if (err || !response) {
				logger.error(`[SVC-SCRAPER] Error scraping: ${err.details}`);
				return;
			}

			const after = moment.duration(moment().diff(producedAt)).asMinutes();
			const took = moment.duration(moment().diff(start)).asSeconds();

			const article = mergeArticle(request, response.data);
			logger.debug(`[SVC-SCRAPER] Article scraped after ${after}m, took ${took}s (url: ${request.url})`);
			logger.debug(`[SVC-SCRAPER] Scraped article published at: ${article.content.published_at}, title: ${response.data.content.title}`);
			await worker.produce('enrich', [{ value: JSON.stringify(article) }]);
		});
	});
};

process.on('unhandledRejection', shutdown);
process.on('uncaughtException', shutdown);
process.on('SIGTERM', shutdown);

main();
