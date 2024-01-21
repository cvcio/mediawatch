const logger = require('../logger');

class Worker {
	constructor (consumer, producer) {
		this.consumer = consumer;
		this.producer = producer;
	}

	async initialize () {
		try {
			await this.consumer.connect();
			await this.producer.connect();
		} catch (err) {
			logger.error(`[SVC-SCRAPER] Kafka connection error: ${err.message}`);
			process.exit(1);
		}
	}

	async consume (topic, callback) {
		await this.consumer.subscribe({ topic, fromBeginning: false });
		await this.consumer.run({
			autoCommit: true,
			eachMessage: async ({ topic, message }) => {
				callback(message);

				this.consumer.pause([{ topic }]);
				setTimeout(() => {
					this.consumer.resume([{ topic }]);
				}, 10);
			}
		});
	}

	async produce (topic, messages) {
		await this.producer.send({ topic, messages });
	}
}

module.exports = { Worker };
