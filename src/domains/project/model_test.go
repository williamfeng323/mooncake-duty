package project

import (
	"testing"

	dao "williamfeng323/mooncake-duty/src/infrastructure/db"
	"williamfeng323/mooncake-duty/src/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInsertTeam(t *testing.T) {
	Convey("Given an registered Team model", t, func() {
		db := &dao.Connection{}
		db.InitConnection(nil, utils.GetConf().Mongo)
		db.Register(&Project{})
		Convey("insert team with empty data should return error", func() {
			project := &Project{}
			db.CollectionRegistry["Team"].New(project)
			rst, err := project.CreateProject()
			So(err, ShouldNotBeNil)
			So(rst, ShouldBeNil)
		})
		SkipConvey("insert team with valid data should success", func() {
			project := &Project{}
			db.CollectionRegistry["Team"].New(project)
			project.Name = "Admin"
			project.Description = "Admin team"
			rst, err := project.CreateProject()
			So(err, ShouldBeNil)
			So(rst, ShouldNotBeNil)
			project.DeleteByID(project.ID)
		})
	})
}
