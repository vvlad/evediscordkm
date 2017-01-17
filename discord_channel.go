package evediscordkm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dustin/go-humanize"
)

type DiscordChannel struct {
	ChannelConfig
}

type discordWebhookMessage struct {
	Content string `json:"content"`
}

func (d *DiscordChannel) Process(kill RelevantReport) {
	var message string

	if kill.Type == Win {
		message = d.formatKillMessage(kill)
	} else {
		message = d.formatLossMessage(kill)
	}

	d.sendToDiscord(message)

}

func (d *DiscordChannel) formatKillMessage(kill RelevantReport) string {
	attackers := kill.RelevantAttackers
	otherAttackersCount := len(kill.Killmail.Attackers) - len(attackers)
	extraAttackersCount := len(attackers) - 3
	if extraAttackersCount < 0 {
		extraAttackersCount = 0
	}

	attackerNames := []string{"**" + attackers[0].Character.Name + "**"}
	if len(attackers) > 1 {
		attackerNames = append(attackerNames, attackers[1].Character.Name)
	}
	if len(attackers) > 2 {
		attackerNames = append(attackerNames, attackers[2].Character.Name)
	}
	if otherAttackersCount+extraAttackersCount > 0 {
		attackerNames = append(attackerNames, fmt.Sprintf("%d other(s)", otherAttackersCount+extraAttackersCount))
	}

	attackerSection := attackerNames[0] + " (solo)"
	if len(attackerNames) > 1 {
		commaAttackers := strings.Join(attackerNames[:len(attackerNames)-1], ", ")
		attackerSection = commaAttackers + " and " + attackerNames[len(attackerNames)-1]
	}

	return fmt.Sprintf(
		"%s killed **%s** (*%s*; *%s ISK*) in **%s** -- https://zkillboard.com/kill/%d",
		attackerSection,
		kill.Killmail.Victim.Character.Name,
		kill.Killmail.Victim.ShipType.Name,
		humanize.Commaf(kill.Metadata.TotalValue),
		kill.Killmail.SolarSystem.Name,
		kill.KillID,
	)
}

func (d *DiscordChannel) formatLossMessage(kill RelevantReport) string {
	countNPC := kill.Killmail.AttackerCount - len(kill.Killmail.Attackers)
	NPCSection := ""
	if countNPC > 0 {
		NPCSection = fmt.Sprintf("and %d NPC(s) ", countNPC)
	}
	return fmt.Sprintf(
		"**%s** was killed (*%s*; *%s* ISK) by %d attacker(s) %sin **%s** -- https://zkillboard.com/kill/%d",
		kill.Killmail.Victim.Character.Name,
		kill.Killmail.Victim.ShipType.Name,
		humanize.Commaf(kill.Metadata.TotalValue),
		len(kill.Killmail.Attackers),
		NPCSection,
		kill.Killmail.SolarSystem.Name,
		kill.KillID,
	)
}

func (d *DiscordChannel) sendToDiscord(message string) {
	if d.Config["webhook"] == nil {
		return
	}
	msg := discordWebhookMessage{message}
	data, _ := json.Marshal(msg)
	http.Post(d.Config["webhook"].(string), "application/json", bytes.NewBuffer(data))
}
