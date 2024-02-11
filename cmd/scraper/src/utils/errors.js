const errorCode = (status) => {
	if (status === 403) {
		return 7;
	}
	if (status === 404) {
		return 5;
	}
	if (status === 406) {
		return 7;
	}
	if (status >= 500) {
		return 13;
	}

	return 2;
};

module.exports = { errorCode };
