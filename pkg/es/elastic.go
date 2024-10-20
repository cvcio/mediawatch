package es

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/cvcio/mediawatch/pkg/es/indices"
	"github.com/elastic/go-elasticsearch/v8"
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
			DiscoverNodesOnStart: true,
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
		return fmt.Errorf("cannot create index: %s", res)
	}

	return res.Body.Close()
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

	defer func() { _ = res.Body.Close() }()
	return res.StatusCode == 200
}

// CreateElasticIndexWithLanguages creates indexes for each language provided.
func (es *Elastic) CreateElasticIndexWithLanguages(prefix string, languages []string) error {
	for _, lang := range languages {
		index := prefix + "_" + strings.ToLower(lang)
		template := getIndex(lang)
		if es.CheckIfIndexExists(index) == true {
			continue
		}

		if err := es.CreateIndex(index, template); err != nil {
			return fmt.Errorf("error creating index %s for lang %s: %v", index, lang, err)
		}
	}

	return nil
}

// getIndex return the index mapping by language.
func getIndex(language string) string {
	switch strings.ToLower(language) {
	case "bg":
		return indices.ArticlesBG
	case "de":
		return indices.ArticlesDE
	case "el":
		return indices.ArticlesEL
	case "en":
		return indices.ArticlesEN
	case "fi":
		return indices.ArticlesFI
	case "fr":
		return indices.ArticlesFR
	case "hu":
		return indices.ArticlesHU
	case "it":
		return indices.ArticlesIT
	case "nl":
		return indices.ArticlesNL
	case "no":
		return indices.ArticlesNO
	case "pt":
		return indices.ArticlesPT
	case "ro":
		return indices.ArticlesRO
	case "ru":
		return indices.ArticlesRU
	case "tr":
		return indices.ArticlesTR
	default:
		return indices.ArticlesEN
	}
}
