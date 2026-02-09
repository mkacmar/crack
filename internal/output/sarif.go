package output

import (
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strings"
	"time"

	"github.com/mkacmar/crack/internal/analyzer"
	"github.com/mkacmar/crack/internal/version"
	"github.com/mkacmar/crack/rule"
)

type SARIFReport struct {
	Version string     `json:"version"`
	Schema  string     `json:"$schema"`
	Runs    []SARIFRun `json:"runs"`
}

type SARIFRun struct {
	Tool        SARIFTool         `json:"tool"`
	Invocations []SARIFInvocation `json:"invocations,omitempty"`
	Results     []SARIFResult     `json:"results"`
	Artifacts   []SARIFArtifact   `json:"artifacts,omitempty"`
}

type SARIFInvocation struct {
	CommandLine                string                 `json:"commandLine,omitempty"`
	Arguments                  []string               `json:"arguments,omitempty"`
	ExecutionSuccessful        bool                   `json:"executionSuccessful"`
	StartTimeUtc               string                 `json:"startTimeUtc,omitempty"`
	EndTimeUtc                 string                 `json:"endTimeUtc,omitempty"`
	WorkingDirectory           *SARIFArtifactLocation `json:"workingDirectory,omitempty"`
	ToolExecutionNotifications []SARIFNotification    `json:"toolExecutionNotifications,omitempty"`
}

type SARIFNotification struct {
	Level     string          `json:"level"`
	Message   SARIFMessage    `json:"message"`
	Locations []SARIFLocation `json:"locations,omitempty"`
}

type SARIFTool struct {
	Driver SARIFDriver `json:"driver"`
}

type SARIFDriver struct {
	Name           string      `json:"name"`
	InformationUri string      `json:"informationUri,omitempty"`
	Version        string      `json:"version,omitempty"`
	Rules          []SARIFRule `json:"rules,omitempty"`
}

type SARIFRule struct {
	ID                   string             `json:"id"`
	Name                 string             `json:"name"`
	HelpUri              string             `json:"helpUri,omitempty"`
	FullDescription      SARIFMessage       `json:"fullDescription,omitempty"`
	DefaultConfiguration SARIFConfiguration `json:"defaultConfiguration"`
}

type SARIFConfiguration struct {
	Level string `json:"level"`
}

type SARIFMessage struct {
	Text string `json:"text"`
}

type SARIFResult struct {
	RuleIndex int             `json:"ruleIndex"`
	Kind      string          `json:"kind,omitempty"`
	Level     string          `json:"level,omitempty"`
	Message   SARIFMessage    `json:"message"`
	Locations []SARIFLocation `json:"locations,omitempty"`
}

type SARIFLocation struct {
	PhysicalLocation SARIFPhysicalLocation `json:"physicalLocation"`
}

type SARIFPhysicalLocation struct {
	ArtifactIndex int `json:"artifactIndex"`
}

type SARIFArtifactLocation struct {
	URI string `json:"uri"`
}

type SARIFArtifact struct {
	Location SARIFArtifactLocation `json:"location"`
	Hashes   map[string]string     `json:"hashes,omitempty"`
}

type InvocationInfo struct {
	CommandLine string
	Arguments   []string
	StartTime   time.Time
	EndTime     time.Time
	WorkingDir  string
	Successful  bool
}

type SARIFFormatter struct {
	IncludePassed  bool
	IncludeSkipped bool
	Invocation     *InvocationInfo
}

func (f *SARIFFormatter) Format(report *analyzer.Report, w io.Writer) error {
	output := f.convertToSARIF(report)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	return encoder.Encode(output)
}

func (f *SARIFFormatter) convertToSARIF(report *analyzer.Report) SARIFReport {
	rules, ruleIndex := f.buildRules(report)
	artifacts, artifactIndex := f.buildArtifacts(report)
	results, notifications := f.buildResults(report, ruleIndex, artifactIndex)

	run := SARIFRun{
		Tool: SARIFTool{
			Driver: SARIFDriver{
				Name:           "crack",
				Version:        version.Version,
				InformationUri: "https://github.com/mkacmar/crack",
				Rules:          rules,
			},
		},
		Results:   results,
		Artifacts: artifacts,
	}

	if f.Invocation != nil || len(notifications) > 0 {
		run.Invocations = []SARIFInvocation{f.buildInvocation(notifications)}
	}

	return SARIFReport{
		Version: "2.1.0",
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		Runs:    []SARIFRun{run},
	}
}

func (f *SARIFFormatter) buildRules(report *analyzer.Report) ([]SARIFRule, map[string]int) {
	ruleMap := make(map[string]analyzer.FindingWithSuggestion)
	for _, res := range report.Results {
		for _, finding := range res.Findings {
			if _, exists := ruleMap[finding.RuleID]; !exists {
				ruleMap[finding.RuleID] = finding
			}
		}
	}

	ruleIDs := make([]string, 0, len(ruleMap))
	for id := range ruleMap {
		ruleIDs = append(ruleIDs, id)
	}
	slices.Sort(ruleIDs)

	rules := make([]SARIFRule, 0, len(ruleMap))
	ruleIndex := make(map[string]int)
	for i, id := range ruleIDs {
		finding := ruleMap[id]
		ruleIndex[id] = i
		r := SARIFRule{
			ID:      finding.RuleID,
			Name:    finding.Name,
			HelpUri: ruleHelpURL(finding.Name, finding.RuleID),
			DefaultConfiguration: SARIFConfiguration{
				Level: "warning",
			},
		}
		if finding.Message != "" && finding.Message != finding.Name {
			r.FullDescription = SARIFMessage{Text: finding.Message}
		}
		rules = append(rules, r)
	}

	return rules, ruleIndex
}

func (f *SARIFFormatter) buildResults(report *analyzer.Report, ruleIndex, artifactIndex map[string]int) ([]SARIFResult, []SARIFNotification) {
	sarifResults := make([]SARIFResult, 0)
	notifications := make([]SARIFNotification, 0)

	for _, res := range report.Results {
		fileURI := toFileURI(res.Path)

		if res.Error != nil {
			notifications = append(notifications, SARIFNotification{
				Level: "error",
				Message: SARIFMessage{
					Text: fmt.Sprintf("Scan error: %v", res.Error),
				},
				Locations: []SARIFLocation{
					{PhysicalLocation: SARIFPhysicalLocation{
						ArtifactIndex: artifactIndex[fileURI],
					}},
				},
			})
			continue
		}

		for _, finding := range res.Findings {
			if finding.Status == rule.StatusPassed && !f.IncludePassed {
				continue
			}
			if finding.Status == rule.StatusSkipped && !f.IncludeSkipped {
				continue
			}

			var kind, level string
			switch finding.Status {
			case rule.StatusPassed:
				kind = "pass"
			case rule.StatusSkipped:
				kind = "notApplicable"
			default:
				kind = "fail"
				level = "warning"
			}

			message := finding.Message
			if finding.Suggestion != "" {
				message = message + " " + finding.Suggestion
			}

			sarifResult := SARIFResult{
				RuleIndex: ruleIndex[finding.RuleID],
				Kind:      kind,
				Level:     level,
				Message:   SARIFMessage{Text: message},
				Locations: []SARIFLocation{
					{PhysicalLocation: SARIFPhysicalLocation{
						ArtifactIndex: artifactIndex[fileURI],
					}},
				},
			}

			sarifResults = append(sarifResults, sarifResult)
		}
	}

	return sarifResults, notifications
}

func (f *SARIFFormatter) buildArtifacts(report *analyzer.Report) ([]SARIFArtifact, map[string]int) {
	artifactHashes := make(map[string]string)
	for _, res := range report.Results {
		fileURI := toFileURI(res.Path)
		artifactHashes[fileURI] = res.SHA256
	}

	uris := make([]string, 0, len(artifactHashes))
	for uri := range artifactHashes {
		uris = append(uris, uri)
	}
	slices.Sort(uris)

	artifacts := make([]SARIFArtifact, 0, len(artifactHashes))
	artifactIndex := make(map[string]int)
	for i, uri := range uris {
		artifactIndex[uri] = i
		artifact := SARIFArtifact{
			Location: SARIFArtifactLocation{URI: uri},
		}
		if hash := artifactHashes[uri]; hash != "" {
			artifact.Hashes = map[string]string{"sha-256": hash}
		}
		artifacts = append(artifacts, artifact)
	}
	return artifacts, artifactIndex
}

func (f *SARIFFormatter) buildInvocation(notifications []SARIFNotification) SARIFInvocation {
	var inv SARIFInvocation

	if f.Invocation != nil {
		inv.CommandLine = f.Invocation.CommandLine
		inv.Arguments = f.Invocation.Arguments
		inv.ExecutionSuccessful = f.Invocation.Successful
		if !f.Invocation.StartTime.IsZero() {
			inv.StartTimeUtc = f.Invocation.StartTime.UTC().Format(time.RFC3339)
		}
		if !f.Invocation.EndTime.IsZero() {
			inv.EndTimeUtc = f.Invocation.EndTime.UTC().Format(time.RFC3339)
		}
		if f.Invocation.WorkingDir != "" {
			inv.WorkingDirectory = &SARIFArtifactLocation{URI: toFileURI(f.Invocation.WorkingDir)}
		}
	}

	if len(notifications) > 0 {
		inv.ToolExecutionNotifications = notifications
	}

	return inv
}

func toFileURI(path string) string {
	if strings.HasPrefix(path, "/") {
		return "file://" + path
	}
	return path
}

const wikiBaseURL = "https://github.com/mkacmar/crack/wiki/Rules"

func ruleHelpURL(name, id string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	return fmt.Sprintf("%s#%s-%s", wikiBaseURL, slug, id)
}
