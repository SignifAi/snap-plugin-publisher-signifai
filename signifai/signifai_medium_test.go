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
	"encoding/json"
	plugin "github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"gopkg.in/jarcoal/httpmock.v1"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func validConfig() plugin.Config {
	config := make(plugin.Config)
	config["api"] = "metrics"
	config["token"] = "1234"
	config["host"] = "my.local.host"

	return config
}

func TestSignifAiPublisher(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", updateSend+"/metrics",
		func(req *http.Request) (*http.Response, error) {

			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				t.Fatal(err)
			}
			metricsSubmitted := map[string][]Metric{}

			err = json.Unmarshal(body, &metricsSubmitted)
			if err != nil {
				t.Fatal(err)
			}

			metrics := metricsSubmitted["events"]

			if len(metrics) != 2 {
				t.Fatal("not receiving correct number of metrics")
			}

			if metrics[0].EventSource != "Snap" {
				t.Fatal("event_source is bad %v", metrics[0].EventSource)
			}
			if metrics[0].Host != "my.local.host" {
				t.Fatalf("host is bad %v", metrics[0].Host)
			}
			if metrics[0].Name != "x.y.z" {
				t.Fatal("name is bad %v", metrics[0].Name)
			}

			if metrics[1].Name != "bar" {
				t.Fatal("name is bad %v", metrics[0].Name)
			}

			val, ok := metrics[1].Attributes["domain_name"]
			if !ok {
				t.Fatal("can't find attributes key %v", metrics[1])
			}

			estr, ok := val.(string)
			if !ok {
				t.Fatal("val is not a string")
			}

			if estr != "my.domain.name" {
				t.Fatalf("key is wrong %v", estr)
			}

			resp, err := httpmock.NewJsonResponse(200, "ok")
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	p := Publisher{}

	metrics := []plugin.Metric{
		plugin.Metric{
			Namespace: plugin.NewNamespace("x", "y", "z"),
			Config:    map[string]interface{}{"pw": "123aB"},
			Data:      3,
			Tags:      map[string]string{"hello": "world"},
			Unit:      "int",
			Timestamp: time.Now(),
		},
		plugin.Metric{
			Namespace: plugin.NewNamespace("bar").AddDynamicElement("domain_name", "Domain Name"),
			Config:    map[string]interface{}{"pw": "123aB"},
			Data:      3,
			Tags:      map[string]string{"hello": "world"},
			Unit:      "int",
			Timestamp: time.Now(),
		},
	}
	metrics[1].Namespace[1].Value = "my.domain.name"
	err := p.Publish(metrics, validConfig())
	if err != nil {
		t.Fatal(err)
	}
}

func TestOverTwentyMetrics(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cnt := 0

	httpmock.RegisterResponder("POST", updateSend+"/metrics",
		func(req *http.Request) (*http.Response, error) {

			cnt += 1

			resp, err := httpmock.NewJsonResponse(200, "ok")
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	p := Publisher{}

	metrics := []plugin.Metric{}

	for i := 0; i < 21; i++ {
		metrics = append(metrics, plugin.Metric{
			Namespace: plugin.NewNamespace("x", "y", "z"),
			Config:    map[string]interface{}{"pw": "123aB"},
			Data:      3,
			Tags:      map[string]string{"hello": "world"},
			Unit:      "int",
			Timestamp: time.Now(),
		})
	}

	err := p.Publish(metrics, validConfig())
	if err != nil {
		t.Fatal(err)
	}

	if cnt != 5 {
		t.Fatalf("server should have been sent two requests instead it sent %v", cnt)
	}

}
