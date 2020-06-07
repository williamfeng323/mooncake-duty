package account

import (
	dao "williamfeng323/mooncake-duty/src/infrastructure/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SentHook is the structure to describe the alternative way to send the alarms.
type SentHook struct {
	URL string `json:"url" bson:"url"`
}

// ContactMethods is describing the way to contact user
type ContactMethods struct {
	SentHook  `json:"sentHook,inline" bson:"sentHook,inline"`
	SentSMS   bool `json:"sentSMS,omitempty" bson:"sentSMS,omitempty"`
	SentEmail bool `json:"sentEmail,omitempty" bson:"sentEmail,omitempty"`
	IsEnabled bool `json:"isEnabled" bson:"isEnabled"`
}

// Account struct of the user account
type Account struct {
	dao.BaseModel  `json:",inline" bson:",inline"`
	Email          string               `json:"email" bson:"email" required:"true"`
	Password       string               `json:"password" bson:"password" required:"true"`
	Mobile         string               `json:"mobile,omitempty" bson:"mobile,omitempty"`
	IsAdmin        bool                 `json:"isAdmin" bson:"isAdmin"`
	Avatar         string               `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Teams          []primitive.ObjectID `json:"teams,omitempty" bson:"teams, omitempty"`
	ContactMethods `json:"contactMethods,omitempty" bson:"contactMethods,omitempty"`
}
