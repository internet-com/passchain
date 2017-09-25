// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trusch/passchain/state"
)

// secretUpdateCmd represents the secretUpdate command
var secretUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "update a secret",
	Long:  `Update a secrets value but retain the shares`,
	Run: func(cmd *cobra.Command, args []string) {
		cli := getCli()
		key := getKey()
		sid := viper.GetString("sid")
		if len(args) > 0 {
			sid = args[0]
		}
		data := secretData
		if data == "" && len(args) > 1 {
			data = args[1]
		}
		if sid == "" || data == "" {
			log.Fatal("you must specify --sid and --data")
		}
		oldSecret, err := cli.GetSecret(sid)
		if err != nil {
			log.Fatal(err)
		}
		aesKey, err := key.DecryptString(oldSecret.Shares[cli.AccountID])
		if err != nil {
			log.Fatal(err)
		}
		s := &state.Secret{ID: sid, Value: data, Shares: oldSecret.Shares}
		err = s.EncryptWithKey(aesKey)
		if err != nil {
			log.Fatal(err)
		}
		if err := cli.UpdateSecret(s); err != nil {
			log.Fatal(err)
		}
		log.Printf("updated secret %v", sid)
	},
}

func init() {
	secretCmd.AddCommand(secretUpdateCmd)
	secretUpdateCmd.PersistentFlags().String("data", "", "secret value")
	viper.BindPFlags(secretUpdateCmd.PersistentFlags())
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// secretUpdateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// secretUpdateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
