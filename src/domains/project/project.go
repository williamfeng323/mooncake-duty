package project

import (
	"fmt"
	"strings"
	"time"
	repoimpl "williamfeng323/mooncake-duty/src/infrastructure/db/repo_impl"
	"williamfeng323/mooncake-duty/src/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	db "williamfeng323/mooncake-duty/src/infrastructure/db"
	validatorimpl "williamfeng323/mooncake-duty/src/infrastructure/db/validator_impl"
)

// Tier is the support level of a team member.
type Tier string

// Severity describe the level of the alarm.
type Severity int

const (
	// T1 predefined tier level
	T1 Tier = "T1"
	// T2 predefined tier level
	T2 Tier = "T2"
	// T3 predefined tier level
	T3 Tier = "T3"
)

const (
	// High represents the most serious issue, must be handle asap
	High Severity = iota
	// Medium represents the medium level issue, you might not need to escalate
	Medium
	// Low could be follow up later
	Low
)

// Alarm defines the struct of an alarm
type Alarm struct {
	ID          primitive.ObjectID `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Severity    `json:"severity" bson:"severity"`
}

// Member describes the member info in project
type Member struct {
	MemberID primitive.ObjectID `json:"memberId" bson:"memberId"`
	IsAdmin  bool               `json:"isAdmin" bson:"isAdmin"`
}

// Project is the struct to contain project information.
type Project struct {
	repo         *repoimpl.ProjectRepo
	db.BaseModel `json:",inline" bson:",inline"`
	Name         string   `json:"name" bson:"name" required:"true"`
	Description  string   `json:"description" bson:"description" required:"true"`
	Members      []Member `json:"members" bson:"members"`
	Alarms       []Alarm  `json:"alarms" bson:"alarms"` // AlarmLog would be needed
}

// NewProject returns a project instance with required fields
func NewProject(name string, description string, members ...Member) *Project {
	project := &Project{
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
		Members:     members,
	}
	project.ID = primitive.NewObjectID()
	project.CreatedAt = time.Now()
	project.repo = repoimpl.GetProjectRepo()
	return project
}

// Create verifies and inserts the project into database
func (prj *Project) Create() error {
	if prj.repo == nil {
		return fmt.Errorf("Project does not initialized")
	}
	validator := validatorimpl.NewDefaultValidator()
	errs := validator.Verify(prj)
	if len(errs) != 0 {
		return fmt.Errorf("Save the account failed due to: %v", errs)
	}
	ctxFind, cancelFind := utils.GetDefaultCtx()
	defer cancelFind()
	rst := prj.repo.FindOne(ctxFind, bson.M{"name": prj.Name})
	foundProject := &Project{}
	err2 := rst.Decode(foundProject)
	print(err2)
	if !foundProject.ID.IsZero() {
		return AlreadyExistError{}
	}
	ctxInsert, cancelInsert := utils.GetDefaultCtx()
	defer cancelInsert()
	_, err := prj.repo.InsertOne(ctxInsert, prj)
	return err
}
