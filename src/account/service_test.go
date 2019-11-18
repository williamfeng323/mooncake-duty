package account

import (
	"testing"
	"williamfeng323/mooncake-duty/src/dao"

	"go.mongodb.org/mongo-driver/bson/primitive"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateAccount(t *testing.T) {
	Convey("Giving a db connection", t, func() {
		acct := &Account{}
		conn := dao.GetConnection()

		Convey("Should panic while no register Account", func() {
			So(func() { createAccount("", "") }, ShouldPanic)
		})
		Convey("Should return error When email or password does not provided", func() {
			conn.Register(acct)
			rst, err := createAccount("", "")
			So(rst, ShouldBeNil)
			So(err, ShouldBeError)
		})
		Convey("Should return inserted account objectId when email/password valid", func() {
			rst, err := createAccount("test123", "rstAbc.")
			So(rst, ShouldNotBeNil)
			So(err, ShouldBeNil)
			stt := rst.InsertedID.(primitive.ObjectID).Hex()
			id, err := primitive.ObjectIDFromHex(stt)
			conn.CollectionRegistry["Account"].New(acct)
			acct.DeleteByID(id)
		})
		Convey("Should return error When email duplicate with existing account", func() {
			initAcct, err := createAccount("test123", "rstAbcd")
			stt := initAcct.InsertedID.(primitive.ObjectID).Hex()
			id, err := primitive.ObjectIDFromHex(stt)
			rst, err := createAccount("test123", "rstAbcd")
			So(rst, ShouldBeNil)
			So(err, ShouldBeError)
			conn.CollectionRegistry["Account"].New(acct)
			acct.DeleteByID(id)
		})
	})
}

func TestSignIn(t *testing.T) {
	Convey("Giving a db connection and init user", t, func() {
		acct := &Account{}
		conn := dao.GetConnection()
		conn.Register(acct)
		conn.CollectionRegistry["Account"].New(acct)
		acct.Email = "test@test.com"
		acct.Password = "password"
		acct.InsertAccount()
		Convey("Should return error when user not found", func() {
			act, err := signIn("sss@ss.com", "abc")
			So(act, ShouldBeEmpty)
			So(err, ShouldNotBeNil)
		})
		Convey("Should return error when user/password does valid", func() {
			act, err := signIn("test@test.com", "abc")
			So(act, ShouldBeEmpty)
			So(err, ShouldNotBeNil)
		})
		FocusConvey("Should return jwt token when user/password valid", func() {
			act, err := signIn("test@test.com", "password")
			So(act, ShouldNotBeEmpty)
			So(err, ShouldBeNil)
		})
		acct.DeleteByID(acct.ID)
	})
}

func TestRefresh(t *testing.T) {
	Convey("Giving a db connection and init user", t, func() {
		acct := &Account{}
		conn := dao.GetConnection()
		conn.Register(acct)
		conn.CollectionRegistry["Account"].New(acct)
		acct.Email = "test@test.com"
		acct.Password = "password"
		acct.InsertAccount()
		Convey("Should return refreshed token when called with valid token", func() {
			tokenString, _ := signIn("test@test.com", "password")
			refreshedToken, err := refresh(tokenString)
			So(refreshedToken, ShouldNotBeEmpty)
			So(err, ShouldBeEmpty)
		})
		Convey("Should return error when token invalid(any reason even timeout)", func() {
			refreshedToken, err := refresh("whateverTokenString")
			So(refreshedToken, ShouldBeEmpty)
			So(err, ShouldBeError)
		})
		acct.DeleteByID(acct.ID)
	})
}
