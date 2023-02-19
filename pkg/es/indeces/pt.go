package indeces

var ArticlesPT = `
{
    "settings": {
        "index": {
            "number_of_shards": 2,
            "number_of_replicas": 2
        },
        "analysis": {
            "filter": {
                "portuguese_stop": {
                    "type": "stop",
                    "stopwords": "_portuguese_"
                },
                "portuguese_stemmer": {
                    "type": "stemmer",
                    "language": "light_portuguese"
                }
            },
            "analyzer": {
                "portuguese_analyzer": {
                    "tokenizer": "standard",
                    "char_filter": [
                        "html_strip"
                    ],
                    "filter": [
                        "lowercase",
                        "portuguese_stop",
                        "portuguese_stemmer"
                    ]
                }
            }
        }
    },
    "alias": {
        "mediawatch_pt-{now/M{yyyy.MM}}": {}
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
                        "analyzer": "portuguese_analyzer"
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
                        "analyzer": "portuguese_analyzer"
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
                        "analyzer": "portuguese_analyzer"
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
                        "analyzer": "portuguese_analyzer"
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
                        "type": "keyword",
                        "analyzer": "portuguese_analyzer"
                    },
                    "quotes": {
                        "type": "text",
                        "analyzer": "portuguese_analyzer"
                    },
                    "sentences": {
                        "type": "text",
                        "analyzer": "portuguese_analyzer"
                    },
                    "stopWords": {
                        "type": "keyword"
                    },
                    "summary": {
                        "type": "text",
                        "analyzer": "portuguese_analyzer"
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
