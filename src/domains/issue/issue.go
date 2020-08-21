package issue

import (
	"fmt"
	"time"
	"williamfeng323/mooncake-duty/src/domains/project"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"williamfeng323/mooncake-duty/src/infrastructure/db"
	repoimpl "williamfeng323/mooncake-duty/src/infrastructure/db/repo_impl"
	validatorimpl "williamfeng323/mooncake-duty/src/infrastructure/db/validator_impl"
	"williamfeng323/mooncake-duty/src/utils"
)

// Tier represents the notification tier
type Tier int

// IssueStatus represents the possible issue statuses
type IssueStatus int

const (
	// T1 maps to T1 members in shift
	T1 Tier = iota + 1
	// T2 maps to T2 members in shift
	T2
	// T3 maps to T3 members in shift
	T3
)

const (
	// Init issue just created
	Init IssueStatus = iota
	// Acknowledged issue acknowledged. Cannot change back from Resolved
	Acknowledged
	// Resolved issue handled. Issue can skip acknowledged status and goes to resolved directly.
	Resolved
)

func (iStatus IssueStatus) String() string {
	statusArray := []string{"Init", "Acknowledged", "Resolved"}
	return statusArray[iStatus]
}

// Issue describe the triggered alert created issues
type Issue struct {
	repo                *repoimpl.IssueRepo
	db.BaseModel        `json:",inline" bson:",inline"`
	ProjectID           primitive.ObjectID `json:"projectId" bson:"projectId" required:"true"`
	IssueKey            string             `json:"issueKey" bson:"issueKey" required:"true"`
	Status              IssueStatus        `json:"status" bson:"status"`
	AcknowledgedAt      time.Time          `json:"acknowledgedAt" bson:"acknowledgedAt"`
	AcknowledgedBy      string             `json:"acknowledgedBy" bson:"acknowledgedBy"`
	ResolvedAt          time.Time          `json:"resolvedAt" bson:"resolvedAt"`
	ResolvedBy          string             `json:"resolvedBy" bson:"resolvedBy"`
	T1NotifiedAt        []time.Time        `json:"t1NotifiedAt" bson:"t1LastNotifiedAt"`
	T1NotificationCount int                `json:"t1NotificationCount" bson:"t1NotificationCount"`
	T2NotifiedAt        []time.Time        `json:"t2NotifiedAt" bson:"t2LastNotifiedAt"`
	T2NotificationCount int                `json:"t2NotificationCount" bson:"t2NotificationCount"`
	T3NotifiedAt        []time.Time        `json:"t3NotifiedAt" bson:"t3LastNotifiedAt"`
	T3NotificationCount int                `json:"t3NotificationCount" bson:"t3NotificationCount"`
}

// GetNotificationTier returns the proper notifier base on the current alert status
func (i *Issue) GetNotificationTier() (Tier, error) {
	projRepo := repoimpl.GetProjectRepo()
	findProjCtx, findProjCancel := utils.GetDefaultCtx()
	defer findProjCancel()
	projRst := projRepo.FindOne(findProjCtx, bson.M{"_id": i.ProjectID})
	if projRst.Err() != nil {
		return 0, project.NotFoundError{}
	}
	proj := &project.Project{}
	err := projRst.Decode(proj)
	if err != nil {
		return 0, err
	}
	callPerTiers := proj.CallsPerTier
	if i.T2NotificationCount >= callPerTiers {
		return T3, nil
	} else if i.T1NotificationCount >= callPerTiers {
		return T2, nil
	}
	return T1, nil
}

// Create verifies and inserts the issue into database
func (i *Issue) Create() error {
	validator := validatorimpl.NewDefaultValidator()
	errs := validator.Verify(i)
	if len(errs) != 0 {
		return fmt.Errorf("Save the issue failed due to: %v", errs)
	}
	ctxInsert, cancelInsert := utils.GetDefaultCtx()
	defer cancelInsert()
	_, err := repoimpl.GetIssueRepo().InsertOne(ctxInsert, i)
	return err
}

// UpdateStatus update the issue status
func (i *Issue) UpdateStatus(status IssueStatus, by string) error {
	if len(by) == 0 {
		return fmt.Errorf("The one who update the status cannot be empty")
	}
	if i.Status > status {
		return fmt.Errorf("The status cannot set back")
	}
	switch status {
	case Acknowledged:
		i.Status = status
		i.AcknowledgedAt = time.Now()
		i.AcknowledgedBy = by
	case Resolved:
		if i.Status == Init {
			i.AcknowledgedAt = time.Now()
			i.AcknowledgedBy = by
		}
		i.Status = status
		i.ResolvedAt = time.Now()
		i.ResolvedBy = by
	}
	i.UpdatedAt = time.Now()
	updCtx, updCtxCancel := utils.GetDefaultCtx()
	defer updCtxCancel()

	var inInterface bson.M
	inrec, _ := bson.Marshal(i)
	bson.Unmarshal(inrec, &inInterface)
	i.repo.UpdateOne(updCtx, bson.M{"_id": i.ID}, bson.M{"$set": inInterface})
	return nil
}

// IsDuplicate distinguishes if 2 issues are the same
func (i *Issue) IsDuplicate(iss *Issue) bool {
	return i.IssueKey == iss.IssueKey
}

// NewIssue validate projectID existence and returns issue
func NewIssue(projectID primitive.ObjectID, key string) (*Issue, error) {
	projRepo := repoimpl.GetProjectRepo()
	findProjCtx, findProjCancel := utils.GetDefaultCtx()
	defer findProjCancel()
	projRst := projRepo.FindOne(findProjCtx, bson.M{"_id": projectID})
	if projRst.Err() != nil {
		return nil, project.NotFoundError{}
	}
	iss := &Issue{
		ProjectID: projectID,
		IssueKey:  key,
		Status:    Init,
	}
	iss.repo = repoimpl.GetIssueRepo()
	iss.ID = primitive.NewObjectID()
	iss.CreatedAt = time.Now()
	return iss, nil
}
