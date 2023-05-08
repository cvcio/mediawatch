const logger = require('../logger');
const { parsers } = require('../parsers');
const { toUpperCase, normalizeString, trimRight } = require('../utils/strings');
const { extract } = require('ascraper'); // require('/home/andefined/js/misc/npm/ascraper/lib'); //
const { errorCode } = require('../utils/errors');

const moment = require('moment');
moment.suppressDeprecationWarnings = true;

class ScrapeService {
    constructor(mongo, proxy) {
        this.passages = [];
        this.mongo = mongo;
		this.proxy = proxy;
        this.GetPassages();
    }

    async GetPassages(err, client) {
        if (!this.mongo) return;

        const db = this.mongo.db(process.env.MONGODB_DB);
        const collection = db.collection('passages');

        const passages = await collection.find({}).limit(5000).toArray();

        this.passages = passages.filter(m => m.type === 'trim');
        logger.debug(`[SVC-SCRAPER] (${this.passages.length}) passages loaded`);
    }

    Scrape = (req, callback) => {
        let request = req.request;
		const feed = JSON.parse(request.feed);

        logger.info(`[SVC-SCRAPER] Scrape - (${feed.hostname}) ${decodeURIComponent(request.url).toString()}`);

        if (request.screen_name && parsers[request.screen_name.toLowerCase()] !== undefined) {
			let parser = parsers[request.screen_name.toLowerCase()];
			let url = parser.url(request.url);

			if (!url) {
                logger.error(`[SVC-SCRAPER] Unable to scrape URL (url error) (${feed.hostname}) ${request.url}`);
                return callback({code: 9, details: `Unable to scrape URL (url error) (${feed.hostname}) ${request.url}`}, null);
            }

			parser.fetchAPI(url)
				.then(res => {
					if (!res) {
						logger.error(`[SVC-SCRAPER] Unable to scrape URL (response empty) (${feed.hostname}) ${request.url}`);
						return callback({code: 9, details: `Unable to scrape URL (response empty) (${feed.hostname}) ${request.url}`}, null);

					}
					logger.debug(`[SVC-SCRAPER] Data ${JSON.stringify(res)}`);

					let article = res;

					if (article.publishedAt && moment.utc(request.crawled_at).isBefore(moment.utc(article.publishedAt))) {
						article.publishedAt = request.crawled_at;
					}

					if (article.text === '' || article.title === '') {
						logger.error(`[SVC-SCRAPER] Unable to scrape URL (malformed data) (${feed.hostname}) ${request.url}`);
						return callback({code: 9, details: `Unable to scrape URL (malformed data) (${feed.hostname}) ${request.url}`}, null);
					}

					return callback(null, {
						status: 'success',
						code: 200,
						message: '',
						data: {
							content: {
								title: article.title.trim(),
								body: normalizeString(article.body),
								authors: (typeof article.authors === 'string' && article.authors.length > 2) ?
									toUpperCase(article.authors).split(',').map(m => m.trim()) : [],
								published_at: article.publishedAt ?
									moment(article.publishedAt).format('YYYY-MM-DDTHH:mm:ssZZ') : moment(request.crawled_at).format('YYYY-MM-DDTHH:mm:ssZZ'),
								tags: (typeof article.tags === 'string' && article.tags.length > 2) ?
									toUpperCase(article.tags).split(',').map(m => m.trim()) : [],
								description: normalizeString(article.description),
								image: article.image || ''
							}
						}
					});
				})
				.catch(err => {
					logger.error(`[SVC-SCRAPER] Error while scraping: (${errorCode(err.response ? err.response.status : 500)}) ${err.message} - (${feed.hostname}) ${request.url}`);
                    return callback({code: errorCode(err.response ? err.response.status : 500), details: err.message}, null);
				});
        } else {
            extract(decodeURIComponent(request.url).toString())
                .then(res => {
                    logger.debug(`[SVC-SCRAPER] Data ${JSON.stringify(res)}`);

                    let article = res;
                    if (article.date && moment(request.crawled_at).isBefore(moment(article.date))) {
                        article.date = request.crawled_at;
                    }

                    if (article.text === '' || article.title === '') {
                        logger.error(`[SVC-SCRAPER] Unable to scrape URL (malformed data) (${feed.hostname}) ${request.url}`);
						return callback({code: 9, details: `Unable to scrape URL (malformed data) (${feed.hostname}) ${request.url}`}, null);
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
                                authors: (typeof article.author === 'string' && article.author.length > 2) ?
                                    toUpperCase(article.author).split(',').map(m => m.trim()) : [],
                                published_at: article.date ?
                                    moment(article.date).format('YYYY-MM-DDTHH:mm:ssZZ') : moment(request.crawled_at).format('YYYY-MM-DDTHH:mm:ssZZ'),
                                tags: (typeof article.keywords === 'string' && article.keywords.length > 2) ?
                                    toUpperCase(article.keywords).split(',').map(m => m.trim()) : [],
								description: normalizeString(article.description),
								image: article.image || ''
                            }
                        }
                    });
                })
                .catch(err => {
					// if (err.response && err.response.status == 403 && this.proxy) {
					// 	return this.RetryWithProxy(request, feed, callback);
					// }
                    logger.error(`[SVC-SCRAPER] Error while scraping: (${errorCode(err.response ? err.response.status : 500)}) ${err.message} - (${feed.hostname}) ${request.url}`);
                    return callback({code: errorCode(err.response ? err.response.status : 500), details: err.message}, null);
                });
        }
    }
    SimpleScrape = (req, callback) => {
        let now = moment();
        let request = req.request;
        request.feed = JSON.parse(request.feed);
        logger.info(`[SVC-SCRAPER] SimpleScrape - (${request.feed.screen_name}) ${decodeURIComponent(request.url).toString()}`);

        if (parsers[request.feed.screen_name] !== undefined) {
			let parser = parsers[request.feed.screen_name.toLowerCase()];
			let url = parser.url(request.url);

			if (!url) {
                logger.error(`[SVC-SCRAPER] URL Error (${request.feed.screen_name}) ${request.url}`);
                return callback(Error(`URL Error (${request.feed.screen_name}) ${request.url}`), null);
            }

			parser.fetchAPI(url)
				.then(res => {
					if (!res) {
						logger.error(`[SVC-SCRAPER] Unable to scrape URL (${request.feed.screen_name}) ${request.url}`);
						return callback(new Error(`Unable to scrape URL (${request.feed.screen_name}) ${request.url}`), null);
					}
					logger.debug(`[SVC-SCRAPER] Data ${JSON.stringify(res)}`);

					let article = res;

					// if (article.publishedAt && moment.utc(request.crawled_at).isBefore(moment.utc(article.publishedAt))) {
					// 	article.publishedAt = request.crawled_at;
					// }

					if (article.text === '' || article.title === '') {
						logger.error(`[SVC-SCRAPER] Unable to scrape URL (${request.feed.screen_name}) ${request.url}`);
						return callback(new Error(`Unable to scrape URL (${request.feed.screen_name}) ${request.url}`), null);
					}

					return callback(null, {
						status: 'success',
						code: 200,
						message: '',
						data: {
							content: {
								title: article.title.trim(),
								body: normalizeString(article.body),
								authors: (typeof article.authors === 'string' && article.authors.length > 2) ?
									toUpperCase(article.authors).split(',').map(m => m.trim()) : [],
								published_at: article.publishedAt ?
									moment(article.publishedAt).format('YYYY-MM-DDTHH:mm:ssZZ') : moment(request.crawled_at).format('YYYY-MM-DDTHH:mm:ssZZ'),
								tags: (typeof article.tags === 'string' && article.tags.length > 2) ?
									toUpperCase(article.tags).split(',').map(m => m.trim()) : [],
								description: normalizeString(article.description),
								image: article.image || ''
							}
						}
					});
				})
				.catch(err => {
					logger.error('Error while scraping:', err.message, `(${request.feed.screen_name}) ${request.url}`);
                    return callback({code: errorCode(err.response ? err.response.status : 500), details: err.message}, null);
				});
        } else {
            // Run Scraper
            extract(decodeURIComponent(request.url).toString())
                .then(res => {
                    let article = res;
                    if (article.date && now.isBefore(moment(article.date))) {
                        article.date = null;
                    }
                    if (article.body === '' || article.title === '') {
                        logger.error(`[SVC-SCRAPER] Unable to scrape URL (${request.feed.screen_name}) ${request.url}`);
                        return callback(new Error(`Unable to scrape URL (${request.feed.screen_name}) ${request.url}`), null);
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
                                tags: (typeof article.keywords === 'string' && article.keywords.length > 2) ?
                                    toUpperCase(article.keywords).split(',').map(m => m.trim()) : [],
                                authors: (typeof article.author === 'string' && article.author.length > 2) ?
                                    toUpperCase(article.author).split(',').map(m => m.trim()) : [],
                                published_at: article.date ?
                                    moment(article.date).format('YYYY-MM-DDTHH:mm:ssZZ') : null,
								description: normalizeString(article.description),
								image: article.image || ''
                            }
                        }
                    });
                })
                .catch(err => {
                    logger.error(`[SVC-SCRAPER] Error while scraping: ${err.message} - (${request.feed.screen_name}) ${request.url}`);
                    return callback({code: errorCode(err.response ? err.response.status : 500), details: err.message}, null);
                });
        }
    }
    ReloadPassages = (req, callback) => {
        this.GetPassages();
        return callback(null, null);
    }
	RetryWithProxy = (request, feed, callback) => {
		logger.info(`[SVC-SCRAPER] RetryWithProxy - (${feed.hostname}) ${decodeURIComponent(request.url).toString()}`);
		extract(decodeURIComponent(request.url).toString(), this.proxy)
			.then(res => {
				let article = res;
				if (article.date && moment(request.crawled_at).isBefore(moment(article.date))) {
					article.date = request.crawled_at;
				}

				if (article.text === '' || article.title === '') {
					logger.error(`[SVC-SCRAPER] Unable to scrape URL (malformed data) (${feed.hostname}) ${request.url}`);
					return callback({code: 9, details: `Unable to scrape URL (malformed data) (${feed.hostname}) ${request.url}`}, null);
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
							authors: (typeof article.author === 'string' && article.author.length > 2) ?
								toUpperCase(article.author).split(',').map(m => m.trim()) : [],
							published_at: article.date ?
								moment(article.date).format('YYYY-MM-DDTHH:mm:ssZZ') : moment(request.crawled_at).format('YYYY-MM-DDTHH:mm:ssZZ'),
							tags: (typeof article.keywords === 'string' && article.keywords.length > 2) ?
								toUpperCase(article.keywords).split(',').map(m => m.trim()) : [],
							description: normalizeString(article.description),
							image: article.image || ''
						}
					}
				});
			})
			.catch(err => {
				logger.error(`[SVC-SCRAPER] Error while scraping: (${errorCode(err.response ? err.response.status : 500)}) ${err.message} - (${feed.hostname}) ${request.url}`);
				return callback({code: errorCode(err.response ? err.response.status : 500), details: err.message}, null);
			});
	}
};

module.exports = { ScrapeService };
