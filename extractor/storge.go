package extractor

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.com/target-smart-data-ai-searsh/task-connector-be/common"
	"gitlab.com/target-smart-data-ai-searsh/task-connector-be/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//Insert fn

//InsertWorkspacs insert new work space in database.
func (s *Server) InsertWorkspacs(data map[string]interface{}) error {
	ws := model.ParseValueToWorkspace(data)

	if data["created_by"] != nil {
		var createdBy int64
		common.GetMapInt64Value(data, "created_by", &createdBy)
		collection := s.mongoClt.Database("wsTasks").Collection("users")

		cur := collection.FindOne(context.Background(), bson.M{"id": createdBy})
		if cur.Err() != nil {
			return errors.Wrap(cur.Err(), "failed to get list of notifications")
		}
		var user model.User

		err := cur.Decode(&user)
		if err != nil && err != mongo.ErrNoDocuments {
			return errors.Wrap(err, "failed to decode result")
		}

		if err != nil && err != mongo.ErrNoDocuments {
			return errors.Wrap(err, "fail to find user ")
		}

		ws.CreatedBy = user
	}

	collection := s.mongoClt.Database("wsTasks").Collection("workspaces")

	_, err := collection.InsertOne(context.Background(), ws)
	if err != nil {
		return errors.Wrap(err, "failed to insert a new user")
	}

	s.SendOneWorkspaceToPubSub(ws)

	return nil
}

//InsertEvent insert new event in database.
func (s *Server) InsertEvent(data map[string]interface{}) error {
	var wsId string
	common.GetMapStringValue(data, "ws_id", &wsId)

	event := model.ParseValueToEvent(data)

	collection := s.mongoClt.Database("wsTasks").Collection("workspaces")
	cur := collection.FindOne(context.Background(), bson.M{"id": wsId})
	if cur.Err() != nil {
		return errors.Wrap(cur.Err(), "failed to get workspace")
	}
	var workspace model.Workspace

	err := cur.Decode(&workspace)
	if err != nil && err != mongo.ErrNoDocuments {
		return errors.Wrap(err, "failed to decode result")
	}

	if len(workspace.Events.Events) > 0 {
		workspace.Events.Events = append(workspace.Events.Events, event)
		workspace.Events.TotalCount = int64(len(workspace.Events.Events))
	} else {
		workspace.Events.Events = append(workspace.Events.Events, event)
		workspace.Events.TotalCount = 1
	}

	_, err = collection.ReplaceOne(context.Background(), bson.M{"id": wsId}, workspace)
	if err != nil {
		return errors.Wrap(err, "failed to replace workspace")
	}

	s.SendOneWorkspaceToPubSub(workspace)

	return nil
}

//InsertStream insert new steam.
func (s *Server) InsertStream(data map[string]interface{}) error {
	var wsId string
	common.GetMapStringValue(data, "ws_id", &wsId)

	stream := model.ParseValueToStream(data)

	collection := s.mongoClt.Database("wsTasks").Collection("workspaces")
	cur := collection.FindOne(context.Background(), bson.M{"id": wsId})
	if cur.Err() != nil {
		return errors.Wrap(cur.Err(), "failed to get list of notifications")
	}
	var workspace model.Workspace

	err := cur.Decode(&workspace)
	if err != nil && err != mongo.ErrNoDocuments {
		return errors.Wrap(err, "failed to decode result")
	}

	if len(workspace.Streams.Streams) > 0 {
		workspace.Streams.Streams = append(workspace.Streams.Streams, stream)
		workspace.Streams.TotalCount = int64(len(workspace.Streams.Streams))
	} else {
		workspace.Streams.Streams = append(workspace.Streams.Streams, stream)
		workspace.Streams.TotalCount = 1
	}

	_, err = collection.ReplaceOne(context.Background(), bson.M{"id": wsId}, workspace)
	if err != nil {
		return errors.Wrap(err, "failed to replace workspace")
	}

	s.SendOneWorkspaceToPubSub(workspace)

	return nil
}

//InsertTask insert new task.
func (s *Server) InsertTask(data map[string]interface{}) error {
	var streamId string
	common.GetMapStringValue(data, "stream_id", &streamId)

	collection := s.mongoClt.Database("wsTasks").Collection("workspaces")
	cur := collection.FindOne(context.Background(), bson.D{{"streamList.streams.id", bson.D{{"$lte", streamId}}}})
	if cur.Err() != nil {
		return errors.Wrap(cur.Err(), "failed to get list of workspace using stream_id")
	}
	var workspace model.Workspace

	err := cur.Decode(&workspace)
	if err != nil && err != mongo.ErrNoDocuments {
		return errors.Wrap(err, "failed to decode result")
	}
	var tasks model.Tasks
	var index int
	for k, v := range workspace.Streams.Streams {
		index = k
		if v.ID == streamId {
			tasks = v.Tasks
			break
		}
	}

	if len(tasks.Tasks) == 0 {
		tasks.Tasks = append(tasks.Tasks, model.ParseValueToTask(data))
		tasks.TotalCount = 1
	} else {
		tasks.Tasks = append(tasks.Tasks, model.ParseValueToTask(data))
		tasks.TotalCount = int64(len(tasks.Tasks))
	}

	workspace.Streams.Streams[index].Tasks = tasks

	_, err = collection.ReplaceOne(context.Background(), bson.M{"id": workspace.ID}, workspace)
	if err != nil {
		return errors.Wrap(err, "failed to replace workspace")
	}

	s.SendOneWorkspaceToPubSub(workspace)

	return nil
}

//GetUserById get one user by id
func (s *Server) GetUserById(id string) (*model.User, error) {
	collection := s.mongoClt.Database("wsTasks").Collection("users")

	cur := collection.FindOne(context.Background(), bson.M{"id": id})
	if cur.Err() != nil {
		return nil, errors.Wrap(cur.Err(), "failed to get list of notifications")
	}
	var user model.User

	err := cur.Decode(&user)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, errors.Wrap(err, "failed to decode result")
	}

	if err != nil && err != mongo.ErrNoDocuments {
		return nil, errors.Wrap(err, "fail to find user ")
	}
	return &user, nil
}

//AddAssigneeMapToTask insert new assignee map to task.
func AddAssigneeMapToTask(task *model.Task, assigneeMaps model.TaskAssigneeMap) {
	if len(task.TaskAssigneeMaps.TaskAssigneeMaps) == 0 {
		task.TaskAssigneeMaps = model.TaskAssigneeMaps{
			TotalCount:       1,
			TaskAssigneeMaps: []model.TaskAssigneeMap{assigneeMaps},
		}
	} else {
		task.TaskAssigneeMaps.TaskAssigneeMaps = append(task.TaskAssigneeMaps.TaskAssigneeMaps, assigneeMaps)
		task.TaskAssigneeMaps.TotalCount = int64(len(task.TaskAssigneeMaps.TaskAssigneeMaps))
	}
}

//InsertTaskAssigneeMap insert new task assignee map.
func (s *Server) InsertTaskAssigneeMap(data map[string]interface{}) error {
	var taskID string
	common.GetMapStringValue(data, "task_id", &taskID)

	collection := s.mongoClt.Database("wsTasks").Collection("workspaces")
	cur := collection.FindOne(context.Background(), bson.D{{"streamList.streams.taskList.tasks.id", bson.D{{"$lte", taskID}}}})
	if cur.Err() != nil {
		return errors.Wrap(cur.Err(), "failed to get list of workspace  using task_id")
	}
	var workspace model.Workspace

	err := cur.Decode(&workspace)
	if err != nil && err != mongo.ErrNoDocuments {
		return errors.Wrap(err, "failed to decode result")
	}

	tam := model.ParseValueToTaskAssigneeMap(data)

	if data["user_id"] != nil {
		var userID string
		common.GetMapStringValue(data, "user_id", &userID)
		user, err := s.GetUserById(userID)
		if err != nil {
			return errors.Wrap(err, "failed to get user")
		}
		if user != nil {
			tam.User = *user
		}
	}

	ok := false
	for _, steam := range workspace.Streams.Streams {
		for k, _ := range steam.Tasks.Tasks {
			if steam.Tasks.Tasks[k].ID == taskID {
				AddAssigneeMapToTask(&steam.Tasks.Tasks[k], tam)
				ok = true
				break
			}
		}
		if ok {
			break
		}
	}
	_, err = collection.ReplaceOne(context.Background(), bson.M{"id": workspace.ID}, workspace)
	if err != nil {
		return errors.Wrap(err, "failed to replace workspace")
	}

	_, err = collection.ReplaceOne(context.Background(), bson.M{"id": workspace.ID}, workspace)
	if err != nil {
		return errors.Wrap(err, "failed to replace workspace")
	}

	s.SendOneWorkspaceToPubSub(workspace)

	return nil
}

//AddChecklistToTask insert new assignee map to task.
func AddChecklistToTask(task *model.Task, checklist model.TaskChecklist) {
	if len(task.TaskChecklists.TaskChecklists) == 0 {
		task.TaskChecklists = model.TaskChecklists{
			TotalCount:     1,
			TaskChecklists: []model.TaskChecklist{checklist},
		}
	} else {
		task.TaskChecklists.TaskChecklists = append(task.TaskChecklists.TaskChecklists, checklist)
		task.TaskChecklists.TotalCount = int64(len(task.TaskChecklists.TaskChecklists))
	}
}

//InsertTaskChecklist insert new task check list.
func (s *Server) InsertTaskChecklist(data map[string]interface{}) error {
	var taskID string
	common.GetMapStringValue(data, "task_id", &taskID)

	collection := s.mongoClt.Database("wsTasks").Collection("workspaces")
	cur := collection.FindOne(context.Background(), bson.D{{"streamList.streams.taskList.tasks.id", bson.D{{"$lte", taskID}}}})
	if cur.Err() != nil {
		return errors.Wrap(cur.Err(), "failed to get list of notifications")
	}
	var workspace model.Workspace

	err := cur.Decode(&workspace)
	if err != nil && err != mongo.ErrNoDocuments {
		return errors.Wrap(err, "failed to decode result")
	}

	tcl := model.ParseValueToTaskChecklist(data)

	if data["finish_by"] != nil {

		var userID string
		common.GetMapStringValue(data, "finish_by", &userID)
		user, err := s.GetUserById(userID)
		if err != nil {
			return errors.Wrap(err, "failed to get user")
		}
		if user != nil {
			tcl.FinishBy = *user
		}

	}

	if data["created_by"] != nil {

		var userID string
		common.GetMapStringValue(data, "created_by", &userID)
		user, err := s.GetUserById(userID)
		if err != nil {
			return errors.Wrap(err, "failed to get user")
		}
		if user != nil {
			tcl.CreatedBy = *user
		}

	}

	ok := false
	for _, steam := range workspace.Streams.Streams {
		for k, _ := range steam.Tasks.Tasks {
			if steam.Tasks.Tasks[k].ID == taskID {
				AddChecklistToTask(&steam.Tasks.Tasks[k], tcl)
				ok = true
				break
			}
		}
		if ok {
			break
		}
	}

	_, err = collection.ReplaceOne(context.Background(), bson.M{"id": workspace.ID}, workspace)
	if err != nil {
		return errors.Wrap(err, "failed to replace workspace")
	}

	s.SendOneWorkspaceToPubSub(workspace)

	return nil
}

//AddTracktimeToTask insert new assignee map to task.
func AddTracktimeToTask(task *model.Task, tracktime model.TaskTracktime) {
	if len(task.TaskTracktimes.TaskTracktimes) == 0 {
		task.TaskTracktimes = model.TaskTracktimes{
			TotalCount:     1,
			TaskTracktimes: []model.TaskTracktime{tracktime},
		}
	} else {
		task.TaskTracktimes.TaskTracktimes = append(task.TaskTracktimes.TaskTracktimes, tracktime)
		task.TaskTracktimes.TotalCount = int64(len(task.TaskTracktimes.TaskTracktimes))
	}
}

//InsertTaskTracktime insert new task check list.
func (s *Server) InsertTaskTracktime(data map[string]interface{}) error {

	var taskID string
	common.GetMapStringValue(data, "task_id", &taskID)

	collection := s.mongoClt.Database("wsTasks").Collection("workspaces")
	cur := collection.FindOne(context.Background(), bson.D{{"streamList.streams.taskList.tasks.id", bson.D{{"$lte", taskID}}}})
	if cur.Err() != nil {
		return errors.Wrap(cur.Err(), "failed to get list of notifications")
	}
	var workspace model.Workspace

	err := cur.Decode(&workspace)
	if err != nil && err != mongo.ErrNoDocuments {
		return errors.Wrap(err, "failed to decode result")
	}
	taskTracktime := model.ParseValueToTaskTracktime(data)

	if data["user_id"] != nil {

		var userID string
		common.GetMapStringValue(data, "user_id", &userID)

		user, err := s.GetUserById(userID)
		if err != nil {
			return errors.Wrap(err, "failed to get user")
		}
		if user != nil {
			taskTracktime.User = *user
		}

	}

	ok := false
	for _, steam := range workspace.Streams.Streams {
		for k, _ := range steam.Tasks.Tasks {
			if steam.Tasks.Tasks[k].ID == taskID {
				AddTracktimeToTask(&steam.Tasks.Tasks[k], taskTracktime)
				ok = true
				break
			}
		}
		if ok {
			break
		}
	}

	_, err = collection.ReplaceOne(context.Background(), bson.M{"id": workspace.ID}, workspace)
	if err != nil {
		return errors.Wrap(err, "failed to replace workspace")
	}

	s.SendOneWorkspaceToPubSub(workspace)

	return nil
}

//InsertUser insert new task check list.
func (s *Server) InsertUser(data map[string]interface{}) error {
	user := model.ParseValueToUser(data)

	collection := s.mongoClt.Database("wsTasks").Collection("users")
	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return errors.Wrap(err, "failed to insert a new user")
	}

	return nil
}

//Update fn

//UpdateWorkspacs update existing work space in database.
func (s *Server) UpdateWorkspacs(data map[string]interface{}) error {
	ws := model.ParseValueToWorkspace(data)
	collection := s.mongoClt.Database("wsTasks").Collection("workspaces")
	cur := collection.FindOne(context.Background(), bson.M{"id": ws.ID})
	if cur.Err() != nil {
		return errors.Wrap(cur.Err(), "failed to get list of notifications")
	}
	var workspace model.Workspace

	err := cur.Decode(&workspace)
	if err != nil && err != mongo.ErrNoDocuments {
		return errors.Wrap(err, "failed to decode result")
	}

	workspace.Name = ws.Name
	workspace.Typ = ws.Typ
	workspace.EndDate = ws.EndDate
	workspace.IconURL = ws.IconURL
	workspace.CoverURL = ws.CoverURL
	workspace.Alias = ws.Alias
	workspace.Description = ws.Description
	workspace.IsPrivate = ws.IsPrivate
	workspace.IsOngoing = ws.IsOngoing
	workspace.IsActive = ws.IsActive
	workspace.ChangedDate = ws.ChangedDate
	workspace.DeletedDate = ws.DeletedDate
	workspace.Expired = ws.Expired

	_, err = collection.ReplaceOne(context.Background(), bson.M{"id": workspace.ID}, workspace)
	if err != nil {
		return errors.Wrap(err, "failed to replace workspace")
	}

	s.SendOneWorkspaceToPubSub(workspace)

	return nil
}

//UpdateEvent update existing event in database.
func (s *Server) UpdateEvent(data map[string]interface{}) error {

	var wsID string
	common.GetMapStringValue(data, "ws_id", &wsID)

	collection := s.mongoClt.Database("wsTasks").Collection("workspaces")
	cur := collection.FindOne(context.Background(), bson.M{"id": wsID})
	if cur.Err() != nil {
		return errors.Wrap(cur.Err(), "failed to get list of notifications")
	}
	var workspace model.Workspace

	err := cur.Decode(&workspace)
	if err != nil && err != mongo.ErrNoDocuments {
		return errors.Wrap(err, "failed to decode result")
	}

	event := model.ParseValueToEvent(data)

	_, err = collection.ReplaceOne(context.Background(), bson.M{"id": workspace.ID}, workspace)
	if err != nil {
		return errors.Wrap(err, "failed to replace workspace")
	}

	for k, v := range workspace.Events.Events {
		if v.ID == event.ID {
			workspace.Events.Events[k].Brief = event.Brief
			workspace.Events.Events[k].Detail = event.Detail
			workspace.Events.Events[k].HolderTyp = event.HolderTyp
			workspace.Events.Events[k].Timestamp = event.Timestamp
			break
		}
	}

	_, err = collection.ReplaceOne(context.Background(), bson.M{"id": workspace.ID}, workspace)
	if err != nil {
		return errors.Wrap(err, "failed to replace workspace")
	}

	s.SendOneWorkspaceToPubSub(workspace)

	return nil
}

//UpdateStream update existing steam.
func (s *Server) UpdateStream(data map[string]interface{}) error {

	var wsID string
	common.GetMapStringValue(data, "ws_id", &wsID)

	collection := s.mongoClt.Database("wsTasks").Collection("workspaces")
	cur := collection.FindOne(context.Background(), bson.M{"id": wsID})
	if cur.Err() != nil {
		return errors.Wrap(cur.Err(), "failed to get list of notifications")
	}
	var workspace model.Workspace

	err := cur.Decode(&workspace)
	if err != nil && err != mongo.ErrNoDocuments {
		return errors.Wrap(err, "failed to decode result")
	}

	stream := model.ParseValueToStream(data)

	for k, v := range workspace.Streams.Streams {
		if v.ID == stream.ID {
			workspace.Streams.Streams[k].Name = stream.Name
			workspace.Streams.Streams[k].Description = stream.Description
			workspace.Streams.Streams[k].SortedTasks = stream.SortedTasks
			workspace.Streams.Streams[k].ChangedDate = stream.ChangedDate
			workspace.Streams.Streams[k].Sequence = stream.Sequence
			workspace.Streams.Streams[k].MeetingID = stream.MeetingID
			break
		}
	}

	_, err = collection.ReplaceOne(context.Background(), bson.M{"id": workspace.ID}, workspace)
	if err != nil {
		return errors.Wrap(err, "failed to replace workspace")
	}

	s.SendOneWorkspaceToPubSub(workspace)

	return nil
}

//RemoveTaskFromStream function
func RemoveTaskFromStream(stream *model.Stream, taskID string) model.Task {
	var RemovedTask model.Task

	tasks := model.Tasks{
		TotalCount: 0,
		Tasks:      []model.Task{},
	}

	tasks.TasksID = stream.Tasks.TasksID
	for _, v := range stream.Tasks.Tasks {
		if v.ID != taskID {
			tasks.TotalCount++
			tasks.Tasks = append(tasks.Tasks, v)
		} else {
			RemovedTask = v
		}
	}

	stream.Tasks = tasks
	return RemovedTask
}

//UpdateTask update existing task.
func (s *Server) UpdateTask(data map[string]interface{}) error {

	var streamID string
	common.GetMapStringValue(data, "stream_id", &streamID)

	collection := s.mongoClt.Database("wsTasks").Collection("workspaces")
	cur := collection.FindOne(context.Background(), bson.D{{"streamList.streams.id", bson.D{{"$lte", streamID}}}})
	if cur.Err() != nil {
		return errors.Wrap(cur.Err(), "failed to get list of notifications")
	}
	var workspace model.Workspace

	err := cur.Decode(&workspace)
	if err != nil && err != mongo.ErrNoDocuments {
		return errors.Wrap(err, "failed to decode result")
	}

	task := model.ParseValueToTask(data)
	var streamIndex int
	isTaskchangeItStream := false
	var ChangeItTask model.Task

	for index, stream := range workspace.Streams.Streams {
		if stream.ID == streamID {
			streamIndex = index
		}

		for k, v := range stream.Tasks.Tasks {

			stream.Tasks.Tasks[k].Name = task.Name
			stream.Tasks.Tasks[k].Description = task.Description
			stream.Tasks.Tasks[k].Tag = task.Tag
			stream.Tasks.Tasks[k].ChangedDate = task.ChangedDate
			stream.Tasks.Tasks[k].DeletedDate = task.DeletedDate
			stream.Tasks.Tasks[k].Deadline = task.Deadline
			stream.Tasks.Tasks[k].Sequence = task.Sequence
			stream.Tasks.Tasks[k].LastTickTime = task.LastTickTime
			stream.Tasks.Tasks[k].TrackedTime = task.TrackedTime
			stream.Tasks.Tasks[k].IsTicking = task.IsTicking
			stream.Tasks.Tasks[k].EstimateTime = task.EstimateTime
			stream.Tasks.Tasks[k].ActualTime = task.ActualTime

			if v.ID == task.ID && stream.ID == streamID {
				break
			}

			if v.ID == task.ID && stream.ID != streamID {
				isTaskchangeItStream = true

				ChangeItTask = RemoveTaskFromStream(&workspace.Streams.Streams[index], v.ID)
			}
		}
	}

	if isTaskchangeItStream {
		workspace.Streams.Streams[streamIndex].Tasks.TotalCount++
		workspace.Streams.Streams[streamIndex].Tasks.Tasks = append(workspace.Streams.Streams[streamIndex].Tasks.Tasks, ChangeItTask)
	}

	_, err = collection.ReplaceOne(context.Background(), bson.M{"id": workspace.ID}, workspace)
	if err != nil {
		return errors.Wrap(err, "failed to replace workspace")
	}

	s.SendOneWorkspaceToPubSub(workspace)

	return nil
}

//UpdateTaskAssigneeMap update existing task assignee map.
func (s *Server) UpdateTaskAssigneeMap(data map[string]interface{}) error {
	return nil
}

//UpdateChecklistInTask update new checklist.
func UpdateChecklistInTask(task *model.Task, checklist model.TaskChecklist) {
	for k, v := range task.TaskChecklists.TaskChecklists {
		if v.ID == checklist.ID {
			task.TaskChecklists.TaskChecklists[k].Content = checklist.Content
			task.TaskChecklists.TaskChecklists[k].Finished = checklist.Finished
			task.TaskChecklists.TaskChecklists[k].FinishBy = checklist.FinishBy
			task.TaskChecklists.TaskChecklists[k].ChangedDate = checklist.ChangedDate
			task.TaskChecklists.TaskChecklists[k].DeletedDate = checklist.DeletedDate
			break

		}
	}
}

//UpdateTaskChecklist update existing task check list.
func (s *Server) UpdateTaskChecklist(data map[string]interface{}) error {
	var taskID string
	common.GetMapStringValue(data, "task_id", &taskID)

	collection := s.mongoClt.Database("wsTasks").Collection("workspaces")
	cur := collection.FindOne(context.Background(), bson.D{{"streamList.streams.taskList.tasks.id", bson.D{{"$lte", taskID}}}})
	if cur.Err() != nil {
		return errors.Wrap(cur.Err(), "failed to get list of notifications")
	}
	var workspace model.Workspace

	err := cur.Decode(&workspace)
	if err != nil && err != mongo.ErrNoDocuments {
		return errors.Wrap(err, "failed to decode result")
	}

	tcl := model.ParseValueToTaskChecklist(data)

	if data["finish_by"] != nil {

		var userID string
		common.GetMapStringValue(data, "finish_by", &userID)
		user, err := s.GetUserById(userID)
		if err != nil {
			return errors.Wrap(err, "failed to get user")
		}
		if user != nil {
			tcl.FinishBy = *user
		}

	}

	ok := false
	for _, steam := range workspace.Streams.Streams {
		for k, _ := range steam.Tasks.Tasks {
			if steam.Tasks.Tasks[k].ID == taskID {
				UpdateChecklistInTask(&steam.Tasks.Tasks[k], tcl)
				ok = true
				break
			}
		}
		if ok {
			break
		}
	}

	_, err = collection.ReplaceOne(context.Background(), bson.M{"id": workspace.ID}, workspace)
	if err != nil {
		return errors.Wrap(err, "failed to replace workspace")
	}

	s.SendOneWorkspaceToPubSub(workspace)
	return nil
}

//UpdateTracktimeInTask function.
func UpdateTracktimeInTask(task *model.Task, tracktime model.TaskTracktime) {
	for k, v := range task.TaskTracktimes.TaskTracktimes {
		if v.ID == tracktime.ID {
			task.TaskTracktimes.TaskTracktimes[k].Finished = tracktime.Finished
			task.TaskTracktimes.TaskTracktimes[k].EndTime = tracktime.EndTime
			task.TaskTracktimes.TaskTracktimes[k].Duration = tracktime.Duration
			task.TaskTracktimes.TaskTracktimes[k].TrackType = tracktime.TrackType
			task.TaskTracktimes.TaskTracktimes[k].ChangedDate = tracktime.ChangedDate
			task.TaskTracktimes.TaskTracktimes[k].DeletedDate = tracktime.DeletedDate
			break
		}
	}
}

//UpdateTaskTracktime update existing task check list.
func (s *Server) UpdateTaskTracktime(data map[string]interface{}) error {
	var taskID string
	common.GetMapStringValue(data, "task_id", &taskID)

	collection := s.mongoClt.Database("wsTasks").Collection("workspaces")
	cur := collection.FindOne(context.Background(), bson.D{{"streamList.streams.taskList.tasks.id", bson.D{{"$lte", taskID}}}})
	if cur.Err() != nil {
		return errors.Wrap(cur.Err(), "failed to get list of notifications")
	}
	var workspace model.Workspace

	err := cur.Decode(&workspace)
	if err != nil && err != mongo.ErrNoDocuments {
		return errors.Wrap(err, "failed to decode result")
	}
	taskTracktime := model.ParseValueToTaskTracktime(data)

	ok := false
	for _, steam := range workspace.Streams.Streams {
		for k, _ := range steam.Tasks.Tasks {
			if steam.Tasks.Tasks[k].ID == taskID {
				UpdateTracktimeInTask(&steam.Tasks.Tasks[k], taskTracktime)
				ok = true
				break
			}
		}
		if ok {
			break
		}
	}

	_, err = collection.ReplaceOne(context.Background(), bson.M{"id": workspace.ID}, workspace)
	if err != nil {
		return errors.Wrap(err, "failed to replace workspace")
	}

	s.SendOneWorkspaceToPubSub(workspace)
	return nil
}
