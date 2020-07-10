package issue

import (
	"time"
	"williamfeng323/mooncake-duty/src/domains/project"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	repoimpl "williamfeng323/mooncake-duty/src/infrastructure/db/repo_impl"
	"williamfeng323/mooncake-duty/src/utils"
)

// Tier represents the notification tier
type Tier int

const (
	// T1 maps to T1 members in shift
	T1 Tier = iota + 1
	// T2 maps to T2 members in shift
	T2
	// T3 maps to T3 members in shift
	T3
)

// Issue describe the triggered alert created issues
type Issue struct {
	ProjectID           primitive.ObjectID `json:"projectId" bson:"projectId" required:"true"`
	Service             string             `json:"service" bson:"service" required:"true"`
	CreatedAt           time.Time          `json:"createdAt" bson:"createdAt" `
	AcknowledgedAt      time.Time          `json:"acknowledgedAt" bson:"acknowledgedAt"`
	ResolvedAt          time.Time          `json:"resolvedAt" bson:"resolvedAt"`
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

// NewIssue validate projectID existence and returns issue
func NewIssue(projectID primitive.ObjectID, service string) (*Issue, error) {
	projRepo := repoimpl.GetProjectRepo()
	findProjCtx, findProjCancel := utils.GetDefaultCtx()
	defer findProjCancel()
	projRst := projRepo.FindOne(findProjCtx, bson.M{"_id": projectID})
	if projRst.Err() != nil {
		return nil, project.NotFoundError{}
	}
	return &Issue{ProjectID: projectID, Service: service, CreatedAt: time.Now()}, nil
}
