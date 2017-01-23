package evediscordkm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
)

type Channel interface {
	Process(RelevantReport)
}

type DebugChannel struct {
	ChannelConfig
}

func (d *DebugChannel) Process(kill RelevantReport) {
	victim := kill.Killmail.Victim
	attackers := make([]string, 0)
	for _, attacker := range kill.Killmail.Attackers {
		attackers = append(attackers, fmt.Sprintf("%s (%s)", attacker.Character.Name, attacker.Corporation.Name))
	}

	message := fmt.Sprintf("[%s] %s (%s) by [%s]", kill.Killmail.SolarSystem.Name, victim.Character.Name, victim.Corporation.Name, strings.Join(attackers, ","))
	if kill.Type == Win {
		color.Green("Win  :" + message)
	} else {
		color.Red("Loss :" + message)
	}

}

type RecordChannel struct {
	ChannelConfig
}

func (d *RecordChannel) Process(kill RelevantReport) {
	dir := fmt.Sprintf("%s", d.Config["path"])
	path := fmt.Sprintf("%s/%v.json", dir, kill.KillID)
	os.MkdirAll(dir, 0700)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		response := zkbRedisQResponse{Package: kill.KillReport}
		if data, _ := json.MarshalIndent(response, "", "  "); data != nil {
			ioutil.WriteFile(path, data, 0644)
		}
	}
}

func CreateChannel(config ChannelConfig) Channel {
	switch config.Type {
	case "discord":
		return &DiscordChannel{config}
	case "debug":
		return &DebugChannel{config}
	case "record":
		return &RecordChannel{config}
	}
	return nil
}
