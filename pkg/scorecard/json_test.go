// Copyright 2021 OpenSSF Scorecard Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package scorecard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/xeipuuv/gojsonschema"

	"github.com/ossf/scorecard/v5/checker"
	"github.com/ossf/scorecard/v5/config"
	"github.com/ossf/scorecard/v5/finding"
	"github.com/ossf/scorecard/v5/log"
)

func jsonMockDocRead() *mockDoc {
	d := map[string]mockCheck{
		"Check-Name": {
			name:        "Check-Name",
			risk:        "High",
			short:       "short description for Check-Name",
			description: "not used",
			url:         "https://github.com/ossf/scorecard/blob/main/docs/checks.md#check-name",
			tags:        []string{"not-used1", "not-used2"},
			remediation: []string{"not-used1", "not-used2"},
		},
		"Check-Name2": {
			name:        "Check-Name2",
			risk:        "Medium",
			short:       "short description for Check-Name2",
			description: "not used",
			url:         "https://github.com/ossf/scorecard/blob/main/docs/checks.md#check-name2",
			tags:        []string{"not-used1", "not-used2"},
			remediation: []string{"not-used1", "not-used2"},
		},
		"Check-Name3": {
			name:        "Check-Name3",
			risk:        "Low",
			short:       "short description for Check-Name3",
			description: "not used",
			url:         "https://github.com/ossf/scorecard/blob/main/docs/checks.md#check-name3",
			tags:        []string{"not-used1", "not-used2"},
			remediation: []string{"not-used1", "not-used2"},
		},
	}

	m := mockDoc{checks: d}
	return &m
}

//nolint:gocognit
func TestJSONOutput(t *testing.T) {
	t.Parallel()

	repoCommit := "68bc59901773ab4c051dfcea0cc4201a1567ab32"
	scorecardCommit := "ccbc59901773ab4c051dfcea0cc4201a1567abdd"
	scorecardVersion := "1.2.3"
	repoName := "org/name"
	date, err := time.Parse(time.RFC3339, "2023-03-02T10:30:43-06:00")
	if err != nil {
		t.Fatalf("time.Parse: %v", err)
	}
	t.Logf("date: %v", date)

	checkDocs := jsonMockDocRead()

	//nolint:govet
	tests := []struct {
		name            string
		expected        string
		showDetails     bool
		showAnnotations bool
		logLevel        log.Level
		result          Result
	}{
		{
			name:        "check-1",
			showDetails: true,
			expected:    "./testdata/check1.json",
			logLevel:    log.DebugLevel,
			result: Result{
				Repo: RepoInfo{
					Name:      repoName,
					CommitSHA: repoCommit,
				},
				Scorecard: ScorecardInfo{
					Version:   scorecardVersion,
					CommitSHA: scorecardCommit,
				},
				Date: date,
				Checks: []checker.CheckResult{
					{
						Details: []checker.CheckDetail{
							{
								Type: checker.DetailWarn,
								Msg: checker.LogMessage{
									Text:    "warn message",
									Path:    "src/file1.cpp",
									Type:    finding.FileTypeSource,
									Offset:  5,
									Snippet: "if (bad) {BUG();}",
								},
							},
						},
						Score:  5,
						Reason: "half score reason",
						Name:   "Check-Name",
					},
				},
				Metadata: []string{},
			},
		},
		{
			name:            "check-1 annotations",
			showDetails:     true,
			showAnnotations: true,
			expected:        "./testdata/check1_annotations.json",
			logLevel:        log.DebugLevel,
			result: Result{
				Repo: RepoInfo{
					Name:      repoName,
					CommitSHA: repoCommit,
				},
				Scorecard: ScorecardInfo{
					Version:   scorecardVersion,
					CommitSHA: scorecardCommit,
				},
				Date: date,
				Checks: []checker.CheckResult{
					{
						Details: []checker.CheckDetail{
							{
								Type: checker.DetailWarn,
								Msg: checker.LogMessage{
									Text:    "warn message",
									Path:    "src/file1.cpp",
									Type:    finding.FileTypeSource,
									Offset:  5,
									Snippet: "if (bad) {BUG();}",
								},
							},
						},
						Score:  5,
						Reason: "half score reason",
						Name:   "Check-Name",
					},
				},
				Config: config.Config{
					Annotations: []config.Annotation{
						{
							Checks: []string{"Check-Name"},
							Reasons: []config.ReasonGroup{
								{Reason: "test-data"},
								{Reason: "remediated"},
							},
						},
					},
				},
				Metadata: []string{},
			},
		},
		{
			name:        "check-2",
			showDetails: true,
			expected:    "./testdata/check2.json",
			logLevel:    log.DebugLevel,
			result: Result{
				Repo: RepoInfo{
					Name:      repoName,
					CommitSHA: repoCommit,
				},
				Scorecard: ScorecardInfo{
					Version:   scorecardVersion,
					CommitSHA: scorecardCommit,
				},
				Date: date,
				Checks: []checker.CheckResult{
					{
						Details: []checker.CheckDetail{
							{
								Type: checker.DetailWarn,
								Msg: checker.LogMessage{
									Text:   "warn message",
									Path:   "bin/binary.elf",
									Type:   finding.FileTypeBinary,
									Offset: 0,
								},
							},
						},
						Score:  checker.MinResultScore,
						Reason: "min score reason",
						Name:   "Check-Name",
					},
				},
				Metadata: []string{},
			},
		},
		{
			name:        "check-3",
			showDetails: true,
			expected:    "./testdata/check3.json",
			logLevel:    log.InfoLevel,
			result: Result{
				Repo: RepoInfo{
					Name:      repoName,
					CommitSHA: repoCommit,
				},
				Scorecard: ScorecardInfo{
					Version:   scorecardVersion,
					CommitSHA: scorecardCommit,
				},
				Date: date,
				Checks: []checker.CheckResult{
					{
						Details: []checker.CheckDetail{
							{
								Type: checker.DetailWarn,
								Msg: checker.LogMessage{
									Text:   "warn message",
									Path:   "bin/binary.elf",
									Type:   finding.FileTypeBinary,
									Offset: 0,
								},
							},
						},
						Score:  checker.MinResultScore,
						Reason: "min result reason",
						Name:   "Check-Name",
					},
					{
						Details: []checker.CheckDetail{
							{
								Type: checker.DetailWarn,
								Msg: checker.LogMessage{
									Text:    "warn message",
									Path:    "src/doc.txt",
									Type:    finding.FileTypeText,
									Offset:  3,
									Snippet: "some text",
								},
							},
						},
						Score:  checker.MinResultScore,
						Reason: "min result reason",
						Name:   "Check-Name2",
					},
					{
						Details: []checker.CheckDetail{
							{
								Type: checker.DetailInfo,
								Msg: checker.LogMessage{
									Text:    "info message",
									Path:    "some/path.js",
									Type:    finding.FileTypeSource,
									Offset:  3,
									Snippet: "if (bad) {BUG();}",
								},
							},
							{
								Type: checker.DetailWarn,
								Msg: checker.LogMessage{
									Text:    "warn message",
									Path:    "some/path.py",
									Type:    finding.FileTypeSource,
									Offset:  3,
									Snippet: "if (bad) {BUG2();}",
								},
							},
							{
								Type: checker.DetailDebug,
								Msg: checker.LogMessage{
									Text:    "debug message",
									Path:    "some/path.go",
									Type:    finding.FileTypeSource,
									Offset:  3,
									Snippet: "if (bad) {BUG5();}",
								},
							},
						},
						Score:  checker.InconclusiveResultScore,
						Reason: "inconclusive reason",
						Name:   "Check-Name3",
					},
				},
				Metadata: []string{},
			},
		},
		{
			name:        "check-4",
			showDetails: true,
			expected:    "./testdata/check4.json",
			logLevel:    log.DebugLevel,
			result: Result{
				Repo: RepoInfo{
					Name:      repoName,
					CommitSHA: repoCommit,
				},
				Scorecard: ScorecardInfo{
					Version:   scorecardVersion,
					CommitSHA: scorecardCommit,
				},
				Date: date,
				Checks: []checker.CheckResult{
					{
						Details: []checker.CheckDetail{
							{
								Type: checker.DetailWarn,
								Msg: checker.LogMessage{
									Text:   "warn message",
									Path:   "bin/binary.elf",
									Type:   finding.FileTypeBinary,
									Offset: 0,
								},
							},
						},
						Score:  checker.MinResultScore,
						Reason: "min result reason",
						Name:   "Check-Name",
					},
					{
						Details: []checker.CheckDetail{
							{
								Type: checker.DetailWarn,
								Msg: checker.LogMessage{
									Text:    "warn message",
									Path:    "src/doc.txt",
									Type:    finding.FileTypeText,
									Offset:  3,
									Snippet: "some text",
								},
							},
						},
						Score:  checker.MinResultScore,
						Reason: "min result reason",
						Name:   "Check-Name2",
					},
					{
						Details: []checker.CheckDetail{
							{
								Type: checker.DetailInfo,
								Msg: checker.LogMessage{
									Text:    "info message",
									Path:    "some/path.js",
									Type:    finding.FileTypeSource,
									Offset:  3,
									Snippet: "if (bad) {BUG();}",
								},
							},
							{
								Type: checker.DetailWarn,
								Msg: checker.LogMessage{
									Text:    "warn message",
									Path:    "some/path.py",
									Type:    finding.FileTypeSource,
									Offset:  3,
									Snippet: "if (bad) {BUG2();}",
								},
							},
							{
								Type: checker.DetailDebug,
								Msg: checker.LogMessage{
									Text:    "debug message",
									Path:    "some/path.go",
									Type:    finding.FileTypeSource,
									Offset:  3,
									Snippet: "if (bad) {BUG5();}",
								},
							},
						},
						Score:  checker.InconclusiveResultScore,
						Reason: "inconclusive reason",
						Name:   "Check-Name3",
					},
				},
				Metadata: []string{},
			},
		},
		{
			name:        "check-5",
			showDetails: true,
			expected:    "./testdata/check5.json",
			logLevel:    log.WarnLevel,
			result: Result{
				Repo: RepoInfo{
					Name:      repoName,
					CommitSHA: repoCommit,
				},
				Scorecard: ScorecardInfo{
					Version:   scorecardVersion,
					CommitSHA: scorecardCommit,
				},
				Date: date,
				Checks: []checker.CheckResult{
					{
						Details: []checker.CheckDetail{
							{
								Type: checker.DetailWarn,
								Msg: checker.LogMessage{
									Text:    "warn message",
									Path:    "src/file1.cpp",
									Type:    finding.FileTypeSource,
									Offset:  5,
									Snippet: "if (bad) {BUG();}",
								},
							},
						},
						Score:  6,
						Reason: "six score reason",
						Name:   "Check-Name",
					},
				},
				Metadata: []string{},
			},
		},
		{
			name:        "check-6",
			showDetails: true,
			expected:    "./testdata/check6.json",
			logLevel:    log.WarnLevel,
			result: Result{
				Repo: RepoInfo{
					Name:      repoName,
					CommitSHA: repoCommit,
				},
				Scorecard: ScorecardInfo{
					Version:   scorecardVersion,
					CommitSHA: scorecardCommit,
				},
				Date: date,
				Checks: []checker.CheckResult{
					{
						Details: []checker.CheckDetail{
							{
								Type: checker.DetailWarn,
								Msg: checker.LogMessage{
									Text: "warn message",
									Path: "https://domain.com/something",
									Type: finding.FileTypeURL,
								},
							},
						},
						Score:  6,
						Reason: "six score reason",
						Name:   "Check-Name",
					},
				},
				Metadata: []string{},
			},
		},
	}

	// Load the JSON schema.
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd: %s", err)
	}
	schemaLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", path.Join(cwd, "json.v2.schema")))
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		t.Fatalf("gojsonschema.NewSchema: %s", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var content []byte
			var err error
			content, err = os.ReadFile(tt.expected)
			if err != nil {
				t.Fatalf("cannot read file: %v", err)
			}

			var expected bytes.Buffer
			n, err := expected.Write(content)
			if err != nil {
				t.Fatalf("%s: cannot write buffer: %v", tt.name, err)
			}
			if n != len(content) {
				t.Fatalf("%s: write %d bytes but expected %d", tt.name, n, len(content))
			}

			var result bytes.Buffer
			o := &AsJSON2ResultOption{
				Details:     tt.showDetails,
				LogLevel:    tt.logLevel,
				Annotations: tt.showAnnotations,
			}
			err = tt.result.AsJSON2(&result, checkDocs, o)
			if err != nil {
				t.Fatalf("%s: AsJSON2: %v", tt.name, err)
			}

			// TODO: add indentation to AsJSON2() and remove
			// the calls to Unmarshall() and Marshall() below.

			// Unmarshall expected output.
			var js JSONScorecardResultV2
			if err := json.Unmarshal(expected.Bytes(), &js); err != nil {
				t.Fatalf("%s: json.Unmarshal: %s", tt.name, err)
			}

			// Marshall.
			var es bytes.Buffer
			encoder := json.NewEncoder(&es)
			if err := encoder.Encode(js); err != nil {
				t.Fatalf("%s: Encode: %s", tt.name, err)
			}

			// Compare outputs.
			r := bytes.Compare(result.Bytes(), es.Bytes())
			if r != 0 {
				t.Fatalf("%s: invalid result %d", tt.name, r)
			}

			// Validate schema.
			docLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", path.Join(cwd, tt.expected)))
			rr, err := schema.Validate(docLoader)
			if err != nil {
				t.Fatalf("%s: Validate error: %s", tt.name, err.Error())
			}

			if !rr.Valid() {
				s := ""
				for _, desc := range rr.Errors() {
					s += fmt.Sprintf("- %s\n", desc)
				}
				t.Fatalf("%s: invalid format: %s", tt.name, s)
			}
		})
	}
}

func TestExperimentalFromJSON2_time(t *testing.T) {
	t.Parallel()
	//nolint:govet
	tests := []struct {
		name    string
		result  string
		want    time.Time
		wantErr bool
	}{
		{
			name:   "main RFC3339 format",
			result: `{"date":"2006-01-02T15:04:05+00:00","repo":{"name":"github.com/foo/bar","commit":"HEAD"},"score":-1.0,"metadata":null}`,
			want:   time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC),
		},
		{
			name:   "backup 2006-01-02 format",
			result: `{"date":"2023-09-26","repo":{"name":"github.com/foo/bar","commit":"HEAD"},"score":-1.0,"metadata":null}`,
			want:   time.Date(2023, time.September, 26, 0, 0, 0, 0, time.UTC),
		},
		{
			name:    "not RFC3339 or 2006-01-02 format",
			result:  `{"date":"January 1, 2023","repo":{"name":"github.com/foo/bar","commit":"HEAD"},"score":-1.0,"metadata":null}`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, _, err := ExperimentalFromJSON2(strings.NewReader(tt.result))
			if tt.wantErr != (err != nil) {
				t.Fatalf("got: %v, wantedErr: %v", err, tt.wantErr)
			}
			if !got.Date.Equal(tt.want) {
				t.Errorf("got: %v, wanted: %v", got.Date, tt.want)
			}
		})
	}
}
