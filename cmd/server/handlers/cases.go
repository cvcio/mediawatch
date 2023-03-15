package handlers

import (
	"net/http"

	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/es"
	"github.com/cvcio/mediawatch/pkg/neo"
	"github.com/cvcio/mediawatch/pkg/web"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

type Cases struct {
	log   *logrus.Logger
	db    *db.MongoDB
	es    *es.Elastic
	neo4j *neo.Neo
}

func NewCasesHandler(log *logrus.Logger, db *db.MongoDB, es *es.Elastic, neo4j *neo.Neo) *Cases {
	return &Cases{
		log:   log,
		db:    db,
		es:    es,
		neo4j: neo4j,
	}
}

func (handler *Cases) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Render(w, r, web.ErrInternalError)
	}
}
func (handler *Cases) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Render(w, r, web.ErrInternalError)
	}
}
func (handler *Cases) Count() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Render(w, r, web.ErrInternalError)
	}
}
