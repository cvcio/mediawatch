package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// Config struct holds all the configuration elements for our apps
type Config struct {
	Env string `envconfig:"ENV" default:"production"`
	Log struct {
		Level string `envconfig:"LOG_LEVEL" default:"debug"`
		Path  string `envconfig:"LOG_PATH" default:"logs"`
	}
	Api struct {
		Host            string        `default:"0.0.0.0" envconfig:"API_HOST"`
		Port            string        `default:"8000" envconfig:"API_PORT"`
		ReadTimeout     time.Duration `default:"10s" envconfig:"API_READ_TIMEOUT"`
		WriteTimeout    time.Duration `default:"20s" envconfig:"API_WRITE_TIMEOUT"`
		ShutdownTimeout time.Duration `default:"10s" envconfig:"API_SHUTDOWN_TIMEOUT"`
		Debug           bool          `default:"false" envconfig:"API_DEBUG"`
	}
	Prometheus struct {
		Host            string        `default:"0.0.0.0" envconfig:"PROM_HOST"`
		Port            string        `default:"9000" envconfig:"PROM_PORT"`
		ReadTimeout     time.Duration `default:"10s" envconfig:"PROM_READ_TIMEOUT"`
		WriteTimeout    time.Duration `default:"20s" envconfig:"PROM_WRITE_TIMEOUT"`
		ShutdownTimeout time.Duration `default:"10s" envconfig:"PROM_SHUTDOWN_TIMEOUT"`
	}
	Service struct {
		Network         string        `envconfig:"SERVICE_NETWORK" default:"tcp"`
		Host            string        `envconfig:"SERVICE_HOST" default:"0.0.0.0"`
		Port            string        `envconfig:"SERVICE_PORT" default:"50050"`
		ReadTimeout     time.Duration `envconfig:"SERVICE_READ_TIMEOUT" default:"10s"`
		WriteTimeout    time.Duration `envconfig:"SERVICE_WRITE_TIMEOUT" default:"20s"`
		ShutdownTimeout time.Duration `envconfig:"SERVICE_SHUTDOWN_TIMEOUT" default:"10s"`
		DomainName      string        `envconfig:"SERVICE_DOMAIN_NAME" default:"mediawatch.io"`
	}
	Twitter struct {
		TwitterConsumerKey       string `envconfig:"TWITTER_CONSUMER_KEY" default:""`
		TwitterConsumerSecret    string `envconfig:"TWITTER_CONSUMER_SECRET" default:""`
		TwitterAccessToken       string `envconfig:"TWITTER_ACCESS_TOKEN" default:""`
		TwitterAccessTokenSecret string `envconfig:"TWITTER_ACCESS_TOKEN_SECRET" default:""`
		TwitterRuleTag           string `envconfig:"TWITTER_RULE_TAG" default:"mediawatch-listener-v2"`
	}
	Scrape struct {
		Host string `envconfig:"SVC_SCRAPER" default:"0.0.0.0:50050"`
	}
	Enrich struct {
		Host string `envconfig:"SVC_ENRICH" default:"0.0.0.0:50030"`
	}
	Compare struct {
		Host string `envconfig:"SVC_COMPARE" default:"0.0.0.0:50040"`
	}
	Mongo struct {
		URL         string        `envconfig:"MONGO_URL" default:"mongodb://localhost:27017"`
		Path        string        `envconfig:"MONGO_PATH" default:"mediawatch"`
		User        string        `envconfig:"MONGO_USER" default:""`
		Pass        string        `envconfig:"MONGO_PASS" default:""`
		DialTimeout time.Duration `envconfig:"DIAL_TIMEOUT" default:"5s"`
	}
	Elasticsearch struct {
		Host  string `envconfig:"ES_HOST" default:"http://localhost:9200"`
		User  string `envconfig:"ES_USER" default:""`
		Pass  string `envconfig:"ES_PASS" default:""`
		Index string `envconfig:"ES_INDEX" default:"mediawatch_articles"`
	}
	Neo struct {
		URL  string `envconfig:"NEO_URL" default:"http://neo4j:neo4j@localhost:7474/db/data/"`
		BOLT string `envconfig:"NEO_BOLT" default:"bolt://localhost:7687"`
		User string `envconfig:"NEO_USER" default:"neo4j"`
		Pass string `envconfig:"NEO_PASS" default:"neo4j"`
	}
	Redis struct {
		Host string `envconfig:"REDIS_HOST" default:"localhost"`
		Port string `envconfig:"REDIS_PORT" default:"6379"`
	}
	Kafka struct {
		Enable               bool   `envconfig:"KAFKA_ENABLE" default:"false"`
		Broker               string `envconfig:"KAFKA_BROKER" default:"kafka:9092"`
		WorkerTopic          string `envconfig:"KAFKA_TOPIC_WORKER" default:"worker"`
		CompareTopic         string `envconfig:"KAFKA_TOPIC_COMPARE" default:"compare"`
		ConsumerGroupWorker  string `envconfig:"KAFKA_CONSUMER_GROUP_WORKERS" default:"mw-worker"`
		ConsumerGroupCompare string `envconfig:"KAFKA_CONSUMER_GROUP_COMPARE" default:"mw-compare"`
		Version              string `envconfig:"KAFKA_VERSION" default:"2.5.0"`
		ConsumerTopic        string `envconfig:"KAFKA_CONSUMER_TOPIC" default:""`
		ProducerTopic        string `envconfig:"KAFKA_PRODUCER_TOPIC" default:""`
		ConsumerGroup        string `envconfig:"KAFKA_CONSUMER_GROUP" default:""`
		ProducerGroup        string `envconfig:"KAFKA_PRODUCER_GROUP" default:""`
		WorkerOffsetOldest   bool   `envconfig:"KAFKA_WORKER_OFFSET_OLDEST" default:"false"`
		AckBefore            string `envconfig:"KAFKA_ACK_BEFORE" default:""`
	}
	SMTP struct {
		Server   string `envconfig:"SMTP_SERVER" default:"smtp"`
		Port     int    `envconfig:"SMTP_PORT" default:"587"`
		User     string `envconfig:"SMTP_USER" default:"no-reply@mediawatch.io"`
		From     string `envconfig:"SMTP_FROM" default:"no-reply@mediawatch.io"`
		FromName string `envconfig:"SMTP_FROM_NAME" default:"MediaWatch"`
		Pass     string `envconfig:"SMTP_PASS" default:""`
		Reply    string `envconfig:"SMTP_REPLY" default:"press@mediawatch.io"`
	}
	Twillio struct {
		Enabled bool   `envconfig:"TWILIO_ENABLED" default:"true"`
		SID     string `envconfig:"TWILIO_SID" default:""`
		Token   string `envconfig:"TWILIO_TOKEN" default:""`
	}
	Auth struct {
		Enabled         bool   `envconfig:"AUTH_ENABLED" default:"true"`
		Authorizer      string `envconfig:"AUTH_AUTHORIZER" default:"mediawatch"`
		BaseCallbackURL string `envconfig:"AUTH_BASE_CALLBACK_URL" default:"http://localhost:8080"`
		Debug           bool   `envconfig:"AUTH_DEBUG" default:"false" `
		Domain          string `envconfig:"DOMAIN_NAME" default:"mediawatch.io"`
		Hash            string `envconfig:"HASH" default:"123"`
		KeyID           string `envconfig:"KEY_ID" default:"0123456789abcdef"`
		PrivateKeyFile  string `envconfig:"PRIVATE_KEY_FILE" default:"private.pem"`
		Algorithm       string `envconfig:"ALGORITHM" default:"RS256"`
	}
	Google struct {
		Enabled      bool   `envconfig:"GOOGLE_ENABLED" default:"true"`
		CallBackURL  string `envconfig:"GOOGLE_AUTH_CB_URL" default:"http://localhost:8000/auth/authorize/google/callback"`
		ClientID     string `envconfig:"GOOGLE_AUTH_CLIENT_ID" default:"1"`
		ClientSecret string `envconfig:"GOOGLE_AUTH_CLIENT_SECRET" default:"1"`
	}
	Github struct {
		Enabled      bool   `envconfig:"GITHUB_ENABLED" default:"true"`
		CallBackURL  string `envconfig:"GITHUB_AUTH_CB_URL" default:"http://localhost:8000/auth/authorize/github/callback"`
		ClientID     string `envconfig:"GITHUB_AUTH_CLIENT_ID" default:""`
		ClientSecret string `envconfig:"GITHUB_AUTH_CLIENT_SECRET" default:""`
	}
	Stripe struct {
		Enabled bool   `envconfig:"STRIPE_ENABLED" default:"true"`
		Key     string `envconfig:"STRIPE_KEY" default:""`
	}
	Proxy struct {
		Enabled   bool   `envconfig:"PROXY_ENABLED" default:"true"`
		Host      string `envconfig:"PROXY_HOST" default:""`
		Port      string `envconfig:"PROXY_PORT" default:""`
		UserName  string `envconfig:"PROXY_USERNAME" default:""`
		Password  string `envconfig:"PROXY_PASSWORD" default:""`
		ProxyList string `envconfig:"PROXY_LIST" default:""`
		Scheme    string `envconfig:"PROXY_SCHEME" default:"https"`
	}
	Streamer struct {
		Init     bool          `envconfig:"STREAMER_INIT" default:"false"`
		Type     string        `envconfig:"STREAMER_TYPE" default:"rss"`
		Lang     string        `envconfig:"STREAMER_LANG" default:"el"`
		Size     int           `envconfig:"STREAMER_SIZE" default:"3000"`
		Chunks   int           `envconfig:"STREAMER_CHUNKS" default:"100"`
		Interval time.Duration `envconfig:"STREAMER_INTERVAL" default:"20m"`
	}
	Langs []string `envconfig:"LANGS" default:"el"`
}

func NewConfig() *Config {
	return new(Config)
}

func NewConfigFromEnv() *Config {
	cfg := NewConfig()

	err := envconfig.Process("", cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

func (c *Config) GetKafkaBrokers() []string {
	return strings.Split(c.Kafka.Broker, ",")
}

func (c *Config) GetMongoURL() string {
	return fmt.Sprintf("%s/%s", c.Mongo.URL, c.Mongo.Path)
}

func (c *Config) GetElasticsearchURL() string {
	return c.Elasticsearch.Host
}

func (c *Config) GetRedisURL() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}

func (c *Config) GetApiURL() string {
	return fmt.Sprintf("%s:%s", c.Api.Host, c.Api.Port)
}

func (c *Config) GetServiceURL() string {
	return fmt.Sprintf("%s:%s", c.Service.Host, c.Service.Port)
}

func (c *Config) GetPrometheusURL() string {
	return fmt.Sprintf("%s:%s", c.Prometheus.Host, c.Prometheus.Port)
}

func (c *Config) GetProxyURL() string {
	if c.Proxy.UserName == "" && c.Proxy.Password == "" {
		return fmt.Sprintf("%s://%s:%s", c.Proxy.Scheme, c.Proxy.Host, c.Proxy.Port)
	}

	return fmt.Sprintf("%s://%s:%s@%s:%s", c.Proxy.Scheme, c.Proxy.UserName, c.Proxy.Password, c.Proxy.Host, c.Proxy.Port)
}

func (c *Config) GetProxyList() []string {
	return strings.Split(c.Proxy.ProxyList, ",")
}

func (c *Config) ExternalAuths() (map[string]*oauth2.Config, error) {
	var auths = make(map[string]*oauth2.Config)
	g, err := c.EnableGoogle()
	if err == nil {
		auths["google"] = g
	}
	return auths, nil
}

func (c *Config) EnableGoogle() (*oauth2.Config, error) {
	if c.Google.ClientID == "" || c.Google.ClientSecret == "" {
		return nil, fmt.Errorf("Invalid google creds")
	}
	return &oauth2.Config{
		RedirectURL:  c.Google.CallBackURL,
		ClientID:     c.Google.ClientID,
		ClientSecret: c.Google.ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
			"openid",
		},
		Endpoint: google.Endpoint,
	}, nil
}

func (c *Config) EnableGithub() (*oauth2.Config, error) {
	if c.Github.ClientID == "" || c.Github.ClientSecret == "" {
		return nil, fmt.Errorf("Invalid Github creds")
	}
	return &oauth2.Config{
		RedirectURL:  c.Github.CallBackURL,
		ClientID:     c.Github.ClientID,
		ClientSecret: c.Github.ClientSecret,
		Scopes: []string{
			"https://www.github.com/auth/userinfo.profile",
			"https://www.github.com/auth/userinfo.email",
			"openid",
		},
		Endpoint: github.Endpoint,
	}, nil
}
