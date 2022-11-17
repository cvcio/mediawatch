package main

import (
	"context"
	"log"
	"time"

	accountsv2 "github.com/cvcio/mediawatch/internal/mediawatch/accounts/v2"
	commonv2 "github.com/cvcio/mediawatch/internal/mediawatch/common/v2"
	"github.com/cvcio/mediawatch/models/deprecated/account"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	collectionv1 := client.Database("mediawatch").Collection("accounts")
	ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := collectionv1.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var acc account.Account
		err := cur.Decode(&acc)
		if err != nil {
			log.Fatal(err)
		}
		newAcc := &accountsv2.Account{
			Id:          acc.ID.Hex(),
			GroupId:     acc.Organization,
			CreatedAt:   timestamppb.New(acc.CreatedAt),
			UpdatedAt:   timestamppb.New(acc.UpdatedAt),
			LastLoginAt: timestamppb.New(acc.LastLoginAt),
			Status:      commonv2.Status_STATUS_PENDING,
			Role:        accountsv2.Role_ROLE_USER,
			Email:       acc.Email,
			Mfa:         &accountsv2.MFA{Enabled: false},
			Source:      "mediawatch.io",
			Authorizer:  "mediawatch",
			Profile: &accountsv2.Profile{
				FirstName:  acc.FirstName,
				LastName:   acc.LastName,
				UserName:   acc.ScreenName,
				Country:    acc.Country,
				Language:   acc.Language,
				Industry:   acc.Industry,
				Occupation: acc.Occupation,
				Mobile:     acc.Mobile,
			},
			Deleted: acc.Deleted,
		}
		log.Println(newAcc)

		// res, err := collection.DeleteOne(ctx, bson.M{"_id": acc.ID})
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// fmt.Println("DeleteOne Result TYPE:", reflect.TypeOf(res))

		// insert, err := collection.InsertOne(ctx, &newAcc)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// log.Println(insert)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
}
