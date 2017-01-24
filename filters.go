package evediscordkm

import (
	"sort"
	"strconv"
)

func (f Filter) isRelevantKill(report RelevantReport) bool {
	return report.RelevantVictim != nil || len(report.RelevantAttackers) > 0
}

func (f Filter) isRelevantSystem(report RelevantReport) bool {
	return len(f.Constraints.Systems) == 0 || entityIsRelevant(report.Killmail.SolarSystem, f.Constraints.Systems)
}

func (f Filter) isRelevantType(report RelevantReport) bool {
	return f.Constraints.Type == "" ||
		(report.Type == Win && f.Constraints.Type == "win") ||
		(report.Type == Loss && f.Constraints.Type == "loss")
}

func (f Filter) isRelevantCost(report RelevantReport) bool {
	return report.Metadata.TotalValue >= f.Constraints.MinimumCost
}

func (f Filter) relevantAttackers(kill KillReport) []Attacker {
	attackers := []Attacker{}
	for _, att := range kill.Killmail.Attackers {
		if entityIsRelevant(att.Corporation, f.Constraints.Corporations) ||
			entityIsRelevant(att.Alliance, f.Constraints.Alliances) ||
			entityIsRelevant(att.Character, f.Constraints.Characters) {
			attackers = append(attackers, att)
		}
	}
	sort.Sort(byDamage(attackers))
	return attackers
}

func (f Filter) relevantVictim(kill KillReport) *Victim {
	victim := kill.Killmail.Victim
	if entityIsRelevant(victim.Corporation, f.Constraints.Corporations) ||
		entityIsRelevant(victim.Alliance, f.Constraints.Alliances) ||
		entityIsRelevant(victim.Character, f.Constraints.Characters) {
		return &victim
	}
	return nil

}

func entityIsRelevant(entity idNameEntity, tokens []string) bool {
	for _, token := range tokens {
		if token == strconv.Itoa(entity.ID) || token == entity.Name {
			return true
		}
	}
	return false
}
