package shift

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TempShift describes the temporary shift period that overrides the default shift
// which generates base on the shift details
type TempShift struct {
	StartDate time.Time
	EndDate   time.Time
	T1Members []*primitive.ObjectID `json:"t1Members" bson:"t1Members"`
	T2Members []*primitive.ObjectID `json:"t2Members" bson:"t2Members"`
	T3Members []*primitive.ObjectID `json:"t3Members" bson:"t3Members"`
}
