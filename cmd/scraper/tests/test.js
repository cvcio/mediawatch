const { ScrapeService } = require('../src/services');
const feeds = require('../data/v2_feeds.json');
const moment = require('moment');
const ONE = false;

const service = new ScrapeService(null);
let data = feeds.data.data;
for (var i = 0; i < data.length; i++) {
	// console.log(data[i]);
    if (ONE && i > 0) break;
    // let f = data[i];
    // if (f.testURL !== '' && f.user_name === 'efsyntakton') { // && f.screen_name === 'documentonews'
    //     let req = {
    //         request: {
    //             feed: f,
    //             url: f.testURL,
	// 			screen_name: f.user_name,
	// 			lang: 'el',
	// 			crawled_at: moment().format('YYYY-MM-DDTHH:mm:ss')
    //         }
    //     };
    //     console.log(req);
    //     service.Scrape(req, (err, data) => {
    //         if (err) {
	// 			console.dir(err);
    //             console.error(f.screen_name, err.message);
    //             return;
    //         }
    //         console.log(f.screen_name, data.data.content.title);
    //     });
    // }
    let f = data[i];
	if (f.test && f.test.url !== '' && f.stream && f.stream.requires_proxy) {
        console.log(f);
	}
    if (f.test && f.test.url !== '' && f.user_name === 'efsyntakton') { // && f.screen_name === 'documentonews'
        let req = {
            request: {
                feed: f,
                url: f.testURL,
				screen_name: f.user_name,
				lang: 'el',
				crawled_at: moment().format('YYYY-MM-DDTHH:mm:ss')
            }
        };
        // service.Scrape(req, (err, data) => {
        //     if (err) {
		// 		console.dir(err);
        //         console.error(f.screen_name, err.message);
        //         return;
        //     }
        //     console.log(f.screen_name, data.data.content.title);
        // });
    }
}
