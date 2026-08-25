package output

import (
	"bufio"
	"encoding/json/jsontext"
	json "encoding/json/v2"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"go.kacmar.sk/crack/internal/suggestions"
	"go.kacmar.sk/crack/internal/version"
	"go.kacmar.sk/crack/rule"
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

const (
	sarifSchema  = "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json"
	sarifVersion = "2.1.0"
	toolInfoURI  = "https://github.com/mkacmar/crack"
)

// SARIFWriter streams a SARIF 2.1.0 document, emitting each result as it arrives.
// Results are written first and the rules and artifacts they reference by index are written on Close, so only those tables are held in memory rather than every finding.
// JSON object members are unordered, so a consumer resolves the indices regardless of the order they appear in.
type SARIFWriter struct {
	buf *bufio.Writer
	enc *jsontext.Encoder
	err error

	IncludePassed  bool
	IncludeSkipped bool

	rules         []SARIFRule
	ruleIndex     map[string]int
	artifacts     []SARIFArtifact
	artifactIndex map[string]int
	notifications []SARIFNotification
}

// NewSARIFWriter begins a SARIF document on w and opens the results array.
// Close must be called to finish the document.
func NewSARIFWriter(w io.Writer, includePassed, includeSkipped bool) *SARIFWriter {
	buf := bufio.NewWriter(w)
	s := &SARIFWriter{
		buf:            buf,
		enc:            jsontext.NewEncoder(buf, jsontext.WithIndent("  ")),
		IncludePassed:  includePassed,
		IncludeSkipped: includeSkipped,
		ruleIndex:      make(map[string]int),
		artifactIndex:  make(map[string]int),
	}
	s.tokens(
		jsontext.BeginObject,
		jsontext.String("version"), jsontext.String(sarifVersion),
		jsontext.String("$schema"), jsontext.String(sarifSchema),
		jsontext.String("runs"), jsontext.BeginArray,
		jsontext.BeginObject,
		jsontext.String("results"), jsontext.BeginArray,
	)
	return s
}

// Notify records a tool execution notification, written under the invocation on Close.
func (s *SARIFWriter) Notify(level, message string) {
	s.notifications = append(s.notifications, SARIFNotification{
		Level:   level,
		Message: SARIFMessage{Text: message},
	})
}

// Write emits the results for one file and registers the rules and artifact they reference.
func (s *SARIFWriter) Write(res DecoratedFileResult) error {
	if s.err != nil {
		return s.err
	}

	artifact := s.registerArtifact(res)

	if res.Error != nil {
		s.notifications = append(s.notifications, SARIFNotification{
			Level:     "error",
			Message:   SARIFMessage{Text: fmt.Sprintf("Scan error: %v", res.Error)},
			Locations: []SARIFLocation{{PhysicalLocation: SARIFPhysicalLocation{ArtifactIndex: artifact}}},
		})
		return s.err
	}

	for _, finding := range res.Findings {
		if finding.Status == rule.StatusPassed && !s.IncludePassed {
			continue
		}
		if finding.Status == rule.StatusSkipped && !s.IncludeSkipped {
			continue
		}
		s.value(s.result(finding, artifact))
	}
	return s.err
}

// Close writes the tables the results refer to and finishes the document.
func (s *SARIFWriter) Close(inv *InvocationInfo) error {
	if s.err != nil {
		return s.err
	}

	s.tokens(jsontext.EndArray)

	s.tokens(jsontext.String("tool"))
	s.value(SARIFTool{Driver: SARIFDriver{
		Name:           "crack",
		Version:        version.Version,
		InformationUri: toolInfoURI,
		Rules:          s.rules,
	}})

	s.tokens(jsontext.String("artifacts"))
	s.value(s.artifacts)

	if inv != nil || len(s.notifications) > 0 {
		s.tokens(jsontext.String("invocations"))
		s.value([]SARIFInvocation{buildInvocation(inv, s.notifications)})
	}

	s.tokens(jsontext.EndObject, jsontext.EndArray, jsontext.EndObject)

	if s.err != nil {
		return s.err
	}
	return s.buf.Flush()
}

func (s *SARIFWriter) tokens(ts ...jsontext.Token) {
	for _, t := range ts {
		if s.err != nil {
			return
		}
		s.err = s.enc.WriteToken(t)
	}
}

func (s *SARIFWriter) value(v any) {
	if s.err != nil {
		return
	}
	s.err = json.MarshalEncode(s.enc, v)
}

func (s *SARIFWriter) result(finding suggestions.DecoratedFinding, artifact int) SARIFResult {
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
		message += " " + finding.Suggestion
	}

	return SARIFResult{
		RuleIndex: s.registerRule(finding),
		Kind:      kind,
		Level:     level,
		Message:   SARIFMessage{Text: message},
		Locations: []SARIFLocation{{PhysicalLocation: SARIFPhysicalLocation{ArtifactIndex: artifact}}},
	}
}

func (s *SARIFWriter) registerArtifact(res DecoratedFileResult) int {
	uri := toFileURI(res.Path)
	if i, ok := s.artifactIndex[uri]; ok {
		return i
	}
	artifact := SARIFArtifact{Location: SARIFArtifactLocation{URI: uri}}
	if res.Identity.SHA256 != "" {
		artifact.Hashes = map[string]string{"sha-256": res.Identity.SHA256}
	}
	s.artifactIndex[uri] = len(s.artifacts)
	s.artifacts = append(s.artifacts, artifact)
	return s.artifactIndex[uri]
}

func (s *SARIFWriter) registerRule(finding suggestions.DecoratedFinding) int {
	if i, ok := s.ruleIndex[finding.RuleID]; ok {
		return i
	}
	r := SARIFRule{
		ID:                   finding.RuleID,
		Name:                 finding.Name,
		HelpUri:              ruleHelpURL(finding.Name, version.Version),
		DefaultConfiguration: SARIFConfiguration{Level: "warning"},
	}
	if finding.Message != "" && finding.Message != finding.Name {
		r.FullDescription = SARIFMessage{Text: finding.Message}
	}
	s.ruleIndex[finding.RuleID] = len(s.rules)
	s.rules = append(s.rules, r)
	return s.ruleIndex[finding.RuleID]
}

func buildInvocation(info *InvocationInfo, notifications []SARIFNotification) SARIFInvocation {
	var inv SARIFInvocation

	if info != nil {
		inv.CommandLine = info.CommandLine
		inv.Arguments = info.Arguments
		inv.ExecutionSuccessful = info.Successful
		if !info.StartTime.IsZero() {
			inv.StartTimeUtc = info.StartTime.UTC().Format(time.RFC3339)
		}
		if !info.EndTime.IsZero() {
			inv.EndTimeUtc = info.EndTime.UTC().Format(time.RFC3339)
		}
		if info.WorkingDir != "" {
			inv.WorkingDirectory = &SARIFArtifactLocation{URI: toFileURI(info.WorkingDir)}
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

const repoBaseURL = "https://github.com/mkacmar/crack"

var semverTag = regexp.MustCompile(`^v\d+\.\d+\.\d+$`)

func ruleHelpURL(name, ver string) string {
	ref := "main"
	if semverTag.MatchString(ver) {
		ref = ver
	}
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	return fmt.Sprintf("%s/blob/%s/docs/rules.md#%s", repoBaseURL, ref, slug)
}
