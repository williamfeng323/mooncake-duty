package project

import (
	"williamfeng323/mooncake-duty/src/domains/account"
	repoimpl "williamfeng323/mooncake-duty/src/infrastructure/db/repo_impl"
	"williamfeng323/mooncake-duty/src/utils"

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

// SetMembers sets the users to the project
func (ps *Service) SetMembers(projectName string, members ...Member) ([]Member, []Member, error) {
	type em struct {
		Members []Member `bson:"members"`
	}
	existingMembers := em{[]Member{}}
	failedMembers := []Member{}
	findProjectCtx, findProjectCancel := utils.GetDefaultCtx()
	defer findProjectCancel()
	prj := ps.repo.FindOne(findProjectCtx, bson.M{"name": projectName})
	if prj.Err() != nil {
		return nil, members, prj.Err()
	}
	acctRepo := repoimpl.GetAccountRepo()
	for k, v := range members {
		findAcctCtx, findAcctCtxCancel := utils.GetDefaultCtx()
		defer findAcctCtxCancel()
		tempAcct := account.Account{}
		rst := acctRepo.FindOne(findAcctCtx, bson.M{"_id": v.MemberID})
		err := rst.Decode(&tempAcct)
		if err == nil {
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
	_, err = ps.repo.UpdateOne(psUpdateCtx, bson.M{"name": projectName}, bson.M{"$set": bson.M{"members": bsonMembers["members"]}})
	if err != nil {
		return nil, members, err
	}
	return existingMembers.Members, failedMembers, nil
}
