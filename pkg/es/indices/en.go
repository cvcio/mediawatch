package indices

var ArticlesEN = `
{
    "settings": {
        "index": {
            "number_of_shards": 1,
            "number_of_replicas": 0
        }
    },
    "alias": {
        "mediawatch_en-{now/M{yyyy.MM}}": {}
    },
    "mappings": {
        "properties": {
            "content": {
                "properties": {
                    "authors": {
                        "type": "keyword"
                    },
                    "body": {
                        "type": "text"
                    },
                    "categories": {
                        "type": "keyword"
                    },
                    "editedAt": {
                        "type": "date",
                        "ignore_malformed": true
                    },
                    "excerpt": {
                        "type": "text"
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
                        "type": "text"
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
                        "type": "text"
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
                        "type": "keyword"
                    },
                    "quotes": {
                        "type": "text"
                    },
                    "sentences": {
                        "type": "text"
                    },
                    "stopWords": {
                        "type": "keyword"
                    },
                    "summary": {
                        "type": "text"
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
