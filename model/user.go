package model

import (
	"gitlab.com/target-smart-data-ai-searsh/task-connector-be/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a row from 'user'.
type User struct {
	UserID      primitive.ObjectID `json:"userId" bson:"_id,omitempty"`
	ID          string             `json:"id"  bson:"id"`
	Subject     string             `json:"subject"  bson:"subject"`
	CreatedDate string             `json:"createdDate"  bson:"createdDate"`
	ChangedDate string             `json:"changedDate"  bson:"changedDate"`
	DeletedDate string             `json:"deletedDate"  bson:"deletedDate"`
}

//Users struct
type Users struct {
	WSsID      primitive.ObjectID `json:"UsersId" bson:"_id,omitempty"` //_id
	TotalCount int64              `json:"totalCount" bson:"totalCount"` //totalCount
	Users      []User             `json:"Users" bson:"Users"`           //allWorkspaces
}

//AllUsers struct
type AllUsers struct {
	Data Users `json:"allUsers"` //allWorkspaces
}

//ParseValueToUser function parse to user object.
func ParseValueToUser(data map[string]interface{}) User {
	var user User
	common.GetMapStringValue(data, "id", &user.ID)
	common.GetMapStringValue(data, "subject", &user.Subject)
	common.GetMapStringValue(data, "created_date", &user.CreatedDate)
	common.GetMapStringValue(data, "changed_date", &user.ChangedDate)
	common.GetMapStringValue(data, "deleted_date", &user.DeletedDate)
	return user
}
