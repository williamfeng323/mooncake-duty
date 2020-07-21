package account

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	repoimpl "williamfeng323/mooncake-duty/src/infrastructure/db/repo_impl"
	"williamfeng323/mooncake-duty/src/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service exposes the account service
type Service struct {
	repo *repoimpl.AccountRepo
}

// SetRepo set the account repository to the service
func (as *Service) SetRepo(repo *repoimpl.AccountRepo) {
	as.repo = repo
}

// SignIn validate the email password pair and return token
func (as *Service) SignIn(email string, password string) (string, error) {
	ctx, cancel := utils.GetDefaultCtx()
	defer cancel()
	result := as.repo.FindOne(ctx, as.repo.EmailFilter(email))
	acct := &Account{}
	result.Decode(acct)
	if acct.ID.IsZero() {
		return "", NotFoundError{}
	}
	decryptedPassword, err := utils.Decrypt(acct.Password)
	if err != nil {
		return "", err
	}
	if string(decryptedPassword) != password {
		return "", fmt.Errorf("Password and account does not match")
	}
	token, err := utils.SignToken(acct.Email)
	return token, err
}

// Refresh refreshes the token
func (as *Service) Refresh(tokenString string) (string, error) {
	claims, err := utils.VerifyToken(tokenString, true)
	if err != nil {
		return "", err
	}
	refreshedToken, err := utils.SignToken(claims.Audience)
	if err != nil {
		return "", err
	}
	return refreshedToken, nil
}

// Register creates the basic account
func (as *Service) Register(email string, password string, isAdmin bool) (string, *Account, error) {
	acct, err := NewAccount(email, password)
	acct.IsAdmin = isAdmin
	if err != nil {
		return "", nil, err
	}
	_, err = acct.Save(false)
	if err != nil {
		return "", nil, err
	}
	token, err := utils.SignToken(acct.Email)
	return token, acct, err
}

// UpdateContactMethods update the ways to send notification
// you empty value will not be update into the account.
func (as *Service) UpdateContactMethods(id string, cm ContactMethods, email string, mobile string) error {

	account, err := as.GetAccountByID(id)
	if err != nil {
		return err
	}
	valueSet := bson.M{}
	if email != "" {
		valueSet["email"] = email
	}
	if mobile != "" {
		valueSet["mobile"] = mobile
	}
	if cm.SentEmail && account.Email == "" {
		return fmt.Errorf("Email must be set before you active send email notification")
	}
	if cm.SentSMS && account.Mobile == "" && mobile == "" {
		return fmt.Errorf("Mobile must be set before you active send email notification")
	}
	if !reflect.DeepEqual(cm, reflect.Zero(reflect.TypeOf(cm)).Interface()) {
		var cmMap bson.D
		originalCm, _ := bson.Marshal(account.ContactMethods)
		bson.Unmarshal(originalCm, &cmMap)
		updatedCm, _ := bson.Marshal(cm)
		bson.Unmarshal(updatedCm, &cmMap)
		valueSet["contactMethods"] = cmMap
		valueSet["updatedAt"] = time.Now().UTC()
	}
	if len(valueSet) == 0 {
		return fmt.Errorf("Nothing needs to be updated")
	}
	ctx, cancel := utils.GetDefaultCtx()
	defer cancel()
	_, err = as.repo.UpdateOne(ctx, bson.M{"_id": account.ID}, bson.M{"$set": valueSet})
	return err
}

// GrantSystemAdmin grants system admin auth to user
func (as *Service) GrantSystemAdmin(email string) error {
	account, err := as.GetAccountByEmail(email)
	if err != nil {
		return err
	}
	ctx, cancel := utils.GetDefaultCtx()
	defer cancel()
	_, err = as.repo.UpdateOne(ctx, bson.M{"_id": account.ID}, bson.M{"$set": bson.M{"isAdmin": true, "updatedAt": time.Now().UTC()}})
	return err
}

// DeactivateAccount logically deletes the account
func (as *Service) DeactivateAccount(email string) error {
	account, err := as.GetAccountByEmail(email)
	if err != nil {
		return err
	}
	ctx, cancel := utils.GetDefaultCtx()
	defer cancel()
	_, err = as.repo.UpdateOne(ctx, bson.M{"_id": account.ID}, bson.M{"$set": bson.M{"deleted": true, "updatedAt": time.Now().UTC()}})
	return err
}

// GetAccountByEmail returns the existing account
func (as *Service) GetAccountByEmail(email string) (*Account, error) {
	ctx, cancel := utils.GetDefaultCtx()
	defer cancel()
	rst := as.repo.FindOne(ctx, as.repo.EmailFilter(email))
	account := &Account{}
	rst.Decode(account)
	if account.ID.IsZero() {
		return nil, NotFoundError{}
	}
	return account, nil
}

// GetAccountByID returns the existing account
func (as *Service) GetAccountByID(id string) (*Account, error) {
	ctx, cancel := utils.GetDefaultCtx()
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, NotFoundError{}
	}
	rst := as.repo.FindOne(ctx, bson.M{"_id": objectID})
	account := &Account{}
	rst.Decode(account)
	if account.ID.IsZero() {
		return nil, NotFoundError{}
	}
	return account, nil
}

var accountService *Service
var accountServiceLock sync.RWMutex

// GetAccountService returns a singleton account service instance
func GetAccountService() *Service {
	accountServiceLock.Lock()
	defer accountServiceLock.Unlock()
	if accountService == nil {
		accountService = &Service{}
		accountService.SetRepo(repoimpl.GetAccountRepo())
	}
	return accountService
}
