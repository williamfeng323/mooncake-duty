package web

import (
	"fmt"
	"net/http"
	"williamfeng323/mooncake-duty/src/domains/issue"
	project "williamfeng323/mooncake-duty/src/domains/project"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
)

//RegisterIssueRoute the account APIs to root.
func RegisterIssueRoute(router *gin.Engine) {
	issueRoutes := router.Group("/issues")
	{
		issueRoutes.GET("", issueList)
		issueRoutes.GET("/:id")
	}
}

type issueFilter struct {
	IssueKey  string            `form:"issueKey,omitempty"`
	ProjectID string            `form:"projectId" binding:"required"`
	Status    issue.IssueStatus `form:"status,omitempty"`
}

func issueList(c *gin.Context) {
	params := issueFilter{}
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	prjID, err := primitive.ObjectIDFromHex(params.ProjectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ProjectID"})
	}
	issues, err := issue.GetIssueService().GetIssueLists(prjID, params.IssueKey, params.Status)
	if err != nil {
		if _, ok := err.(project.NotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Project %s not found", params.ProjectID)})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"issues": issues})
}
