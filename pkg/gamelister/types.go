package gamelister

type Game struct {
	Server string            `json:"server"`
	Info   map[string]string `json:"info"`
	Status map[string]string `json:"status"`
	Ping   int64             `json:"ping"`
}
