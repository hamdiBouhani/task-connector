package extractor

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"cloud.google.com/go/pubsub"
	"gitlab.com/target-smart-data-ai-searsh/task-connector-be/model"
)

//SendToPubSub send data to pub-sub.
func (s *Server) SendToPubSub(data model.Workspaces) {

	if data.TotalCount > 0 {
		for _, ws := range data.Workspaces {
			reqBodyBytes := new(bytes.Buffer)
			json.NewEncoder(reqBodyBytes).Encode(ws)
			msg := reqBodyBytes.Bytes()

			//r.Logger.Printf("INDEX :%d sent: %s\n", index, string(msg))

			message := &pubsub.Message{
				Data: msg,
			}

			ctx := context.Background()

			res := s.Session.PSConn.Topic.Publish(ctx, message)

			msgID, err := res.Get(ctx)
			if err != nil {
				log.Fatal(err)
			}
			s.Logger.Printf("Message is stored in topic(%s) msgId :%s\n", s.Session.PSConn.TopicName, msgID)

		}
	}

}

//SendToPubSub send data to pub-sub.
func (s *Server) SendOneWorkspaceToPubSub(ws model.Workspace) {

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(ws)
	msg := reqBodyBytes.Bytes()

	//r.Logger.Printf("INDEX :%d sent: %s\n", index, string(msg))

	message := &pubsub.Message{
		Data: msg,
	}

	ctx := context.Background()

	res := s.Session.PSConn.Topic.Publish(ctx, message)

	msgID, err := res.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
	s.Logger.Printf("Message is stored in topic(%s) msgId :%s\n", s.Session.PSConn.TopicName, msgID)

}
