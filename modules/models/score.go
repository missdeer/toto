package models

type Standing struct {
	Ranking  int `json:",string"`
	Name     string
	Teamid   int `json:",string"`
	Appear   int `json:",string"`
	Win      int `json:",string"`
	Draw     int `json:",string"`
	Lose     int `json:",string"`
	Goal     int `json:",string"`
	Conceded int `json:",string"`
	Point    int `json:",string"`
}

type StandingElement struct {
	Id        int `json:",string"`
	Name      string
	Standings []Standing
}

type PlayerRank struct {
	Ranking  int `json:",string"`
	Name     string
	Team     string
	Teamid   int `json:",string"`
	Goal     int `json:",string"`
	Appear   int `json:",string,omitempty"`
	Period   int `json:",string,omitempty"`
	Penalty  int `json:",string,omitempty"`
	Epicycle int `json:",string,omitempty"`
	Point    int `json:",string,omitempty"`
}

type PlayerrankElement struct {
	Id         int `json:",string"`
	Name       string
	Playerrank []PlayerRank
}

type Assistant struct {
	Ranking int `json:",string"`
	Name    string
	Team    string
	Assists int `json:",string"`
}

type AssistantElement struct {
	Id               int `json:",string"`
	Name             string
	Playerassistrank []Assistant
}

type Card struct {
	Name            string
	Team            string
	YellowCardCount int `json:",string"`
	RedCardCount    int `json:",string"`
}

type CardrankElement struct {
	Id        int `json:",string"`
	Name      string
	Standings []Card
}

type FootballScore struct {
	Standings  []StandingElement
	Playerrank []PlayerrankElement
	Assistrank []AssistantElement
	Cardrank   []CardrankElement
}
