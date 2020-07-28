package project

import (
	"fmt"
	"sync"
	"time"
	"williamfeng323/mooncake-duty/src/domains/account"
	repoimpl "williamfeng323/mooncake-duty/src/infrastructure/db/repo_impl"
	"williamfeng323/mooncake-duty/src/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
)

// Service provides the services to project domain
type Service struct {
	repo *repoimpl.ProjectRepo
}

// SetRepo set the account repository to the service
func (ps *Service) SetRepo(repo *repoimpl.ProjectRepo) {
	ps.repo = repo
}

// CreateProject creates the basic project
func (ps *Service) CreateProject(name string, description string, members ...Member) (*Project, error) {
	if len(members) > 0 {
		for _, member := range members {
			_, err := account.GetAccountService().GetAccountByID(member.MemberID.Hex())
			if err != nil {
				return nil, err
			}
		}
	}
	proj := NewProject(name, description, members...)
	err := proj.Create()
	if err != nil {
		return nil, err
	}
	return proj, nil
}

// SetMembers sets the users to the project
func (ps *Service) SetMembers(id primitive.ObjectID, members ...Member) ([]Member, []Member, error) {
	type em struct {
		Members []Member `bson:"members"`
	}
	existingMembers := em{[]Member{}}
	failedMembers := []Member{}
	findProjectCtx, findProjectCancel := utils.GetDefaultCtx()
	defer findProjectCancel()
	prj := ps.repo.FindOne(findProjectCtx, bson.M{"_id": id})
	if prj.Err() != nil {
		return nil, members, NotFoundError{}
	}
	acctRepo := repoimpl.GetAccountRepo()
	for k, v := range members {
		findAcctCtx, findAcctCtxCancel := utils.GetDefaultCtx()
		defer findAcctCtxCancel()
		rst := acctRepo.FindOne(findAcctCtx, bson.M{"_id": v.MemberID})
		if rst.Err() == nil {
			existingMembers.Members = append(existingMembers.Members, members[k])
		} else {
			failedMembers = append(failedMembers, members[k])
		}
	}
	if len(existingMembers.Members) == 0 {
		return nil, failedMembers, account.NotFoundError{}
	}
	encMembers, err := bson.Marshal(&existingMembers)
	if err != nil {
		return nil, members, err
	}
	bsonMembers := bson.M{}
	bson.Unmarshal(encMembers, bsonMembers)
	psUpdateCtx, psUpdateCancel := utils.GetDefaultCtx()
	defer psUpdateCancel()
	_, err = ps.repo.UpdateOne(psUpdateCtx,
		bson.M{"_id": id}, bson.M{"$set": bson.M{"members": bsonMembers["members"], "updatedAt": time.Now().UTC()}})
	if err != nil {
		return nil, members, err
	}
	return existingMembers.Members, failedMembers, nil
}

// SetNameOrDescription update the name or description of the project.
func (ps *Service) SetNameOrDescription(id primitive.ObjectID, name string, desc string) error {
	if len(name) == 0 && len(desc) == 0 {
		return fmt.Errorf("Nothing to update")
	}

	updValue := bson.M{"updatedAt": time.Now().UTC()}
	if len(name) > 0 {
		updValue["name"] = name
	}
	if len(desc) > 0 {
		updValue["description"] = desc
	}
	updProjectCtx, updProjectCancel := utils.GetDefaultCtx()
	defer updProjectCancel()
	rst, err := ps.repo.UpdateOne(updProjectCtx, bson.M{"_id": id}, bson.M{"$set": updValue})
	if err != nil {
		return err
	}
	if rst.MatchedCount == 0 {
		return NotFoundError{}
	}
	return nil
}

// GetProjectByID find project by its id.
func (ps *Service) GetProjectByID(id primitive.ObjectID) (*Project, error) {
	ctx, cancelCtx := utils.GetDefaultCtx()
	defer cancelCtx()
	rst := ps.repo.FindOne(ctx, bson.M{"_id": id})
	if rst.Err() != nil {
		return nil, NotFoundError{}
	}
	prj := Project{}
	err := rst.Decode(&prj)
	if err != nil {
		return nil, err
	}
	return &prj, nil
}

// GetProjectByName find project by its name.
func (ps *Service) GetProjectByName(projectName string) (*Project, error) {
	ctx, cancelCtx := utils.GetDefaultCtx()
	defer cancelCtx()
	rst := ps.repo.FindOne(ctx, bson.M{"name": projectName})
	if rst.Err() != nil {
		return nil, NotFoundError{}
	}
	prj := Project{}
	err := rst.Decode(&prj)
	if err != nil {
		return nil, err
	}
	return &prj, nil
}

var projectService *Service
var projectServiceLock sync.RWMutex

// GetProjectService returns a singleton project service instance
func GetProjectService() *Service {
	projectServiceLock.Lock()
	defer projectServiceLock.Unlock()
	if projectService == nil {
		projectService = &Service{}
		projectService.SetRepo(repoimpl.GetProjectRepo())
	}
	return projectService
}
