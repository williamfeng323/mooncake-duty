package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSignedToken(t *testing.T) {
	Convey("Should return signed valid token when string provided", t, func() {
		tokenString, _ := SignToken("tester@test.com")
		claim, err := VerifyToken(tokenString, true)
		So(tokenString, ShouldNotBeEmpty)
		So(claim, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

func TestVerifyToken(t *testing.T) {
	Convey("When encrypted with HMAC signed method", t, func() {
		tokenString, _ := SignToken("tester@test.com")
		Convey("return parsed claims when token valid", func() {
			claim, err := VerifyToken(tokenString, false)
			So(claim, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
	})
	Convey("When token string is invalid", t, func() {
		tokenString := "ACatIsSoStupid"
		Convey("should return error", func() {
			claim, err := VerifyToken(tokenString, false)
			So(claim, ShouldBeNil)
			So(err, ShouldBeError)
		})
	})
	Convey("When encrypted with HMAC signed method", t, func() {
		jwtExpireIn = -10
		tokenString, _ := SignToken("tester@test.com")
		Convey("should not return parsed claims when token valid when need to verify expiration", func() {
			claim, err := VerifyToken(tokenString, false)
			So(claim, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
		Convey("return parsed claims when token valid when skip verify expiration", func() {
			claim, err := VerifyToken(tokenString, true)
			So(claim, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
		jwtExpireIn = GetConf().JWTExpireIn
	})
}
