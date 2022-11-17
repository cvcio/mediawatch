package subscription

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Subscription is a product with a plan for a specific user
type Subscription struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	Deleted   bool      `json:"-" bson:"deleted"`

	// Stripe Ids
	CustomerID string `json:"customerId" bson:"customerId"`
	ProductID  string `json:"productId" bson:"productId"` // Embed ?
	PriceID    string `json:"priceId" bson:"priceId"`     // Embed ?

	UserID         string `json:"userId" bson:"userId"`                 // Embed ?
	OrganizationID string `json:"organizationId" bson:"organizationId"` // Embed ?
}

type SubsciptionSessionResponse struct {
	Status int    `json:"status"`
	URL    string `json:"url"`
}
