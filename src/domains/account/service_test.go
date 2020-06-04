package account

import (
	"testing"
	dao "williamfeng323/mooncake-duty/src/infrastructure/db"

	"go.mongodb.org/mongo-driver/bson/primitive"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateAccount(t *testing.T) {
	SkipConvey("Giving a db connection", t, func() {
		acct := &Account{}
		conn := dao.GetConnection()

		Convey("Should panic while no register Account", func() {
			So(func() { createAccount("", "", true) }, ShouldPanic)
		})
		Convey("Should return error When email or password does not provided", func() {
			conn.Register(acct)
			rst, err := createAccount("", "", true)
			So(rst, ShouldBeNil)
			So(err, ShouldBeError)
		})
		Convey("Should return inserted account objectId when email/password valid", func() {
			rst, err := createAccount("test123", "rstAbc.", true)
			So(rst, ShouldNotBeNil)
			So(err, ShouldBeNil)
			stt := rst.InsertedID.(primitive.ObjectID).Hex()
			id, err := primitive.ObjectIDFromHex(stt)
			conn.CollectionRegistry["Account"].New(acct)
			acct.DeleteByID(id)
		})
		Convey("Should return error When email duplicate with existing account", func() {
			initAcct, err := createAccount("test123", "rstAbcd", true)
			stt := initAcct.InsertedID.(primitive.ObjectID).Hex()
			id, err := primitive.ObjectIDFromHex(stt)
			rst, err := createAccount("test123", "rstAbcd", true)
			So(rst, ShouldBeNil)
			So(err, ShouldBeError)
			conn.CollectionRegistry["Account"].New(acct)
			acct.DeleteByID(id)
		})
	})
}

func TestFindByID(t *testing.T) {
	SkipConvey("Giving a db connection and init user", t, func() {
		acct := &Account{}
		conn := dao.GetConnection()
		conn.Register(acct)
		conn.CollectionRegistry["Account"].New(acct)
		acct.Email = "test@test.com"
		acct.Password = "password"
		rst, _ := acct.InsertAccount()
		id := rst.InsertedID.(primitive.ObjectID).Hex()
		Convey("should return the account document when find a valid id", func() {
			acct2, err := getAccountByID(id)
			So(acct2, ShouldNotBeZeroValue)
			So(err, ShouldBeNil)
		})
		Convey("should return nil account document when cannot find related id", func() {
			fakeID := primitive.NewObjectID()
			acct2, err := getAccountByID(fakeID.Hex())
			So(acct2, ShouldBeZeroValue)
			So(err, ShouldNotBeNil)
		})
		acct.DeleteByID(acct.ID)
	})
}

func TestUpdateAccount(t *testing.T) {
	SkipConvey("Giving a db connection an init user", t, func() {
		acct := &Account{}
		conn := dao.GetConnection()
		conn.Register(acct)
		conn.CollectionRegistry["Account"].New(acct)
		acct.Email = "test@test.com"
		acct.Password = "password"
		acct.InsertAccount()
		Convey("Should updated document correctly", func() {
			rst, err := updateAccount(acct.ID, "https://avatar/test/1234567", "1234567")
			So(rst, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(rst.ModifiedCount, ShouldEqual, 1)
			act, _ := getAccountByID(acct.ID.Hex())
			So(act.Mobile, ShouldEqual, "1234567")
			So(act.Avatar, ShouldEqual, "https://avatar/test/1234567")
		})
		acct.DeleteByID(acct.ID)
	})
}

func TestSignIn(t *testing.T) {
	SkipConvey("Giving a db connection and init user", t, func() {
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
	SkipConvey("Giving a db connection and init user", t, func() {
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
