package main

import (
	"FallGuys66/live/douyu/DMconfig/config"
	"FallGuys66/live/douyu/DYtype"
	"FallGuys66/live/douyu/client"
)

func main() {
	webSocketClient := client.DyBarrageWebSocketClient{
		Config: config.SpiderConfig,
		MsgBreakers: DYtype.CodeBreakershandler{
			IsLive: false,
		},
	}
	webSocketClient.Init()
	webSocketClient.Start()
}
