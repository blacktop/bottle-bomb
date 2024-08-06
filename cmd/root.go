/*
Copyright © 2024 blacktop

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

const (
	brewAPI   = "https://formulae.brew.sh/api/formula/%s.json"
	bottleAPI = "https://ghcr.io/v2/homebrew/core/%s/manifests/%s" // 1st %s is the formula name; 2nd %s is the version
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bottle-bomb",
	Short: "Download a homebrew bottle and install it",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.New(os.Stderr)

		req, err := http.NewRequest("GET", fmt.Sprintf(brewAPI, args[0]), nil)
		if err != nil {
			panic(err)
		}
		// req.Header.Add("Authorization", "Bearer QQ==")
		// req.Header.Set("Accept", "application/vnd.oci.image.index.v1+json")
		req.Header.Set("Accept", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		var formula Formula
		if err := json.Unmarshal(body, &formula); err != nil {
			panic(err)
		}
		// fmt.Println(formula)

		if len(formula.Dependencies) > 0 {
			for _, dep := range formula.Dependencies {
				logger.Warn("Dependencies", "dep", dep)
			}
		}

		if formula.Bottle.Stable.Files.Arm64Sonoma == nil {
			logger.Fatal("No bottle for macOS Sonoma (arm64)")
		}

		dlreq, err := http.NewRequest("GET", formula.Bottle.Stable.Files.Arm64Sonoma.URL, nil)
		if err != nil {
			panic(err)
		}
		dlreq.Header.Add("Authorization", "Bearer QQ==")
		resp, err = http.DefaultClient.Do(dlreq)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		logger.Info("Creating", "file", fmt.Sprintf("%s.tar.gz", args[0]))
		f, err := os.Create(args[0] + ".tar.gz")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		_, err = io.Copy(f, resp.Body)
		if err != nil {
			panic(err)
		}
		// var bottle Bottle
		// if err := json.Unmarshal(body, &bottle); err != nil {
		// 	panic(err)
		// }
		// fmt.Println(bottle)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bottle-bomb.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
