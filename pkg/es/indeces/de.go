package indeces

var ArticlesDE = `
{
    "settings": {
        "index": {
            "number_of_shards": 2,
            "number_of_replicas": 2
        },
        "analysis": {
            "filter": {
                "german_stop": {
                    "type": "stop",
                    "stopwords": "_german_"
                },
                "german_stemmer": {
                    "type": "stemmer",
                    "language": "light_german"
                }
            },
            "analyzer": {
                "german_analyzer": {
                    "tokenizer": "standard",
                    "char_filter": [
                        "html_strip"
                    ],
                    "filter": [
                        "lowercase",
                        "german_stop",
                        "german_normalization",
                        "german_stemmer"
                    ]
                }
            }
        }
    },
    "alias": {
        "mediawatch_de-{now/M{yyyy.MM}}": {}
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
                        "analyzer": "german_analyzer"
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
                        "analyzer": "german_analyzer"
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
                        "analyzer": "german_analyzer"
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
                        "analyzer": "german_analyzer"
                    },
                    "entities": {
                        "properties": {
                            "entity_text": {
                                "type": "keyword"
                            },
                            "entity_type": {
                                "type": "keyword"
                            }
                        }
                    },
                    "keywords": {
                        "type": "keyword",
                        "analyzer": "german_analyzer"
                    },
                    "quotes": {
                        "type": "text",
                        "analyzer": "german_analyzer"
                    },
                    "sentences": {
                        "type": "text",
                        "analyzer": "german_analyzer"
                    },
                    "stopWords": {
                        "type": "keyword"
                    },
                    "summary": {
                        "type": "text",
                        "analyzer": "german_analyzer"
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
