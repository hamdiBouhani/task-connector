package extractor

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/machinebox/graphql"
	"github.com/pkg/errors"
	"gitlab.com/target-smart-data-ai-searsh/task-connector-be/model"
)

//TokenResp struct.
type TokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

//Snapshot create snapdhot of database
func (s *Server) Snapshot(c *gin.Context) {

	s.Logger.Info(fmt.Sprintf("\n\n\n\n\n Task Job Running at %v . \n\n\n\n\n", time.Now()))

	var tokenResp TokenResp

	payload := strings.NewReader("grant_type=client_credentials&client_id=" + s.clientID + "&client_secret=" + s.clientSecret + "&scope=openid email groups profile offline_access")

	reqToken, _ := http.NewRequest("POST", s.oidcAddr, payload)
	reqToken.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(reqToken)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(body, &tokenResp)

	client := graphql.NewClient(s.WsAddr)

	// make a request
	req := graphql.NewRequest(ALLWorkspaceQuery)
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)

	// run it and capture the response
	var respData model.AllWorkspaces
	if err := client.Run(context.Background(), req, &respData); err != nil {
		log.Fatal(err)
	}

	file, _ := json.MarshalIndent(respData, "", " ")
	_ = ioutil.WriteFile("ws.json", file, 0644)

	//Get All users data
	UsersReq := graphql.NewRequest(ALLUsersQuery)
	UsersReq.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	var usersData model.AllUsers
	if err := client.Run(context.Background(), UsersReq, &usersData); err != nil {
		log.Fatal(err)
	}

	// usersFile, _ := json.MarshalIndent(usersData, "", " ")
	// _ = ioutil.WriteFile("users.json", usersFile, 0644)

	go s.SendToPubSub(respData.Data)
	go s.SaveAllWorkspaces(respData.Data)
	go s.SaveAllUsers(usersData.Data)
}

//SavaAllWorkspaces save data into db.
func (s *Server) SaveAllWorkspaces(data model.Workspaces) error {

	if data.TotalCount > 0 {
		for _, ws := range data.Workspaces {
			// Get a handle for your collection
			collection := s.mongoClt.Database("wsTasks").Collection("workspaces")

			_, err := collection.InsertOne(context.Background(), ws)
			if err != nil {
				return errors.Wrap(err, "failed to insert a new user")
			}
		}
	}
	return nil
}

//SaveAllUsers save data into db.
func (s *Server) SaveAllUsers(data model.Users) error {

	if data.TotalCount > 0 {
		for _, user := range data.Users {
			// Get a handle for your collection
			collection := s.mongoClt.Database("wsTasks").Collection("users")

			_, err := collection.InsertOne(context.Background(), user)
			if err != nil {
				return errors.Wrap(err, "failed to insert a new user")
			}
		}
	}
	return nil
}
