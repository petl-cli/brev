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

var whatsappCampaignsUpdateCampaignByIdCmd = &cobra.Command{
	Use:   "update-campaign-by-id",
	Short: "Update a WhatsApp campaign",
	RunE:  runWhatsappCampaignsUpdateCampaignById,
}

var whatsappCampaignsUpdateCampaignByIdFlags struct {
	campaignId     int
	campaignName   string
	campaignStatus string
	rescheduleFor  string
	body           string
}

func init() {
	whatsappCampaignsUpdateCampaignByIdCmd.Flags().IntVar(&whatsappCampaignsUpdateCampaignByIdFlags.campaignId, "campaign-id", 0, "id of the campaign")
	whatsappCampaignsUpdateCampaignByIdCmd.MarkFlagRequired("campaign-id")
	whatsappCampaignsUpdateCampaignByIdCmd.Flags().StringVar(&whatsappCampaignsUpdateCampaignByIdFlags.campaignName, "campaign-name", "", "Name of the campaign")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	whatsappCampaignsUpdateCampaignByIdCmd.Flags().StringVar(&whatsappCampaignsUpdateCampaignByIdFlags.campaignStatus, "campaign-status", "", "Status of the campaign")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	whatsappCampaignsUpdateCampaignByIdCmd.Flags().StringVar(&whatsappCampaignsUpdateCampaignByIdFlags.rescheduleFor, "reschedule-for", "", "Reschedule the sending UTC date-time (YYYY-MM-DDTHH:mm:ss.SSSZ) of campaign. **Prefer to pass your timezone in date-time format for accurate result.For example: **2017-06-01T12:30:00+02:00** Use this field to update the scheduledAt of any existing draft or scheduled WhatsApp campaign. ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	whatsappCampaignsUpdateCampaignByIdCmd.Flags().StringVar(&whatsappCampaignsUpdateCampaignByIdFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	whatsappCampaignsCmd.AddCommand(whatsappCampaignsUpdateCampaignByIdCmd)
}

func runWhatsappCampaignsUpdateCampaignById(cmd *cobra.Command, args []string) error {
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
			Name:        "campaign-name",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Name of the campaign",
		})
		flags = append(flags, flagSchema{
			Name:        "campaign-status",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Status of the campaign",
		})
		flags = append(flags, flagSchema{
			Name:        "reschedule-for",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Reschedule the sending UTC date-time (YYYY-MM-DDTHH:mm:ss.SSSZ) of campaign. **Prefer to pass your timezone in date-time format for accurate result.For example: **2017-06-01T12:30:00+02:00** Use this field to update the scheduledAt of any existing draft or scheduled WhatsApp campaign. ",
		})
		flags = append(flags, flagSchema{
			Name:        "recipients",
			Type:        "object",
			Required:    false,
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
			Status:      "204",
			ContentType: "",
			Description: "WhatsApp campaign has been deleted",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "update-campaign-by-id",
			"description": "Update a WhatsApp campaign",
			"http": map[string]any{
				"method": "PUT",
				"path":   "/whatsappCampaigns/{campaignId}",
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
	pathParams["campaignId"] = fmt.Sprintf("%v", whatsappCampaignsUpdateCampaignByIdFlags.campaignId)

	req := &httpclient.Request{
		Method:      "PUT",
		Path:        httpclient.SubstitutePath("/whatsappCampaigns/{campaignId}", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if whatsappCampaignsUpdateCampaignByIdFlags.body != "" {
		if err := json.Unmarshal([]byte(whatsappCampaignsUpdateCampaignByIdFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("campaign-name") {
		bodyMap["campaignName"] = whatsappCampaignsUpdateCampaignByIdFlags.campaignName
	}
	if cmd.Flags().Changed("campaign-status") {
		bodyMap["campaignStatus"] = whatsappCampaignsUpdateCampaignByIdFlags.campaignStatus
	}
	if cmd.Flags().Changed("reschedule-for") {
		bodyMap["rescheduleFor"] = whatsappCampaignsUpdateCampaignByIdFlags.rescheduleFor
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
