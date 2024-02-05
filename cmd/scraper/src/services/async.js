const moment = require('moment');
const { extract } = require('ascraper'); // require('/home/andefined/js/misc/npm/ascraper/lib'); //

const logger = require('../logger');
const { parsers } = require('../parsers');

const { toUpperCase, normalizeString, trimRight } = require('../utils/strings');
// const { errorCode } = require('../utils/errors');
const { getProxyUrl } = require('../utils/proxy');

moment.suppressDeprecationWarnings = true;

// Async Scrape scrapes a given URL and returns the extracted data.
// It also checks if the feed has a custom parser, located under parsers package.
const Scrape = async (req, passages) => {
	// eslint-disable-next-line prefer-destructuring
	const request = req.request;
	const feed =		typeof request.feed === 'string'
		? JSON.parse(request.feed)
		: request.feed;

	logger.debug(`[SVC-SCRAPER] Scrape - (${feed.hostname}) ${decodeURIComponent(request.url).toString()}`);

	const proxy = feed.stream.requires_proxy ? getProxyUrl() : null;

	// Check if feed has a custom parser
	if (feed.userName && parsers[feed.userName.toLowerCase()] !== undefined) {
		const parser = parsers[feed.userName.toLowerCase()];
		const url = parser.url(request.url);

		if (!url) {
			logger.error(`[SVC-SCRAPER] Unable to scrape URL (url error) (${feed.hostname}) ${request.url}`);
			throw new Error(`Unable to scrape URL (url error) (${feed.hostname}) ${request.url}`);
		}

		const res = await parser.fetchAPI(url);
		const article = res;

		const crawled_at_utc = moment.utc(request.crawled_at);
		const publishedAt_utc = moment.utc(article.publishedAt);

		if (article.publishedAt && crawled_at_utc.isBefore(publishedAt_utc)) {
			article.publishedAt = request.crawled_at;
		}

		if (article.text === '' || article.title === '') {
			throw new Error(`Unable to scrape URL (empty response) (${feed.hostname}) ${request.url}`);
		}

		return {
			status: 'success',
			code: 200,
			message: '',
			data: {
				content: {
					title: article.title.trim(),
					body: normalizeString(article.body),
					authors:
						typeof article.authors === 'string'
						&& article.authors.length > 2
							? toUpperCase(article.authors)
								.split(',')
								.map(m => m.trim())
							: [],
					published_at: article.publishedAt
						? moment(article.publishedAt).format('YYYY-MM-DDTHH:mm:ssZZ')
						: moment(request.crawled_at).format('YYYY-MM-DDTHH:mm:ssZZ'),
					tags:
						typeof article.tags === 'string'
						&& article.tags.length > 2
							? toUpperCase(article.tags)
								.split(',')
								.map(m => m.trim())
							: [],
					description: normalizeString(article.description),
					image: article.image || ''
				}
			}
		};
	}

	const url = decodeURIComponent(request.url).toString();
	if (!url) {
		logger.error(`[SVC-SCRAPER] Unable to scrape URL (url error) (${feed.hostname}) ${request.url}`);
		throw new Error(`Unable to scrape URL (url error) (${feed.hostname}) ${request.url}`);
	}

	let res = null;
	try {
		res = await extract(url, proxy);
	} catch (e) {
		if (!feed.stream.requires_proxy && e.code === 403 && !req.retry) {
			feed.stream.requires_proxy = true;
			logger.warn(`[SVC-SCRAPER] Proxy required for (${feed.hostname}) ${request.url}`);
			return Scrape({ request, feed, retry: true }, passages);
		}
		logger.error(`[SVC-SCRAPER] Unable to scrape URL (extract error) (${feed.hostname}) ${e.code} ${e.message}`);
		throw e;
	}

	const article = res;
	if (article.date && moment(request.crawled_at).isBefore(moment(article.date))) {
		article.date = request.crawled_at;
	}

	if (article.text === '' || article.title === '') {
		throw new Error(`Unable to scrape URL (empty response) (${feed.hostname}) ${request.url}`);
	}

	article.text = trimRight(article.text, passages);

	return {
		status: 'success',
		code: 200,
		message: '',
		data: {
			content: {
				title: article.title,
				body: normalizeString(article.text),
				authors:
					typeof article.author === 'string'
					&& article.author.length > 2
						? toUpperCase(article.author)
							.split(',')
							.map(m => m.trim())
						: [],
				published_at: article.date
					? moment(article.date).format('YYYY-MM-DDTHH:mm:ssZZ')
					: moment(request.crawled_at).format('YYYY-MM-DDTHH:mm:ssZZ'),
				tags:
					typeof article.keywords === 'string'
					&& article.keywords.length > 2
						? toUpperCase(article.keywords)
							.split(',')
							.map(m => m.trim())
						: [],
				description: normalizeString(article.description) || '',
				image: article.image || ''
			}
		}
	};
};

module.exports = Scrape;
