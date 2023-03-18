package handlers

import (
	"context"

	"github.com/bufbuild/connect-go"
	commonv2 "github.com/cvcio/mediawatch/internal/mediawatch/common/v2"
	feedsv2 "github.com/cvcio/mediawatch/internal/mediawatch/feeds/v2"
	"github.com/cvcio/mediawatch/internal/mediawatch/feeds/v2/feedsv2connect"
	"github.com/cvcio/mediawatch/models/feed"
	"github.com/cvcio/mediawatch/pkg/auth"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/es"
	"go.uber.org/zap"
)

// FeedsHandler implements feeds connect service
type FeedsHandler struct {
	log           *zap.SugaredLogger
	mg            *db.MongoDB
	elastic       *es.Elastic
	authenticator *auth.JWTAuthenticator
	// Embed the unimplemented server
	feedsv2connect.UnimplementedFeedServiceHandler
}

// NewFeedsHandler returns a new FeedsHandler service.
func NewFeedsHandler(cfg *config.Config, log *zap.SugaredLogger, mg *db.MongoDB, elastic *es.Elastic, authenticator *auth.JWTAuthenticator) *FeedsHandler {
	return &FeedsHandler{log: log, mg: mg, elastic: elastic, authenticator: authenticator}
}

// CreateFeed creates a new feed.
func (h *FeedsHandler) CreateFeed(ctx context.Context, req *connect.Request[feedsv2.Feed]) (*connect.Response[feedsv2.Feed], error) {
	return connect.NewResponse(&feedsv2.Feed{}), nil
}

// GetFeed returns a single feed.
func (h *FeedsHandler) GetFeed(ctx context.Context, req *connect.Request[feedsv2.QueryFeed]) (*connect.Response[feedsv2.Feed], error) {
	return connect.NewResponse(&feedsv2.Feed{}), nil
}

// GetFeeds return a list of feeds.
func (h *FeedsHandler) GetFeeds(ctx context.Context, req *connect.Request[feedsv2.QueryFeed]) (*connect.Response[feedsv2.FeedList], error) {
	return connect.NewResponse(&feedsv2.FeedList{}), nil
}

// UpdateFeed updates a single feed.
func (h *FeedsHandler) UpdateFeed(ctx context.Context, req *connect.Request[feedsv2.Feed]) (*connect.Response[commonv2.ResponseWithMessage], error) {
	// TODO: parse claims

	if err := feed.Update(ctx, h.mg, req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&commonv2.ResponseWithMessage{
		Status:  "OK",
		Message: "UPDATE",
	}), nil
}

// UpdateFeedWithFields updates a single feed with given feilds.
func (h *FeedsHandler) UpdateFeedWithFields(ctx context.Context, req *connect.Request[feedsv2.FeedWithFields]) (*connect.Response[commonv2.ResponseWithMessage], error) {
	return connect.NewResponse(&commonv2.ResponseWithMessage{}), nil
}

// DeleteFeed deletes a single article.
func (h *FeedsHandler) DeleteFeed(ctx context.Context, req *connect.Request[feedsv2.Feed]) (*connect.Response[commonv2.ResponseWithMessage], error) {
	return connect.NewResponse(&commonv2.ResponseWithMessage{}), nil
}

// GetFeedsStreamList returns a stream list.
func (h *FeedsHandler) GetFeedsStreamList(ctx context.Context, req *connect.Request[feedsv2.QueryFeed]) (*connect.Response[feedsv2.FeedStreamList], error) {
	return connect.NewResponse(&feedsv2.FeedStreamList{}), nil
}
