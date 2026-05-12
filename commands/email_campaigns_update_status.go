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

var emailCampaignsUpdateStatusCmd = &cobra.Command{
	Use:   "update-status",
	Short: "Update an email campaign status",
	RunE:  runEmailCampaignsUpdateStatus,
}

var emailCampaignsUpdateStatusFlags struct {
	campaignId int
	status     string
	body       string
}

func init() {
	emailCampaignsUpdateStatusCmd.Flags().IntVar(&emailCampaignsUpdateStatusFlags.campaignId, "campaign-id", 0, "Id of the campaign")
	emailCampaignsUpdateStatusCmd.MarkFlagRequired("campaign-id")
	emailCampaignsUpdateStatusCmd.Flags().StringVar(&emailCampaignsUpdateStatusFlags.status, "status", "", "Note:- **replicateTemplate** status will be available **only for template type campaigns.** ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateStatusCmd.Flags().StringVar(&emailCampaignsUpdateStatusFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	emailCampaignsCmd.AddCommand(emailCampaignsUpdateStatusCmd)
}

func runEmailCampaignsUpdateStatus(cmd *cobra.Command, args []string) error {
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
			Name:        "status",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Note:- **replicateTemplate** status will be available **only for template type campaigns.** ",
		})

		type responseSchema struct {
			Status      string `json:"status"`
			ContentType string `json:"content_type,omitempty"`
			Description string `json:"description,omitempty"`
		}
		var responses []responseSchema
		responses = append(responses, responseSchema{
			Status:      "204",
			ContentType: "",
			Description: "The campaign status has been updated successfully",
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
			"command":     "update-status",
			"description": "Update an email campaign status",
			"http": map[string]any{
				"method": "PUT",
				"path":   "/emailCampaigns/{campaignId}/status",
			},
			"input": map[string]any{
				"flags":         flags,
				"body_flag":     true,
				"body_required": true,
			},
			"output": map[string]any{
				"responses": responses,
			},
			"semantics": map[string]any{
				"safe":         false,
				"idempotent":   true,
				"reversible":   true,
				"side_effects": []string{"mutates_resource"},
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
	pathParams["campaignId"] = fmt.Sprintf("%v", emailCampaignsUpdateStatusFlags.campaignId)

	req := &httpclient.Request{
		Method:      "PUT",
		Path:        httpclient.SubstitutePath("/emailCampaigns/{campaignId}/status", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if emailCampaignsUpdateStatusFlags.body != "" {
		if err := json.Unmarshal([]byte(emailCampaignsUpdateStatusFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("status") {
		bodyMap["status"] = emailCampaignsUpdateStatusFlags.status
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
