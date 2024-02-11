package handlers

import (
	"context"
	"encoding/json"
	"time"

	"connectrpc.com/connect"
	"github.com/cvcio/mediawatch/models/feed"
	"github.com/cvcio/mediawatch/pkg/auth"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/es"
	"github.com/cvcio/mediawatch/pkg/interceptors"
	commonv2 "github.com/cvcio/mediawatch/pkg/mediawatch/common/v2"
	feedsv2 "github.com/cvcio/mediawatch/pkg/mediawatch/feeds/v2"
	"github.com/cvcio/mediawatch/pkg/mediawatch/feeds/v2/feedsv2connect"
	scrapev2 "github.com/cvcio/mediawatch/pkg/mediawatch/scrape/v2"
	"github.com/cvcio/mediawatch/pkg/redis"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// FeedsHandler implements feeds connect service
type FeedsHandler struct {
	log           *zap.SugaredLogger
	mg            *db.MongoDB
	elastic       *es.Elastic
	authenticator *auth.JWTAuthenticator
	rdb           *redis.RedisClient
	// Embed the unimplemented server
	feedsv2connect.UnimplementedFeedServiceHandler
}

// NewFeedsHandler returns a new FeedsHandler service.
func NewFeedsHandler(cfg *config.Config, log *zap.SugaredLogger, mg *db.MongoDB, elastic *es.Elastic, authenticator *auth.JWTAuthenticator, rdb *redis.RedisClient) (*FeedsHandler, error) {
	if err := feed.EnsureIndex(context.Background(), mg); err != nil {
		return nil, err
	}
	return &FeedsHandler{log: log, mg: mg, elastic: elastic, authenticator: authenticator, rdb: rdb}, nil
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
	h.log.Debugf("GetFeed Request Message: %+v", req.Msg)
	// TODO: parse claims and authorization tokens

	data, err := feed.Get(ctx, h.mg,
		feed.Id(req.Msg.Id),
		feed.Hostname(req.Msg.Hostname),
		feed.UserName(req.Msg.UserName),
	)

	if err != nil {
		errorMessage := connect.NewError(connect.CodeInternal, errors.Errorf("unable to retrieve feed"))
		h.log.Errorf("Internal: %s", err.Error())
		return nil, errorMessage
	}

	data.Stream.State = commonv2.State_STATE_UNSPECIFIED
	state, err := h.rdb.Get("feed:status:" + data.Id)
	if err != nil {
		data.Stream.State = commonv2.State_STATE_OK
	}

	if state == "offline" {
		data.Stream.State = commonv2.State_STATE_NOT_OK
	}

	return connect.NewResponse(data), nil
}

// GetFeeds return a list of feeds.
func (h *FeedsHandler) GetFeeds(ctx context.Context, req *connect.Request[feedsv2.QueryFeed]) (*connect.Response[feedsv2.FeedList], error) {
	h.log.Debugf("GetFeed Request Message: %+v", req.Msg)
	// TODO: parse claims and authorization tokens

	data, err := feed.List(ctx, h.mg,
		feed.Q(req.Msg.Q),
		feed.StreamStatus(int(req.Msg.StreamStatus.Number())),
		feed.StreamType(int(req.Msg.StreamType.Number())),
		feed.Lang(req.Msg.Lang),
		feed.Country(req.Msg.Country),
		feed.Limit(int(req.Msg.Limit)),
		feed.Offset(int(req.Msg.Offset)),
		feed.SortOrder(int(req.Msg.SortOrder)),
		feed.SortKey(req.Msg.SortKey),
	)

	if err != nil {
		errorMessage := connect.NewError(connect.CodeInternal, errors.Errorf("unable to retrieve feed"))
		h.log.Errorf("Internal: %s", err.Error())
		return nil, errorMessage
	}

	for _, f := range data.Data {
		f.Stream.State = commonv2.State_STATE_UNSPECIFIED
		state, err := h.rdb.Get("feed:status:" + f.Id)
		if err != nil {
			f.Stream.State = commonv2.State_STATE_OK
		}

		if state == "offline" {
			f.Stream.State = commonv2.State_STATE_NOT_OK
		}
	}

	return connect.NewResponse(data), nil
}

// UpdateFeed updates a single feed.
func (h *FeedsHandler) UpdateFeed(ctx context.Context, req *connect.Request[feedsv2.Feed]) (*connect.Response[commonv2.ResponseWithMessage], error) {
	h.log.Debugf("UpdateFeed Request Message: %+v", req.Msg)
	// TODO: parse claims and authorization tokens

	if err := feed.Update(ctx, h.mg, req.Msg); err != nil {
		errorMessage := connect.NewError(connect.CodeInternal, errors.Errorf("unable to update feed"))
		h.log.Errorf("Internal: %s", err.Error())
		return nil, errorMessage
	}

	return connect.NewResponse(&commonv2.ResponseWithMessage{
		Status:  "ok",
		Message: "feed updated",
	}), nil
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

// TestFeed tests a feed.
func (h *FeedsHandler) TestFeed(ctx context.Context, req *connect.Request[feedsv2.Feed]) (*connect.Response[feedsv2.FeedTest], error) {
	h.log.Debugf("TestFeed Request Message: %+v", req.Msg)
	cfg := config.NewConfigFromEnv()

	// Create the gRPC service clients
	// Parse Server Options
	var grpcOptions []grpc.DialOption
	grpcOptions = append(
		grpcOptions,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptors.TimeoutInterceptor(60*time.Second)),
	)
	// Create gRPC Scrape Connection
	scrapeGRPC, err := grpc.Dial(cfg.Scrape.Host, grpcOptions...)
	if err != nil {
		errorMessage := connect.NewError(connect.CodeInternal, errors.Errorf("unable to connect to scrape service"))
		h.log.Errorf("Internal: %s", err.Error())
		return nil, errorMessage
	}

	defer scrapeGRPC.Close()

	// Create gRPC Scrape client
	scrape := scrapev2.NewScrapeServiceClient(scrapeGRPC)

	// scraper client
	feedString, _ := json.Marshal(req.Msg)

	// create the scrape request
	scrapeReq := scrapev2.ScrapeRequest{
		Feed: string(feedString),
		Url:  req.Msg.Test.Url,
		Lang: req.Msg.Localization.Lang,
	}
	// scrape the article
	scrapeResp, err := scrape.Scrape(context.Background(), &scrapeReq)
	if err != nil {
		errorMessage := connect.NewError(connect.CodeInternal, errors.Errorf("unable to scrape url: %s. Error: %s", req.Msg.Test.Url, err.Error()))
		h.log.Errorf("Internal: %s", err.Error())
		return nil, errorMessage
	}

	return connect.NewResponse(&feedsv2.FeedTest{
		Title:       scrapeResp.Data.Content.Title,
		Body:        scrapeResp.Data.Content.Body,
		Authors:     scrapeResp.Data.Content.Authors,
		Tags:        scrapeResp.Data.Content.Tags,
		PublishedAt: scrapeResp.Data.Content.PublishedAt,
		Description: scrapeResp.Data.Content.Description,
		Image:       scrapeResp.Data.Content.Image,
		Status:      scrapeResp.Status,
	}), nil
}
