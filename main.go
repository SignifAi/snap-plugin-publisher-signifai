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

package main

import (
	"github.com/signifai/snap-plugin-publisher-signifai/signifai"
	"github.com/signifai/snap-plugin-lib-go/v1/plugin"
        "google.golang.org/grpc"
)

const (
	pluginName    = "signifai-publisher"
	pluginVersion = 1
)

func main() {
	plugin.StartPublisher(signifai.New(), pluginName, pluginVersion, plugin.GRPCServerOptions(grpc.MaxMsgSize(20 * 1024 * 1024)))
}
