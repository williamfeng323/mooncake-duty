package project

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

func createProject(project Project) (*mongo.InsertOneResult, error) {
	if len(strings.TrimSpace(project.Name)) == 0 || len(strings.TrimSpace(project.Description)) == 0 {
		return nil, fmt.Errorf("filed Name or Description could not be empty")
	}
	return nil, nil
}
