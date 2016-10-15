// Copyright (c) 2016 Pulcy.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/op/go-logging"
	api "github.com/pulcy/robin-api"
	"github.com/spf13/cobra"
)

var (
	projectName    = "robinctl"
	projectVersion = "dev"
	projectBuild   = "dev"
)

const (
	defaultAPIURL = "http://localhost:8056"
)

var (
	cmdMain = &cobra.Command{
		Use:   projectName,
		Short: "Control a Robin load-balancer",
		Long:  "Control a Robin load-balancer",
		Run:   UsageFunc,
	}
	log         = logging.MustGetLogger(cmdMain.Use)
	globalFlags struct {
		apiURL string
		quiet  bool
	}
)

func init() {
	cmdMain.PersistentFlags().StringVar(&globalFlags.apiURL, "api-url", defaultAPIURL, "URL of the Robin API")
	cmdMain.PersistentFlags().BoolVarP(&globalFlags.quiet, "quiet", "q", false, "Suppress informational output")
	logging.SetFormatter(logging.MustStringFormatter("[%{level:-5s}] %{message}"))
}

func main() {
	cmdMain.Execute()
}

func UsageFunc(cmd *cobra.Command, args []string) {
	cmd.Usage()
}

func newAPIClient() api.API {
	if globalFlags.apiURL == "" {
		Exitf("api-url cannot be empty")
	}
	apiURL, err := url.Parse(globalFlags.apiURL)
	if err != nil {
		Exitf("api-url cannot be parsed: %v", err)
	}
	result, err := api.NewClient(apiURL)
	if err != nil {
		Exitf("api-client cannot be created: %v", err)
	}
	return result
}

func Infof(format string, args ...interface{}) {
	if !globalFlags.quiet {
		if !strings.HasSuffix(format, "\n") {
			format = format + "\n"
		}
		fmt.Printf(format, args...)
	}
}

func Exitf(format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}
	fmt.Printf(format, args...)
	os.Exit(1)
}
