package model

import (
	"gitlab.com/target-smart-data-ai-searsh/task-connector-be/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Workspace represents a row from 'workspace'.
type Workspace struct {
	WsID        primitive.ObjectID `json:"wsId" bson:"_id,omitempty"`      //_id
	ID          string             `json:"id" bson:"id"`                   // id
	Name        string             `json:"name" bson:"name"`               // name
	Typ         string             `json:"typ" bson:"typ"`                 // typ
	StartDate   string             `json:"startDate" bson:"startDate"`     // start_date
	EndDate     string             `json:"endDate" bson:"endDate"`         // end_date
	IconURL     string             `json:"iconURL" bson:"iconURL"`         // icon_url
	CoverURL    string             `json:"coverURL" bson:"coverURL"`       // cover_url
	Alias       string             `json:"alias" bson:"alias"`             // alias
	Description string             `json:"description" bson:"description"` // description
	IsPrivate   bool               `json:"isPrivate" bson:"isPrivate"`     // is_private
	IsOngoing   bool               `json:"isOngoing" bson:"isOngoing"`     // is_ongoing
	IsActive    bool               `json:"isActive" bson:"isActive"`       // is_active
	CreatedBy   User               `json:"createdBy" bson:"createdBy"`     // created_by
	Metadata    string             `json:"metadata" bson:"metadata"`       // metadata
	CreatedDate string             `json:"createdDate" bson:"createdDate"` // created_date
	ChangedDate string             `json:"changedDate" bson:"changedDate"` // changed_date
	DeletedDate string             `json:"deletedDate" bson:"deletedDate"` // deleted_date
	Expired     bool               `json:"expired" bson:"expired"`         // expired
	Streams     Streams            `json:"streams" bson:"streamList"`      //Streams
	Events      Events             `json:"events" bson:"eventList"`        //events
}

type Workspaces struct {
	WSsID      primitive.ObjectID `json:"wssId" bson:"_id,omitempty"`   //_id
	TotalCount int64              `json:"totalCount" bson:"totalCount"` //totalCount
	Workspaces []Workspace        `json:"workspaces" bson:"workspaces"` //allWorkspaces
}

type AllWorkspaces struct {
	Data Workspaces `json:"allWorkspaces"` //allWorkspaces
}

func ParseValueToWorkspace(data map[string]interface{}) Workspace {
	var ws Workspace
	common.GetMapStringValue(data, "id", &ws.ID)
	common.GetMapStringValue(data, "name", &ws.Name)
	common.GetMapStringValue(data, "typ", &ws.Typ)
	common.GetMapStringValue(data, "start_date", &ws.StartDate)
	common.GetMapStringValue(data, "end_date", &ws.EndDate)
	common.GetMapStringValue(data, "icon_url", &ws.IconURL)
	common.GetMapStringValue(data, "cover_url", &ws.CoverURL)
	common.GetMapStringValue(data, "alias", &ws.Alias)
	common.GetMapStringValue(data, "description", &ws.Description)
	common.GetMapBoolValue(data, "is_private", &ws.IsPrivate)
	common.GetMapBoolValue(data, "is_ongoing", &ws.IsOngoing)
	common.GetMapBoolValue(data, "is_active", &ws.IsActive)
	common.GetMapStringValue(data, "metadata", &ws.Metadata)
	common.GetMapStringValue(data, "created_date", &ws.CreatedDate)
	common.GetMapStringValue(data, "changed_date", &ws.ChangedDate)
	common.GetMapStringValue(data, "deleted_date", &ws.DeletedDate)
	common.GetMapBoolValue(data, "expired", &ws.Expired)

	return ws
}

type WsTasks struct {
	WsID  primitive.ObjectID `bson:"_id,omitempty"` //_id
	ID    string             `bson:"id"`            // id
	Tasks []string           `bson:"tasks"`
}
