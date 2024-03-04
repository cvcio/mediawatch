require('dotenv').config({ path: '.env' });

const moment = require('moment');
const axios = require('axios');
const logger = require('../src/logger');
const { Scrape } = require('../src/services');

const test_feed_scrapability_async = async f => {
	const req = {
		request: {
			doc_id: 'test',
			feed: JSON.stringify(f),
			url: f.test.url,
			screen_name: f.user_name,
			lang: 'el',
			crawled_at: moment().format('YYYY-MM-DDTHH:mm:ss')
		}
	};
	try {
		const data = await Scrape(req, []);
		return { hostname: f.hostname, status: 'ok', message: data };
	} catch (e) {
		console.error(e);
		return {
			hostname: f.hostname,
			status: 'error',
			message: e.message,
			code: e.code
		};
	}
};

const test = async () => {
	const options = {
		method: 'POST',
		url: 'http://localhost:8000/mediawatch.feeds.v2.FeedService/GetFeeds',
		headers: { 'Content-Type': 'application/json' },
		data: { limit: 3000 }
	};

	const results = [];
	const feeds = await axios.request(options);

	// eslint-disable-next-line no-restricted-syntax
	for (const feed of feeds.data.data) {
		if (feed.test && feed.test.url !== '' && results.length < 2000) {
			if (feed.hostname === 'efsyn.gr') {
				// feed.stream.requires_proxy = true;
				results.push(test_feed_scrapability_async(feed));
			}
		}
	}

	const done = await Promise.all(results);
	// eslint-disable-next-line no-restricted-syntax
	for (const result of done) {
		if (result.status === 'error') {
			logger.info(`${result.code} - ${result.hostname} - ${result.message}`);
		} else {
			logger.info(`200 - ${result.hostname}`);
			console.log(result.message);
		}
	}
};

test().catch(e => {
	console.error(e);
	process.exit(1);
});
