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

var metricsURL = "/metrics"

type Metric struct {
	EventSource string                 `json:"event_source"`
	Host        string                 `json:"host,omitempty"`
	Service     string                 `json:"service,omitempty"`
	Application string                 `json:"application,omitempty"`
	Name        string                 `json:"name"`
	Value       interface{}            `json:"value"`
	Type        string                 `json:"type,omitempty"`
	Timestamp   int64                  `json:"timestamp,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
}
