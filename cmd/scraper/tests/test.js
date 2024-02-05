require('dotenv').config({ path: '.env' });

const moment = require('moment');
const axios = require('axios');

const logger = require('../src/logger');
const { Scrape } = require('../src/services');

const test_feed_scrapability_async = async f => {
	// f.stream.requires_proxy = true;
	const req = {
		request: {
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
			// if (feed.hostname === 'epirusonline.gr') {
				// feed.stream.requires_proxy = true;
				results.push(test_feed_scrapability_async(feed));
			// }
		}
	}

	const done = await Promise.all(results);
	// eslint-disable-next-line no-restricted-syntax
	for (const result of done) {
		if (result.status === 'error') {
			logger.info(`${result.code} - ${result.hostname} - ${result.message}`);
		} else {
			logger.info(`200 - ${result.hostname}`);
		}
	}
};

test().catch(e => {
	console.error(e);
	process.exit(1);
});

// const options = {
// 	method: 'POST',
// 	url: 'http://localhost:8000/mediawatch.feeds.v2.FeedService/GetFeeds',
// 	headers: { 'Content-Type': 'application/json' },
// 	data: { limit: 3000 }
// };

// axios
// 	.request(options)
// 	.then(res => {
// 		for (let i = 0; i < res.data.data.length; i++) {
// 			const feed = res.data.data[i];
// 			// if (feed.hostname === 'efsyn.gr') {
// 			if (feed.test && feed.test.url !== '') {
// 				// if (feed.hostname === 'efsyn.gr') {
// 				test_feed_scrapability_service(feed);
// 			}
// 		}
// 	})
// 	.catch(error => {
// 		console.error(error);
// 	});

// for (var i = 0; i < data.length; i++) {
// 	// console.log(data[i]);
//     if (ONE && i > 0) break;
//     // let f = data[i];
//     // if (f.testURL !== '' && f.user_name === 'efsyntakton') { // && f.screen_name === 'documentonews'
//     //     let req = {
//     //         request: {
//     //             feed: f,
//     //             url: f.testURL,
// 	// 			screen_name: f.user_name,
// 	// 			lang: 'el',
// 	// 			crawled_at: moment().format('YYYY-MM-DDTHH:mm:ss')
//     //         }
//     //     };
//     //     console.log(req);
//     //     service.Scrape(req, (err, data) => {
//     //         if (err) {
// 	// 			console.dir(err);
//     //             console.error(f.screen_name, err.message);
//     //             return;
//     //         }
//     //         console.log(f.screen_name, data.data.content.title);
//     //     });
//     // }
//     let f = data[i];
//     // if (f.test && f.test.url !== '' && f.stream && f.stream.requires_proxy) {
//     if (f.test && f.test.url !== '' && f.user_name === 'EFSYNTAKTON') {
//         let req = {
//             request: {
//                 feed: JSON.stringify(f),
//                 url: f.test.url,
// 				screen_name: f.user_name,
// 				lang: 'el',
// 				crawled_at: moment().format('YYYY-MM-DDTHH:mm:ss')
//             }
//         };
// 		// console.log(req);
//         // service.Scrape(req, (err, data) => {
// 		// 	console.log(data)
// 		// });

// 		console.log(req)
//         service.RetryWithProxy(req.request, f, (err, data) => {
// 			console.log(data)
// 		});
//     }
// }
