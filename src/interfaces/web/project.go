package web

import (
	"fmt"
	"net/http"
	"williamfeng323/mooncake-duty/src/domains/account"
	project "williamfeng323/mooncake-duty/src/domains/project"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
)

// CreateProject defines the parameters to create a project
type CreateProject struct {
	Name        string           `json:"name" binding:"required"`
	Description string           `json:"description" binding:"required"`
	Members     []project.Member `json:"members,omitempty" binding:"omitempty"`
}

// ProjectFilter defines the parameters to filter the project in GET /projects
type ProjectFilter struct {
	Name string `form:"name" binding:"required"`
}

// ProjectID defines the url wildcard parameter
type ProjectID struct {
	ID string `uri:"id" binding:"required"`
}

// UpdateProjectParams update project name, description or members in project
type UpdateProjectParams struct {
	Name        string           `json:"name,omitempty"`
	Description string           `json:"description,omitempty"`
	Members     []project.Member `json:"members,omitempty"`
}

//RegisterProjectRoute the project APIs to root.
func RegisterProjectRoute(router *gin.Engine) {
	projectsRoutes := router.Group("/projects")
	{
		projectsRoutes.GET("", func(c *gin.Context) {
			params := ProjectFilter{}
			if err := c.ShouldBindQuery(&params); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			prj, err := project.GetProjectService().GetProjectByName(params.Name)
			if err != nil {
				if _, ok := err.(project.NotFoundError); ok {
					c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Project %s not found", params.Name)})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"project": prj})
			return
		})
		projectsRoutes.PUT("", func(c *gin.Context) {
			params := CreateProject{}
			if err := c.ShouldBindJSON(&params); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			proj, err := project.GetProjectService().CreateProject(params.Name, params.Description, params.Members...)
			if err != nil {
				if _, ok := err.(account.NotFoundError); ok {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid member account"})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"project": proj})
			return
		})
		projectsRoutes.GET("/:id", func(c *gin.Context) {
			param := ProjectID{}
			if err := c.ShouldBindUri(&param); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			id, e := primitive.ObjectIDFromHex(param.ID)
			if e != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
				return
			}
			prj, err := project.GetProjectService().GetProjectByID(id)
			if err != nil {
				if _, ok := err.(project.NotFoundError); ok {
					c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Project ID %s not found", param.ID)})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"project": prj})
			return
		})
		projectsRoutes.POST("/:id", func(c *gin.Context) {
			idParam := ProjectID{}
			updParams := UpdateProjectParams{}
			if err := c.ShouldBindUri(&idParam); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			id, e := primitive.ObjectIDFromHex(idParam.ID)
			if e != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": project.NotFoundError{}.Error()})
				return
			}
			if err := c.ShouldBindJSON(&updParams); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if len(updParams.Name) == 0 && len(updParams.Description) == 0 && updParams.Members == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("Empty request body")})
				return
			}
			rsp := gin.H{}
			if len(updParams.Name) > 0 || len(updParams.Description) > 0 {
				err := project.GetProjectService().SetNameOrDescription(id, updParams.Name, updParams.Description)
				if err != nil {
					if _, ok := err.(project.NotFoundError); ok {
						c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Project ID %s not found", idParam.ID)})
						return
					}
					rsp["error"] = err.Error()
					c.JSON(http.StatusInternalServerError, rsp)
					return
				}
			}
			if len(updParams.Members) > 0 {
				suc, failed, err := project.GetProjectService().SetMembers(id, updParams.Members...)
				rsp["updatedMembers"] = suc
				rsp["failedMembers"] = failed
				if err != nil {
					rsp["error"] = err.Error()
				}
			}
			c.JSON(http.StatusOK, rsp)
		})
	}
}
