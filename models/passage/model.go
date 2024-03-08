package passage

import (
	passagesv2 "github.com/cvcio/mediawatch/pkg/mediawatch/passages/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Passage to Trim
// Passage to Trim
type Passage struct {
	ID       primitive.ObjectID     `bson:"_id" json:"id"`
	Text     string                 `json:"text" bson:"text"`
	Language string                 `json:"language" bson:"language"`
	Type     passagesv2.PassageType `json:"type" bson:"type"`
}
