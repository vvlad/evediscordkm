package evediscordkm

type Constraints struct {
	Alliances    []string
	Corporations []string
	Type         string
	Systems      []string
	Characters   []string
	MinimumCost  float64 `yaml:"minimum_cost"`
}

type filterFunction func(RelevantReport) bool
type Filter struct {
	Constraints
	Channel Channel
	filters []filterFunction
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

func NewFilter(constraints Constraints, channel Channel) Filter {
	filter := Filter{
		Constraints: constraints,
		Channel:     channel,
		filters:     make([]filterFunction, 0),
	}
	filter.filters = []filterFunction{
		filter.isRelevantSystem,
		filter.isRelevantCost,
		filter.isRelevantKill,
		filter.isRelevantType,
	}
	return filter
}

func NewRelevantReport(kill KillReport, victim *Victim, attackers []Attacker) RelevantReport {

	report := RelevantReport{
		KillReport:        kill,
		RelevantVictim:    victim,
		RelevantAttackers: attackers,
		Type:              Unknown,
	}
	switch {
	case len(report.RelevantAttackers) > 0:
		report.Type = Win
	case report.RelevantVictim != nil:
		report.Type = Loss
	}
	return report
}

func (f Filter) Process(kill KillReport) {
	report := NewRelevantReport(kill, f.relevantVictim(kill), f.relevantAttackers(kill))

	for _, relevant := range f.filters {
		if !relevant(report) {
			return
		}
	}
	// FIXME: add channel pool
	go f.Channel.Process(report)
}
