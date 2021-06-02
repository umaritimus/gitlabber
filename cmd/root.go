/*
Copyright © 2021 Andy Dorfman <github.com/umaritimus>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultConfigFilename = ".gitlabber"
	envPrefix             = "GITLABBER"
)

var cfgFile string
var port int
var token string
var secret string
var version int
var url string
var project string
var logLevel string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gitlabber",
	Short: "Gitlab Bot and API Dispatcher",
	Long: `
Dispatcher utility to process gitlab webhooks from the gitlab and trigger corresponding pipeline actions, e.g.:

- On Merge Request change, if Approval conditions are satisfied, trigger Rebase/Merge and production deployment
- On Merge Request comment, if a pipeline trigger request is detected, rerun pipeline
- On Merge Request comment, if a status trigger is detected, send a slack status notification
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./"+defaultConfigFilename+".toml", "Path to the Configuration File [🏁]")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 443, "Listen port [🏁]")
	rootCmd.PersistentFlags().StringVarP(&secret, "secret", "", "", "Gitlabber Authentication token [🏁]")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "logLevel", "l", "info", "Gitlabber log level [🏁]")
	rootCmd.PersistentFlags().StringVarP(&url, "url", "", "https://gitlab.com/api/v4", "Gitlab api url [🏁]")
	rootCmd.PersistentFlags().StringVarP(&project, "project", "", "", "*Gitlab project [🏁]")
	rootCmd.PersistentFlags().IntVarP(&version, "version", "v", 4, "Gitlab api version [🏁]")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "", "", "*Gitlab token [🚩]")
	rootCmd.MarkPersistentFlagRequired("token")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix(envPrefix)

	// Search for a configuration file in the current directory
	viper.AddConfigPath(".")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(defaultConfigFilename)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	rootCmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --api-version => %GITLABBER_API_VERSION%
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			viper.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			rootCmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
