package issue

import (
	"context"
	"sync"
	"williamfeng323/mooncake-duty/src/domains/project"
	repoimpl "williamfeng323/mooncake-duty/src/infrastructure/db/repo_impl"
	"williamfeng323/mooncake-duty/src/utils"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service contains the issue related services
type Service struct {
	repo *repoimpl.IssueRepo
}

// SetRepo set the issue repository to the service
func (iss *Service) SetRepo(repo *repoimpl.IssueRepo) {
	iss.repo = repo
}

// CreateNewIssue create new issue in db.
func (iss *Service) CreateNewIssue(prjID primitive.ObjectID, issueKey string) (*Issue, error) {
	i, err := NewIssue(prjID, issueKey)
	if err != nil {
		return nil, err
	}
	err = i.Create()
	if err != nil {
		return nil, err
	}
	return i, nil
}

// GetIssueLists return the issues filtered by projectId, issueKey and status, if status is -1,
// it means no criteria on issue status.
func (iss *Service) GetIssueLists(prjID primitive.ObjectID, issueKey string, status IssueStatus) ([]Issue, error) {
	if prj, _ := project.GetProjectService().GetProjectByID(prjID); prjID.IsZero() || prj == nil {
		return nil, project.NotFoundError{}
	}
	filter := []bson.M{bson.M{"projectId": prjID}}
	if len(issueKey) != 0 {
		filter = append(filter, bson.M{"issueKey": issueKey})
	}
	if status.Valid() {
		filter = append(filter, bson.M{"status": status})
	}

	fCtx, fCancel := utils.GetDefaultCtx()
	defer fCancel()
	rst, err := iss.repo.Find(fCtx, bson.M{"$and": filter})
	if err != nil {
		return nil, err
	}

	issues := []Issue{}
	if err := rst.All(context.Background(), &issues); err != nil {
		return nil, err
	}
	return issues, nil
}

// GetIssueByID returns the issue found by ID
func (iss *Service) GetIssueByID(id primitive.ObjectID) (*Issue, error) {
	fCtx, fCancel := utils.GetDefaultCtx()
	defer fCancel()
	rst := iss.repo.FindOne(fCtx, bson.M{"_id": id})
	if rst.Err() != nil {
		return nil, NotFoundError{}
	}
	i := Issue{}
	err := rst.Decode(&i)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

var issueService *Service
var issueServiceLock sync.RWMutex

// GetIssueService returns a singleton issue service instance
func GetIssueService() *Service {
	issueServiceLock.Lock()
	defer issueServiceLock.Unlock()
	if issueService == nil {
		issueService = &Service{}
		issueService.SetRepo(repoimpl.GetIssueRepo())
	}
	return issueService
}
