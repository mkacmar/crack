package output

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/mkacmar/crack/internal/result"
	"github.com/mkacmar/crack/internal/rule"
)

type SARIFReport struct {
	Version string     `json:"version"`
	Schema  string     `json:"$schema"`
	Runs    []SARIFRun `json:"runs"`
}

type SARIFRun struct {
	Tool      SARIFTool       `json:"tool"`
	Results   []SARIFResult   `json:"results"`
	Artifacts []SARIFArtifact `json:"artifacts,omitempty"`
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
	ShortDescription     SARIFMessage       `json:"shortDescription"`
	FullDescription      SARIFMessage       `json:"fullDescription,omitempty"`
	DefaultConfiguration SARIFConfiguration `json:"defaultConfiguration"`
	Help                 SARIFMessage       `json:"help,omitempty"`
}

type SARIFConfiguration struct {
	Level string `json:"level"`
}

type SARIFMessage struct {
	Text string `json:"text"`
}

type SARIFResult struct {
	RuleID    string          `json:"ruleId"`
	Kind      string          `json:"kind,omitempty"`
	Level     string          `json:"level,omitempty"`
	Message   SARIFMessage    `json:"message"`
	Locations []SARIFLocation `json:"locations,omitempty"`
	Fixes     []SARIFFix      `json:"fixes,omitempty"`
}

type SARIFLocation struct {
	PhysicalLocation SARIFPhysicalLocation `json:"physicalLocation"`
}

type SARIFPhysicalLocation struct {
	ArtifactLocation SARIFArtifactLocation `json:"artifactLocation"`
}

type SARIFArtifactLocation struct {
	URI string `json:"uri"`
}

type SARIFArtifact struct {
	Location SARIFArtifactLocation `json:"location"`
	Hashes   map[string]string     `json:"hashes,omitempty"`
}

type SARIFFix struct {
	Description SARIFMessage `json:"description"`
}

type SARIFFormatter struct {
	IncludePassed  bool
	IncludeSkipped bool
}

func (f *SARIFFormatter) Format(report *result.ScanResults, w io.Writer) error {
	output := f.convertToSARIF(report)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	return encoder.Encode(output)
}

func (f *SARIFFormatter) convertToSARIF(report *result.ScanResults) SARIFReport {
	ruleMap := make(map[string]rule.ProcessedResult)
	for _, res := range report.Results {
		for _, check := range res.Results {
			if _, exists := ruleMap[check.RuleID]; !exists {
				ruleMap[check.RuleID] = check
			}
		}
	}

	rules := make([]SARIFRule, 0, len(ruleMap))
	for _, check := range ruleMap {
		rule := SARIFRule{
			ID:   check.RuleID,
			Name: check.Name,
			ShortDescription: SARIFMessage{
				Text: check.Name,
			},
			DefaultConfiguration: SARIFConfiguration{
				Level: "warning",
			},
		}
		if check.Message != "" && check.Message != check.Name {
			rule.FullDescription = SARIFMessage{Text: check.Message}
		}
		rules = append(rules, rule)
	}

	sarifResults := make([]SARIFResult, 0)
	artifactHashes := make(map[string]string)

	for _, res := range report.Results {
		fileURI := toFileURI(res.Path)
		artifactHashes[fileURI] = res.SHA256

		if res.Error != nil {
			sarifResults = append(sarifResults, SARIFResult{
				RuleID: "scan-error",
				Kind:   "fail",
				Level:  "error",
				Message: SARIFMessage{
					Text: fmt.Sprintf("Scan error: %v", res.Error),
				},
				Locations: []SARIFLocation{
					{PhysicalLocation: SARIFPhysicalLocation{
						ArtifactLocation: SARIFArtifactLocation{URI: fileURI},
					}},
				},
			})
			continue
		}

		for _, check := range res.Results {
			if check.Status == rule.StatusPassed && !f.IncludePassed {
				continue
			}
			if check.Status == rule.StatusSkipped && !f.IncludeSkipped {
				continue
			}

			var kind, level string
			switch check.Status {
			case rule.StatusPassed:
				kind = "pass"
			case rule.StatusSkipped:
				kind = "notApplicable"
			default:
				kind = "fail"
				level = "warning"
			}

			sarifResult := SARIFResult{
				RuleID:  check.RuleID,
				Kind:    kind,
				Level:   level,
				Message: SARIFMessage{Text: check.Message},
				Locations: []SARIFLocation{
					{PhysicalLocation: SARIFPhysicalLocation{
						ArtifactLocation: SARIFArtifactLocation{URI: fileURI},
					}},
				},
			}

			if check.Suggestion != "" {
				sarifResult.Fixes = []SARIFFix{
					{Description: SARIFMessage{Text: check.Suggestion}},
				}
			}

			sarifResults = append(sarifResults, sarifResult)
		}
	}

	artifactList := make([]SARIFArtifact, 0, len(artifactHashes))
	for uri, hash := range artifactHashes {
		artifact := SARIFArtifact{
			Location: SARIFArtifactLocation{URI: uri},
		}
		if hash != "" {
			artifact.Hashes = map[string]string{"sha-256": hash}
		}
		artifactList = append(artifactList, artifact)
	}

	return SARIFReport{
		Version: "2.1.0",
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		Runs: []SARIFRun{
			{
				Tool: SARIFTool{
					Driver: SARIFDriver{
						Name:           "crack",
						InformationUri: "https://github.com/mkacmar/crack",
						Rules:          rules,
					},
				},
				Results:   sarifResults,
				Artifacts: artifactList,
			},
		},
	}
}

func toFileURI(path string) string {
	if strings.HasPrefix(path, "/") {
		return "file://" + url.PathEscape(path)
	}
	return path
}
