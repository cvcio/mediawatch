module.exports = {
    root: true,
    env: {
        node: true
    },
    extends: [
        'airbnb-base'
    ],
    rules: {
        'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'off',
        // allow debugger during development
        'space-before-function-paren': ['error', 'always'],
        // allow paren-less arrow functions
        'arrow-parens': 0,
        // allow async-await
        'generator-star-spacing': 0,
        // allow debugger during development
        'no-debugger': process.env.NODE_ENV === 'production' ? 2 : 0,
        'no-mixed-spaces-and-tabs': 0,
        // tabs
        'indent': [4, 'tab'],
        'no-tabs': 0,
        'semi': ['error', 'always'],
        'comma-dangle': 0,
        'consistent-return': 0,
        'function-paren-newline': ['error', 'never'],
        'implicit-arrow-linebreak': ['off'],
        'no-param-reassign': 0,
        'no-underscore-dangle': 0,
        'no-shadow': 0,
        'no-console': 0,
        'no-plusplus': 0,
        'no-unused-expressions': 0,
    },
    parserOptions: {
        parser: 'babel-eslint',
        sourceType: 'module',
        ecmaVersion: 8
    }
};