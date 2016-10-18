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
	"sort"
	"strconv"
	"strings"

	api "github.com/pulcy/robin-api"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

type field struct {
	Header string
	Get    func(d fieldData) string
}

type fieldData struct {
	ID       string
	Frontend api.FrontendRecord
	Selector api.FrontendSelectorRecord
}

const (
	defaultFields = "id,port,domain,path-prefix,frontend-port"
)

var (
	cmdList = &cobra.Command{
		Use:     "ls",
		Short:   "Show frontends",
		Run:     cmdListRun,
		Example: fmt.Sprintf("%s ls --fields=%s", projectName, strings.Join(allFieldIDs(), ",")),
	}
	listFlags struct {
		fields string
	}
	allFields = map[string]field{
		"id":            field{"ID", func(d fieldData) string { return d.ID }},
		"service":       field{"Service", func(d fieldData) string { return d.Frontend.Service }},
		"port":          field{"Port", func(d fieldData) string { return strconv.Itoa(d.Selector.ServicePort) }},
		"mode":          field{"Mode", func(d fieldData) string { return def(d.Frontend.Mode, "http") }},
		"domain":        field{"Domain", func(d fieldData) string { return d.Selector.Domain }},
		"path-prefix":   field{"Path-prefix", func(d fieldData) string { return d.Selector.PathPrefix }},
		"frontend-port": field{"Frontend-port", func(d fieldData) string { return strconv.Itoa(d.Selector.FrontendPort) }},
		"private":       field{"Private", func(d fieldData) string { return strconv.FormatBool(d.Selector.Private) }},
		"weight":        field{"Weight", func(d fieldData) string { return strconv.Itoa(d.Selector.Weight) }},
		"ssl-cert":      field{"SSL cert", func(d fieldData) string { return d.Selector.SslCert }},
	}
)

func init() {
	cmdList.Flags().StringVar(&listFlags.fields, "fields", defaultFields, "The field to include in the list")

	cmdMain.AddCommand(cmdList)
}

func cmdListRun(cmd *cobra.Command, args []string) {
	c := newAPIClient()
	result, err := c.All()
	if err != nil {
		Exitf("Failed to list frontends: %v", err)
	}
	if len(result) == 0 {
		Infof("No frontends found.")
	} else {
		fieldIDs := strings.Split(listFlags.fields, ",")
		var header []string
		for _, fieldID := range fieldIDs {
			f, ok := allFields[fieldID]
			if !ok {
				Exitf("Unknown field '%s'", fieldID)
			}
			header = append(header, f.Header)
		}
		lines := []string{
			strings.Join(header, "|"),
		}
		for id, rec := range result {
			for _, sel := range rec.Selectors {
				data := fieldData{
					ID:       id,
					Frontend: rec,
					Selector: sel,
				}
				var cols []string
				for _, fieldID := range fieldIDs {
					f, ok := allFields[fieldID]
					if !ok {
						Exitf("Unknown field '%s'", fieldID)
					}
					cols = append(cols, f.Get(data))
				}
				lines = append(lines, strings.Join(cols, "|"))
			}
		}
		sort.Strings(lines[1:])
		fmt.Println(columnize.SimpleFormat(lines))
	}
}

func allFieldIDs() []string {
	var list []string
	for k := range allFields {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}

func def(value, defaultValue string) string {
	if value != "" {
		return value
	}
	return defaultValue
}
