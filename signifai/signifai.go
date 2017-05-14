package signifai

import (
	"bytes"
	"errors"

	"encoding/json"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"log"
	"net/http"

	"strings"
)

var MissingHostServiceApplication = errors.New("Your Configuration is Missing a Host, Service, or Application Field")
var MissingAPI = errors.New("Your Configuration is Missing an API field - {metrics, deployments, incidents}")
var MissingAuth = errors.New("Your Configuration is Missing a token fields")

var updateSend = "https://collectors.signifai.io/v1"

func (p Publisher) postit(list []interface{}) error {
	postList := make(map[string][]interface{})
	postList["events"] = list
	jbytes, err := json.Marshal(postList)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", updateSend+"/"+p.api, bytes.NewBuffer(jbytes))
	req.Header.Set("Authorization", "Bearer "+p.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
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
				p.addMetrics(mts[x : total-1])
			} else {
				p.addMetrics(mts[x : x+batch])
			}
			x += batch
		}

	} else {
		p.addMetrics(mts)
	}

	return nil
}

func (p Publisher) addMetrics(mts []plugin.Metric) {
	var list []interface{}

	for _, m := range mts {

		var statics []string
		var attributes = make(map[string]interface{})
		for _, element := range m.Namespace {
			if element.IsDynamic() {
				attributes[element.Name] = element.Value
			} else {
				statics = append(statics, element.Value)
			}
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
	}
}
