package config

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/jackc/pgx"
	"github.com/prometheus/common/log"
	"gitlab.com/target-smart-data-ai-searsh/task-connector-be/types"
)

//InitSession  Initialize the database Session configuration
func InitSession(dbName, pgUser, pgPass, pgHost string, pgPort int) (*types.Session, error) {
	// Initialize the database configuration

	session := types.Session{}
	bdConfig := pgx.ConnConfig{}

	bdConfig.Database = dbName
	bdConfig.Host = pgHost
	bdConfig.Port = uint16(pgPort)
	bdConfig.User = pgUser
	bdConfig.Password = pgPass

	// - creates a db connection
	// create a regular pg connection for use by transactions
	log.Info("Creating regular connection to db")
	pgConn, err := pgx.Connect(bdConfig)
	if err != nil {
		return nil, err
	}

	session.PGConn = pgConn

	// - creates a replication connection
	replConn, err := pgx.ReplicationConnect(bdConfig)
	if err != nil {
		return nil, err
	}

	session.ReplConn = replConn

	// delete all existing slots
	err = deleteAllSlots(&session)
	if err != nil {
		log.Errorf("could not delete replication slots : %v", err)
	}

	// - creates a new replication slot
	slotName := generateSlotName()
	session.SlotName = slotName

	log.Info("Creating replication slot ", slotName)
	consistentPoint, snapshotName, err := session.ReplConn.CreateReplicationSlotEx(slotName, "wal2json")
	if err != nil {
		return nil, err
	}

	log.Infof("Created replication slot \"%s\" with consistent point LSN = %s, snapshot name = %s",
		slotName, consistentPoint, snapshotName)

	lsn, _ := pgx.ParseLSN(consistentPoint)

	session.RestartLSN = lsn
	session.SnapshotName = snapshotName
	return &session, nil
}

// generates a random slot name which can be remembered
func generateSlotName() string {
	// list of random words
	strs := []string{
		"gigantic",
		"scold",
		"greasy",
		"shaggy",
		"wasteful",
		"few",
		"face",
		"pet",
		"ablaze",
		"mundane",
	}

	rand.Seed(time.Now().Unix())

	// generate name such as delta_gigantic20
	name := fmt.Sprintf("delta_%s%d", strs[rand.Intn(len(strs))], rand.Intn(100))

	return name
}

// delete all old slots that were created by us
func deleteAllSlots(session *types.Session) error {
	rows, err := session.PGConn.Query("SELECT slot_name FROM pg_replication_slots")
	if err != nil {
		return err
	}
	for rows.Next() {
		var slotName string
		rows.Scan(&slotName)

		// only delete slots created by this program
		if !strings.Contains(slotName, "delta_") {
			continue
		}

		log.Infof("Deleting replication slot %s", slotName)
		err = session.ReplConn.DropReplicationSlot(slotName)
		//_,err = session.PGConn.Exec(fmt.Sprintf("SELECT pg_drop_replication_slot(\"%s\")", slotName))
		if err != nil {
			log.With("could not delete slot ", slotName).Error(err)
		}
	}
	return nil
}
