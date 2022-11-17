// api: '"http://www.amna.gr/feeds/getarticle.php?id="+url.split(\'/\').splice(-2, 2)[0]+"&infolevel=ADVANCED"',
// ex. https://www.amna.gr/feeds/twitterupload.php?id=544853&cat=home&title=Bouli-Apochi-apo-tis-onomastikes-psifofories-apofasise-o-SYRIZA--Kanei-logo-gia-parabiasi-tou-Kanonismou
// ex. https://www.amna.gr/home/article/525014/Stin-Patra-ektaktos-o-M-Chrusochoidis-me-ton-gg-Politikis-Prostasias-BPapageorgiou
// https://www.amna.gr/home/article/624073/Komision-Schedio-gia-doruforiko-sustima-sundesimotitas-pou-tha-epitrepei-se-oli-tin-Europi--prosbasi-upsilis-tachutitas-sto-diadiktuo

const { getParameterByName } = require('../utils/strings');
const { URL } = require('url');
const axios = require('axios');

const url = (u) => {
    const id = getParameterByName('id', u) || u.split('/').splice(-2, 2)[0];
    return id ? `https://www.amna.gr/feeds/getarticle.php?id=${id}&infolevel=ADVANCED` : null;
};

const fetchAPI = async (link) => {
	const url = new URL(link);
	try {
		axios.defaults.headers.get['Content-Type'] = 'application/json;charset=utf-8;text/html;text/plain';
		axios.defaults.headers.get['Access-Control-Allow-Origin'] = '*';
		axios.defaults.headers.get['User-Agent'] = 'MediaWatch Bot/2.0 (mediawatch.io)';

		const html = await axios({
			method: 'get',
			url: url.href,
			responseType:'json',
			insecureHTTPParser: true
		});
		if (html.status >= 400) {
			throw new Error(`Error Not Found or Not Authorized: ${html.status} ${html.statusText}`);
		}

		return parseJSON(html.data);
	} catch (err) {
		throw new Error(`Error Fetching URL: ${err.message}`);
	}
};

const parseJSON = (body) => {
	let article = body;

	let image = article.photo1.split('/');
	image[image.length - 1] = 'w' + image[image.length - 1];
	image = image.join('/');

	return {
		title: article.title,
		body: article.text,
		authors: article.author,
		publishedAt: article.c_daytime,
		tags: article.tags,
		description: article.capelo,
		image: image.replace('..', 'https://www.amna.gr')
	};
};

module.exports = { url, fetchAPI };
