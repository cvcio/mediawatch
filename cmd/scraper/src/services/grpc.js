const moment = require('moment');
const { extract } = require('ascraper'); // require('/home/andefined/js/misc/npm/ascraper/lib'); //
const logger = require('../logger');
const { parsers } = require('../parsers');
const { toUpperCase, normalizeString, trimRight } = require('../utils/strings');
const { errorCode } = require('../utils/errors');
const { getProxyUrl } = require('../utils/proxy');

moment.suppressDeprecationWarnings = true;

class ScrapeService {
	constructor (passages) {
		this.passages = passages.filter(m => m.type === 'trim');
		logger.info(`[SVC-SCRAPER] (${this.passages.length}) passages loaded`);
	}

	Scrape (req, callback) {
		const { request } = req;
		const feed =			typeof request.feed === 'string'
			? JSON.parse(request.feed)
			: request.feed;

		logger.info(`[SVC-SCRAPER] Scrape - (${feed.hostname}) ${decodeURIComponent(request.url).toString()}`);

		if (
			request.screen_name
			&& parsers[request.screen_name.toLowerCase()] !== undefined
		) {
			const parser = parsers[request.screen_name.toLowerCase()];
			const url = parser.url(request.url);

			if (!url) {
				logger.error(`[SVC-SCRAPER] Unable to scrape URL (url error) (${feed.hostname}) ${request.url}`);
				return callback({
					code: 9,
					details: `Unable to scrape URL (url error) (${feed.hostname}) ${request.url}`
				},
				null);
			}

			parser
				.fetchAPI(url)
				.then(res => {
					if (!res) {
						logger.error(`[SVC-SCRAPER] Unable to scrape URL (response empty) (${feed.hostname}) ${request.url}`);
						return callback({
							code: 9,
							details: `Unable to scrape URL (response empty) (${feed.hostname}) ${request.url}`
						},
						null);
					}
					// logger.debug(`[SVC-SCRAPER] Data ${JSON.stringify(res)}`);

					const article = res;

					if (
						article.publishedAt
						&& moment
							.utc(request.crawled_at)
							.isBefore(moment.utc(article.publishedAt))
					) {
						article.publishedAt = request.crawled_at;
					}

					if (article.text === '' || article.title === '') {
						logger.error(`[SVC-SCRAPER] Unable to scrape URL (malformed data) (${feed.hostname}) ${request.url}`);
						return callback({
							code: 9,
							details: `Unable to scrape URL (malformed data) (${feed.hostname}) ${request.url}`
						},
						null);
					}

					return callback(null, {
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
					});
				})
				.catch(err => {
					logger.error(`[SVC-SCRAPER] Error while scraping: (${errorCode(err.response ? err.response.status : 500)}) ${err.message} - (${feed.hostname}) ${request.url}`);
					return callback({
						code: errorCode(err.response ? err.response.status : 500),
						details: err.message
					},
					null);
				});
		} else {
			extract(decodeURIComponent(request.url).toString(),
				feed.stream.requires_proxy ? getProxyUrl() : null)
				.then(res => {
					const article = res;
					if (
						article.date
						&& moment(request.crawled_at).isBefore(moment(article.date))
					) {
						article.date = request.crawled_at;
					}

					if (article.text === '' || article.title === '') {
						logger.error(`[SVC-SCRAPER] Unable to scrape URL (malformed data) (${feed.hostname}) ${request.url}`);
						return callback({
							code: 9,
							details: `Unable to scrape URL (malformed data) (${feed.hostname}) ${request.url}`
						},
						null);
					}
					article.text = trimRight(article.text, this.passages);
					return callback(null, {
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
								description: normalizeString(article.description),
								image: article.image || ''
							}
						}
					});
				})
				.catch(err => {
					if (
						err.code
						&& err.code === 403
						&& getProxyUrl()
					) {
						return this.RetryWithProxy(request, feed, callback);
					}
					logger.error(`[SVC-SCRAPER] Error while scraping: (${errorCode(err.code ? err.code : 500)}) ${err.message} - (${feed.hostname}) ${request.url}`);
					return callback({
						code: errorCode(err.response ? err.response.status : 500),
						details: err.message
					},
					null);
				});
		}
	}

	SimpleScrape (req, callback) {
		const now = moment();
		const { request } = req;
		request.feed = JSON.parse(request.feed);
		logger.info(`[SVC-SCRAPER] SimpleScrape - (${
			request.feed.screen_name
		}) ${decodeURIComponent(request.url).toString()}`);

		if (parsers[request.feed.screen_name] !== undefined) {
			const parser = parsers[request.feed.screen_name.toLowerCase()];
			const url = parser.url(request.url);

			if (!url) {
				logger.error(`[SVC-SCRAPER] URL Error (${request.feed.screen_name}) ${request.url}`);
				return callback(Error(`URL Error (${request.feed.screen_name}) ${request.url}`),
					null);
			}

			parser
				.fetchAPI(url)
				.then(res => {
					if (!res) {
						logger.error(`[SVC-SCRAPER] Unable to scrape URL (${request.feed.screen_name}) ${request.url}`);
						return callback(new Error(`Unable to scrape URL (${request.feed.screen_name}) ${request.url}`),
							null);
					}
					logger.debug(`[SVC-SCRAPER] Data ${JSON.stringify(res)}`);

					const article = res;

					if (article.text === '' || article.title === '') {
						logger.error(`[SVC-SCRAPER] Unable to scrape URL (${request.feed.screen_name}) ${request.url}`);
						return callback(new Error(`Unable to scrape URL (${request.feed.screen_name}) ${request.url}`),
							null);
					}

					return callback(null, {
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
					});
				})
				.catch(err => {
					logger.error('Error while scraping:',
						err.message,
						`(${request.feed.screen_name}) ${request.url}`);
					return callback({
						code: errorCode(err.response ? err.response.status : 500),
						details: err.message
					},
					null);
				});
		} else {
			// Run Scraper
			extract(decodeURIComponent(request.url).toString())
				.then(res => {
					const article = res;
					if (article.date && now.isBefore(moment(article.date))) {
						article.date = null;
					}
					if (article.body === '' || article.title === '') {
						logger.error(`[SVC-SCRAPER] Unable to scrape URL (${request.feed.screen_name}) ${request.url}`);
						return callback(new Error(`Unable to scrape URL (${request.feed.screen_name}) ${request.url}`),
							null);
					}

					article.text = trimRight(article.text, this.passages);
					return callback(null, {
						status: 'success',
						code: 200,
						message: '',
						data: {
							content: {
								title: article.title,
								body: normalizeString(article.text),
								tags:
									typeof article.keywords === 'string'
									&& article.keywords.length > 2
										? toUpperCase(article.keywords)
											.split(',')
											.map(m => m.trim())
										: [],
								authors:
									typeof article.author === 'string'
									&& article.author.length > 2
										? toUpperCase(article.author)
											.split(',')
											.map(m => m.trim())
										: [],
								published_at: article.date
									? moment(article.date).format('YYYY-MM-DDTHH:mm:ssZZ')
									: null,
								description: normalizeString(article.description),
								image: article.image || ''
							}
						}
					});
				})
				.catch(err => {
					logger.error(`[SVC-SCRAPER] Error while scraping: ${err.message} - (${request.feed.screen_name}) ${request.url}`);
					return callback({
						code: errorCode(err.code ? err.code : 500),
						details: err.message
					},
					null);
				});
		}
	}

	ReloadPassages (req, callback) {
		logger.info('[SVC-SCRAPER] ReloadPassages: ', this.passages.length);
		return callback({ code: errorCode(500), details: 'Unimplemented method' },
			null);
	}

	RetryWithProxy (request, feed, callback) {
		logger.info(`[SVC-SCRAPER] RetryWithProxy - (${
			feed.hostname
		}) ${decodeURIComponent(request.url).toString()}`);
		extract(decodeURIComponent(request.url).toString(), getProxyUrl())
			.then(res => {
				const article = res;
				if (
					article.date
					&& moment(request.crawled_at).isBefore(moment(article.date))
				) {
					article.date = request.crawled_at;
				}

				if (article.text === '' || article.title === '') {
					logger.error(`[SVC-SCRAPER] Unable to scrape URL (malformed data) (${feed.hostname}) ${request.url}`);
					return callback({
						code: 9,
						details: `Unable to scrape URL (malformed data) (${feed.hostname}) ${request.url}`
					},
					null);
				}
				article.text = trimRight(article.text, this.passages);
				logger.debug(`[SVC-SCRAPER] DONEWITHPROXY - (${feed.hostname}) ${article.title}`);
				return callback(null, {
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
							description: normalizeString(article.description),
							image: article.image || ''
						}
					}
				});
			})
			.catch(err => {
				logger.error(`[SVC-SCRAPER] Error while scraping: (${errorCode(err.code ? err.code : 500)}) ${err.message} - (${feed.hostname}) ${request.url}`);
				return callback({
					code: errorCode(err.code ? err.code : 500),
					details: err.message
				},
				null);
			});
	}
}

module.exports = ScrapeService;
