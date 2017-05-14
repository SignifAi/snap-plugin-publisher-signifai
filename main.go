package main

import (
	"github.com/SignifAi/snap-plugin-publisher-signifai/signifai"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

const (
	pluginName    = "signifai-publisher"
	pluginVersion = 1
)

func main() {
	plugin.StartPublisher(signifai.New(), pluginName, pluginVersion)
}
