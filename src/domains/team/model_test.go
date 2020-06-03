package team

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
		db.Register(&Team{})
		Convey("insert team with empty data should return error", func() {
			team := &Team{}
			db.CollectionRegistry["Team"].New(team)
			rst, err := team.InsertTeam()
			So(err, ShouldNotBeNil)
			So(rst, ShouldBeNil)
		})
		SkipConvey("insert team with valid data should success", func() {
			team := &Team{}
			db.CollectionRegistry["Team"].New(team)
			team.Name = "Admin"
			team.Description = "Admin team"
			rst, err := team.InsertTeam()
			So(err, ShouldBeNil)
			So(rst, ShouldNotBeNil)
			team.DeleteByID(team.ID)
		})
	})
}
