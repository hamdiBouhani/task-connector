package types

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx"
)

type PubSubCon struct {
	TopicName      string
	ProjectID      string
	ServiceAccount string

	Topic  *pubsub.Topic
	Client *pubsub.Client
}

// Session stores the context, active db and ws connections, and replication slot state
type Session struct {
	Ctx        context.Context
	CancelFunc context.CancelFunc

	ReplConn *pgx.ReplicationConn
	PGConn   *pgx.Conn

	WSConn *websocket.Conn

	PSConn *PubSubCon

	SlotName     string
	SnapshotName string
	RestartLSN   uint64 //The pg_lsn data type can be used to store LSN (Log Sequence Number) data which is a pointer to a location in the XLOG. This type is a representation of XLogRecPtr and an internal system type of PostgreSQL.
}

// Cancel the currently running session
// Recreate replication connection
func (s *Session) ResetSession(dbName, pgUser, pgPass, pgHost string, pgPort int) error {

	bdConfig := pgx.ConnConfig{}

	bdConfig.Database = dbName
	bdConfig.Host = pgHost
	bdConfig.Port = uint16(pgPort)
	bdConfig.User = pgUser
	bdConfig.Password = pgPass

	var err error
	// cancel the currently running session
	if s.CancelFunc != nil {
		s.CancelFunc()
	}

	// close websocket connection
	if s.WSConn != nil {
		//err = session.WSConn.Close()
		if err != nil {
			return err
		}
	}

	// create new context
	ctx, cancelFunc := context.WithCancel(context.Background())
	s.Ctx = ctx
	s.CancelFunc = cancelFunc

	// create the replication connection
	if s.ReplConn != nil {
		if s.ReplConn.IsAlive() {
			// reuse the existing connection (or close it nonetheless?)
			return nil
		}
	}

	replConn, err := pgx.ReplicationConnect(bdConfig)
	if err != nil {
		return err
	}
	s.ReplConn = replConn

	return nil

}

// SnapshotDataJSON is the struct that binds with an incoming request for snapshot data
type SnapshotDataJSON struct {
	// SlotName is the name of the replication slot for which the snapshot data needs to be fetched
	// (not used as of now, will be useful in multi client setup)
	SlotName string `json:"slotName" binding:"omitempty"`

	Table   string   `json:"table" binding:"required"`
	Offset  *uint    `json:"offset" binding:"exists"`
	Limit   *uint    `json:"limit" binding:"exists"`
	OrderBy *OrderBy `json:"order_by" binding:"exists"`
}

//OrderBy is the struct
type OrderBy struct {
	Column string `json:"column" binding:"exists"`
	Order  string `json:"order" binding:"exists"`
	// Nulls TODO
}

//Wal2JSONEvent is the struct
type Wal2JSONEvent struct {
	NextLSN string `json:"nextlsn"`
	Change  []map[string]interface{}
}
