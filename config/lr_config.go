package config

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/jackc/pgx"
	"gitlab.com/target-smart-data-ai-searsh/task-connector-be/types"
)

// LRListenAck listens on the websocket for ack messages
// The commited LSN is extracted and is updated to the server
func LRListenAck(session *types.Session, wsErr chan<- error) {
	jsonMsg := make(map[string]string)
	for {
		log.Info("Listening for WS message")
		//_, msg, err := session.WSConn.ReadMessage()
		err := session.WSConn.ReadJSON(&jsonMsg)
		if err != nil {
			log.WithError(err).Error("Error reading from websocket")
			wsErr <- err // send the error to the channel to terminate connection
			return
		}
		log.Info("Received WS message: ", jsonMsg)
		lsn := jsonMsg["lsn"]
		lrAckLSN(session, lsn)
	}
}

// LRAckLSN will set the flushed LSN value and trigger a StandbyStatus update
func lrAckLSN(session *types.Session, restartLSNStr string) error {
	restartLSN, err := pgx.ParseLSN(restartLSNStr)
	if err != nil {
		return err
	}

	session.RestartLSN = restartLSN
	return SendStandbyStatus(session)
}

// SendStandbyStatus sends a StandbyStatus object with the current RestartLSN value to the server
func SendStandbyStatus(session *types.Session) error {
	standbyStatus, err := pgx.NewStandbyStatus(session.RestartLSN)
	if err != nil {
		return fmt.Errorf("unable to create StandbyStatus object: %s", err)
	}
	log.Info(standbyStatus)
	standbyStatus.ReplyRequested = 0
	log.Info("Sending Standby Status with LSN ", pgx.FormatLSN(session.RestartLSN))
	err = session.ReplConn.SendStandbyStatus(standbyStatus)
	if err != nil {
		return fmt.Errorf("unable to send StandbyStatus object: %s", err)
	}

	return nil
}

var statusHeartbeatIntervalSeconds = 10

//SendPeriodicHeartbeats send periodic keep alive hearbeats to the server so that the connection isn't dropped
func SendPeriodicHeartbeats(session *types.Session) {
	for {
		select {
		case <-session.Ctx.Done():
			// context closed; stop sending heartbeats
			return
		case <-time.Tick(time.Duration(statusHeartbeatIntervalSeconds) * time.Second):
			{
				// send hearbeat message at every statusHeartbeatIntervalSeconds interval
				log.Info("Sending periodic status heartbeat")
				err := SendStandbyStatus(session)
				if err != nil {
					log.WithError(err).Error("Failed to send status heartbeat")
				}
			}
		}
	}

}
