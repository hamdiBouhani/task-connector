package extractor

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"gitlab.com/target-smart-data-ai-searsh/task-connector-be/config"
	"gitlab.com/target-smart-data-ai-searsh/task-connector-be/types"
)

func (s *Server) StreamChanges(c *gin.Context) {

	s.Session.ResetSession(s.DbName, s.PgUser, s.PgPass, s.PgHost, s.PgPort)
	go func() {
		wsErr := make(chan error, 1)

		s.Logger.Infof("Starting replication for slot '%s' from LSN %s", s.Session.SlotName, pgx.FormatLSN(s.Session.RestartLSN))
		err := s.Session.ReplConn.StartReplication(s.Session.SlotName, s.Session.RestartLSN, -1, "\"include-lsn\" 'on'", "\"pretty-print\" 'off'")
		if err != nil {
			e := fmt.Sprintf("Could not Start replication for slot '%s' from LSN %s", s.Session.SlotName, pgx.FormatLSN(s.Session.RestartLSN))
			s.Logger.WithError(err).Error(e)

		}

		// start sending periodic status heartbeats to postgres
		go config.SendPeriodicHeartbeats(s.Session)

		for {

			if !s.Session.ReplConn.IsAlive() {
				e := fmt.Sprintf("Looks like the connection is dead")
				s.Logger.WithError(s.Session.ReplConn.CauseOfDeath()).Error(e)
			}
			s.Logger.Info("Waiting for LR message")

			ctx := s.Session.Ctx
			message, err := s.Session.ReplConn.WaitForReplicationMessage(ctx)
			if err != nil {
				// check whether the error is because of the context being cancelled
				if ctx.Err() != nil {
					// context cancelled, exit
					s.Logger.Warn("Websocket closed")
					return
				}

				s.Logger.WithError(err).Errorf("%s", reflect.TypeOf(err))
			}

			if message.WalMessage != nil {
				if message == nil {
					s.Logger.Error("Message nil")
					continue
				}

				walData := message.WalMessage.WalData
				s.Logger.Infof("Received replication message: %s", string(walData))

				var wData types.WalData
				if err := json.Unmarshal(walData, &wData); err != nil {
					e := fmt.Sprintf("faild to decode waldata:%s", string(walData))
					s.Logger.WithError(s.Session.ReplConn.CauseOfDeath()).Error(e)
				}

				for _, value := range wData.Change {
					if value.Kind == "delete" {
						continue
					}

					if value.Kind == "insert" {
						switch value.Table {
						case "workspace":
							err := s.InsertWorkspacs(value.GetValue())
							if err != nil {
								e := fmt.Sprintf("Could not insert workspace data")
								s.Logger.WithError(err).Error(e)
							}
						case "stream":
							err := s.InsertStream(value.GetValue())
							if err != nil {
								e := fmt.Sprintf("Could not insert workspace data")
								s.Logger.WithError(err).Error(e)
							}
						case "event":
							err := s.InsertEvent(value.GetValue())
							if err != nil {
								e := fmt.Sprintf("Could not insert workspace data")
								s.Logger.WithError(err).Error(e)
							}
						case "task":
							err := s.InsertTask(value.GetValue())
							if err != nil {
								e := fmt.Sprintf("Could not insert workspace data")
								s.Logger.WithError(err).Error(e)
							}
						case "task_assignee_map":
							err := s.InsertTaskAssigneeMap(value.GetValue())
							if err != nil {
								e := fmt.Sprintf("Could not insert workspace data")
								s.Logger.WithError(err).Error(e)
							}
						case "task_checklist":
							err := s.InsertTaskChecklist(value.GetValue())
							if err != nil {
								e := fmt.Sprintf("Could not insert workspace data")
								s.Logger.WithError(err).Error(e)
							}
						case "task_tracktime":
							err := s.InsertTaskTracktime(value.GetValue())
							if err != nil {
								e := fmt.Sprintf("Could not insert workspace data")
								s.Logger.WithError(err).Error(e)
							}
						case "user":
							err := s.InsertUser(value.GetValue())
							if err != nil {
								e := fmt.Sprintf("Could not insert workspace data")
								s.Logger.WithError(err).Error(e)
							}
						}
					}

					if value.Kind == "update" {
						switch value.Table {
						case "workspace":
							err := s.UpdateWorkspacs(value.GetValue())
							if err != nil {
								e := fmt.Sprintf("Could not Update workspace data")
								s.Logger.WithError(err).Error(e)
							}
						case "stream":
							err := s.UpdateStream(value.GetValue())
							if err != nil {
								e := fmt.Sprintf("Could not Update workspace data")
								s.Logger.WithError(err).Error(e)
							}
						case "event":
							err := s.UpdateEvent(value.GetValue())
							if err != nil {
								e := fmt.Sprintf("Could not Update workspace data")
								s.Logger.WithError(err).Error(e)
							}
						case "task":
							err := s.UpdateTask(value.GetValue())
							if err != nil {
								e := fmt.Sprintf("Could not Update workspace data")
								s.Logger.WithError(err).Error(e)
							}
						case "task_assignee_map":
							err := s.UpdateTaskAssigneeMap(value.GetValue())
							if err != nil {
								e := fmt.Sprintf("Could not update workspace data")
								s.Logger.WithError(err).Error(e)
							}
						case "task_checklist":
							err := s.UpdateTaskChecklist(value.GetValue())
							if err != nil {
								e := fmt.Sprintf("Could not update workspace data")
								s.Logger.WithError(err).Error(e)
							}
						case "task_tracktime":
							err := s.UpdateTaskTracktime(value.GetValue())
							if err != nil {
								e := fmt.Sprintf("Could not update workspace data")
								s.Logger.WithError(err).Error(e)
							}

						}
					}

				}

			}

		}

		select {
		case <-wsErr: // ws closed
			s.Logger.Warn("Cancelling context.")
			// cancel session context
			s.Session.CancelFunc()

			err = s.Session.ReplConn.Close()
			if err != nil {
				s.Logger.WithError(err).Error("Could not close replication connection")
			}

		}

	}()

}
