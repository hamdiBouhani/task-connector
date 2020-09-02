package model

import (
	"gitlab.com/target-smart-data-ai-searsh/task-connector-be/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Streams struct {
	StreamsID  primitive.ObjectID `json:"streamsId" bson:"_id,omitempty"`
	TotalCount int64              `json:"totalCount" bson:"totalCount"` //totalCount
	Streams    []Stream           `json:"streams" bson:"streams"`       //Streams
}

// Stream represents a row from 'stream'.
type Stream struct {
	StreamID    primitive.ObjectID `json:"streamId" bson:"_id,omitempty"`
	ID          string             `json:"id" bson:"id"`                   // id
	Name        string             `json:"name" bson:"name"`               // name
	Description string             `json:"description" bson:"description"` // description
	SortedTasks string             `json:"sortedTasks" bson:"sortedTasks"` // sorted_tasks
	CreatedDate string             `json:"createdDate" bson:"createdDate"` // created_date
	ChangedDate string             `json:"changedDate" bson:"changedDate"` // changed_date
	DeletedDate string             `json:"deletedDate" bson:"deletedDate"` // deleted_date
	Sequence    string             `json:"sequence" bson:"sequence"`       // sequence
	MeetingID   string             `json:"meetingID" bson:"meetingID"`     // meeting_id
	Tasks       Tasks              `json:"tasks" bson:"taskList"`          //tasks
}

//ParseValueToStream function parse to user object.
func ParseValueToStream(data map[string]interface{}) Stream {
	var stream Stream
	common.GetMapStringValue(data, "id", &stream.ID)
	common.GetMapStringValue(data, "name", &stream.Name)
	common.GetMapStringValue(data, "description", &stream.Description)
	common.GetMapStringValue(data, "sorted_tasks", &stream.SortedTasks)
	common.GetMapStringValue(data, "created_date", &stream.CreatedDate)
	common.GetMapStringValue(data, "changed_date", &stream.ChangedDate)
	common.GetMapStringValue(data, "deleted_date", &stream.DeletedDate)
	common.GetMapStringValue(data, "sequence", &stream.Sequence)
	common.GetMapStringValue(data, "meeting_id", &stream.MeetingID)
	return stream
}
