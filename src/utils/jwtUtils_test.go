package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSignedToken(t *testing.T) {
	Convey("Should return signed valid token when string provided", t, func() {
		tokenString, _ := SignToken("tester@test.com")
		claim, err := VerifyToken(tokenString)
		So(tokenString, ShouldNotBeEmpty)
		So(claim, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

func TestVerifyToken(t *testing.T) {
	Convey("When encrypted with HMAC signed method", t, func() {
		tokenString, _ := SignToken("tester@test.com")
		Convey("return parsed claims when token valid", func() {
			claim, err := VerifyToken(tokenString)
			So(claim, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
	})
	Convey("When token string is invalid", t, func() {
		tokenString := "ACatIsSoStupid"
		Convey("should return error", func() {
			claim, err := VerifyToken(tokenString)
			So(claim, ShouldBeNil)
			So(err, ShouldBeError)
		})
	})
}
