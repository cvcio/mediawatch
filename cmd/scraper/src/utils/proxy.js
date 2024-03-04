const { URL, format } = require('url');

const getProxy = () => {
	if (!process.env.PROXY_ENABLED) return null;

	const proxyList = process.env.PROXY_LIST.split(',');
	const [host, port] = proxyList[Math.floor(Math.random() * proxyList.length)].split(':');

	const proxy = {
		protocol: 'http',
		host,
		port: parseInt(port, 10)
	};

	if (process.env.PROXY_USERNAME !== '' && process.env.PROXY_PASSWORD !== '') {
		proxy.auth = {
			username: process.env.PROXY_USERNAME,
			password: process.env.PROXY_PASSWORD
		};
	}

	return proxy;
};

const getProxyUrl = () => {
	if (!process.env.PROXY_ENABLED) return null;

	const proxy = {
		protocol: process.env.PROXY_SCHEME || 'http',
		host: process.env.PROXY_HOST || '',
		port: parseInt(process.env.PROXY_PORT),
		username: process.env.PROXY_USERNAME || '',
		password: process.env.PROXY_PASSWORD || '',
	};

	let url = null;

	if (proxy.username !== '' && proxy.password !== '') {
		url = new URL(`${proxy.protocol}://${proxy.username}:${proxy.password}@${proxy.host}:${proxy.port}`);
	} else {
		url = new URL(`${proxy.protocol}://${proxy.host}:${proxy.port}`);
	}

	return format(url);
};

module.exports = { getProxy, getProxyUrl };
