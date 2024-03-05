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
    "alias": {
        "mediawatch_el-{now/M{yyyy.MM}}": {}
    },
    "mappings": {
        "properties": {
            "content": {
                "properties": {
                    "authors": {
                        "type": "keyword"
                    },
                    "body": {
                        "type": "text",
                        "analyzer": "greek_analyzer"
                    },
                    "categories": {
                        "type": "keyword"
                    },
                    "editedAt": {
                        "type": "date",
                        "ignore_malformed": true
                    },
                    "excerpt": {
                        "type": "text",
                        "analyzer": "greek_analyzer"
                    },
                    "image": {
                        "type": "keyword"
                    },
                    "publishedAt": {
                        "type": "date",
                        "ignore_malformed": true
                    },
                    "sources": {
                        "type": "keyword"
                    },
                    "tags": {
                        "type": "keyword"
                    },
                    "title": {
                        "type": "text",
                        "analyzer": "greek_analyzer"
                    }
                }
            },
            "crawledAt": {
                "type": "date",
                "ignore_malformed": true
            },
            "docId": {
                "type": "keyword"
            },
            "lang": {
                "type": "keyword"
            },
            "nlp": {
                "properties": {
                    "claims": {
                        "type": "text",
                        "analyzer": "greek_analyzer"
                    },
                    "entities": {
                        "properties": {
                            "text": {
                                "type": "keyword"
                            },
                            "type": {
                                "type": "keyword"
                            }
                        }
                    },
                    "keywords": {
                        "type": "text",
                        "analyzer": "greek_analyzer"
                    },
                    "quotes": {
                        "type": "text",
                        "analyzer": "greek_analyzer"
                    },
                    "sentences": {
                        "type": "text",
                        "analyzer": "greek_analyzer"
                    },
                    "stopWords": {
                        "type": "keyword"
                    },
                    "summary": {
                        "type": "text",
                        "analyzer": "greek_analyzer"
                    },
                    "tokens": {
                        "type": "keyword"
                    },
                    "topics": {
                        "type": "keyword"
                    }
                }
            },
            "screen_name": {
                "type": "keyword"
            },
            "tweet_id": {
                "type": "long"
            },
            "tweet_id_str": {
                "type": "keyword"
            },
            "url": {
                "type": "keyword"
            }
        }
    }
}
`
