package issue

import (
	"sync"
	repoimpl "williamfeng323/mooncake-duty/src/infrastructure/db/repo_impl"

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

// GetIssueLists ...
func (iss *Service) GetIssueLists(status string) {

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
