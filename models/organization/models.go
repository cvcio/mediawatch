package organization

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Organization Model
type Organization struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	Deleted   bool      `json:"-" bson:"deleted"`

	Name       string    `json:"name" bson:"name"`
	Industry   string    `json:"industry" bson:"industry"`
	Type       string    `json:"type,omitempty" bson:"type,omitempty"`
	Size       string    `json:"size,omitempty" bson:"size,omitempty"`
	Email      string    `json:"email" bson:"email,omitempty"`
	Phone      string    `json:"phone" bson:"phone,omitempty"`
	ScreenName string    `json:"screenName,omitempty" bson:"screenName,omitempty"`
	URL        string    `json:"url,omitempty" bson:"url,omitempty"`
	Country    string    `json:"country,omitempty" bson:"country,omitempty"`
	City       string    `json:"city,omitempty" bson:"city,omitempty"`
	Language   string    `bson:"language,omitempty" json:"language,omitempty"`
	Timezone   string    `bson:"timezone,omitempty" json:"timezone,omitempty"`
	Seats      int       `bson:"seats" json:"seats"`
	Members    []*Member `json:"members" bson:"members"`
}

type Member struct {
	ID     string `json:"id,omitempty"`
	Nonce  string `json:"nonce,omitempty"`
	Status string `json:"status,omitempty"`
	Role   string `json:"role,omitempty"`
	Email  string `json:"email,omitempty"`
}

type NewOrg struct {
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"required"`
	Country string `json:"country" validate:"required"`
}

type UpdOrg struct {
	Name       *string `json:"name,omitempty" bson:"name,omitempty"`
	Industry   *string `json:"industry,omitempty" bson:"industry,omitempty"`
	Type       *string `json:"type,omitempty" bson:"type,omitempty"`
	Size       *string `json:"size,omitempty" bson:"size,omitempty"`
	Email      *string `json:"email,omitempty" bson:"email,omitempty"`
	Phone      *string `json:"phone,omitempty" bson:"phone,omitempty"`
	ScreenName *string `json:"screenName,omitempty" bson:"screenName,omitempty"`
	URL        *string `json:"url,omitempty" bson:"url,omitempty"`
	Country    *string `json:"country,omitempty" bson:"country,omitempty"`
	City       *string `json:"city,omitempty" bson:"city,omitempty"`
	Language   *string `bson:"language,omitempty" json:"language,omitempty"`
	Timezone   *string `bson:"timezone,omitempty" json:"timezone,omitempty"`

	Members []*Member `json:"members" bson:"members"`
}
