package es

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/cvcio/mediawatch/pkg/es/indeces"
	"github.com/elastic/go-elasticsearch/v7"
)

// Elastic struct implements elasticsearch client.
type Elastic struct {
	Client *elasticsearch.Client
}

// NewElasticsearch returns a new Elastic struct.
func NewElasticsearch(host, user, pass string) (*Elastic, error) {
	es, err := elasticsearch.NewClient(
		elasticsearch.Config{
			Addresses: []string{host},
			Username:  user,
			Password:  pass,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	)

	if err != nil {
		return nil, err
	}

	if _, err := es.Info(); err != nil {
		return nil, err
	}

	return &Elastic{es}, nil
}

// CreateIndex creates a new index.
func (es *Elastic) CreateIndex(name, template string) error {
	res, err := es.Client.Indices.Create(
		name,
		es.Client.Indices.Create.WithBody(strings.NewReader(template)),
	)
	if err != nil {
		return err
	}
	if res.IsError() {
		return fmt.Errorf("Cannot create index: %s", res)
	}
	res.Body.Close()
	return nil
}

// CheckIfIndexExists check if index exists.
func (es *Elastic) CheckIfIndexExists(name string) bool {
	res, err := es.Client.Indices.Exists([]string{name})
	if err != nil {
		return false
	}

	if res.IsError() {
		return false
	}

	res.Body.Close()
	if res.StatusCode != 200 {
		return false
	}
	return true
}

// CreateElasticIndexWithLanguages creates indexes for each language provided.
func (es *Elastic) CreateElasticIndexWithLanguages(prefix string, languages []string) error {
	for _, lang := range languages {
		index := prefix + "_" + strings.ToLower(lang) + "*"
		template := getIndex(lang)
		if es.CheckIfIndexExists(index) == true {
			continue
		}

		if err := es.CreateIndex(index, template); err != nil {
			return fmt.Errorf("Error creating index %s for lang %s: %v", index, lang, err)
		}
	}

	return nil
}

func getIndex(language string) string {
	switch strings.ToLower(language) {
	case "bg":
		return indeces.ArticlesBG
	case "de":
		return indeces.ArticlesDE
	case "el":
		return indeces.ArticlesEL
	case "en":
		return indeces.ArticlesEN
	case "fi":
		return indeces.ArticlesFI
	case "fr":
		return indeces.ArticlesFR
	case "hu":
		return indeces.ArticlesHU
	case "it":
		return indeces.ArticlesIT
	case "nl":
		return indeces.ArticlesNL
	case "no":
		return indeces.ArticlesNO
	case "pt":
		return indeces.ArticlesPT
	case "ro":
		return indeces.ArticlesRO
	case "ru":
		return indeces.ArticlesRU
	case "tr":
		return indeces.ArticlesTR
	default:
		return indeces.ArticlesEN
	}
}
