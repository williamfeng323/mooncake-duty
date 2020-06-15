package shift

import (
	"time"
	"williamfeng323/mooncake-duty/src/infrastructure/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type weekStart int

const (
	// Sun week starts from Sunday
	Sun weekStart = iota
	// Mon week starts from Monday
	Mon
)

type shiftRecurrence int

const (
	// Daily means the shift rotates daily
	Daily shiftRecurrence = iota
	// Weekly means the shift rotates weekly
	Weekly
	// BiWeekly means the shift rotates bi-weekly
	BiWeekly
	// Monthly means the shift rotates monthly
	Monthly
)

// Shift describes the on call shift
type Shift struct {
	db.BaseModel    `json:",inline" bson:",inline"`
	ProjectID       primitive.ObjectID   `json:"projectId" bson:"projectId"`
	WeekStart       weekStart            `json:"weekStart" bson:"weekStart"`
	ShiftFirstDate  time.Time            `json:"shiftFirstDate" bson:"shiftFirstDate"`
	ShiftRecurrence shiftRecurrence      `json:"shiftRecurrence" bson:"shiftRecurrence"`
	T1Members       []primitive.ObjectID `json:"t1Members" bson:"t1Members"`
	T2Members       []primitive.ObjectID `json:"t2Members" bson:"t2Members"`
	T3Members       []primitive.ObjectID `json:"t3Members" bson:"t3Members"`
}

// NewShift creates an basic shift instance which has no detail member information
func NewShift(projectID primitive.ObjectID, weekStart weekStart, shiftFirstDate time.Time, recurrence shiftRecurrence) (*Shift, error) {
	return nil, nil
}

// DefaultShifts returns the default shifts in this period according to its info
func (shift *Shift) DefaultShifts() (*primitive.ObjectID, *primitive.ObjectID, *primitive.ObjectID) {
	return nil, nil, nil
}
