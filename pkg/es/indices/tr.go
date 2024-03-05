package indices

var ArticlesTR = `
{
    "settings": {
        "index": {
            "number_of_shards": 1,
            "number_of_replicas": 0
        },
        "analysis": {
            "filter": {
                "turkish_stop": {
                    "type": "stop",
                    "stopwords": "_turkish_"
                },
                "turkish_stemmer": {
                    "type": "stemmer",
                    "language": "turkish"
                }
            },
            "analyzer": {
                "turkish_analyzer": {
                    "tokenizer": "standard",
                    "char_filter": [
                        "html_strip"
                    ],
                    "filter": [
                        "apostrophe",
                        "lowercase",
                        "turkish_stop",
                        "turkish_stemmer"
                    ]
                }
            }
        }
    },
    "alias": {
        "mediawatch_tr-{now/M{yyyy.MM}}": {}
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
                        "analyzer": "turkish_analyzer"
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
                        "analyzer": "turkish_analyzer"
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
                        "analyzer": "turkish_analyzer"
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
                        "analyzer": "turkish_analyzer"
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
                        "analyzer": "turkish_analyzer"
                    },
                    "quotes": {
                        "type": "text",
                        "analyzer": "turkish_analyzer"
                    },
                    "sentences": {
                        "type": "text",
                        "analyzer": "turkish_analyzer"
                    },
                    "stopWords": {
                        "type": "keyword"
                    },
                    "summary": {
                        "type": "text",
                        "analyzer": "turkish_analyzer"
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
