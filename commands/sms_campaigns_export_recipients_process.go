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

var smsCampaignsExportRecipientsProcessCmd = &cobra.Command{
	Use:   "export-recipients-process",
	Short: "Export an SMS campaign's recipients",
	RunE:  runSmsCampaignsExportRecipientsProcess,
}

var smsCampaignsExportRecipientsProcessFlags struct {
	campaignId     int
	notifyUrl      string
	recipientsType string
	body           string
}

func init() {
	smsCampaignsExportRecipientsProcessCmd.Flags().IntVar(&smsCampaignsExportRecipientsProcessFlags.campaignId, "campaign-id", 0, "id of the campaign")
	smsCampaignsExportRecipientsProcessCmd.MarkFlagRequired("campaign-id")
	smsCampaignsExportRecipientsProcessCmd.Flags().StringVar(&smsCampaignsExportRecipientsProcessFlags.notifyUrl, "notify-url", "", "URL that will be called once the export process is finished. For reference, https://help.brevo.com/hc/en-us/articles/360007666479")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	smsCampaignsExportRecipientsProcessCmd.Flags().StringVar(&smsCampaignsExportRecipientsProcessFlags.recipientsType, "recipients-type", "", "Filter the recipients based on how they interacted with the campaign")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	smsCampaignsExportRecipientsProcessCmd.Flags().StringVar(&smsCampaignsExportRecipientsProcessFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	smsCampaignsCmd.AddCommand(smsCampaignsExportRecipientsProcessCmd)
}

func runSmsCampaignsExportRecipientsProcess(cmd *cobra.Command, args []string) error {
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
			Description: "id of the campaign",
		})
		flags = append(flags, flagSchema{
			Name:        "notify-url",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "URL that will be called once the export process is finished. For reference, https://help.brevo.com/hc/en-us/articles/360007666479",
		})
		flags = append(flags, flagSchema{
			Name:        "recipients-type",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Filter the recipients based on how they interacted with the campaign",
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
			"command":     "export-recipients-process",
			"description": "Export an SMS campaign's recipients",
			"http": map[string]any{
				"method": "POST",
				"path":   "/smsCampaigns/{campaignId}/exportRecipients",
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
				"side_effects": []string{"creates_resource"},
				"impact":       "medium",
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
	pathParams["campaignId"] = fmt.Sprintf("%v", smsCampaignsExportRecipientsProcessFlags.campaignId)

	req := &httpclient.Request{
		Method:      "POST",
		Path:        httpclient.SubstitutePath("/smsCampaigns/{campaignId}/exportRecipients", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if smsCampaignsExportRecipientsProcessFlags.body != "" {
		if err := json.Unmarshal([]byte(smsCampaignsExportRecipientsProcessFlags.body), &bodyMap); err != nil {
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
		bodyMap["notifyURL"] = smsCampaignsExportRecipientsProcessFlags.notifyUrl
	}
	if cmd.Flags().Changed("recipients-type") {
		bodyMap["recipientsType"] = smsCampaignsExportRecipientsProcessFlags.recipientsType
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
