package es

var indexArticles = `
{
  "settings": {
      "index": {
          "number_of_shards": 2,
          "number_of_replicas": 2
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
                  "keywords": {
                      "type": "text",
                      "fields": {
                          "keyword": {
                              "type": "keyword"
                          }
                      }
                  },
                  "stopWords": {
                      "type": "text"
                  },
                  "sentences": {
                      "type": "text",
                      "analyzer": "greek_analyzer"
                  },
                  "summary": {
                      "type": "text",
                      "analyzer": "greek_analyzer"
                  },
                  "entities": {
                      "properties": {
                          "entity_text": {
                              "type": "text",
                              "fields": {
                                  "keyword": {
                                      "type": "keyword"
                                  }
                              }
                          },
                          "entity_type": {
                              "type": "text",
                              "fields": {
                                  "keyword": {
                                      "type": "keyword"
                                  }
                              }
                          }
                      }
                  },
                  "tokens": {
                      "type": "text"
                  }
              }
          },
          "nodeId": {
              "type": "long"
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
