package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cvcio/mediawatch/internal/startup"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/es"
	mailer "github.com/cvcio/mediawatch/pkg/mailer/v1"
	"github.com/cvcio/mediawatch/pkg/mid"
	"github.com/cvcio/mediawatch/pkg/neo"
	"github.com/cvcio/mediawatch/pkg/redis"
	"github.com/cvcio/mediawatch/pkg/twillio"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RegisterRoutes registers all the routes for the application.
func RegisterRoutes(
	cfg *config.Config,
	log *zap.Logger,
	mongoDB *db.MongoDB,
	elasticClient *es.Elastic,
	neoClient *neo.Neo,
	redisClient *redis.RedisClient,
) http.Handler {
	// ========================================
	// Create a server
	router := gin.New()
	if cfg.Env != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	router.RedirectTrailingSlash = true
	router.RedirectFixedPath = true
	router.ForwardedByClientIP = true
	router.UseH2C = true

	router.Use(gin.Recovery())
	router.Use(ginzap.Ginzap(log, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(log, true))
	router.Use(mid.EnableCORS(cfg.Env))

	// ========================================
	// Create authenticator
	authenticator, err := startup.GetAuthenticator(cfg)
	if err != nil {
		log.Sugar().Fatalf("main: Authenticator: %v", err)
	}

	// Create GoogleAuth
	externalAuths, err := cfg.ExternalAuths()
	if err != nil {
		log.Error("main: ExternalAuths: %v", zap.Error(err))
	}
	log.Info("main: Created auth keys")

	storeService := persistence.NewRedisCache(cfg.GetRedisURL(), "", time.Minute)
	log.Info("main: Created cache store")

	// Create mail service
	mailService := mailer.New(
		cfg.SMTP.Server,
		cfg.SMTP.Port,
		cfg.SMTP.User,
		cfg.SMTP.Pass,
		cfg.SMTP.From,
		cfg.SMTP.FromName,
		cfg.SMTP.Reply,
	)
	log.Info("main: Created mail service")

	smsService := twillio.New(cfg.Twillio.SID, cfg.Twillio.Token, "MediaWatch")
	log.Info("main: Created Twillio SMS service")

	accountsHandler, err := NewAccountsHandler(log, mongoDB, authenticator, externalAuths, mailService, smsService)
	if err != nil {
		log.Sugar().Fatal(err)
	}
	log.Info("main: Created accounts handler")

	organizationsHandler, err := NewOrganizationsHandler(log, mongoDB, mailService)
	if err != nil {
		log.Sugar().Fatal(err)
	}
	log.Info("main: Created organizations handler")

	articlesHandler, err := NewArticlesHandler(log, mongoDB, elasticClient, neoClient)
	if err != nil {
		log.Sugar().Fatal(err)
	}
	log.Info("main: Created articles handler")

	casesHandler, err := NewCasesHandler(log, mongoDB, elasticClient, neoClient)
	if err != nil {
		log.Sugar().Fatal(err)
	}
	log.Info("main: Created cases handler")

	passagesHandler, err := NewPassagesHandler(log, mongoDB)
	if err != nil {
		log.Sugar().Fatal(err)
	}
	log.Info("main: Created passages handler")

	// r := Reports{
	// 	DB:  db,
	// 	log: log,
	// }

	// s := Subscription{
	// 	DB:        db,
	// 	log:       log,
	// 	StripeKey: cfg.Stripe.Key,
	// }

	authMiddleware := mid.Auth{
		Authenticator: authenticator,
	}

	v2 := router.Group("/v2")

	v2Accounts := v2.Group("/accounts")

	v2Accounts.GET("/", authMiddleware.Authenticate(), accountsHandler.List)
	v2Accounts.POST("/", accountsHandler.Create)
	v2Accounts.GET("/:id", authMiddleware.Authenticate(), accountsHandler.Get)
	v2Accounts.PUT("/:id", authMiddleware.Authenticate(), accountsHandler.Update)
	v2Accounts.DELETE("/:id", authMiddleware.Authenticate(), accountsHandler.Delete)

	v2Auth := v2.Group("/auth")

	v2Auth.POST("/login", accountsHandler.Login)
	v2Auth.POST("/register", accountsHandler.Create)
	v2Auth.PUT("/register/:id", accountsHandler.Update)
	v2Auth.POST("/verify/:id", accountsHandler.Verify)
	v2Auth.POST("/token", accountsHandler.Token)
	v2Auth.POST("/reset", accountsHandler.Reset)
	v2Auth.POST("/reset/verify/:id", accountsHandler.ResetVerify)

	if len(accountsHandler.oauth2Providers) > 0 {
		for k := range accountsHandler.oauth2Providers {
			if k == "google" {
				v2Auth.GET(fmt.Sprintf("/v2/auth/%s", k), accountsHandler.OAuthGoogle)
				v2Auth.GET(fmt.Sprintf("/v2/auth/%s/callback", k), accountsHandler.OAuthGoogleCB)
			}
		}
	}

	v2Orgs := v2.Group("/orgs")

	v2Orgs.Use(authMiddleware.Authenticate())
	v2Orgs.GET("/", organizationsHandler.List)
	v2Orgs.POST("/", organizationsHandler.Create)
	v2Orgs.GET("/:id", organizationsHandler.Get)
	v2Orgs.PUT("/:id", organizationsHandler.Update)
	v2Orgs.DELETE("/:id", organizationsHandler.Delete)

	v2Orgs.PUT("/:id/members", organizationsHandler.UpsertMember)
	// v2Orgs.PUT("/:id/members/:memberId", organizationsHandler.UpsertMember)
	v2Orgs.DELETE("/:id/members/:memberId", organizationsHandler.RemoveMember)
	v2Orgs.PUT("/:id/invitations/:memberId/:method", organizationsHandler.ProcessInvitation)

	// router.GET("/feeds", f.List)
	// router.POST("/feeds", f.Create)
	// router.GET("/feeds/:id", f.Get)
	// router.PUT("/feeds/:id", f.Update)
	// router.DELETE("/feeds/:id", f.Delete)

	v2.GET("/articles", authMiddleware.Authenticate(), articlesHandler.List)

	v2Cases := v2.Group("/cases")

	v2Cases.GET("/", authMiddleware.Authenticate(), casesHandler.List)
	v2Cases.GET("/:id/count", authMiddleware.Authenticate(), cache.SiteCache(storeService, time.Minute), casesHandler.Count)

	// router.GET("/reports", r.List)
	// router.POST("/reports", r.Create)
	// router.GET("/reports/:id", r.Get)
	// router.PUT("/reports/:id", r.Update)
	// router.DELETE("/reports/:id", r.Delete)

	// router.GET("/subscriptions", s.List)
	// router.POST("/subscriptions", s.Create)
	// router.GET("/subscriptions/:id", s.Get)
	// router.PUT("/subscriptions/:id", s.Update)
	// router.DELETE("/subscriptions/:id", s.Delete)

	v2Passages := v2.Group("/passages")

	v2Passages.GET("/", authMiddleware.Authenticate(), passagesHandler.GetPassages)
	v2Passages.POST("/", authMiddleware.Authenticate(), passagesHandler.Create)
	v2Passages.DELETE("/:id", authMiddleware.Authenticate(), passagesHandler.DeletePassage)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	return router
}
