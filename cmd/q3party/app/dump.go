package app

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/finarfin/q3party/pkg/gamelister"
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
	dumpCmd.PersistentFlags().String("output", filepath.Join(os.TempDir(), "q3servers.json"), "Output path")
	viper.BindPFlag("servers", dumpCmd.PersistentFlags().Lookup("servers"))
	viper.BindPFlag("output", dumpCmd.PersistentFlags().Lookup("output"))
}

func runDump(cmd *cobra.Command, args []string) {
	log.Trace("Dump started")

	servers := viper.GetStringSlice("servers")
	repo, err := gamelister.NewGameLister(cmd.Context(), "68", servers)
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

		if repo.IsBusy() == true {
			time.Sleep(1 * time.Second)
			continue
		}

		break
	}

	games, err := repo.List()
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	b, err := json.Marshal(games)
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	outputFile := viper.GetString("output")
	err = ioutil.WriteFile(outputFile, b, 666)
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	log.Info("Done")
}
