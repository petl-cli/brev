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

var transactionalSmsSendSmsMessageToMobileCmd = &cobra.Command{
	Use:   "send-sms-message-to-mobile",
	Short: "Send SMS message to a mobile number",
	RunE:  runTransactionalSmsSendSmsMessageToMobile,
}

var transactionalSmsSendSmsMessageToMobileFlags struct {
	sender             string
	recipient          string
	content            string
	type_              string
	tag                string
	webUrl             string
	unicodeEnabled     bool
	organisationPrefix string
	body               string
}

func init() {
	transactionalSmsSendSmsMessageToMobileCmd.Flags().StringVar(&transactionalSmsSendSmsMessageToMobileFlags.sender, "sender", "", "Name of the sender. **The number of characters is limited to 11 for alphanumeric characters and 15 for numeric characters** ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalSmsSendSmsMessageToMobileCmd.Flags().StringVar(&transactionalSmsSendSmsMessageToMobileFlags.recipient, "recipient", "", "Mobile number to send SMS with the country code")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalSmsSendSmsMessageToMobileCmd.Flags().StringVar(&transactionalSmsSendSmsMessageToMobileFlags.content, "content", "", "Content of the message. If more than **160 characters** long, will be sent as multiple text messages ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalSmsSendSmsMessageToMobileCmd.Flags().StringVar(&transactionalSmsSendSmsMessageToMobileFlags.type_, "type", "", "Type of the SMS. Marketing SMS messages are those sent typically with marketing content. Transactional SMS messages are sent to individuals and are triggered in response to some action, such as a sign-up, purchase, etc.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalSmsSendSmsMessageToMobileCmd.Flags().StringVar(&transactionalSmsSendSmsMessageToMobileFlags.tag, "tag", "", "Tag of the message")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalSmsSendSmsMessageToMobileCmd.Flags().StringVar(&transactionalSmsSendSmsMessageToMobileFlags.webUrl, "web-url", "", "Webhook to call for each event triggered by the message (delivered etc.)")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalSmsSendSmsMessageToMobileCmd.Flags().BoolVar(&transactionalSmsSendSmsMessageToMobileFlags.unicodeEnabled, "unicode-enabled", false, "Format of the message. It indicates whether the content should be treated as unicode or not. ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalSmsSendSmsMessageToMobileCmd.Flags().StringVar(&transactionalSmsSendSmsMessageToMobileFlags.organisationPrefix, "organisation-prefix", "", "A recognizable prefix will ensure your audience knows who you are. Recommended by U.S. carriers. This will be added as your Brand Name before the message content. **Prefer verifying maximum length of 160 characters including this prefix in message content to avoid multiple sending of same sms.**")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalSmsSendSmsMessageToMobileCmd.Flags().StringVar(&transactionalSmsSendSmsMessageToMobileFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	transactionalSmsCmd.AddCommand(transactionalSmsSendSmsMessageToMobileCmd)
}

func runTransactionalSmsSendSmsMessageToMobile(cmd *cobra.Command, args []string) error {
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
			Name:        "sender",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Name of the sender. **The number of characters is limited to 11 for alphanumeric characters and 15 for numeric characters** ",
		})
		flags = append(flags, flagSchema{
			Name:        "recipient",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Mobile number to send SMS with the country code",
		})
		flags = append(flags, flagSchema{
			Name:        "content",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Content of the message. If more than **160 characters** long, will be sent as multiple text messages ",
		})
		flags = append(flags, flagSchema{
			Name:        "type",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Type of the SMS. Marketing SMS messages are those sent typically with marketing content. Transactional SMS messages are sent to individuals and are triggered in response to some action, such as a sign-up, purchase, etc.",
		})
		flags = append(flags, flagSchema{
			Name:        "tag",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Tag of the message",
		})
		flags = append(flags, flagSchema{
			Name:        "web-url",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Webhook to call for each event triggered by the message (delivered etc.)",
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

		type responseSchema struct {
			Status      string `json:"status"`
			ContentType string `json:"content_type,omitempty"`
			Description string `json:"description,omitempty"`
		}
		var responses []responseSchema
		responses = append(responses, responseSchema{
			Status:      "201",
			ContentType: "application/json",
			Description: "SMS has been sent successfully to the recipient",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})
		responses = append(responses, responseSchema{
			Status:      "402",
			ContentType: "application/json",
			Description: "You don't have enough credit to send your SMS. Please update your plan",
		})

		schema := map[string]any{
			"command":     "send-sms-message-to-mobile",
			"description": "Send SMS message to a mobile number",
			"http": map[string]any{
				"method": "POST",
				"path":   "/transactionalSMS/sms",
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
		Path:        httpclient.SubstitutePath("/transactionalSMS/sms", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if transactionalSmsSendSmsMessageToMobileFlags.body != "" {
		if err := json.Unmarshal([]byte(transactionalSmsSendSmsMessageToMobileFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("sender") {
		bodyMap["sender"] = transactionalSmsSendSmsMessageToMobileFlags.sender
	}
	if cmd.Flags().Changed("recipient") {
		bodyMap["recipient"] = transactionalSmsSendSmsMessageToMobileFlags.recipient
	}
	if cmd.Flags().Changed("content") {
		bodyMap["content"] = transactionalSmsSendSmsMessageToMobileFlags.content
	}
	if cmd.Flags().Changed("type") {
		bodyMap["type"] = transactionalSmsSendSmsMessageToMobileFlags.type_
	}
	if cmd.Flags().Changed("tag") {
		bodyMap["tag"] = transactionalSmsSendSmsMessageToMobileFlags.tag
	}
	if cmd.Flags().Changed("web-url") {
		bodyMap["webUrl"] = transactionalSmsSendSmsMessageToMobileFlags.webUrl
	}
	if cmd.Flags().Changed("unicode-enabled") {
		bodyMap["unicodeEnabled"] = transactionalSmsSendSmsMessageToMobileFlags.unicodeEnabled
	}
	if cmd.Flags().Changed("organisation-prefix") {
		bodyMap["organisationPrefix"] = transactionalSmsSendSmsMessageToMobileFlags.organisationPrefix
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
