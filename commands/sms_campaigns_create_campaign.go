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

var smsCampaignsCreateCampaignCmd = &cobra.Command{
	Use:   "create-campaign",
	Short: "Creates an SMS campaign",
	RunE:  runSmsCampaignsCreateCampaign,
}

var smsCampaignsCreateCampaignFlags struct {
	name                   string
	sender                 string
	content                string
	scheduledAt            string
	unicodeEnabled         bool
	organisationPrefix     string
	unsubscribeInstruction string
	body                   string
}

func init() {
	smsCampaignsCreateCampaignCmd.Flags().StringVar(&smsCampaignsCreateCampaignFlags.name, "name", "", "Name of the campaign")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	smsCampaignsCreateCampaignCmd.Flags().StringVar(&smsCampaignsCreateCampaignFlags.sender, "sender", "", "Name of the sender. **The number of characters is limited to 11 for alphanumeric characters and 15 for numeric characters** ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	smsCampaignsCreateCampaignCmd.Flags().StringVar(&smsCampaignsCreateCampaignFlags.content, "content", "", "Content of the message. The **maximum characters used per SMS is 160**, if used more than that, it will be counted as more than one SMS ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	smsCampaignsCreateCampaignCmd.Flags().StringVar(&smsCampaignsCreateCampaignFlags.scheduledAt, "scheduled-at", "", "UTC date-time on which the campaign has to run (YYYY-MM-DDTHH:mm:ss.SSSZ). **Prefer to pass your timezone in date-time format for accurate result.** ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	smsCampaignsCreateCampaignCmd.Flags().BoolVar(&smsCampaignsCreateCampaignFlags.unicodeEnabled, "unicode-enabled", false, "Format of the message. It indicates whether the content should be treated as unicode or not. ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	smsCampaignsCreateCampaignCmd.Flags().StringVar(&smsCampaignsCreateCampaignFlags.organisationPrefix, "organisation-prefix", "", "A recognizable prefix will ensure your audience knows who you are. Recommended by U.S. carriers. This will be added as your Brand Name before the message content. **Prefer verifying maximum length of 160 characters including this prefix in message content to avoid multiple sending of same sms.**")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	smsCampaignsCreateCampaignCmd.Flags().StringVar(&smsCampaignsCreateCampaignFlags.unsubscribeInstruction, "unsubscribe-instruction", "", "Instructions to unsubscribe from future communications. Recommended by U.S. carriers. Must include **STOP** keyword. This will be added as instructions after the end of message content. **Prefer verifying maximum length of 160 characters including this instructions in message content to avoid multiple sending of same sms.**")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	smsCampaignsCreateCampaignCmd.Flags().StringVar(&smsCampaignsCreateCampaignFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	smsCampaignsCmd.AddCommand(smsCampaignsCreateCampaignCmd)
}

func runSmsCampaignsCreateCampaign(cmd *cobra.Command, args []string) error {
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
			Description: "Name of the campaign",
		})
		flags = append(flags, flagSchema{
			Name:        "sender",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Name of the sender. **The number of characters is limited to 11 for alphanumeric characters and 15 for numeric characters** ",
		})
		flags = append(flags, flagSchema{
			Name:        "content",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Content of the message. The **maximum characters used per SMS is 160**, if used more than that, it will be counted as more than one SMS ",
		})
		flags = append(flags, flagSchema{
			Name:        "recipients",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "scheduled-at",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "UTC date-time on which the campaign has to run (YYYY-MM-DDTHH:mm:ss.SSSZ). **Prefer to pass your timezone in date-time format for accurate result.** ",
		})
		flags = append(flags, flagSchema{
			Name:        "unicode-enabled",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Format of the message. It indicates whether the content should be treated as unicode or not. ",
		})
		flags = append(flags, flagSchema{
			Name:        "organisation-prefix",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "A recognizable prefix will ensure your audience knows who you are. Recommended by U.S. carriers. This will be added as your Brand Name before the message content. **Prefer verifying maximum length of 160 characters including this prefix in message content to avoid multiple sending of same sms.**",
		})
		flags = append(flags, flagSchema{
			Name:        "unsubscribe-instruction",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Instructions to unsubscribe from future communications. Recommended by U.S. carriers. Must include **STOP** keyword. This will be added as instructions after the end of message content. **Prefer verifying maximum length of 160 characters including this instructions in message content to avoid multiple sending of same sms.**",
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
			"command":     "create-campaign",
			"description": "Creates an SMS campaign",
			"http": map[string]any{
				"method": "POST",
				"path":   "/smsCampaigns",
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

	req := &httpclient.Request{
		Method:      "POST",
		Path:        httpclient.SubstitutePath("/smsCampaigns", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if smsCampaignsCreateCampaignFlags.body != "" {
		if err := json.Unmarshal([]byte(smsCampaignsCreateCampaignFlags.body), &bodyMap); err != nil {
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
		bodyMap["name"] = smsCampaignsCreateCampaignFlags.name
	}
	if cmd.Flags().Changed("sender") {
		bodyMap["sender"] = smsCampaignsCreateCampaignFlags.sender
	}
	if cmd.Flags().Changed("content") {
		bodyMap["content"] = smsCampaignsCreateCampaignFlags.content
	}
	if cmd.Flags().Changed("scheduled-at") {
		bodyMap["scheduledAt"] = smsCampaignsCreateCampaignFlags.scheduledAt
	}
	if cmd.Flags().Changed("unicode-enabled") {
		bodyMap["unicodeEnabled"] = smsCampaignsCreateCampaignFlags.unicodeEnabled
	}
	if cmd.Flags().Changed("organisation-prefix") {
		bodyMap["organisationPrefix"] = smsCampaignsCreateCampaignFlags.organisationPrefix
	}
	if cmd.Flags().Changed("unsubscribe-instruction") {
		bodyMap["unsubscribeInstruction"] = smsCampaignsCreateCampaignFlags.unsubscribeInstruction
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
