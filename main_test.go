package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuditRepository(t *testing.T) {

	finding146 := Finding{"1.9.2", false, false}
	finding534 := Finding{"2.2.0", false, false}
	finding550 := Finding{"2.0.3", false, false}
	resolve146 := Resolve{146, false}
	resolve534 := Resolve{534, true}
	resolve550 := Resolve{550, false}
	resolves := []Resolve{resolve146, resolve534, resolve550}
	actions := []Action{Action{"", resolves}}
	advisories := map[int]Advisory{
		146: Advisory{146, "moderate", []Finding{finding146}},
		534: Advisory{534, "low", []Finding{finding534}},
		550: Advisory{550, "critical", []Finding{finding550}},
	}
	vulnerabilities := Vulnerabilities{0, 1, 1, 0, 1}
	metadata := Metadata{vulnerabilities, 527, 39, 0, 566}
	var auditReport AuditReportBody
	auditReport = AuditReportBody{
		Actions:    actions,
		Advisories: advisories,
		Metadata:   metadata,
	}

	t.Run("AuditRepositoryLowProdCriticalDev", func(t *testing.T) {

		prodVulnLevel := Level(Low)
		devVulnLevel := Level(Critical)
		failBuild, hasVulns := checkVulnerabilities(auditReport, prodVulnLevel, devVulnLevel)
		assert.True(t, failBuild)
		assert.True(t, hasVulns)
	})

	t.Run("AuditRepositoryLowProdLowDev", func(t *testing.T) {

		prodVulnLevel := Level(Low)
		devVulnLevel := Level(Low)
		failBuild, hasVulns := checkVulnerabilities(auditReport, prodVulnLevel, devVulnLevel)
		assert.True(t, failBuild)
		assert.True(t, hasVulns)
	})

	t.Run("AuditRepositoryHighProdNoneDev", func(t *testing.T) {

		prodVulnLevel := Level(Critical)
		devVulnLevel := Level(None)
		failBuild, hasVulns := checkVulnerabilities(auditReport, prodVulnLevel, devVulnLevel)
		assert.True(t, failBuild)
		assert.True(t, hasVulns)
	})

	t.Run("AuditRepositoryWithoutVulnerabilities", func(t *testing.T) {

		auditReportNoVulns := AuditReportBody{
			Advisories: map[int]Advisory{},
			Metadata:   Metadata{},
		}
		prodVulnLevel := Level(Low)
		devVulnLevel := Level(None)
		failBuild, hasVulns := checkVulnerabilities(auditReportNoVulns, prodVulnLevel, devVulnLevel)
		assert.False(t, failBuild)
		assert.False(t, hasVulns)
	})
}

func TestGetSlackClient(t *testing.T) {
	slackChannels := "builds,builds-estafette"
	slackWorkspace := "estafette"
	slackCredentialsJSON := `[{
"name": "estafette",
"type": "idk",
"additionalProperties": {
  "workspace": "estafette",
  "webhook": "https://estafette.slack.com/webhook"
  }
}]`

	t.Run("GetSlackIntegration", func(t *testing.T) {

		slackEnabled, slackWebhookClient := getSlackIntegration(&slackChannels, &slackCredentialsJSON, &slackWorkspace)

		assert.True(t, slackEnabled)
		assert.NotNil(t, slackWebhookClient)
	})

	t.Run("GetSlackIntegrationNoChannels", func(t *testing.T) {
		noChannels := ""

		slackEnabled, slackWebhookClient := getSlackIntegration(&noChannels, &slackCredentialsJSON, &slackWorkspace)

		assert.False(t, slackEnabled)
		assert.Nil(t, slackWebhookClient)
	})

	t.Run("GetSlackIntegrationNoWebhookURLWithCredentials", func(t *testing.T) {

		slackEnabled, slackWebhookClient := getSlackIntegration(&slackChannels, &slackCredentialsJSON, &slackWorkspace)

		assert.True(t, slackEnabled)
		assert.NotNil(t, slackWebhookClient)
	})
}
