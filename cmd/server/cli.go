package main

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "A http redirecting server",
	Long: `redirect is a pure http 'redirecting server', i.e. all requests will receive a http redirect status. 
redirect uses hostname and request URL to determine the redirect target. 
It can be used as a backend for URL shorteners or domain forwarders.

In the default settings, all redirects will be saved to a file 'redirects.json', the server will listen on port 8080 and provide no admin API.runServer
To run in this setup, an initial 'redirects.json' is required. To generate one, the server can be started with 'server -f' and stopped again, this will produce an empty save file.`,
	Version: Build,
	Run: func(cmd *cobra.Command, args []string) {
		config.adminAddress = viper.GetString("api")
		config.debug = viper.GetBool("debug")
		config.listenAddress = viper.GetString("listen")
		config.redirectFile = viper.GetString("storage")
		config.redirectFileIgnoreErr = viper.GetBool("force")
		config.redirectNoSave = viper.GetBool("volatile")

		runServer()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.server.yaml)")
	rootCmd.PersistentFlags().StringVarP(&config.listenAddress, "listen", "l", ":8080", "Sets listen address (ip:port) for redirector; empty ip for all interfaces")
	rootCmd.PersistentFlags().StringVarP(&config.adminAddress, "api", "a", "", "Enable HTTP API on a specific hostname (listen address has to cover this hostname)")
	rootCmd.PersistentFlags().StringVarP(&config.redirectFile, "storage", "s", "redirects.json", "Save file for the redirector (loaded at start of server, saved at closing of server)")
	rootCmd.PersistentFlags().BoolVarP(&config.redirectFileIgnoreErr, "force", "f", false, "Ignore load errors when opening redirector save file (starts with empty redirector), this can be useful for first setup of server")
	rootCmd.PersistentFlags().BoolVar(&config.redirectNoSave, "volatile", false, "Do not save redirects when closing server")
	rootCmd.PersistentFlags().BoolVar(&config.debug, "debug", false, "Enable debut output")

	viper.BindPFlags(rootCmd.PersistentFlags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".server" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".server")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
