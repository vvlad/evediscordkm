package evediscordkm

import (
	"sort"
	"strconv"
)

type Constraints struct {
	Alliances    []string
	Corporations []string
	Type         string
	Systems      []string
}

type Filter struct {
	Constraints
	Channel Channel
}

type ReportType int

const (
	Unknown ReportType = iota
	Win
	Loss
)

type RelevantReport struct {
	KillReport
	RelevantVictim    *Victim
	RelevantAttackers []Attacker
	Type              ReportType
}

func (f *Filter) Process(kill KillReport) {
	if len(f.Constraints.Systems) == 0 ||
		entityIsRelevant(kill.Killmail.SolarSystem, f.Constraints.Systems) {

		report := RelevantReport{
			KillReport:        kill,
			RelevantVictim:    f.relevantVictim(kill),
			RelevantAttackers: f.relevantAttackers(kill),
			Type:              Unknown,
		}

		switch {
		case len(report.RelevantAttackers) > 0:
			report.Type = Win
		case report.RelevantVictim != nil:
			report.Type = Loss
		}

		if f.matchesTypeFilter(report) {
			go f.Channel.Process(report)
		}
	}
}

func (f *Filter) matchesTypeFilter(report RelevantReport) bool {
	return (f.Constraints.Type == "" && report.Type != Unknown) ||
		(report.Type == Win && f.Constraints.Type == "win") ||
		(report.Type == Loss && f.Constraints.Type == "loss")
}

func (f *Filter) relevantAttackers(kill KillReport) []Attacker {
	attackers := []Attacker{}
	for _, att := range kill.Killmail.Attackers {
		if entityIsRelevant(att.Corporation, f.Constraints.Corporations) ||
			entityIsRelevant(att.Alliance, f.Constraints.Alliances) {
			attackers = append(attackers, att)
		}
	}
	sort.Sort(byDamage(attackers))
	return attackers

}

func (f *Filter) relevantVictim(kill KillReport) *Victim {
	victim := kill.Killmail.Victim
	if entityIsRelevant(victim.Corporation, f.Constraints.Corporations) ||
		entityIsRelevant(victim.Alliance, f.Constraints.Alliances) {
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
