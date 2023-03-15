const toUpperCase = (str) => {
	if (!str) {
		return '';
	}
	str = str.normalize('NFD').replace(/[\u0300-\u036f]/g, '');
	return str.toUpperCase();
};

const normalizeString = (str) => {
	let text = str.trim();

	text = text.replace(/<[^>]+>/g, '');
	text = text.replace(/<(?:.|\n)*?>/gm, '');
	text = text.replace(/<br>/gi, '\n');
	text = text.replace(/<p.*>/gi, '\n');
	text = text.replace(/<a.*href="(.*?)".*>(.*?)<\/a>/gi, " $2 (Link->$1) ");
	text = text.replace(/<(?:.|\s)*?>/g, '');
	text = text.replace(/[\r\n]+/g, '\n\n');
	text = text.replace(/ +/g, ' ');

	return text;
};

const getParameterByName = (name, url) => {
	name = name.replace(/[[\]]/g, '\\$&');
	var regex = new RegExp('[?&]' + name + '(=([^&#]*)|&|#|$)'),
		results = regex.exec(url);
	if (!results) return null;
	if (!results[2]) return '';
	return decodeURIComponent(results[2].replace(/\+/g, ' '));
};

const trimRight = (str, passages) => {
	let text = str;
	for (var i = 0; i < passages.length; i++) {
		let t = passages[i];
		if (text.match(new RegExp(t.text))) {
			text = text.split(new RegExp(t.text))[0];
		}
	}
	return text;
};



export {
	trimRight,
	toUpperCase,
	normalizeString,
	getParameterByName
};
