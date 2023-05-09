module.exports = {
    'plugins': [
        '@babel/plugin-transform-regenerator',
        '@babel/plugin-transform-runtime',
    ],
    presets: [
        [
            '@babel/preset-env',
            {
                targets: {
                    esmodules: true,
                },
            },
        ],
    ],
}