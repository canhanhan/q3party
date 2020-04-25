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
	apiCmd.PersistentFlags().Bool("mock", false, "Create mock server")
	apiCmd.PersistentFlags().String("region", "", "AWS region")
	apiCmd.PersistentFlags().String("access-id", "", "AWS access ID")
	apiCmd.PersistentFlags().String("access-secret", "", "AWS access secret")
	apiCmd.PersistentFlags().String("bucket", "", "AWS bucket")
	viper.BindPFlag("servers", apiCmd.PersistentFlags().Lookup("servers"))
	viper.BindPFlag("mock", apiCmd.PersistentFlags().Lookup("mock"))
	viper.BindPFlag("region", apiCmd.PersistentFlags().Lookup("region"))
	viper.BindPFlag("accessId", apiCmd.PersistentFlags().Lookup("access-id"))
	viper.BindPFlag("accessSecret", apiCmd.PersistentFlags().Lookup("access-secret"))
	viper.BindPFlag("bucket", apiCmd.PersistentFlags().Lookup("bucket"))
}

func runServer(cmd *cobra.Command, args []string) {
	log.Trace("Server started")

	servers := viper.GetStringSlice("servers")
	var gameRepo apiserver.GameRepository
	if viper.GetBool("mock") {
		var err error
		gameRepo, err = apiserver.NewMockGameRepository(cmd.Context(), "testdata/games.json")
		if err != nil {
			cmd.PrintErr(err)
			return
		}
	} else {
		var err error
		gameRepo, err = apiserver.NewQ3GameRepository(cmd.Context(), "68", servers)
		if err != nil {
			cmd.PrintErr(err)
			return
		}
	}

	gs, err := apiserver.NewGameService(gameRepo)
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	var listRepo apiserver.ListRepository
	if viper.GetBool("mock") {
		var err error
		listRepo, err = apiserver.NewMockListRepository("testdata/lists.json")
		if err != nil {
			cmd.PrintErr(err)
			return
		}
	} else {
		region := viper.GetString("region")
		accessID := viper.GetString("accessId")
		accessSecret := viper.GetString("accessSecret")
		bucket := viper.GetString("bucket")

		var err error
		listRepo, err = apiserver.NewS3ListRepository(region, accessID, accessSecret, bucket)
		if err != nil {
			cmd.PrintErr(err)
			return
		}
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
