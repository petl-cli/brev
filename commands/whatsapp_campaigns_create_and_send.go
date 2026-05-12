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

var whatsappCampaignsCreateAndSendCmd = &cobra.Command{
	Use:   "create-and-send",
	Short: "Create and Send a WhatsApp campaign",
	RunE:  runWhatsappCampaignsCreateAndSend,
}

var whatsappCampaignsCreateAndSendFlags struct {
	name        string
	templateId  int
	scheduledAt string
	body        string
}

func init() {
	whatsappCampaignsCreateAndSendCmd.Flags().StringVar(&whatsappCampaignsCreateAndSendFlags.name, "name", "", "Name of the WhatsApp campaign creation")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	whatsappCampaignsCreateAndSendCmd.Flags().IntVar(&whatsappCampaignsCreateAndSendFlags.templateId, "template-id", 0, "Id of the WhatsApp template in **approved** state")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	whatsappCampaignsCreateAndSendCmd.Flags().StringVar(&whatsappCampaignsCreateAndSendFlags.scheduledAt, "scheduled-at", "", "Sending UTC date-time (YYYY-MM-DDTHH:mm:ss.SSSZ). **Prefer to pass your timezone in date-time format for accurate result.For example: **2017-06-01T12:30:00+02:00** ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	whatsappCampaignsCreateAndSendCmd.Flags().StringVar(&whatsappCampaignsCreateAndSendFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	whatsappCampaignsCmd.AddCommand(whatsappCampaignsCreateAndSendCmd)
}

func runWhatsappCampaignsCreateAndSend(cmd *cobra.Command, args []string) error {
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
			Name:        "name",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Name of the WhatsApp campaign creation",
		})
		flags = append(flags, flagSchema{
			Name:        "template-id",
			Type:        "integer",
			Required:    true,
			Location:    "body",
			Description: "Id of the WhatsApp template in **approved** state",
		})
		flags = append(flags, flagSchema{
			Name:        "scheduled-at",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Sending UTC date-time (YYYY-MM-DDTHH:mm:ss.SSSZ). **Prefer to pass your timezone in date-time format for accurate result.For example: **2017-06-01T12:30:00+02:00** ",
		})
		flags = append(flags, flagSchema{
			Name:        "recipients",
			Type:        "object",
			Required:    true,
			Location:    "body",
			Description: "Segment ids and List ids to include/exclude from campaign",
		})

		type responseSchema struct {
			Status      string `json:"status"`
			ContentType string `json:"content_type,omitempty"`
			Description string `json:"description,omitempty"`
		}
		var responses []responseSchema
		responses = append(responses, responseSchema{
			Status:      "201",
			ContentType: "application/json",
			Description: "successfully created",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "create-and-send",
			"description": "Create and Send a WhatsApp campaign",
			"http": map[string]any{
				"method": "POST",
				"path":   "/whatsappCampaigns",
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

	req := &httpclient.Request{
		Method:      "POST",
		Path:        httpclient.SubstitutePath("/whatsappCampaigns", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if whatsappCampaignsCreateAndSendFlags.body != "" {
		if err := json.Unmarshal([]byte(whatsappCampaignsCreateAndSendFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("name") {
		bodyMap["name"] = whatsappCampaignsCreateAndSendFlags.name
	}
	if cmd.Flags().Changed("template-id") {
		bodyMap["templateId"] = whatsappCampaignsCreateAndSendFlags.templateId
	}
	if cmd.Flags().Changed("scheduled-at") {
		bodyMap["scheduledAt"] = whatsappCampaignsCreateAndSendFlags.scheduledAt
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
