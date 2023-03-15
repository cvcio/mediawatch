package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/cvcio/mediawatch/models/deprecated/feed"
	"github.com/cvcio/mediawatch/pkg/db"
)

const FILE = "/home/andefined/py/mediawatch/import-feeds/mediawatch-feeds-20230311.csv"

func main() {
	dConn, err := db.NewMongoDB("mongodb://localhost:27017", "mediawatch", 10*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(FILE)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	records := csv.NewReader(file)
	for {
		// Read each record from csv
		record, err := records.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if record[0] != "_id" {
			// _id, 0
			// createdAt, 1
			// updatedAt, 2
			// deleted, 3
			// country, 4
			// lang, 5
			// name, 6
			// screen_name, 7
			// twitter_id, 8
			// business_type, 9
			// tier, 10
			// content_type, 11
			// url, 12
			// rss, 13
			// stream_type, 14
			// business_owner, 15
			// locality, 16
			// political_orientation, 17
			// political_stance, 18
			// status, 19
			// testURL, 20
			// twitter_id_str, 21
			// registered, 22
			// hostname, 23
			// registry_id, 24
			// meta_classes.api, 25
			// meta_classes.feed_type, 26
			// profile_image_url, 27
			if record[0] != "" {
				f := new(feed.UpdateFeed)
				f.Name = &record[6]
				f.ScreenName = &record[7]
				f.Tier = &record[10]
				f.BusinessType = &record[9]
				f.BusinessOwner = &record[15]
				f.Registered = &record[24]
				f.PoliticalStance = &record[18]
				f.PoliticalOrientation = &record[17]
				f.ContentType = &record[11]
				f.Country = &record[4]
				f.Locality = &record[16]
				f.Lang = &record[5]
				f.URL = &record[12]
				f.Status = &record[19]
				f.TestURL = &record[20]
				f.TwitterProfileImage = &record[27]

				if record[13] != "-" {
					f.RSS = &record[13]
				}

				f.StreamType = &record[14]

				twitterID, _ := strconv.ParseInt(record[8], 10, 64)
				f.TwitterID = &twitterID
				f.TwitterIDStr = &record[21]

				j, _ := json.Marshal(f)
				log.Printf("%v", string(j))

				err = feed.Update(context.Background(), dConn, record[0], f, time.Now())
				if err != nil {
					fmt.Printf("COULDNT UPDATE (%s)\n", record[6])
				}
			} else {
				f := new(feed.Feed)
				f.Name = record[6]
				f.ScreenName = record[7]
				f.Tier = record[10]
				f.BusinessType = record[9]
				f.BusinessOwner = record[15]
				f.Registered = record[24]
				f.PoliticalStance = record[18]
				f.PoliticalOrientation = record[17]
				f.ContentType = record[11]
				f.Country = record[4]
				f.Locality = record[16]
				f.Lang = record[5]
				f.URL = record[12]
				f.Status = record[19]
				f.TestURL = record[20]
				f.TwitterProfileImage = record[27]

				if record[13] != "-" {
					f.RSS = record[13]
				}

				f.StreamType = record[14]

				twitterID, _ := strconv.ParseInt(record[8], 10, 64)
				f.TwitterID = twitterID
				f.TwitterIDStr = record[21]

				_, err = feed.Create(context.Background(), dConn, f, time.Now())
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}

}
