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
	"encoding/base64"
	"encoding/json"
	"os"
	"time"

	terminate "github.com/pulcy/go-terminate"
	api "github.com/pulcy/robin-api"
	"github.com/spf13/cobra"
)

var (
	cmdAdd = &cobra.Command{
		Use:   "add",
		Short: "Add frontends",
		Run:   cmdAddRun,
	}
	addFlags struct {
		ID       string
		JSON     string
		Wait     bool
		Frontend api.FrontendRecord
		Selector api.FrontendSelectorRecord
	}
)

func init() {
	cmdAdd.Flags().StringVar(&addFlags.ID, "id", "", "ID of the new frontend")
	cmdAdd.Flags().StringVar(&addFlags.JSON, "json", "", "Frontend formatted at json (can be base64 encoded)")
	cmdAdd.Flags().BoolVar(&addFlags.Wait, "wait", false, "If set, the process will add the frontend and then wait until termination. Just before termination the frontend will be removed.")

	// Frontend flags
	cmdAdd.Flags().StringVar(&addFlags.Frontend.Service, "service", "", "Service that is providing the frontend")
	cmdAdd.Flags().StringVar(&addFlags.Frontend.Mode, "mode", "", "Mode of the frontend (http|tcp)")
	cmdAdd.Flags().StringVar(&addFlags.Frontend.HttpCheckPath, "http-check-path", "", "Path for HTTP health checks")
	cmdAdd.Flags().StringVar(&addFlags.Frontend.HttpCheckMethod, "http-check-method", "", "HTTP method for health checks")
	cmdAdd.Flags().BoolVar(&addFlags.Frontend.Sticky, "sticky", false, "If set, requests will be send to same server (when possible)")
	cmdAdd.Flags().BoolVar(&addFlags.Frontend.Backup, "backup", false, "If set, requests will only be send to this frontend when other frontends with same selectors are down")

	// Frontend selector flags
	cmdAdd.Flags().IntVar(&addFlags.Selector.Weight, "weight", 0, "A value between 0-100 used for prioritizing frontends (100 most important)")
	cmdAdd.Flags().StringVar(&addFlags.Selector.Domain, "domain", "", "The domain to select for this frontend")
	cmdAdd.Flags().StringVar(&addFlags.Selector.PathPrefix, "path-prefix", "", "The path-prefix to select for this frontend")
	cmdAdd.Flags().IntVar(&addFlags.Selector.ServicePort, "port", 0, "The port on the service to forward requests for this frontend to")
	cmdAdd.Flags().IntVar(&addFlags.Selector.FrontendPort, "frontend-port", 0, "The port on lb host for this frontend to listen on")
	cmdAdd.Flags().BoolVar(&addFlags.Selector.Private, "private", false, "If set, this will be a cluster-local frontend")
	cmdAdd.Flags().StringVar(&addFlags.Selector.SslCert, "ssl-cert", "", "Name of the SSL certificate file to use for this frontend")

	cmdMain.AddCommand(cmdAdd)
}

func cmdAddRun(cmd *cobra.Command, args []string) {
	c := newAPIClient()
	id := addFlags.ID
	if id == "" && len(args) == 1 {
		id = args[0]
	}
	if id == "" {
		id = addFlags.Frontend.Service
	}

	var f api.FrontendRecord
	if addFlags.JSON != "" {
		// Use json
		data := addFlags.JSON
		// Try to decode from base64 first
		raw, err := base64.StdEncoding.DecodeString(data)
		if err == nil {
			data = string(raw)
		}
		if err := json.Unmarshal([]byte(data), &f); err != nil {
			Exitf("Failed to parse JSON: %v", err)
		}
	} else {
		// Use flags
		f = addFlags.Frontend
		f.Selectors = []api.FrontendSelectorRecord{addFlags.Selector}
	}

	if err := c.Add(id, f); err != nil {
		Exitf("Failed to add frontend '%s': %v", id, err)
	}

	if addFlags.Wait {
		onClose := func() {
			log.Infof("Removing frontend '%s'", id)
			if err := c.Remove(id); err != nil {
				Exitf("Failed to remove frontend '%s': %v", id, err)
			}
			os.Exit(0)
		}
		t := terminate.NewTerminator(log.Infof, onClose)
		t.OsExitDelay = time.Minute
		t.ListenSignals()
	}
}
