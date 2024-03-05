module.exports = {
	root: true,
	env: {
		node: true,
		es2021: true
	},
	extends: [
		'eslint:recommended'
	],
	rules: {
		'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'off',
		'no-debugger': process.env.NODE_ENV === 'production' ? 2 : 0,
		'space-before-function-paren': ['error', 'always'],
		'arrow-parens': 0,
		'generator-star-spacing': 0,
		'no-mixed-spaces-and-tabs': 0,
		indent: ['error', 'tab'],
		semi: ['error', 'always'],
		'no-tabs': 0,
		'comma-dangle': 0,
		'consistent-return': 0,
		'function-paren-newline': ['error', 'never'],
		'implicit-arrow-linebreak': ['off'],
		'no-param-reassign': 0,
		'no-underscore-dangle': 0,
		'no-shadow': 0,
		'no-plusplus': 0,
		'no-unused-expressions': 0,
		// 'import/prefer-default-export': ['error'],
		'max-len': [0, { code: 120 }],
	},
	parserOptions: {
		ecmaVersion: 'latest',
		sourceType: 'module'
	}
};
