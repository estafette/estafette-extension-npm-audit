package main

import (
	"bytes"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"

	"github.com/sethgrid/pester"
)

// SlackWebhookClient is used to send messages to slack using a webhook
type SlackWebhookClient interface {
	SendMessage(string, string, string, string) error
}

type slackWebhookClientImpl struct {
	webhookURL string
}

// NewSlackWebhookClient returns a new SlackWebhookClient
func NewSlackWebhookClient(webhookURL string) SlackWebhookClient {
	return &slackWebhookClientImpl{
		webhookURL: webhookURL,
	}
}

// GetAccessToken returns an access token to access the Bitbucket api
func (sc *slackWebhookClientImpl) SendMessage(target, title, link, message string) (err error) {

	var requestBody io.Reader

	slackMessageBody := SlackMessageBody{
		Channel:  target,
		Username: "Estafette CI",
		Attachments: []SlackMessageAttachment{
			SlackMessageAttachment{
				Fallback:  message,
				Title:     title,
				TitleLink: link,
				Text:      message,
				Color:     "warning",
				MarkdownIn: []string{
					"text",
				},
			},
		},
	}

	data, err := json.Marshal(slackMessageBody)
	if err != nil {
		log.Info().Msgf("Failed marshalling SlackMessageBody: %v. Error: %v", slackMessageBody, err)
		return
	}
	requestBody = bytes.NewReader(data)

	// create client, in order to add headers
	client := pester.New()
	client.MaxRetries = 3
	client.Backoff = pester.ExponentialJitterBackoff
	client.KeepLog = true
	request, err := http.NewRequest("POST", sc.webhookURL, requestBody)
	if err != nil {
		log.Info().Msgf("Failed creating http client: %v", err)
		return
	}

	// add headers
	request.Header.Add("Content-type", "application/json")

	// perform actual request
	response, err := client.Do(request)
	if err != nil {
		log.Info().Msgf("Failed performing http request to Slack: %v", err)
		return
	}

	defer response.Body.Close()

	return
}
