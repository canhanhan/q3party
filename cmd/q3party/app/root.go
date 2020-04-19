package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

func Execute() error {
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "q3master",
	Short: "Quake 3 Master Server",
	Long:  `Quake 3 master server`,
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path")

	rootCmd.AddCommand(proxyCmd)
	//rootCmd.AddCommand(listCmd)
	//rootCmd.AddCommand(staticCmd)
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Get working directory
		wd, err := os.Getwd()
		if err != nil {
			er(err)
		}
		viper.AddConfigPath(wd)

		// Get executable directory
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			er(err)
		}
		viper.AddConfigPath(dir)
		viper.SetConfigName("q3party")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		println(err.Error())
	}
}
