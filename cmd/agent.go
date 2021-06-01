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

	guuid "github.com/google/uuid"
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
		fmt.Println("gitlabber agent is listening on port", port)
		fmt.Println("use the following integration secret:\n\n", secret, "\n\n")
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

	secret, _ := cmd.Flags().GetString("secret")

	if secret == "" {
		cmd.Flags().Set("secret", fmt.Sprintf("%v", guuid.New()))
	}

	return nil
}
