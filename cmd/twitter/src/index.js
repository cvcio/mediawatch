require('dotenv').config();

const { Kafka } = require('kafkajs');
const url = require('url');
const { v4: uuidv4 } = require('uuid');
const { ETwitterStreamEvent, TwitterApi } = require('twitter-api-v2');
const logger = require('./logger');

const token = process.env.TWITTER_BEARER_TOKEN;
const clientId = process.env.KAFKA_CONSUMER_GROUP_WORKERS;
const brokers = [process.env.KAFKA_BROKER];
const topic = process.env.KAFKA_TOPIC_WORKER;

const twitterClient = new TwitterApi(token);
const kafka = new Kafka({ clientId, brokers });
const producer = kafka.producer();

const isValidUrl = (s) => {
	if (s.includes('twitter.com')) return false;
	if (s.includes('live24.gr')) return false;
	if (s.includes('indymedia.org')) return false;
	if (s.includes('youtube.com')) return false;
	if (s.includes('.jpg') || s.includes('.png')) return false;

	const parts = url.parse(s);
	if (parts.path.length <= 1) return false;

	return true;
};

const getUserNameFromTweet = (authorId, users) => {
	for (let i = 0; i < users.length; i++) {
		if (authorId === users[i].id) return users[i].username;
	}

	return '';
};

const getMessage = (data, url, username) => {
	if (process.env.VERSION === 'v1') {
		return {
			docId: uuidv4(),
			url,
			tweet_id: data.id,
			twitter_user_id: data.author_id,
			twitter_user_id_str: data.author_id.toString(),
			screen_name: username,
			user_name: username,
			created_at: new Date(data.created_at),
			created_at_str: data.created_at

		};
	}
	return {
		docId: uuidv4(),
		type: 'twitter',
		url,
		tweet_id: data.id,
		twitter_user_id: data.author_id,
		created_at: data.created_at,
		user_name: username
	};
};

const produce = async (msg) => {
	try {
		// send a message to the configured topic with
		// the key and value formed from the current value of `i`
		await producer.send({
			topic,
			messages: [
				{
					value: JSON.stringify(msg)
				},
			],
			acks: 0
		});

		logger.info(`New Tweet: ${msg.user_name} ${msg.url}`);
	} catch (err) {
		logger.error(`could not write message ${err}`);
	}
};

const main = async () => {
	logger.info('Starting Service');
	await producer.connect();

	const stream = await twitterClient.v2.searchStream({
		'tweet.fields':
			'created_at,id,author_id,lang,entities,in_reply_to_user_id'.split(','),
		'user.fields': 'id,name,profile_image_url,url,username,verified'.split(','),
		expansions: 'author_id,attachments.media_keys'.split(',')
	});

	stream.on(ETwitterStreamEvent.ConnectionError,
		err => logger.error('Connection error!', err));

	stream.on(ETwitterStreamEvent.ConnectionClosed,
		() => logger.info('Connection has been closed.'));

	stream.on(ETwitterStreamEvent.Data,
		data => {
			// console.debug(data);
			// stop
			if (!data.data) return;
			if (!data.data.entities) return;
			if (!data.data.entities.urls) return;
			if (!data.includes) return;
			if (data.data.in_reply_to_user_id) return;
			if (data.data.entities.urls.length <= 0) return;

			const matchingRules = data.matching_rules.map(m => m.tag);
			if (!matchingRules.includes(`mediawatch-listener-${process.env.VERSION}`)) return;

			// process
			const urls = JSON.parse(JSON.stringify(data.data.entities));
			for (let i = 0; i < urls.urls.length; i++) {
				const u = urls.urls[i];
				if (u.expanded_url && isValidUrl(u.expanded_url)) {
					const username = getUserNameFromTweet(data.data.author_id, data.includes.users);
					if (username !== '') {
						produce(getMessage(data.data, u.expanded_url, username));
					}
				}
			}
		});

	stream.on(ETwitterStreamEvent.DataKeepAlive,
		() => logger.debug('Twitter has a keep-alive packet.'));
};

const shutdown = err => {
	if (err) {
		logger.error(`Service error: ${err.message}`);
	}
	process.exit(err ? 1 : 0);
};

process.on('unhandledRejection', shutdown);
process.on('uncaughtException', shutdown);
process.on('SIGTERM', shutdown);

main(0);
