package project

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"strings"
	db "williamfeng323/mooncake-duty/src/infrastructure/db"
)

// Permission the predefined permissions strings
type Permission string

// Tier is the support level of a team member.
type Tier string

// Severity describe the level of the alarm.
type Severity int

const (
	// Create predefined permission mode, create permission
	Create Permission = "create"
	// Read predefined permission mode, read permission
	Read Permission = "read"
	// Update predefined permission mode, update permission
	Update Permission = "update"
	// Delete predefined permission mode, deleted permission
	Delete Permission = "delete"
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
	return project
}
