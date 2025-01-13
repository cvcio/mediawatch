package es

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/olivere/elastic/v7"
)

type ES struct {
	Client *elastic.Client
	URL    string
}

// NewElastic Client from the given configuration.
//
// Deprecated: Use NewElasticsearch instead.
func NewElastic(host, user, pass string) (*ES, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c := &http.Client{Transport: tr}
	esclient, err := elastic.NewClient(
		elastic.SetURL(host),
		elastic.SetBasicAuth(user, pass),
		elastic.SetSniff(false),
		elastic.SetHttpClient(c),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
	)
	if err != nil {
		return nil, err
	}
	return &ES{Client: esclient, URL: host}, nil
}

// // CreateElasticIndex runs an elasticsearch container.
// func CreateElasticIndex(client *ES) error {
// 	// log.Println("Init indexes")
// 	mappings := map[string]string{
// 		"articles": indexArticles,
// 	}
//
// 	ctx := context.Background()
// 	for k, v := range mappings {
// 		createIndex, err := client.Client.CreateIndex(k).BodyJson(v).Do(ctx)
// 		if err != nil {
// 			log.Printf("Error creating mapping %s from file %s: %v", k, v, err)
// 			continue
// 		}
// 		if !createIndex.Acknowledged {
// 			log.Printf("Error mapping %s from file %s not acknowledged", k, v)
// 			continue
// 		}
// 	}
//
// 	return nil
// }
//
// // CreateElasticIndexArticles with index name
// func CreateElasticIndexArticles(client *ES, indexes []string) error {
// 	ctx := context.Background()
// 	for _, v := range indexes {
// 		// Test if already exists
// 		indexExists, err := client.Client.IndexExists(v).Do(context.TODO())
// 		if err != nil {
// 			log.Println(err)
// 			continue
// 		}
// 		if indexExists {
// 			continue
// 		}
// 		// Create Index
// 		createIndex, err := client.Client.CreateIndex(v).Body(indexArticles).Do(ctx)
// 		if err != nil {
// 			log.Printf("Error creating index %s from file %s", v, err)
// 			continue
// 		}
// 		if !createIndex.Acknowledged {
// 			log.Printf("Error mapping %s not acknowledged", v)
// 			continue
// 		}
// 	}
//
// 	return nil
// }
