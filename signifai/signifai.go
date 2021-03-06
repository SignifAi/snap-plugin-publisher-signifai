/*
http://www.apache.org/licenses/LICENSE-2.0.txt
Copyright 2017 SignifAI Inc
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package signifai

import (
	"bytes"
	"errors"

	"encoding/json"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"log"
	"net/http"

	"strings"
	"time"
)

var MissingHostServiceApplication = errors.New("Your Configuration is Missing a Host, Service, or Application Field")
var MissingAPI = errors.New("Your Configuration is Missing an API field - {metrics, deployments, incidents}")
var MissingAuth = errors.New("Your Configuration is Missing a token fields")

var updateSend = "https://collectors.signifai.io/v1"

const SIGNIFAI_AGENT_VERSION = "1.0.7"

func (p Publisher) postit(list []interface{}) error {
	if list == nil || len(list) == 0 {
		// Nothing to publish? Don't publish, it's a success
		return nil
	}

	postList := make(map[string][]interface{})
	postList["events"] = list
	jbytes, err := json.Marshal(postList)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", updateSend+"/"+p.api, bytes.NewBuffer(jbytes))
	req.Header.Set("Authorization", "Bearer "+p.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 1 * time.Second}
	var resp *http.Response
	for attempts := 0; attempts < 8; attempts++ {
		resp, err = client.Do(req)
		if err != nil {
			// XXX: There _has_ to be a better way to check timeout in Go
			if !strings.Contains(err.Error(), "Client.Timeout exceeded") {
				return err
			}
		} else {
			break
		}
	}

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

// Publisher is a publisher to the SignifAi System
type Publisher struct {
	api         string // api to use {metrics, incidents, deployment}
	token       string // JWT token to use
	host        string // host that is being collected from
	service     string // service that is being collected from
	application string // application that is being collected from
	initialized bool   // indicates that we've initialized the plugin
}

func New() *Publisher {
	return new(Publisher)
}

// GetConfigPolicy returns the configuration Policy needed for using
// this plugin
//
// we have quite a few optional parameters here
func (p *Publisher) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	policy := plugin.NewConfigPolicy()
	policy.AddNewStringRule([]string{""}, "api", true)
	policy.AddNewStringRule([]string{""}, "token", false)
	policy.AddNewStringRule([]string{""}, "host", false)
	policy.AddNewStringRule([]string{""}, "service", false)
	policy.AddNewStringRule([]string{""}, "application", false)

	return *policy, nil
}

// prob. want to refactor me
// the default Get* functions from plugin do assertations along w/nil
// chks
func (p *Publisher) setConfig(cfg plugin.Config) error {

	if p.initialized {
		return nil
	}

	// mandatory
	api, err := cfg.GetString("api")
	if err != nil {
		log.Println(err)
		return err
	}
	p.api = api

	token, err := cfg.GetString("token")
	if err != nil {
		if err != plugin.ErrConfigNotFound {
			log.Println(err)
			return err
		}
	} else {
		p.token = token
	}

	host, err := cfg.GetString("host")
	if err != nil {
		if err != plugin.ErrConfigNotFound {
			log.Println(err)
			return err
		}
	} else {
		p.host = host
	}

	service, err := cfg.GetString("service")
	if err != nil {
		if err != plugin.ErrConfigNotFound {
			log.Println(err)
			return err
		}
	} else {
		p.service = service
	}

	application, err := cfg.GetString("application")
	if err != nil {
		if err != plugin.ErrConfigNotFound {
			log.Println(err)
			return err
		}
	} else {
		p.application = application
	}

	if p.api == "" {
		return MissingAPI
	}

	if p.host == "" && p.application == "" && p.service == "" {
		return MissingHostServiceApplication
	}

	if p.token == "" {
		return MissingAuth
	}

	p.initialized = true

	return nil
}

// Publish publishes snap metrics to SignifAI API
func (p *Publisher) Publish(mts []plugin.Metric, cfg plugin.Config) error {
	err := p.setConfig(cfg)
	if err != nil {
		return err
	}

	batch := 5
	total := len(mts)

	// break it up first into batches of 5
	if total > batch {
		loops := total / batch
		if total%batch != 0 {
			loops += 1
		}

		x := 0
		for i := 0; i < loops; i++ {
			if i+1 == loops {
				err = p.addMetrics(mts[x : total])
			} else {
				err = p.addMetrics(mts[x : x+batch])
			}
			if err != nil {
				return err
			}
			x += batch
		}

	} else {
		return p.addMetrics(mts)
	}

	return nil
}

func (p Publisher) addMetrics(mts []plugin.Metric) (error) {
	var list []interface{}

	for _, m := range mts {

		var statics []string
		var attributes = make(map[string]interface{})
                attributes["agent_version"] = SIGNIFAI_AGENT_VERSION
		for _, element := range m.Namespace {
			if element.IsDynamic() {
				attributes[element.Name] = element.Value
			} else {
				statics = append(statics, element.Value)
			}
		}

		for tag, value := range m.Tags {
			attributes["tag/" + tag] = value
		}

		switch p.api {
		case "metrics":
			o := Metric{
				EventSource: "Snap",
				Name:        strings.Join(statics, "."),
				Value:       m.Data,
				Timestamp:   m.Timestamp.Unix(),
				Attributes:  attributes,
			}

			if p.host != "" {
				o.Host = p.host
			}

			if p.service != "" {
				o.Service = p.service
			}

			if p.application != "" {
				o.Application = p.application
			}
			list = append(list, o)

		default:
			log.Println("sorry - incidents && deployments are not support yet.")
		}

	}

	err := p.postit(list)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
