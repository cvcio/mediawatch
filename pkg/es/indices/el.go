package indices

var ArticlesEL = `
{
    "settings": {
        "index": {
            "number_of_shards": 1,
            "number_of_replicas": 0
        },
        "analysis": {
            "filter": {
                "greek_stop": {
                    "type": "stop",
                    "stopwords": "_greek_"
                },
                "greek_lowercase": {
                    "type": "lowercase",
                    "language": "greek"
                },
                "greek_stemmer": {
                    "type": "stemmer",
                    "language": "greek"
                }
            },
            "analyzer": {
                "greek_analyzer": {
                    "tokenizer": "standard",
                    "char_filter": [
                        "html_strip"
                    ],
                    "filter": [
                        "greek_lowercase",
                        "greek_stop",
                        "greek_stemmer",
                        "asciifolding"
                    ]
                }
            }
        }
    },
    "aliases": {
        "<mediawatch_articles_el-{now/M{yyyy.MM}}>": {}
    }
}
`
