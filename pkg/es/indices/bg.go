package indices

var ArticlesBG = `
{
    "settings": {
        "index": {
            "number_of_shards": 1,
            "number_of_replicas": 0
        },
        "analysis": {
            "filter": {
                "bulgarian_stop": {
                    "type": "stop",
                    "stopwords": "_bulgarian_"
                },
                "bulgarian_stemmer": {
                    "type": "stemmer",
                    "language": "bulgarian"
                }
            },
            "analyzer": {
                "bulgarian_analyzer": {
                    "tokenizer": "standard",
                    "char_filter": [
                        "html_strip"
                    ],
                    "filter": [
                        "lowercase",
                        "bulgarian_stop",
                        "bulgarian_stemmer"
                    ]
                }
            }
        }
    },
    "alias": {
        "mediawatch_bg-{now/M{yyyy.MM}}": {}
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
                        "analyzer": "bulgarian_analyzer"
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
                        "analyzer": "bulgarian_analyzer"
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
                        "analyzer": "bulgarian_analyzer"
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
                        "analyzer": "bulgarian_analyzer"
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
                        "analyzer": "bulgarian_analyzer"
                    },
                    "quotes": {
                        "type": "text",
                        "analyzer": "bulgarian_analyzer"
                    },
                    "sentences": {
                        "type": "text",
                        "analyzer": "bulgarian_analyzer"
                    },
                    "stopWords": {
                        "type": "keyword"
                    },
                    "summary": {
                        "type": "text",
                        "analyzer": "bulgarian_analyzer"
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
