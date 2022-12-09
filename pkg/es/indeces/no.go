package indeces

var ArticlesNO = `
{
    "settings": {
        "index": {
            "number_of_shards": 2,
            "number_of_replicas": 2
        },
        "analysis": {
            "filter": {
                "norwegian_stop": {
                    "type": "stop",
                    "stopwords": "_norwegian_"
                },
                "norwegian_stemmer": {
                    "type": "stemmer",
                    "language": "norwegian"
                }
            },
            "analyzer": {
                "norwegian_analyzer": {
                    "tokenizer": "standard",
                    "char_filter": [
                        "html_strip"
                    ],
                    "filter": [
                        "lowercase",
                        "norwegian_stop",
                        "norwegian_stemmer"
                    ]
                }
            }
        }
    },
    "alias": {
        "mediawatch_no-{now/M{yyyy.MM}}": {}
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
                        "analyzer": "norwegian_analyzer"
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
                        "analyzer": "norwegian_analyzer"
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
                        "analyzer": "norwegian_analyzer"
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
                        "analyzer": "norwegian_analyzer"
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
                        "analyzer": "norwegian_analyzer"
                    },
                    "quotes": {
                        "type": "text",
                        "analyzer": "norwegian_analyzer"
                    },
                    "sentences": {
                        "type": "text",
                        "analyzer": "norwegian_analyzer"
                    },
                    "stopWords": {
                        "type": "keyword"
                    },
                    "summary": {
                        "type": "text",
                        "analyzer": "norwegian_analyzer"
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
