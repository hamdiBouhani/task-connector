package model

import (
	"gitlab.com/target-smart-data-ai-searsh/task-connector-be/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tasks struct {
	TasksID    primitive.ObjectID `json:"tasksId" bson:"_id,omitempty"` //_id
	TotalCount int64              `json:"totalCount" bson:"totalCount"` //totalCount
	Tasks      []Task             `json:"tasks" bson:"tasks"`           //tasks
}

// Task represents a row from 'task'.
type Task struct {
	TaskID           primitive.ObjectID `json:"taskId" bson:"_id,omitempty"`              //_id
	ID               string             `json:"id" bson:"id"`                             // id
	Name             string             `json:"name" bson:"name"`                         // name
	Description      string             `json:"description" bson:"description"`           // description
	Tag              string             `json:"tag" bson:"tag"`                           // tag
	CreatedDate      string             `json:"createdDate" bson:"createdDate"`           // created_date
	ChangedDate      string             `json:"changedDate" bson:"changedDate"`           // changed_date
	DeletedDate      string             `json:"deletedDate" bson:"deletedDate"`           // deleted_date
	Deadline         string             `json:"deadline" bson:"deadline"`                 // deadline
	Sequence         string             `json:"sequence" bson:"sequence"`                 // sequence
	LastTickTime     string             `json:"lastTickTime" bson:"lastTickTime"`         // last_tick_time
	TrackedTime      string             `json:"trackedTime" bson:"trackedTime"`           // tracked_time
	IsTicking        bool               `json:"isTicking" bson:"isTicking"`               // is_ticking
	EstimateTime     int64              `json:"estimateTime" bson:"estimateTime"`         // estimate_time
	ActualTime       int64              `json:"actualTime" bson:"actualTime"`             // actual_time
	TaskAssigneeMaps TaskAssigneeMaps   `json:"taskAssigneeMaps" bson:"taskAssigneeMaps"` // taskAssigneeMaps
	TaskChecklists   TaskChecklists     `json:"taskChecklists" bson:"taskChecklists"`     //taskChecklists
	TaskTracktimes   TaskTracktimes     `json:"taskTracktimes" bson:"taskTracktimes"`     //taskTracktimes
}

//ParseValueToTask function parse to task object.
func ParseValueToTask(data map[string]interface{}) Task {
	var task Task
	common.GetMapStringValue(data, "id", &task.ID)
	common.GetMapStringValue(data, "name", &task.Name)
	common.GetMapStringValue(data, "description", &task.Description)
	common.GetMapStringValue(data, "tag", &task.Tag)
	common.GetMapStringValue(data, "created_date", &task.CreatedDate)
	common.GetMapStringValue(data, "changed_date", &task.ChangedDate)
	common.GetMapStringValue(data, "deleted_date", &task.DeletedDate)
	common.GetMapStringValue(data, "deadline", &task.Deadline)
	common.GetMapStringValue(data, "sequence", &task.Sequence)
	common.GetMapStringValue(data, "last_tick_time", &task.LastTickTime)
	common.GetMapStringValue(data, "tracked_time", &task.TrackedTime)
	common.GetMapBoolValue(data, "is_ticking", &task.IsTicking)
	common.GetMapInt64Value(data, "estimate_time", &task.EstimateTime)
	common.GetMapInt64Value(data, "actual_time", &task.ActualTime)
	return task
}

// TaskAssigneeMaps represents a list of asignne maps
type TaskAssigneeMaps struct {
	TaskAssigneeMapsID primitive.ObjectID `json:"taskAssigneeMapsId" bson:"_id,omitempty"`  //_id
	TotalCount         int64              `json:"totalCount" bson:"totalCount"`             //totalCount
	TaskAssigneeMaps   []TaskAssigneeMap  `json:"taskAssigneeMaps" bson:"taskAssigneeMaps"` //taskAssigneeMaps
}

// TaskAssigneeMap represents a row from 'task_assignee_map'.
type TaskAssigneeMap struct {
	TaskAssigneeMapID primitive.ObjectID `json:"taskAssigneeMapId" bson:"_id,omitempty"` //_id
	ID                string             `json:"id" bson:"id"`                           // id
	User              User               `json:"user" bson:"user" `                      // user_id
	CreatedDate       string             `json:"createdDate" bson:"createdDate"`         // created_date
	ChangedDate       string             `json:"changedDate" bson:"changedDate"`         // changed_date
	DeletedDate       string             `json:"deletedDate" bson:"deletedDate"`         // deleted_date
}

//ParseValueToTaskAssigneeMap function parse to TaskAssigneeMap object.
func ParseValueToTaskAssigneeMap(data map[string]interface{}) TaskAssigneeMap {
	var task TaskAssigneeMap
	common.GetMapStringValue(data, "id", &task.ID)
	common.GetMapStringValue(data, "created_date", &task.CreatedDate)
	common.GetMapStringValue(data, "changed_date", &task.ChangedDate)
	common.GetMapStringValue(data, "deleted_date", &task.DeletedDate)

	return task
}

// TaskAttachment represents a list of attachment per task
type TaskAttachments struct {
	TaskAttachmentsID primitive.ObjectID `json:"taskAttachmentsId" bson:"_id,omitempty"` //_id
	TotalCount        int64              `json:"totalCount" bson:"totalCount"`           //totalCount
	TaskAttachments   []TaskAttachment   `json:"taskAttachments" bson:"taskAttachments"` //taskAttachments
}

// TaskAttachment represents a row from 'task_attachment'.
type TaskAttachment struct {
	TaskAttachmentID primitive.ObjectID `json:"taskAttachmentId" bson:"_id,omitempty"` //_id
	ID               string             `json:"id" bson:"id"`                          // id
	Name             string             `json:"name" bson:"name"`                      // name
	URL              string             `json:"url" bson:"url"`                        // url
	Size             int64              `json:"size" bson:"size"`                      // size
	Uploader         int64              `json:"uploader" bson:"uploader"`              // uploader
	TaskID           int64              `json:"task_id" bson:"task_id"`                // task_id
	CreatedDate      string             `json:"created_date" bson:"created_date"`      // created_date
	ChangedDate      string             `json:"changed_date" bson:"changed_date"`      // changed_date
	DeletedDate      string             `json:"deleted_date" bson:"deleted_date"`      // deleted_date
}

type TaskChecklists struct {
	TaskChecklistsID primitive.ObjectID `json:"taskChecklistsId" bson:"_id,omitempty"`
	TotalCount       int64              `json:"totalCount" bson:"totalCount"`         //totalCount
	TaskChecklists   []TaskChecklist    `json:"taskChecklists" bson:"taskChecklists"` //taskChecklists
}

// TaskChecklist represents a row from 'task_checklist'.
type TaskChecklist struct {
	TaskChecklistID primitive.ObjectID `json:"taskChecklistId" bson:"_id,omitempty"`
	ID              string             `json:"id" bson:"id"`                   // id
	Content         string             `json:"content" bson:"content"`         // content
	Finished        bool               `json:"finished" bson:"finished"`       // finished
	FinishBy        User               `json:"finishBy" bson:"finishBy"`       // finish_by
	CreatedDate     string             `json:"createdDate" bson:"createdDate"` // created_date
	ChangedDate     string             `json:"changedDate" bson:"changedDate"` // changed_date
	DeletedDate     string             `json:"deletedDate" bson:"deletedDate"` // deleted_date
	CreatedBy       User               `json:"createdBy" bson:"createdBy"`     // created_by
}

//ParseValueToTaskChecklist function parse to TaskChecklist object.
func ParseValueToTaskChecklist(data map[string]interface{}) TaskChecklist {
	var task TaskChecklist
	common.GetMapStringValue(data, "id", &task.ID)
	common.GetMapStringValue(data, "content", &task.Content)
	common.GetMapBoolValue(data, "finished", &task.Finished)
	common.GetMapStringValue(data, "created_date", &task.CreatedDate)
	common.GetMapStringValue(data, "changed_date", &task.ChangedDate)
	common.GetMapStringValue(data, "deleted_date", &task.DeletedDate)

	return task
}

type TaskTracktimes struct {
	TaskTracktimesID primitive.ObjectID `json:"taskTracktimesId" bson:"_id,omitempty"`
	TotalCount       int64              `json:"totalCount" bson:"totalCount"`         //totalCount
	TaskTracktimes   []TaskTracktime    `json:"taskTracktimes" bson:"taskTracktimes"` //taskTracktimes
}

// TaskTracktime represents a row from 'task_tracktime'.
type TaskTracktime struct {
	TaskTracktimeID primitive.ObjectID `json:"taskTracktimeId" bson:"_id,omitempty"`
	ID              string             `json:"id" bson:"id"`                   // id
	StartTime       string             `json:"startTime" bson:"startTime"`     // start_time
	EndTime         string             `json:"endTime" bson:"endTime"`         // end_time
	Duration        string             `json:"duration" bson:"duration"`       // duration
	Finished        bool               `json:"finished" bson:"finished"`       // finished
	TrackType       string             `json:"trackType" bson:"trackType"`     // track_type
	CreatedDate     string             `json:"createdDate" bson:"createdDate"` // created_date
	ChangedDate     string             `json:"changedDate" bson:"changedDate"` // changed_date
	DeletedDate     string             `json:"deletedDate" bson:"deletedDate"` // deleted_date
	User            User               `json:"user" bson:"user"`               // user
}

//ParseValueToTaskTracktime function parse to TasTaskTracktimekChecklist object.
func ParseValueToTaskTracktime(data map[string]interface{}) TaskTracktime {
	var task TaskTracktime
	common.GetMapStringValue(data, "id", &task.ID)
	common.GetMapStringValue(data, "start_time", &task.StartTime)
	common.GetMapStringValue(data, "end_time", &task.EndTime)
	common.GetMapStringValue(data, "duration", &task.Duration)
	common.GetMapBoolValue(data, "finished", &task.Finished)
	common.GetMapStringValue(data, "created_date", &task.CreatedDate)
	common.GetMapStringValue(data, "changed_date", &task.ChangedDate)
	common.GetMapStringValue(data, "deleted_date", &task.DeletedDate)

	return task
}

type TaskDTO struct {
	TaskID       primitive.ObjectID `json:"taskId" bson:"_id,omitempty"`
	ID           string             `json:"id" bson:"id"`                     // id
	Name         string             `json:"name" bson:"name"`                 // name
	Description  string             `json:"description" bson:"description"`   // description
	Tag          string             `json:"tag" bson:"tag"`                   // tag
	CreatedDate  string             `json:"createdDate" bson:"createdDate"`   // created_date
	ChangedDate  string             `json:"changedDate" bson:"changedDate"`   // changed_date
	DeletedDate  string             `json:"deletedDate" bson:"deletedDate"`   // deleted_date
	Deadline     string             `json:"deadline" bson:"deadline"`         // deadline
	Sequence     string             `json:"sequence" bson:"sequence"`         // sequence
	LastTickTime string             `json:"lastTickTime" bson:"lastTickTime"` // last_tick_time
	TrackedTime  string             `json:"trackedTime" bson:"trackedTime"`   // tracked_time
	IsTicking    bool               `json:"isTicking" bson:"isTicking"`       // is_ticking
	EstimateTime int64              `json:"estimateTime" bson:"estimateTime"` // estimate_time
	ActualTime   int64              `json:"actualTime" bson:"actualTime"`     // actual_time
}

//ParseValueToTask function parse to task object.
func ParseValueToTaskDTO(data map[string]interface{}) TaskDTO {
	var task TaskDTO
	common.GetMapStringValue(data, "id", &task.ID)
	common.GetMapStringValue(data, "name", &task.Name)
	common.GetMapStringValue(data, "description", &task.Description)
	common.GetMapStringValue(data, "tag", &task.Tag)
	common.GetMapStringValue(data, "created_date", &task.CreatedDate)
	common.GetMapStringValue(data, "changed_date", &task.ChangedDate)
	common.GetMapStringValue(data, "deleted_date", &task.DeletedDate)
	common.GetMapStringValue(data, "deadline", &task.Deadline)
	common.GetMapStringValue(data, "sequence", &task.Sequence)
	common.GetMapStringValue(data, "last_tick_time", &task.LastTickTime)
	common.GetMapStringValue(data, "tracked_time", &task.TrackedTime)
	common.GetMapBoolValue(data, "is_ticking", &task.IsTicking)
	common.GetMapInt64Value(data, "estimate_time", &task.EstimateTime)
	common.GetMapInt64Value(data, "actual_time", &task.ActualTime)
	return task
}
