package output

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/mkacmar/crack/internal/model"
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
	Level     string          `json:"level"`
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
}

type SARIFFix struct {
	Description SARIFMessage `json:"description"`
}

type SARIFFormatter struct {
	IncludePassed bool
}

func (f *SARIFFormatter) Format(report *model.ScanResults, w io.Writer) error {
	output := f.convertToSARIF(report)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	return encoder.Encode(output)
}

func (f *SARIFFormatter) convertToSARIF(report *model.ScanResults) SARIFReport {
	ruleMap := make(map[string]model.RuleResult)
	for _, result := range report.Results {
		for _, check := range result.Results {
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
	artifactSet := make(map[string]bool)

	for _, result := range report.Results {
		fileURI := toFileURI(result.Path)
		artifactSet[fileURI] = true

		if result.Error != nil {
			sarifResults = append(sarifResults, SARIFResult{
				RuleID: "scan-error",
				Kind:   "fail",
				Level:  "error",
				Message: SARIFMessage{
					Text: fmt.Sprintf("Scan error: %v", result.Error),
				},
				Locations: []SARIFLocation{
					{PhysicalLocation: SARIFPhysicalLocation{
						ArtifactLocation: SARIFArtifactLocation{URI: fileURI},
					}},
				},
			})
			continue
		}

		for _, check := range result.Results {
			if check.State == model.CheckStatePassed && !f.IncludePassed {
				continue
			}
			if check.State == model.CheckStateSkipped {
				continue
			}

			kind := "fail"
			level := "warning"
			if check.State == model.CheckStatePassed {
				kind = "pass"
				level = "note"
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

	artifactList := make([]SARIFArtifact, 0, len(artifactSet))
	for uri := range artifactSet {
		artifactList = append(artifactList, SARIFArtifact{
			Location: SARIFArtifactLocation{URI: uri},
		})
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
