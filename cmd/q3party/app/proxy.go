package app

import (
	"github.com/finarfin/q3party/pkg/masterclient"
	"github.com/finarfin/q3party/pkg/masterserver"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Starts a master server to proxy upstream servers",
	Run:   runProxy,
}

func init() {
	proxyCmd.PersistentFlags().StringP("bind", "b", "127.0.0.1:27950", "Bind address")
	proxyCmd.PersistentFlags().StringSlice("servers", nil, "Comma seperated list of master servers")
	viper.BindPFlag("bind", proxyCmd.PersistentFlags().Lookup("bind"))
	viper.BindPFlag("servers", proxyCmd.PersistentFlags().Lookup("servers"))
}

func runProxy(cmd *cobra.Command, args []string) {
	log.Trace("Starting proxy command")

	bind := viper.GetString("bind")
	servers := viper.GetStringSlice("servers")
	mc, err := masterclient.NewMasterClient(cmd.Context(), servers)
	if err != nil {
		log.Error(err)
		cmd.PrintErr(err)
		return
	}
	defer mc.Close()

	ms, err := masterserver.NewMasterServer(cmd.Context(), bind, mc)
	if err != nil {
		log.Error(err)
		cmd.PrintErr(err)
		return
	}
	defer ms.Close()

	<-cmd.Context().Done()
}
