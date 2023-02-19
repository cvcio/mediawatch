package indeces

var ArticlesIT = `
{
    "settings": {
        "index": {
            "number_of_shards": 2,
            "number_of_replicas": 2
        },
        "analysis": {
            "filter": {
                "italian_elision": {
                    "type": "elision",
                    "articles": [
                        "c", "l", "all", "dall", "dell",
                        "nell", "sull", "coll", "pell",
                        "gl", "agl", "dagl", "degl", "negl",
                        "sugl", "un", "m", "t", "s", "v", "d"
                    ],
                    "articles_case": true
                },
                "italian_stop": {
                    "type": "stop",
                    "stopwords": "_italian_"
                },
                "italian_stemmer": {
                    "type": "stemmer",
                    "language": "light_italian"
                }
            },
            "analyzer": {
                "italian_analyzer": {
                    "tokenizer": "standard",
                    "char_filter": [
                        "html_strip"
                    ],
                    "filter": [
                        "italian_elision",
                        "lowercase",
                        "italian_stop",
                        "italian_stemmer"
                    ]
                }
            }
        }
    },
    "alias": {
        "mediawatch_it-{now/M{yyyy.MM}}": {}
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
                        "analyzer": "italian_analyzer"
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
                        "analyzer": "italian_analyzer"
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
                        "analyzer": "italian_analyzer"
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
                        "analyzer": "italian_analyzer"
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
                        "analyzer": "italian_analyzer"
                    },
                    "quotes": {
                        "type": "text",
                        "analyzer": "italian_analyzer"
                    },
                    "sentences": {
                        "type": "text",
                        "analyzer": "italian_analyzer"
                    },
                    "stopWords": {
                        "type": "keyword"
                    },
                    "summary": {
                        "type": "text",
                        "analyzer": "italian_analyzer"
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
