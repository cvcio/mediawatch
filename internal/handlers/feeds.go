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
	"github.com/pkg/errors"
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
func NewFeedsHandler(cfg *config.Config, log *zap.SugaredLogger, mg *db.MongoDB, elastic *es.Elastic, authenticator *auth.JWTAuthenticator) (*FeedsHandler, error) {
	if err := feed.EnsureIndex(context.Background(), mg); err != nil {
		return nil, err
	}
	return &FeedsHandler{log: log, mg: mg, elastic: elastic, authenticator: authenticator}, nil
}

// CreateFeed creates a new feed.
func (h *FeedsHandler) CreateFeed(ctx context.Context, req *connect.Request[feedsv2.Feed]) (*connect.Response[feedsv2.Feed], error) {
	h.log.Debugf("CreateFeed Request Message: %+v", req.Msg)
	// TODO: parse claims and authorization tokens

	f, err := feed.Create(ctx, h.mg, req.Msg)
	if err != nil {
		errorMessage := connect.NewError(connect.CodeInternal, errors.Errorf("unable to create feed"))
		h.log.Errorf("Internal: %s", err.Error())
		return nil, errorMessage
	}

	return connect.NewResponse(f), nil
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
	return connect.NewResponse(&commonv2.ResponseWithMessage{}), nil
}

// UpdateFeedWithFields updates a single feed with given feilds.
func (h *FeedsHandler) UpdateFeedWithFields(ctx context.Context, req *connect.Request[feedsv2.FeedWithFields]) (*connect.Response[commonv2.ResponseWithMessage], error) {
	return connect.NewResponse(&commonv2.ResponseWithMessage{}), nil
}

// DeleteFeed deletes a single article.
func (h *FeedsHandler) DeleteFeed(ctx context.Context, req *connect.Request[feedsv2.Feed]) (*connect.Response[commonv2.ResponseWithMessage], error) {
	h.log.Debugf("DeleteFeed Request Message: %+v", req.Msg)
	// TODO: parse claims and authorization tokens

	if err := feed.Delete(ctx, h.mg, req.Msg); err != nil {
		errorMessage := connect.NewError(connect.CodeInternal, errors.Errorf("unable to delete feed"))
		h.log.Errorf("Internal: %s", err.Error())
		return nil, errorMessage
	}

	return connect.NewResponse(&commonv2.ResponseWithMessage{
		Status:  "ok",
		Message: "feed deleted",
	}), nil
}

// GetFeedsStreamList returns a stream list.
func (h *FeedsHandler) GetFeedsStreamList(ctx context.Context, req *connect.Request[feedsv2.QueryFeed]) (*connect.Response[feedsv2.FeedList], error) {
	h.log.Debugf("GetFeedsStreamList Request Message: %+v", req.Msg)
	// TODO: parse claims and authorization tokens

	if req.Msg.StreamType != commonv2.StreamType_STREAM_TYPE_RSS && req.Msg.StreamType != commonv2.StreamType_STREAM_TYPE_TWITTER {
		return nil, connect.NewError(connect.CodeInternal, errors.Errorf("only twitter and rss streams are supported"))
	}

	data, err := feed.GetFeedsStreamList(ctx, h.mg,
		// return only active streams
		feed.StreamStatus(int(commonv2.Status_STATUS_ACTIVE.Number())),
		feed.StreamType(int(req.Msg.StreamType.Number())),
		feed.Lang(req.Msg.Lang),
	)
	if err != nil {
		errorMessage := connect.NewError(connect.CodeInternal, errors.Errorf("unable to retrieve feeds stream list"))
		h.log.Errorf("Internal: %s", err.Error())
		return nil, errorMessage
	}

	return connect.NewResponse(&feedsv2.FeedList{
		Data: data,
	}), nil
}
