package demos

type Booking struct {
	Demos   []Demo
	Players []string
}

type Demo struct {
	Name  string
	Map   string
	URL   string
	Teams TeamNames
}

type TeamNames struct {
	RedTeam string
	BluTeam string
}
