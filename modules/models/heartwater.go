package models

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
	Status      string
	catchStat   string
	HometeamID  string
	GameId      uint
	State       string
	Daxiao      []HeartwaterDaxiao
	League      string
	Living      string
	Bgcolor     string
	Id          uint
	StartTime   string
	LeagueId    string
	GuestteamID string
	HomeTeam    string
	Yazhi       []HeartwaterYazhi
	GuestTeam   string
}
