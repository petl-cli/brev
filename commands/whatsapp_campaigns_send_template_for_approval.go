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

var whatsappCampaignsSendTemplateForApprovalCmd = &cobra.Command{
	Use:   "send-template-for-approval",
	Short: "Send your WhatsApp template for approval",
	RunE:  runWhatsappCampaignsSendTemplateForApproval,
}

var whatsappCampaignsSendTemplateForApprovalFlags struct {
	templateId int
}

func init() {
	whatsappCampaignsSendTemplateForApprovalCmd.Flags().IntVar(&whatsappCampaignsSendTemplateForApprovalFlags.templateId, "template-id", 0, "id of the template")
	whatsappCampaignsSendTemplateForApprovalCmd.MarkFlagRequired("template-id")

	whatsappCampaignsCmd.AddCommand(whatsappCampaignsSendTemplateForApprovalCmd)
}

func runWhatsappCampaignsSendTemplateForApproval(cmd *cobra.Command, args []string) error {
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
			Name:        "template-id",
			Type:        "integer",
			Required:    true,
			Location:    "path",
			Description: "id of the template",
		})

		type responseSchema struct {
			Status      string `json:"status"`
			ContentType string `json:"content_type,omitempty"`
			Description string `json:"description,omitempty"`
		}
		var responses []responseSchema
		responses = append(responses, responseSchema{
			Status:      "200",
			ContentType: "",
			Description: "WhatsApp template sent for approval",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "send-template-for-approval",
			"description": "Send your WhatsApp template for approval",
			"http": map[string]any{
				"method": "POST",
				"path":   "/whatsappCampaigns/template/approval/{templateId}",
			},
			"input": map[string]any{
				"flags":         flags,
				"body_flag":     false,
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
	pathParams["templateId"] = fmt.Sprintf("%v", whatsappCampaignsSendTemplateForApprovalFlags.templateId)

	req := &httpclient.Request{
		Method:      "POST",
		Path:        httpclient.SubstitutePath("/whatsappCampaigns/template/approval/{templateId}", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

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
