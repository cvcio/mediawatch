package passage

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Passage to Trim
type Passage struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Text string             `json:"text" bson:"text"`
	Type string             `json:"type" bson:"type"`
}
