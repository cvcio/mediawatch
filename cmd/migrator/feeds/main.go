package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/cvcio/mediawatch/models/deprecated/feed"
	"github.com/cvcio/mediawatch/pkg/db"
	"github.com/cvcio/mediawatch/pkg/twitter"
)

func main() {
	dConn, err := db.NewMongoDB("mongodb://localhost:27017", "mediawatch", 10*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new twitter client
	twtt, err := twitter.NewAPI("",
		"", "",
		"")
	if err != nil {
		log.Fatalf("Error connecting to twitter: %s", err.Error())
	}
	a, err := os.Open("tmp/exports/MediaWatch Feeds - feeds 20210131.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	ra := csv.NewReader(a)

	for {
		// Read each record from csv
		record, err := ra.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if record[0] != "_id" {
			// _id,name,screen_name,tier,business_type,business_owner,registered,political_stance,political_orientation,content_type,country,locality,lang,url,status,testURL
			if record[0] != "" {
				// update
				f := new(feed.UpdateFeed)
				f.Name = &record[1]
				f.ScreenName = &record[2]
				f.Tier = &record[3]
				f.BusinessType = &record[4]
				f.BusinessOwner = &record[5]
				f.Registered = &record[6]
				f.PoliticalStance = &record[7]
				f.PoliticalOrientation = &record[8]
				f.ContentType = &record[9]
				f.Country = &record[10]
				f.Locality = &record[11]
				f.Lang = &record[12]
				f.URL = &record[13]
				f.Status = &record[14]
				f.TestURL = &record[15]

				user, err := twtt.GetUsersShow(*f.ScreenName, url.Values{})
				if err != nil {
					fmt.Printf("MISSING ACCOUNT (%s)\n", *f.ScreenName)
				} else {
					f.ScreenName = &user.ScreenName
					f.TwitterID = &user.Id
					f.TwitterIDStr = &user.IdStr
					f.Description = &user.Description
					f.TwitterProfileImage = &user.ProfileImageUrlHttps
				}

				err = feed.Update(context.Background(), dConn, record[0], f, time.Now())
				if err != nil {
					fmt.Printf("COULDNT UPDATE (%s)\n", *f.ScreenName)
				}
			}
			// else {
			// 	f := new(feed.Feed)
			// 	f.Name = record[1]
			// 	f.ScreenName = record[2]
			// 	f.BusinessType = record[3]
			// 	f.ContentType = record[4]
			// 	f.URL = record[5]
			// 	f.Status = "pending"
			// 	f.TestURL = record[7]

			// 	user, err := twtt.GetUsersShow(f.ScreenName, url.Values{})
			// 	if err != nil {
			// 		fmt.Println(err)
			// 		continue
			// 	}

			// 	f.TwitterID = user.Id
			// 	f.TwitterIDStr = user.IdStr
			// 	f.TwitterProfileImage = user.ProfileImageUrlHttps

			// 	_, err = feed.Create(context.Background(), dConn, f, time.Now())
			// 	if err != nil {
			// 		fmt.Println(err)
			// 	}
			// }
		}
	}

}
