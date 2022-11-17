package es

var indexRelationships = `

{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"relationships":{
        "properties": {
          "score": {
            "type": "float"
          },
          "source": {
            "properties": {
              "docId": {
                "type": "keyword"
              },
              "publishedAt": {
                "type": "date",
                "ignore_malformed": true,
                "null_value": "NULL"
              },
              "screenName": {
                "type": "keyword"
              },
              "tweetId": {
                "type": "long"
              }
            }
          },
          "target": {
            "properties": {
              "docId": {
                "type": "keyword"
              },
              "publishedAt": {
                "type": "date",
                "ignore_malformed": true,
                "null_value": "NULL"
              },
              "screenName": {
                "type": "keyword"
              },
              "tweetId": {
                "type": "long"
              }
            }
          }
        }
    }
  }
}
`
