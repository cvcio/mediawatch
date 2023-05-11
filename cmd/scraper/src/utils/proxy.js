const getProxy = () => {
	if (!process.env.PROXY_ENABLED) return null;

	const proxy = {
		protocol: 'http',
		host: process.env.PROXY_HOST,
		port: parseInt(process.env.PROXY_PORT, 10)
	};

	if (process.env.PROXY_USERNAME !== '' && process.env.PROXY_USERNAME !== '') {
		proxy.auth = {
			username: process.env.PROXY_USERNAME,
			password: process.env.PROXY_PASSWORD
		};
	}

	return proxy;
};

export { getProxy };
