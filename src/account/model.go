package account

import (
	"williamfeng323/mooncake-duty/src/dao"
)

//Role the role struct type for account.
type Role struct {
	dao.BaseModel
	Name string `json:"name" bson:"name"`
}

// Account struct of the user account
type Account struct {
	dao.BaseModel
	Projects []string `json:"projects" bson:"projects"`
	Email    string   `json:"email" bson:"email"`
	Password string   `json:"password" bson:"password"`
	Roles    []Role   `json:"roles" bson:"roles"`
}
