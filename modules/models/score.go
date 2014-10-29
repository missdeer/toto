package models

type Standing struct {
	Ranking  string
	Name     string
	Teamid   string
	Appear   string
	Win      string
	Draw     string
	Lose     string
	Goal     string
	Conceded string
	Point    string
}

type StandingElement struct {
	Id        string
	Name      string
	Standings []Standing
}

type PlayerRank struct {
	Ranking  string
	Name     string
	Team     string
	Teamid   string
	Goal     string
	Appear   string
	Period   string
	Penalty  string
	Epicycle string
	Point    string
}

type PlayerrankElement struct {
	Id         string
	Name       string
	Playerrank []PlayerRank
}

type Assistant struct {
	Ranking string
	Name    string
	Team    string
	Assists string
}

type AssistantElement struct {
	Id               string
	Name             string
	Playerassistrank []Assistant
}

type Card struct {
	Name            string
	Team            string
	YellowCardCount string
	RedCardCount    string
}

type CardrankElement struct {
	Id        string
	Name      string
	Standings []Card
}

type FootballScore struct {
	Standings  []StandingElement
	Playerrank []PlayerrankElement
	Assistrank []AssistantElement
	Cardrank   []CardrankElement
}
