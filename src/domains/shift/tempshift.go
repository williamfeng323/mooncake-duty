package shift

import (
	"time"
	"williamfeng323/mooncake-duty/src/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TempShift describes the temporary shift period that overrides the default shift
// which generates base on the shift details
type TempShift struct {
	StartDate time.Time
	EndDate   time.Time
	T1Member  primitive.ObjectID `json:"t1Member" bson:"t1Member"`
	T2Member  primitive.ObjectID `json:"t2Member" bson:"t2Member"`
	T3Member  primitive.ObjectID `json:"t3Member" bson:"t3Member"`
}

// NewTempShift inits an temp shift instance with start date and end date
func NewTempShift(startDate time.Time, recurrence shiftRecurrence) *TempShift {
	startDate = utils.ToDateStarted(startDate)
	var endDate time.Time
	switch recurrence {
	case Daily:
		endDate = utils.ToDateEnded(startDate)
	case Weekly:
		endDate = utils.ToDateEnded(startDate.AddDate(0, 0, 6))
	case BiWeekly:
		endDate = utils.ToDateEnded(startDate.AddDate(0, 0, 13))
	case Monthly:
		endDate = time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location()).AddDate(0, 1, 0).Add(-time.Nanosecond)
	}
	return &TempShift{
		StartDate: startDate,
		EndDate:   endDate,
	}
}
