package twitter

import (
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	"go.uber.org/zap"
)

// Demux receives channels or interfaces and type switches them to call the appropriate handle function.
type Demux interface {
	Handle(m interface{})
	HandleChan(c <-chan interface{})
}

// StreamDemux receives messages and type switches them to call functions with typed messages.
type StreamDemux struct {
	All        func(m interface{})
	Tweet      func(tweet anaconda.Tweet)
	Event      func(event anaconda.Event)
	EventTweet func(event anaconda.EventTweet)
	Other      func(m interface{})
}

// NewAPI creates a new anaconda instance.
func NewAPI(consumerkey string, consumersecret string, accesstoken string, accesstokensecret string) (*anaconda.TwitterApi, error) {
	api := anaconda.NewTwitterApiWithCredentials(
		accesstoken,
		accesstokensecret,
		consumerkey,
		consumersecret,
	)
	if _, err := api.VerifyCredentials(); err != nil {
		return nil, err
	}
	return api, nil
}

// NewStreamDemux initializes a new StreamDemux.
func NewStreamDemux() StreamDemux {
	return StreamDemux{
		All:        func(m interface{}) {},
		Tweet:      func(tweet anaconda.Tweet) {},
		Event:      func(event anaconda.Event) {},
		EventTweet: func(event anaconda.EventTweet) {},
		Other:      func(m interface{}) {},
	}
}

// Handle handles messages.
func (d StreamDemux) Handle(m interface{}) {
	d.All(m)

	switch t := m.(type) {
	case anaconda.Tweet:
		d.Tweet(t)
	case anaconda.Event:
		d.Event(t)
	case anaconda.EventTweet:
		d.EventTweet(t)
	default:
		d.Other(t)
	}
}

// HandleChan handles channels.
func (d StreamDemux) HandleChan(c <-chan interface{}) {
	for m := range c {
		d.Handle(m)
	}
}

// Listen struct.
type Listen struct {
	TwitterAPI *anaconda.TwitterApi
	log        *zap.SugaredLogger
	stream     *anaconda.Stream
}

// NewListener return a new Listener service, given a twitter api client
func NewListener(tw *anaconda.TwitterApi, log *zap.SugaredLogger, opts ...ListenOpts) (*Listen, error) {
	s := new(Listen)
	s.TwitterAPI = tw
	s.log = log
	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

type ListenOpts func(*Listen) error

func WithPublicStream(values map[string][]string) ListenOpts {
	return func(s *Listen) error {
		v := url.Values(values)
		s.stream = s.TwitterAPI.PublicStreamFilter(v)
		return nil
	}
}

// TweetListen start the listener and send cached urls to chan
func (s *Listen) TweetListen(f func(anaconda.Tweet)) {
	demux := NewStreamDemux()
	demux.Tweet = f
	demux.HandleChan(s.stream.C)
}

func (s *Listen) EventListen(f func(anaconda.Event)) {
	demux := NewStreamDemux()
	demux.Event = f
	demux.HandleChan(s.stream.C)
}

func (s *Listen) EventTweetListen(f func(anaconda.EventTweet)) {
	demux := NewStreamDemux()
	demux.EventTweet = f
	demux.HandleChan(s.stream.C)
}

func IsTweet(t anaconda.Tweet) bool {
	if t.InReplyToStatusIdStr == "" && t.InReplyToUserIdStr == "" &&
		t.RetweetedStatus == nil && t.QuotedStatus == nil {
		return true
	}
	return false
}

func GetUsersLookup(twtt *anaconda.TwitterApi, users []string) ([]string, error) {
	ids := make([]string, 0)

	for _, uu := range users {
		twitterAccount, err := twtt.GetUsersShow(uu, url.Values{})
		if err != nil {
			continue
		}
		ids = append(ids, twitterAccount.IdStr)
	}

	return ids, nil
}
