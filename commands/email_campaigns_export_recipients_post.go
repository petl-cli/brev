package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/rishimantri795/CLICreator/runtime/httpclient"
	"github.com/rishimantri795/CLICreator/runtime/output"
	"github.com/spf13/cobra"
)

var emailCampaignsExportRecipientsPostCmd = &cobra.Command{
	Use:   "export-recipients-post",
	Short: "Export the recipients of an email campaign",
	RunE:  runEmailCampaignsExportRecipientsPost,
}

var emailCampaignsExportRecipientsPostFlags struct {
	campaignId     int
	notifyUrl      string
	recipientsType string
	body           string
}

func init() {
	emailCampaignsExportRecipientsPostCmd.Flags().IntVar(&emailCampaignsExportRecipientsPostFlags.campaignId, "campaign-id", 0, "Id of the campaign")
	emailCampaignsExportRecipientsPostCmd.MarkFlagRequired("campaign-id")
	emailCampaignsExportRecipientsPostCmd.Flags().StringVar(&emailCampaignsExportRecipientsPostFlags.notifyUrl, "notify-url", "", "Webhook called once the export process is finished. For reference, https://help.brevo.com/hc/en-us/articles/360007666479")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsExportRecipientsPostCmd.Flags().StringVar(&emailCampaignsExportRecipientsPostFlags.recipientsType, "recipients-type", "", "Type of recipients to export for a campaign")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsExportRecipientsPostCmd.Flags().StringVar(&emailCampaignsExportRecipientsPostFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	emailCampaignsCmd.AddCommand(emailCampaignsExportRecipientsPostCmd)
}

func runEmailCampaignsExportRecipientsPost(cmd *cobra.Command, args []string) error {
	// --schema: print full input/output type contract without making any network call.
	if rootFlags.schema {
		type flagSchema struct {
			Name        string `json:"name"`
			Type        string `json:"type"`
			Required    bool   `json:"required"`
			Location    string `json:"location"`
			Description string `json:"description,omitempty"`
		}
		var flags []flagSchema
		flags = append(flags, flagSchema{
			Name:        "campaign-id",
			Type:        "integer",
			Required:    true,
			Location:    "path",
			Description: "Id of the campaign",
		})
		flags = append(flags, flagSchema{
			Name:        "notify-url",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Webhook called once the export process is finished. For reference, https://help.brevo.com/hc/en-us/articles/360007666479",
		})
		flags = append(flags, flagSchema{
			Name:        "recipients-type",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Type of recipients to export for a campaign",
		})

		type responseSchema struct {
			Status      string `json:"status"`
			ContentType string `json:"content_type,omitempty"`
			Description string `json:"description,omitempty"`
		}
		var responses []responseSchema
		responses = append(responses, responseSchema{
			Status:      "202",
			ContentType: "application/json",
			Description: "process id created",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})
		responses = append(responses, responseSchema{
			Status:      "404",
			ContentType: "application/json",
			Description: "Campaign ID not found",
		})

		schema := map[string]any{
			"command":     "export-recipients-post",
			"description": "Export the recipients of an email campaign",
			"http": map[string]any{
				"method": "POST",
				"path":   "/emailCampaigns/{campaignId}/exportRecipients",
			},
			"input": map[string]any{
				"flags":         flags,
				"body_flag":     true,
				"body_required": false,
			},
			"output": map[string]any{
				"responses": responses,
			},
			"semantics": map[string]any{
				"safe":         false,
				"idempotent":   false,
				"reversible":   true,
				"side_effects": []string{"sends_notification"},
				"impact":       "high",
			},
			"requires_auth": true,
		}
		data, _ := json.MarshalIndent(schema, "", "  ")
		fmt.Fprintln(_stdoutCounter, string(data))
		return nil
	}

	cfg, err := rootConfig()
	if err != nil {
		e := output.NetworkError(err)
		e.Write(os.Stderr)
		return output.NewExitError(e)
	}

	client := httpclient.New(cfg.BaseURL, cfg.AuthProvider())
	client.Debug = rootFlags.debug
	client.DryRun = rootFlags.dryRun
	if rootFlags.noRetries {
		client.RetryConfig.MaxRetries = 0
	}

	// Build path params
	pathParams := map[string]string{}
	pathParams["campaignId"] = fmt.Sprintf("%v", emailCampaignsExportRecipientsPostFlags.campaignId)

	req := &httpclient.Request{
		Method:      "POST",
		Path:        httpclient.SubstitutePath("/emailCampaigns/{campaignId}/exportRecipients", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if emailCampaignsExportRecipientsPostFlags.body != "" {
		if err := json.Unmarshal([]byte(emailCampaignsExportRecipientsPostFlags.body), &bodyMap); err != nil {
			_invState.errorType = "parse_error"
			cliErr := &output.CLIError{
				Error:    true,
				Code:     "validation_error",
				Message:  fmt.Sprintf("invalid JSON in --body: %v", err),
				ExitCode: output.ExitValidation,
			}
			cliErr.Write(os.Stderr)
			return output.NewExitError(cliErr)
		}
	}
	// Individual flags overlay onto body (flags take precedence over --body JSON)
	if cmd.Flags().Changed("notify-url") {
		bodyMap["notifyURL"] = emailCampaignsExportRecipientsPostFlags.notifyUrl
	}
	if cmd.Flags().Changed("recipients-type") {
		bodyMap["recipientsType"] = emailCampaignsExportRecipientsPostFlags.recipientsType
	}
	req.Body = bodyMap

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "deadline exceeded") {
			_invState.errorType = "timeout"
		} else {
			_invState.errorType = "network_error"
		}
		e := output.NetworkError(err)
		e.Write(os.Stderr)
		return output.NewExitError(e)
	}

	if resp.StatusCode >= 400 {
		if resp.StatusCode >= 500 {
			_invState.errorType = "http_5xx"
		} else {
			_invState.errorType = "http_4xx"
		}
		_invState.errorCode = resp.StatusCode
		e := output.HTTPError(resp.StatusCode, resp.Body)
		e.Write(os.Stderr)
		return output.NewExitError(e)
	}

	if rootFlags.jq != "" {
		return output.JQFilter(_stdoutCounter, resp.Body, rootFlags.jq)
	}
	return output.Print(_stdoutCounter, resp.Body, output.Format(cfg.OutputFormat))
}
