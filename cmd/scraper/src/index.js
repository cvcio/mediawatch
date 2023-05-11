require('dotenv').config();

const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const { MongoClient } = require('mongodb');
const services = require('./services');
const logger = require('./logger');
const { getProxy } = require('./utils/proxy');

let server;

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
	const mongo = new MongoClient(process.env.MONGODB_URL);
	await mongo.connect();

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
	const service = new services.ScrapeService(mongo, getProxy());

	// Add Services (Endpoints)
	server.addService(scrapeProto.mediawatch.scrape.v2.ScrapeService.service, service);
	server.bindAsync(process.env.SERVER_ADDRESS, grpc.ServerCredentials.createInsecure(), (err) => {
		if (err !== null) {
			shutdown(err);
		}

		server.start();
		logger.info(`[SVC-SCRAPER] gRPC server started at: ${process.env.SERVER_ADDRESS}`);
	});
};

process.on('unhandledRejection', shutdown);
process.on('uncaughtException', shutdown);
process.on('SIGTERM', shutdown);

main();
