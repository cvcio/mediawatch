package indeces

var ArticlesFI = `
{
    "settings": {
        "index": {
            "number_of_shards": 2,
            "number_of_replicas": 2
        },
        "analysis": {
            "filter": {
                "finnish_stop": {
                    "type": "stop",
                    "stopwords": "_finnish_"
                },
                "finnish_stemmer": {
                    "type": "stemmer",
                    "language": "finnish"
                }
            },
            "analyzer": {
                "finnish_analyzer": {
                    "tokenizer": "standard",
                    "char_filter": [
                        "html_strip"
                    ],
                    "filter": [
                        "lowercase",
                        "finnish_stop",
                        "finnish_stemmer"
                    ]
                }
            }
        }
    },
    "alias": {
        "mediawatch_fi-{now/M{yyyy.MM}}": {}
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
                        "analyzer": "finnish_analyzer"
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
                        "analyzer": "finnish_analyzer"
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
                        "analyzer": "finnish_analyzer"
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
                        "analyzer": "finnish_analyzer"
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
                        "analyzer": "finnish_analyzer"
                    },
                    "quotes": {
                        "type": "text",
                        "analyzer": "finnish_analyzer"
                    },
                    "sentences": {
                        "type": "text",
                        "analyzer": "finnish_analyzer"
                    },
                    "stopWords": {
                        "type": "keyword"
                    },
                    "summary": {
                        "type": "text",
                        "analyzer": "finnish_analyzer"
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
