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

	console.log(proxy);

	return proxy;
};

export { getProxy };
