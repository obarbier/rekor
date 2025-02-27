//
// Copyright 2021 The Sigstore Authors.
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

package app

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"text/tabwriter"

	"github.com/pkg/errors"
	"github.com/sigstore/rekor/pkg/api"
	"github.com/spf13/cobra"
)

type versionOptions struct {
	json bool
}

var versionOpts = &versionOptions{}

// verifyCmd represents the verify command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "rekor-server version",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runVersion(versionOpts)
	},
}

func init() {
	versionCmd.PersistentFlags().BoolVarP(&versionOpts.json, "json", "j", false,
		"print JSON instead of text")
	rootCmd.AddCommand(versionCmd)
}

func runVersion(opts *versionOptions) error {
	v := VersionInfo()
	res := v.String()

	if opts.json {
		j, err := v.JSONString()
		if err != nil {
			return errors.Wrap(err, "unable to generate JSON from version info")
		}
		res = j
	}

	fmt.Println(res)
	return nil
}

type Info struct {
	GitVersion   string
	GitCommit    string
	GitTreeState string
	BuildDate    string
	GoVersion    string
	Compiler     string
	Platform     string
}

func VersionInfo() Info {
	// These variables typically come from -ldflags settings and in
	// their absence fallback to the global defaults set above.
	return Info{
		GitVersion:   api.GitVersion,
		GitCommit:    api.GitCommit,
		GitTreeState: api.GitTreeState,
		BuildDate:    api.BuildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String returns the string representation of the version info
func (i *Info) String() string {
	b := strings.Builder{}
	w := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "GitVersion:\t%s\n", i.GitVersion)
	fmt.Fprintf(w, "GitCommit:\t%s\n", i.GitCommit)
	fmt.Fprintf(w, "GitTreeState:\t%s\n", i.GitTreeState)
	fmt.Fprintf(w, "BuildDate:\t%s\n", i.BuildDate)
	fmt.Fprintf(w, "GoVersion:\t%s\n", i.GoVersion)
	fmt.Fprintf(w, "Compiler:\t%s\n", i.Compiler)
	fmt.Fprintf(w, "Platform:\t%s\n", i.Platform)

	w.Flush() // #nosec
	return b.String()
}

// JSONString returns the JSON representation of the version info
func (i *Info) JSONString() (string, error) {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}
