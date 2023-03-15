package commands

import (
	"context"
	"time"

	"github.com/cvcio/mediawatch/models/deprecated/feed"
	"github.com/cvcio/mediawatch/pkg/config"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	// migrateCmd represents the init command
	migrateFeedsCmd = &cobra.Command{
		Use:   "feeds",
		Short: "",
		Long:  ``,
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// Configure
			cfg := config.NewConfig()
			log := logrus.New()

			// Read config from env variables
			err := envconfig.Process("", cfg)
			if err != nil {
				log.Fatalf("Error loading config: %s", err.Error())
			}
			// ============================================== ==============
			// Start Mongo
			log.Debug("Initialize Mongo")
			dbConn, err := db.NewMongoDB(cfg.Mongo.URL, cfg.Mongo.Path, cfg.Mongo.DialTimeout)
			if err != nil {
				log.Fatalf("Register DB: %v", err)
			}
			log.Debug("Connected to Mongo")
			defer dbConn.Close()

			collection := dbConn.Client.Database("mediawatch").Collection("feeds")
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			cur, err := collection.Find(ctx, bson.M{})
			if err != nil {
				log.Fatal(err)
			}
			defer cur.Close(ctx)

			for cur.Next(ctx) {
				var f *feed.Feed
				err := cur.Decode(&f)
				if err != nil {
					log.Fatal(err)
				}

				log.Info(f.Name)
			}
			if err := cur.Err(); err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	migrateCmd.AddCommand(migrateFeedsCmd)
}
