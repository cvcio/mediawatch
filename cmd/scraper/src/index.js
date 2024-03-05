require('dotenv').config();

const os = require('os');
const grpc = require('@grpc/grpc-js');
const { HealthImplementation } = require('grpc-health-check');
const { ScrapeServiceService } = require('@buf/cvcio_mediawatch.grpc_node/mediawatch/scrape/v2/scrape_grpc_pb');
const { MongoClient } = require('mongodb');
const { Kafka, logLevel } = require('kafkajs');
const moment = require('moment');

const { Worker } = require('./worker');
const services = require('./services');
const logger = require('./logger');
const { mergeArticle } = require('./utils/strings');

let server;
const statusMap = { '': 'SERVING', service: 'SERVING' };

// Connect to Kafka
const kafka = new Kafka({
	clientId: `${os.hostname}`,
	brokers: process.env.KAFKA_BROKERS.split(','),
	logLevel: logLevel.ERROR
});
const consumer = kafka.consumer({ groupId: 'mw-scraper' });
const producer = kafka.producer();
const worker = new Worker(consumer, producer);

const shutdown = async (err) => {
	if (err && err.message) {
		logger.error(`[SVC-SCRAPER] gRPC server error: ${err.message}`);
	}
	if (server) {
		server.tryShutdown(() => {
			logger.info('[SVC-SCRAPER] gRPC server closed');
		});
	}
	if (consumer) {
		await consumer.disconnect();
	}
	if (producer) {
		await producer.disconnect();
	}
	process.exit(err ? 1 : 0);
};

const main = async () => {
	logger.info('[SVC-SCRAPER] Starting gRPC server');

	// Connect to mongo and retrieve passages
	const mongo = new MongoClient(process.env.MONGODB_URL);
	mongo.on('serverClosed', event =>
		logger.info(`[SVC-SCRAPER] MongoDB connection closed: ${event.address}`));
	await mongo.connect();
	const db = mongo.db(process.env.MONGODB_DB);
	const collection = db.collection('passages');
	const passages = await collection.find({}).limit(5000).toArray();
	mongo.close();

	// Server Constructor
	server = new grpc.Server();

	const service = new services.ScrapeService(passages);
	const healthImpl = new HealthImplementation(statusMap);
	healthImpl.addToServer(server);

	// Add Services (Endpoints)
	server.addService(ScrapeServiceService, service);
	server.bindAsync(process.env.SERVER_ADDRESS,
		grpc.ServerCredentials.createInsecure(),
		err => {
			if (err !== null) {
				shutdown(err);
			}

			logger.info(`[SVC-SCRAPER] gRPC server started at: ${process.env.SERVER_ADDRESS}`);
		});

	await worker.initialize();

	worker.consume('scrape', async message => {
		const producedAt = moment.unix(message.timestamp / 1000);
		const start = moment();
		const request = JSON.parse(message.value.toString());
		try {
			const response = await services.Scrape({ request }, passages);
			const article = mergeArticle(request, response.data);

			const after = moment.duration(moment().diff(producedAt)).asMinutes();
			const took = moment.duration(moment().diff(start)).asSeconds();

			logger.info(`[SVC-SCRAPER] Article scraped after ${after}m, took ${took}s (id: ${article.doc_id})`);
			logger.debug(`[SVC-SCRAPER] Scraped article published at: ${article.content.published_at}, body: ${article.content.title} (id: ${article.doc_id})`);
			await worker.produce('enrich', [{ value: JSON.stringify(article) }]);
		} catch (err) {
			console.error(err);
			logger.error(`[SVC-SCRAPER] Error scraping: ${err}`);
		}
	});
};

main().catch(e => {
	logger.error(`[SVC-SCRAPER] main: ${e.message}`, e);
	process.exit(1);
});

process.on('unhandledRejection', async err => shutdown(err));
process.on('uncaughtException', async err => shutdown(err));
process.once('SIGTERM', async err => shutdown(err));
process.once('SIGINT', async err => shutdown(err));
