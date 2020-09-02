package model

import (
	"gitlab.com/target-smart-data-ai-searsh/task-connector-be/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Events struct {
	EventsID   primitive.ObjectID `bson:"_id,omitempty"`
	TotalCount int64              `json:"totalCount" bson:"totalCount"` //totalCount
	Events     []Event            `json:"events" bson:"events"`
}

type Event struct {
	EventID   primitive.ObjectID `bson:"_id,omitempty"`
	ID        string             `json:"id" bson:"id"`               // id
	Timestamp string             `json:"timestamp" bson:"timestamp"` //timestamp
	HolderTyp string             `json:"holderTyp" bson:"holderTyp"` //holderTyp
	Brief     string             `json:"brief" bson:"brief"`         //brief
	Detail    string             `json:"detail" bson:"detail"`       //detail
	Task      TaskDTO            `json:"task" bson:"task"`           //task
}

//ParseValueToEvent function parse to event object.
func ParseValueToEvent(data map[string]interface{}) Event {
	var event Event
	common.GetMapStringValue(data, "id", &event.ID)
	common.GetMapStringValue(data, "timestamp", &event.Timestamp)
	common.GetMapStringValue(data, "holder_typ", &event.HolderTyp)
	common.GetMapStringValue(data, "brief", &event.Brief)
	common.GetMapStringValue(data, "detail", &event.Detail)
	return event
}
