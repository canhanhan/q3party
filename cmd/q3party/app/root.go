package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/finarfin/q3party/internal/logging"
	log "github.com/sirupsen/logrus"
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
	rootCmd.AddCommand(apiCmd)
	rootCmd.AddCommand(dumpCmd)
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

	initLogger()
}

func initLogger() {
	log.SetOutput(os.Stdout)
	if viper.IsSet("logfile") {
		file, err := os.OpenFile(viper.GetString("logfile"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			er(err)
		}

		log.SetOutput(file)
	}

	log.SetLevel(log.InfoLevel)
	if viper.IsSet("loglevel") {
		level, err := log.ParseLevel(viper.GetString("loglevel"))
		if err != nil {
			er(err)
		}

		log.SetLevel(level)
	}

	log.AddHook(&logging.SplitHook{})
}
