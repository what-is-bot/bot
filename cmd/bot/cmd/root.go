package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/what-is-bot/bot/internal/bot"
)

var (
	cfgFile string
	cmd     = &cobra.Command{
		Use:   "bot",
		Short: "What is bot is your friendly onboarding helper for Slack",
		Run: func(cmd *cobra.Command, args []string) {
			bot.Start()
		},
	}
)

func Execute() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "configuration file (default is $HOME/.what-is-bot)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
		}

		viper.AddConfigPath(home)
		viper.SetConfigType("json")
		viper.SetConfigName(".what-is-bot")
	}

	viper.AutomaticEnv()
	viper.Set("Verbose", true)

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println(err)
	}
}
