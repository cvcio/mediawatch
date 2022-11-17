const { ScrapeService } = require('../src/services');
const feeds = require('../data/feeds.json');
const ONE = false;

const service = new ScrapeService(null);
let data = feeds.data.data;
for (var i = 0; i < data.length; i++) {
    if (ONE && i > 0) break;
    let f = data[i];
    if (f.testURL !== '' && f.screen_name === 'amna_news') { // && f.screen_name === 'documentonews'
        let req = {
            request: {
                feed: JSON.stringify(f),
                url: f.testURL
            }
        };
        // console.log(req);
        service.SimpleScrape(req, (err, data) => {
            if (err) {
                console.error(f.screen_name, err.message);
                return;
            }
            console.log(f.screen_name, data);
        });
    }
}
