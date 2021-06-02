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
	"gitlabber/controller"
	"net/http"

	guuid "github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// agentCmd represents the agent command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Gitlab Bot server",
	Long:  `Main gitlab api dispatcher service`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
		return preinit(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		bindaddr := controller.InitConfig(cmd)
		router := controller.Router()

		log.Fatal("Error starting the server: ", http.ListenAndServe(bindaddr, router))
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// agentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// agentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func preinit(cmd *cobra.Command) error {

	loglevel, _ := cmd.Flags().GetString("logLevel")
	lvl, err := log.ParseLevel(loglevel)

	if err != nil {
		log.SetLevel(log.DebugLevel)
		log.Error(fmt.Errorf("invalid logging level specified %s : %s . Defaulting to 'debug' \n", loglevel, err))
	} else {
		log.SetLevel(lvl)
		log.Debug(fmt.Sprintf("Setting loglevel to '%s'\n", loglevel))
	}

	secret, _ := cmd.Flags().GetString("secret")

	if secret == "" {
		log.Debug("Creating new secret")
		cmd.Flags().Set("secret", fmt.Sprintf("%v", guuid.New()))
	}

	return nil
}
