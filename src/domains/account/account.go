package account

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	db "williamfeng323/mooncake-duty/src/infrastructure/db"
	repoimpl "williamfeng323/mooncake-duty/src/infrastructure/db/repo_impl"
	validatorimpl "williamfeng323/mooncake-duty/src/infrastructure/db/validator_impl"
	"williamfeng323/mooncake-duty/src/utils"
)

// SentHook is the structure to describe the alternative way to send the alarms.
type SentHook struct {
	URL       string `json:"url" bson:"url"`
	IsEnabled bool   `json:"isEnabled" bson:"isEnabled"`
}

// ContactMethods is describing the way to contact user
type ContactMethods struct {
	SentHook  `json:"sentHook" bson:"sentHook"`
	SentSMS   bool `json:"sentSMS" bson:"sentSMS"`
	SentEmail bool `json:"sentEmail" bson:"sentEmail"`
}

// Account struct of the user account
type Account struct {
	repo           *repoimpl.AccountRepo
	db.BaseModel   `json:",inline" bson:",inline"`
	Email          string               `json:"email" bson:"email" required:"true"` // account email is unique to-do add the unique index
	Password       string               `json:"password" bson:"password" required:"true"`
	Mobile         string               `json:"mobile,omitempty" bson:"mobile,omitempty"`
	IsAdmin        bool                 `json:"isAdmin" bson:"isAdmin"`
	Projects       []primitive.ObjectID `json:"projects,omitempty" bson:"projects, omitempty"`
	ContactMethods `json:"contactMethods,omitempty" bson:"contactMethods,omitempty"`
}

// NewAccount returns an account instance with email and password
func NewAccount(email string, password string) (*Account, error) {
	account := &Account{Email: email}
	encryptedPassword, err := utils.Encrypt(password)
	if err != nil {
		return nil, err
	}
	account.ID = primitive.NewObjectID()
	account.Password = string(encryptedPassword)
	account.CreatedAt = time.Now()
	account.repo = repoimpl.GetAccountRepo()
	return account, nil
}

// Save creates or updates an account document into database, returns the success count and error
func (acct *Account) Save(allowReplace bool) (int, error) {
	validator := validatorimpl.NewDefaultValidator()
	errs := validator.Verify(acct)
	if len(errs) != 0 {
		return 0, fmt.Errorf("Save the account failed due to: %v", errs)
	}
	ctx, cancel := utils.GetDefaultCtx()
	defer cancel()
	rst := acct.repo.FindOne(ctx, acct.repo.EmailFilter(acct.Email))
	tempAcct := &Account{}
	rst.Decode(tempAcct)

	if tempAcct.ID.IsZero() {
		rst, err := acct.repo.InsertOne(ctx, acct)
		if err != nil {
			return 0, err
		}
	} else if allowReplace {
		acct.UpdatedAt = time.Now()
		var inInterface bson.M
		inrec, _ := bson.Marshal(acct)
		bson.Unmarshal(inrec, &inInterface)
		delete(inInterface, "_id")
		rst, err := acct.repo.UpdateOne(ctx, acct.repo.EmailFilter(acct.Email), bson.M{"$set": inInterface})
		if err != nil {
			return 0, err
		}
	} else {
		return 0, AlreadyExistError{}
	}
	return 1, nil
}

// func refresh(tokenString string) (string, error) {
// 	claims, err := utils.VerifyToken(tokenString)
// 	if err != nil {
// 		return "", err
// 	}
// 	refreshedToken, err := utils.SignToken(claims.Audience)
// 	if err != nil {
// 		return "", err
// 	}
// 	return refreshedToken, nil
// }
