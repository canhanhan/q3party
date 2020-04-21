package app

import (
	"encoding/json"
	"io/ioutil"

	"github.com/finarfin/q3party/pkg/apiserver"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Generate a Q3 Party server",
	Run:   runDump,
}

func init() {
	dumpCmd.PersistentFlags().StringSlice("servers", nil, "Comma seperated list of master servers")
	viper.BindPFlag("servers", dumpCmd.PersistentFlags().Lookup("servers"))
}

func runDump(cmd *cobra.Command, args []string) {
	log.Trace("Dump started")

	servers := viper.GetStringSlice("servers")
	repo, err := apiserver.NewQ3GameRepository(cmd.Context(), "68", servers)
	if err != nil {
		cmd.PrintErr(err)
		return
	}
	defer repo.Close()

	for {
		select {
		case <-cmd.Context().Done():
			return
		default:
		}

		games, err := repo.List()
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		if len(games) > 0 {
			b, err := json.Marshal(games)
			if err != nil {
				cmd.PrintErr(err)
				return
			}

			err = ioutil.WriteFile("C:\\temp\\q3.json", b, 666)
			if err != nil {
				cmd.PrintErr(err)
				return
			}

			log.Info("Done")
			return
		}
	}
}
