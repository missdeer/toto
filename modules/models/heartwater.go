package models

import (
	"encoding/json"
	"html/template"
)

type HeartwaterOp struct {
	Home  float32
	Guest float32
	Drawn float32
}

type HeartwaterDaxiao struct {
	HandicapShow string
	Da           float32
	Xiao         float32
	Handicap     string
}

type HeartwaterYazhi struct {
	HomeHandcip  float32
	HandicapShow string
	GuestHandcip float32
	Handicap     string
}

type HeartwaterRecord struct {
	VsString    string
	Turn        string
	Op          HeartwaterOp
	Gy          string
	Status      json.Number
	CatchStat   string
	HometeamID  string
	GameId      json.Number
	State       json.Number
	Daxiao      []HeartwaterDaxiao
	League      string
	Living      string
	Bgcolor     string
	Id          json.Number
	StartTime   string
	LeagueId    string
	GuestteamID string
	HomeTeam    template.HTML
	Yazhi       []HeartwaterYazhi
	GuestTeam   template.HTML
}
