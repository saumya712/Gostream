package domain

type Metrics struct {
	Totallogs    int
	Errorcount   int
	Servicesstat map[string]int
}

func Newmetric() *Metrics {
	return &Metrics{
		Servicesstat: make(map[string]int),
	}
}

func (m *Metrics) Increment(service string, level LogLevel) {
	m.Totallogs++
	m.Servicesstat[service]++

	if level == LevelError {
		m.Errorcount++
	}
}
