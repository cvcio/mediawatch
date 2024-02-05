const moment = require('moment');

const toUpperCase = (str) => {
	if (!str) {
		return '';
	}
	str = str.normalize('NFD').replace(/[\u0300-\u036f]/g, '');
	return str.toUpperCase();
};

const normalizeString = (str) => {
	let text = str.trim();

	text = text.replace(/\b([a-zA-Z=]*[0-9]){24,}[a-zA-Z=]*\b/g, '');
	text = text.replace(/<[^>]+>/g, '');
	text = text.replace(/<(?:.|\n)*?>/gm, '');
	text = text.replace(/<br>/gi, '\n');
	text = text.replace(/<p.*>/gi, '\n');
	text = text.replace(/<a.*href="(.*?)".*>(.*?)<\/a>/gi, ' $2 (Link->$1) ');
	text = text.replace(/<(?:.|\s)*?>/g, '');
	text = text.replace(/[\r\n]+/g, '\n\n');
	text = text.replace(/ +/g, ' ');

	// text = text.replace(/\b([a-zA-Z=]*[0-9]){24,}[a-zA-Z=]*\b/g, '');
	// // text = text.replace(/<(?:.|\s)*?>/g, '');
	// text = text.replace(/<a.*href="(.*?)".*>(.*?)<\/a>/gi, ' $2 (Link->$1) ');
	// text = text.replace(/<br>/gi, '\n');
	// text = text.replace(/<p.*>/gi, '\n');
	// text = text.replace(/[\r\n]+/g, '\n\n');
	// text = text.replace(/<(?:.|\n)*?>/gm, '');
	// text = text.replace(/<[^>]+>/g, '');
	// text = text.replace(/<\/?[^>]+>/gi, '');
	// text = text.replace(/ +/g, ' ');

	return text.trim();
};

const getParameterByName = (name, url) => {
	name = name.replace(/[[\]]/g, '\\$&');
	const regex = new RegExp(`[?&]${name}(=([^&#]*)|&|#|$)`);
	const results = regex.exec(url);
	if (!results) return null;
	if (!results[2]) return '';
	return decodeURIComponent(results[2].replace(/\+/g, ' '));
};

const trimRight = (str, passages) => {
	let text = str;
	for (let i = 0; i < passages.length; i++) {
		const t = passages[i];
		if (text.match(new RegExp(t.text))) {
			// eslint-disable-next-line
			text = text.split(new RegExp(t.text))[0];
		}
	}
	return text;
};

const mergeArticle = (article, data) => {
	if (!article.content) {
		article.content = {};
	}

	if (!article.content.title || article.content.title === '') {
		article.content.title = data.content.title;
	}
	article.content.body = data.content.body;
	article.content.authors = data.content.authors;
	article.content.tags = data.content.tags;
	article.content.excerpt = data.content.description;
	article.content.image = data.content.image;
	article.content.published_at = moment(data.content.date_published).toISOString();

	return article;
};

export {
	trimRight,
	toUpperCase,
	normalizeString,
	getParameterByName,
	mergeArticle
};
