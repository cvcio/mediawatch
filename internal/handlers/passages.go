package handlers

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/cvcio/mediawatch/models/passage"
	"github.com/cvcio/mediawatch/pkg/auth"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/es"
	passagesv2 "github.com/cvcio/mediawatch/pkg/mediawatch/passages/v2"
	"github.com/cvcio/mediawatch/pkg/mediawatch/passages/v2/passagesv2connect"
	"github.com/cvcio/mediawatch/pkg/redis"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type PassagesHandler struct {
	log           *zap.SugaredLogger
	mg            *db.MongoDB
	elastic       *es.Elastic
	authenticator *auth.JWTAuthenticator
	rdb           *redis.RedisClient
	// Embed the unimplemented server
	passagesv2connect.UnimplementedPassageServiceHandler
}

// NewPassagesHandler returns a new PassagesHandler service.
func NewPassagesHandler(cfg *config.Config, log *zap.SugaredLogger, mg *db.MongoDB, elastic *es.Elastic, authenticator *auth.JWTAuthenticator, rdb *redis.RedisClient) (*PassagesHandler, error) {
	if err := passage.EnsureIndex(context.Background(), mg); err != nil {
		return nil, err
	}
	return &PassagesHandler{log: log, mg: mg, elastic: elastic, authenticator: authenticator, rdb: rdb}, nil
}

// GetPassages return a list of passages.
func (h *PassagesHandler) GetPassages(ctx context.Context, req *connect.Request[passagesv2.QueryPassage]) (*connect.Response[passagesv2.PassageList], error) {
	h.log.Debugf("GetPassages Request Message: %+v", req.Msg)
	// TODO: parse claims and authorization tokens

	data, err := passage.List(ctx, h.mg,
		passage.Lang(req.Msg.Language),
	)

	if err != nil {
		errorMessage := connect.NewError(connect.CodeInternal, errors.Errorf("unable to retrieve passages"))
		h.log.Errorf("Internal: %s", err.Error())
		return nil, errorMessage
	}

	return connect.NewResponse(data), nil
}

// CreateFeed creates a new feed.
func (h *PassagesHandler) CreatePassage(ctx context.Context, req *connect.Request[passagesv2.Passage]) (*connect.Response[passagesv2.Passage], error) {
	h.log.Debugf("CreateFeed Request Message: %+v", req.Msg)

	now := time.Now()
	f, err := passage.Create(ctx, h.mg, req.Msg, now)
	if err != nil {
		errorMessage := connect.NewError(connect.CodeInternal, errors.Errorf("unable to create feed"))
		h.log.Errorf("Internal: %s", err.Error())
		return nil, errorMessage
	}

	return connect.NewResponse(f), nil
}
