package evediscordkm

type ReportProessor struct {
	Reports chan KillReport
	filters []Filter
}

func NewReportProcessor(config Config) ReportProessor {
	return ReportProessor{
		Reports: make(chan KillReport),
		filters: createFilters(config),
	}
}

func (r ReportProessor) Run() {
	for killReport := range r.Reports {
		for _, filter := range r.filters {
			filter.Process(killReport)
		}
	}
}

func createFilters(c Config) []Filter {
	filters := make([]Filter, 0)
	for _, channelConfig := range c.Channels {
		if reportChannel := CreateChannel(channelConfig); reportChannel != nil {
			filter := NewFilter(channelConfig.Constraints, reportChannel)
			filters = append(filters, filter)
		}
	}
	return filters
}
