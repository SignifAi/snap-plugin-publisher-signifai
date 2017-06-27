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
	"github.com/signifai/snap-plugin-lib-go/v1/plugin"
	"testing"
)

func TestValidConfig(t *testing.T) {
	p := Publisher{}

	config := make(plugin.Config)
	config["api"] = "metrics"
	config["token"] = "1234"
	config["host"] = "my.local.host"

	err := p.setConfig(config)
	if err != nil {
		t.Fatal(err)
	}

	if p.api != config["api"] {
		t.Fatalf("bad config, %v, %v", p.api, config["api"])
	}

	if p.token != config["token"] {
		t.Fatalf("bad config, %v, %v", p.token, config["token"])
	}

	if p.host != config["host"] {
		t.Fatalf("bad config, %v, %v", p.host, config["host"])
	}

	if !p.initialized {
		t.Fatal("bad config, %v", p.initialized, true)
	}
}

func TestBadConfig(t *testing.T) {
	p := Publisher{}

	config := make(plugin.Config)
	config["api"] = "metrics"
	config["token"] = "1234"

	err := p.setConfig(config)
	if err != MissingHostServiceApplication {
		t.Fatal("mandatory field not erroring")
	}

	config["api"] = ""
	config["service"] = "my-webapp"

	err = p.setConfig(config)
	if err != MissingAPI {
		t.Fatal("mandatory field not erroring")
	}

	config["api"] = "metrics"
	config["token"] = ""

	err = p.setConfig(config)
	if err != MissingAuth {
		t.Fatal("mandatory field not erroring")
	}

}
