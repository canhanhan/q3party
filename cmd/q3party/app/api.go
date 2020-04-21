package app

import (
	"github.com/finarfin/q3party/pkg/apiserver"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var apiCmd = &cobra.Command{
	Use:   "server",
	Short: "Generate a Q3 Party server",
	Run:   runServer,
}

func init() {
	apiCmd.PersistentFlags().StringSlice("servers", nil, "Comma seperated list of master servers")
	viper.BindPFlag("servers", apiCmd.PersistentFlags().Lookup("servers"))
}

func runServer(cmd *cobra.Command, args []string) {
	log.Trace("Server started")

	servers := viper.GetStringSlice("servers")
	gameRepo, err := apiserver.NewQ3GameRepository(cmd.Context(), "68", servers)
	if err != nil {
		cmd.PrintErr(err)
		return
	}
	gs, err := apiserver.NewGameService(gameRepo)
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	listRepo, err := apiserver.NewMockListRepository("testdata/lists.json")
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	ls, err := apiserver.NewListService(listRepo)
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	s, err := apiserver.NewApiServer("127.0.0.1:8080", ls, gs)
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	log.Error(s.Listen())
}
