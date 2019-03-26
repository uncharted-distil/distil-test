//
//   Copyright © 2019 Uncharted Software Inc.
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/unchartedsoftware/plog"
	"github.com/urfave/cli"

	"github.com/uncharted-distil/distil-test/env"
)

type WSMessage struct {
	ID         string `json:"id"`
	ResultID   string `json:"resultId"`
	RequestID  string `json:"requestId"`
	SolutionID string `json:"solutionId"`
	Timestamp  string `json:"timestamp"`
	Progress   string `json:"progress"`
	Error      string `json:"error"`
	Complete   bool   `json:"complete"`
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	config, err := env.LoadConfig()
	if err != nil {
		log.Errorf("%+v", err)
		os.Exit(1)
	}

	app := cli.NewApp()
	app.Name = "distil-test"
	app.Version = "0.1.0"
	app.Usage = "Test distil"
	app.UsageText = "distil-test"
	app.Flags = []cli.Flag{}
	app.Action = func(c *cli.Context) error {

		// wait fot Distil
		err = waitForDistil(config.RetryCount, fmt.Sprintf("http://%s/distil/config", config.AppEndpoint))
		if err != nil {
			log.Errorf("%v", err)
			return cli.NewExitError(errors.Cause(err), 2)
		}

		// initialize client
		u := url.URL{Scheme: "ws", Host: config.AppEndpoint, Path: "/ws"}
		log.Infof("connecting to %s", u.String())

		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Errorf("%v", err)
			return cli.NewExitError(errors.Cause(err), 2)
		}
		defer conn.Close()

		log.Infof("Using app interface at `%s` ", config.AppEndpoint)
		err = conn.WriteMessage(websocket.TextMessage, getRequest(config.Dataset))
		if err != nil {
			log.Errorf("%v", err)
			return cli.NewExitError(errors.Cause(err), 2)
		}

		success := isSuccess(conn)

		if !success {
			log.Errorf("no successful pipelines produced by Distil")
			return cli.NewExitError(errors.Cause(errors.Errorf("no successful pipelines produced by Distil")), 2)
		}
		log.Infof("at least one pipeline was successfully created")

		return nil
	}
	// run app
	app.Run(os.Args)
}

func isSuccess(conn *websocket.Conn) bool {
	// success is defined as one pipeline running to completion without error
	log.Infof("Waiting for messages...")
	success := false
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Errorf("%v", err)
			break
		}
		log.Infof("%s", message)

		var msg WSMessage
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Errorf("%v", err)
			break
		}

		if msg.Complete {
			break
		}

		if msg.Error != "" {
			log.Errorf("error received from Distil")
		}

		if msg.Progress == "REQUEST_COMPLETED" && msg.Error == "" {
			success = true
			//break
		}
	}
	log.Infof("Done reading messages")

	return success
}

func waitForDistil(maxRetries int, url string) error {
	// can be determined by hitting the config endpoint
	log.Infof("waiting for distil to be up...")
	for i := 0; i < maxRetries; i++ {
		// will error if not available
		_, err := http.Get(url)
		if err == nil {
			return nil
		}
		log.Infof("distil not up yet (attempt %d)", i+1)
		time.Sleep(5 * time.Second)
	}

	return errors.Errorf("no response after %d retries", maxRetries)
}

func getRequest(dataset string) []byte {
	return []byte(fmt.Sprintf(`
        {
            "type": "CREATE_SOLUTIONS",
            "id": "0",
            "dataset": "%s",
            "filters": {
                "filters": [],
                "variables": ["cylinders", "displacement", "acceleration"]
            },
            "maxSolutions": 3,
            "maxTime": 1000,
            "metrics": null,
            "subTask": "univariate",
            "target": "class",
            "task": "regression"
        }`, dataset))
}
