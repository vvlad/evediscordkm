package evediscordkm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	zKillboardRedisQUrl = "http://redisq.zkillboard.com/listen.php"
)

type ZKillboard struct {
	Config  ZKillboardConfig
	Reports chan KillReport
}

type zkbRedisQResponse struct {
	Package KillReport `json:"package"`
}

func NewZKillboard(config Config) *ZKillboard {
	return &ZKillboard{
		Config:  config.ZKillboard,
		Reports: make(chan KillReport),
	}
}

func (z *ZKillboard) Run() {
	if z.Config.ReplayFrom != "" {
		if files, err := ioutil.ReadDir(z.Config.ReplayFrom); err == nil {
			for _, file := range files {
				data, _ := ioutil.ReadFile(fmt.Sprintf("%s/%s", z.Config.ReplayFrom, file.Name()))
				response := zkbRedisQResponse{}
				json.Unmarshal(data, &response)
				z.process(response.Package)
			}
		}
	}

	for {
		resp, err := http.Get("https://redisq.zkillboard.com/listen.php")
		if err != nil {
			fmt.Fprintln(os.Stderr, "HTTP error!", err)
			continue
		}
		if resp.StatusCode != 200 {
			fmt.Fprintln(os.Stderr, "HTTP error!", resp.StatusCode, *resp)
			continue
		}
		defer resp.Body.Close()
		data, _ := ioutil.ReadAll(resp.Body)
		response := zkbRedisQResponse{}
		json.Unmarshal(data, &response)
		z.process(response.Package)
	}
}

func (z *ZKillboard) process(report KillReport) {
	if report.KillID == 0 {
		return
	}
	filteredAttackers := []Attacker{}
	for _, att := range report.Killmail.Attackers {
		if att.Character.ID != 0 {
			filteredAttackers = append(filteredAttackers, att)
		}
	}
	if len(filteredAttackers) == 0 {
		return
	}
	report.Killmail.Attackers = filteredAttackers
	z.Reports <- report
}
