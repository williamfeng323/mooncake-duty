package shift

import (
	"math"
	"time"
	"williamfeng323/mooncake-duty/src/domains/project"
	"williamfeng323/mooncake-duty/src/infrastructure/db"
	repoimpl "williamfeng323/mooncake-duty/src/infrastructure/db/repo_impl"
	"williamfeng323/mooncake-duty/src/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

// WeekStart defines when the week start. Sun/Mon
type WeekStart int

const (
	// Sun week starts from Sunday
	Sun WeekStart = iota
	// Mon week starts from Monday
	Mon
)

// String returns the English name of the week start day ("Sun", "Mon", ...).
func (w WeekStart) String() string {
	days := []string{"Sun", "Mon"}
	if Sun <= w && w <= Mon {
		return days[w]
	}
	return ""
}

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

// String returns the English name of the shift recurrence ("Daily", "Weekly", ...).
func (w shiftRecurrence) String() string {
	days := []string{"Daily", "Weekly", "BiWeekly", "Monthly"}
	if Daily <= w && w <= Monthly {
		return days[w]
	}
	return ""
}

// Shift describes the on call shift
type Shift struct {
	db.BaseModel    `json:",inline" bson:",inline"`
	ProjectID       primitive.ObjectID   `json:"projectId" bson:"projectId"`
	WeekStart       WeekStart            `json:"weekStart" bson:"weekStart"`
	ShiftFirstDate  time.Time            `json:"shiftFirstDate" bson:"shiftFirstDate"`
	ShiftRecurrence shiftRecurrence      `json:"shiftRecurrence" bson:"shiftRecurrence"`
	T1Members       []primitive.ObjectID `json:"t1Members" bson:"t1Members"`
	T2Members       []primitive.ObjectID `json:"t2Members" bson:"t2Members"`
	T3Members       []primitive.ObjectID `json:"t3Members" bson:"t3Members"`
}

// NewShift creates an basic shift instance which has no detail member information
func NewShift(projectID primitive.ObjectID, weekStart WeekStart, shiftFirstDate time.Time, recurrence shiftRecurrence) (*Shift, error) {
	projectRepo := repoimpl.GetProjectRepo()
	findProjectCtx, cancel := utils.GetDefaultCtx()
	defer cancel()
	prjResult := projectRepo.FindOne(findProjectCtx, bson.M{"_id": projectID})
	if prjResult.Err() != nil {
		return nil, project.NotFoundError{}
	}
	shift := &Shift{WeekStart: weekStart, ShiftRecurrence: recurrence, ProjectID: projectID}
	shift.ID = primitive.NewObjectID()
	shift.CreatedAt = time.Now()
	switch recurrence {
	case Daily:
		shift.ShiftFirstDate = utils.ToDateStarted(shiftFirstDate)
	case Weekly:
		fallthrough
	case BiWeekly:
		shift.ShiftFirstDate = utils.FirstDateOfWeek(shiftFirstDate, time.Weekday(weekStart))
	case Monthly:
		shift.ShiftFirstDate = utils.ToMonthStarted(shiftFirstDate)
	}
	return shift, nil
}

// DefaultShifts returns the default shifts in this period according to its info
func (sh *Shift) DefaultShifts(startDate time.Time, endDate time.Time) ([]*TempShift, error) {
	startDate = utils.ToDateStarted(startDate)
	endDate = utils.ToDateStarted(endDate)
	if len(sh.T1Members) == 0 || len(sh.T2Members) == 0 {
		return nil, NoMemberError{}
	}
	if startDate.After(endDate) || startDate.Before(sh.ShiftFirstDate) {
		return nil, OutOfScopeError{}
	}
	tempShifts := []*TempShift{}
	var periodDuration float64
	var sinceShiftStarted float64
	var recurrenceMultiplier int
	switch sh.ShiftRecurrence {
	case Daily:
		d := endDate.Sub(startDate)
		periodDuration = math.Ceil(d.Hours() / 24)
		sinceShiftStarted = math.Ceil(startDate.Sub(sh.ShiftFirstDate).Hours() / 24)
		recurrenceMultiplier = 1
	case Weekly:
		startDate = utils.FirstDateOfWeek(startDate, time.Weekday(sh.WeekStart))
		periodDuration = math.Ceil(utils.ToDateEnded(endDate).Sub(startDate).Hours() / (24 * 7))
		sinceShiftStarted = math.Ceil(startDate.Sub(sh.ShiftFirstDate).Hours() / (24 * 7))
		recurrenceMultiplier = 7
	case BiWeekly:
	case Monthly:
	default:
		return nil, nil
	}
	for i := int(sinceShiftStarted); i < int(sinceShiftStarted+periodDuration); i++ {
		periodStartDate := utils.ToDateStarted(sh.ShiftFirstDate.AddDate(0, 0, i*recurrenceMultiplier))
		ts := NewTempShift(periodStartDate, sh.ShiftRecurrence)

		ts.T1Member = getMember(sh.T1Members, i)
		ts.T2Member = getMember(sh.T2Members, i)
		if len(sh.T3Members) != 0 {
			ts.T3Member = getMember(sh.T3Members, i)
		}
		tempShifts = append(tempShifts, ts)
	}
	return tempShifts, nil
}

func getMember(members []primitive.ObjectID, fromStarted int) primitive.ObjectID {
	actualWeeksPassed := fromStarted + 1
	k := actualWeeksPassed % len(members)
	if k == 0 {
		return members[len(members)-1]
	}
	return members[k-1]
}
