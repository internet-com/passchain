/*
 * Copyright (C) 2017 Tino Rusch
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// secretUpdateCmd represents the secretUpdate command
var secretUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "update a secret",
	Long:  `Update a secrets value but retain the shares`,
	Run: func(cmd *cobra.Command, args []string) {
		cli := getCli()
		key := cli.Key
		sid := viper.GetString("sid")
		if sid == "" && len(args) > 0 {
			sid = args[0]
		}
		data := secretData
		if data == "" && len(args) > 1 {
			data = args[1]
		}
		if sid == "" || data == "" {
			log.Fatal("you must specify --sid and --data")
		}
		sec, err := cli.GetSecret(sid)
		if err != nil {
			log.Fatalf("failed to get secret: %v", err)
		}
		aesKey, err := key.DecryptString(sec.Shares[cli.AccountID])
		if err != nil {
			log.Fatalf("failed to decrypt secret: %v", err)
		}
		sec.Value = data
		err = sec.EncryptWithKey(aesKey)
		if err != nil {
			log.Fatal(err)
		}
		if err := cli.UpdateSecret(sec); err != nil {
			log.Fatal(err)
		}
		log.Printf("updated secret %v", sid)
	},
}

func init() {
	secretCmd.AddCommand(secretUpdateCmd)
	secretUpdateCmd.PersistentFlags().String("data", "", "secret value")
	viper.BindPFlags(secretUpdateCmd.PersistentFlags())
}
