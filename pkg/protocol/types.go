package protocol

var OOBHeader = []byte{0xff, 0xff, 0xff, 0xff}
var EOTFooter = []byte{'E', 'O', 'T'}

type GameType string

const (
	FFAGameType      GameType = "ffa"
	TeamPlayGameType GameType = "team"
	TourneyGameType  GameType = "tourney"
	CTFGameType      GameType = "ctf"
)

type GetServersRequest struct {
	Protocol     string
	GameType     GameType
	IncludeFull  bool
	IncludeEmpty bool
	Demo         bool
}

type GetServersResponse struct {
	Servers []string
}

type GetInfoRequest struct {
	Challenge string
}

type GetInfoResponse struct {
	Data map[string]string
}

type GetStatusRequest struct {
	Challenge string
}

type GetStatusResponse struct {
	Data map[string]string
}

type GetChallengeRequest struct{}

type GetChallengeResponse struct {
	Challenge string
}
