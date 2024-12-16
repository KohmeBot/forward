package main

import (
	"github.com/kohmebot/forward/forward"
	"github.com/kohmebot/plugin"
)

func NewPlugin() plugin.Plugin {
	return forward.NewPlugin()
}
