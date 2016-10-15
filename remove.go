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

import "github.com/spf13/cobra"

var (
	cmdRemove = &cobra.Command{
		Use:   "rm",
		Short: "Remove frontends",
		Run:   cmdRemoveRun,
	}
)

func init() {
	cmdMain.AddCommand(cmdRemove)
}

func cmdRemoveRun(cmd *cobra.Command, args []string) {
	c := newAPIClient()
	for _, id := range args {
		if err := c.Remove(id); err != nil {
			Exitf("Failed to remove frontend '%s': %v", id, err)
		}
	}
}
