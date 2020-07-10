package issue

import (
	repoimpl "williamfeng323/mooncake-duty/src/infrastructure/db/repo_impl"
)

// Service contains the issue related services
type Service struct {
	repo *repoimpl.IssueRepo
}

// func (is *Service) CreateIssue (projectID primitive.ObjectID) {

// }
