package main

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
)

// AuditReportBody represents the body coming from npm audit
type AuditReportBody struct {
	Advisories map[int]Advisory `json:"advisories,omitempty"`
	Metadata   Metadata         `json:"metadata,omitempty"`
	Error      CommandError     `json:"error,omitempty"`
}

type Advisory struct {
	Id       int       `json:"id,omitempty"`
	Severity string    `json:"severity,omitempty"`
	Findings []Finding `json:"findings,omitempty"`
}

type Finding struct {
	Version  string `json:"version,omitempty"`
	Dev      bool   `json:"dev,omitempty"`
	Optional bool   `json:"optional,omitempty"`
	Bundled  bool   `json:"bundled,omitempty"`
}

type Metadata struct {
	Vulnerabilities      Vulnerabilities `json:"vulnerabilities,omitempty"`
	Dependencies         int             `json:"dependencies,omitempty"`
	DevDependencies      int             `json:"devDependencies,omitempty"`
	OptionalDependencies int             `json:"optionalDependencies,omitempty"`
	TotalDependencies    int             `json:"totalDependencies,omitempty"`
}

type Vulnerabilities struct {
	Info     int `json:"info,omitempty"`
	Low      int `json:"low,omitempty"`
	Moderate int `json:"moderate,omitempty"`
	High     int `json:"high,omitempty"`
	Critical int `json:"critical,omitempty"`
}

type CommandError struct {
	Code    string `json:"code,omitempty"`
	Summary string `json:"summary,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

func readAuditReport(report string) AuditReportBody {
	var auditReport AuditReportBody
	err := json.Unmarshal([]byte(report), &auditReport)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed unmarshalling audit report")
	}
	return auditReport
}
