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

var webhooksUpdateWebhookByIdCmd = &cobra.Command{
	Use:   "update-webhook-by-id",
	Short: "Update a webhook",
	RunE:  runWebhooksUpdateWebhookById,
}

var webhooksUpdateWebhookByIdFlags struct {
	webhookId   int
	description string
	url         string
	events      []string
	domain      string
	batched     bool
	headers     []string
	body        string
}

func init() {
	webhooksUpdateWebhookByIdCmd.Flags().IntVar(&webhooksUpdateWebhookByIdFlags.webhookId, "webhook-id", 0, "Id of the webhook")
	webhooksUpdateWebhookByIdCmd.MarkFlagRequired("webhook-id")
	webhooksUpdateWebhookByIdCmd.Flags().StringVar(&webhooksUpdateWebhookByIdFlags.description, "description", "", "Description of the webhook")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	webhooksUpdateWebhookByIdCmd.Flags().StringVar(&webhooksUpdateWebhookByIdFlags.url, "url", "", "URL of the webhook")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	webhooksUpdateWebhookByIdCmd.Flags().StringSliceVar(&webhooksUpdateWebhookByIdFlags.events, "events", nil, "- Events triggering the webhook. Possible values for **Transactional** type webhook: #### `sent` OR `request`, `delivered`, `hardBounce`, `softBounce`, `blocked`, `spam`, `invalid`, `deferred`, `click`, `opened`, `uniqueOpened` and `unsubscribed` - Possible values for **Marketing** type webhook: #### `spam`, `opened`, `click`, `hardBounce`, `softBounce`, `unsubscribed`, `listAddition` & `delivered` - Possible values for **Inbound** type webhook: #### `inboundEmailProcessed` ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	webhooksUpdateWebhookByIdCmd.Flags().StringVar(&webhooksUpdateWebhookByIdFlags.domain, "domain", "", "Inbound domain of webhook, used in case of event type `inbound`")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	webhooksUpdateWebhookByIdCmd.Flags().BoolVar(&webhooksUpdateWebhookByIdFlags.batched, "batched", false, "Batching configuration of the webhook, we send batched webhooks if its true")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	webhooksUpdateWebhookByIdCmd.Flags().StringSliceVar(&webhooksUpdateWebhookByIdFlags.headers, "headers", nil, "")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	webhooksUpdateWebhookByIdCmd.Flags().StringVar(&webhooksUpdateWebhookByIdFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	webhooksCmd.AddCommand(webhooksUpdateWebhookByIdCmd)
}

func runWebhooksUpdateWebhookById(cmd *cobra.Command, args []string) error {
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
			Name:        "webhook-id",
			Type:        "integer",
			Required:    true,
			Location:    "path",
			Description: "Id of the webhook",
		})
		flags = append(flags, flagSchema{
			Name:        "description",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Description of the webhook",
		})
		flags = append(flags, flagSchema{
			Name:        "url",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "URL of the webhook",
		})
		flags = append(flags, flagSchema{
			Name:        "events",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "- Events triggering the webhook. Possible values for **Transactional** type webhook: #### `sent` OR `request`, `delivered`, `hardBounce`, `softBounce`, `blocked`, `spam`, `invalid`, `deferred`, `click`, `opened`, `uniqueOpened` and `unsubscribed` - Possible values for **Marketing** type webhook: #### `spam`, `opened`, `click`, `hardBounce`, `softBounce`, `unsubscribed`, `listAddition` & `delivered` - Possible values for **Inbound** type webhook: #### `inboundEmailProcessed` ",
		})
		flags = append(flags, flagSchema{
			Name:        "domain",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Inbound domain of webhook, used in case of event type `inbound`",
		})
		flags = append(flags, flagSchema{
			Name:        "batched",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Batching configuration of the webhook, we send batched webhooks if its true",
		})
		flags = append(flags, flagSchema{
			Name:        "auth",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "Authentication header to be send with the webhook requests",
		})
		flags = append(flags, flagSchema{
			Name:        "headers",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "",
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
			Description: "Webhook updated",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})
		responses = append(responses, responseSchema{
			Status:      "404",
			ContentType: "application/json",
			Description: "Webhook ID not found",
		})

		schema := map[string]any{
			"command":     "update-webhook-by-id",
			"description": "Update a webhook",
			"http": map[string]any{
				"method": "PUT",
				"path":   "/webhooks/{webhookId}",
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
	pathParams["webhookId"] = fmt.Sprintf("%v", webhooksUpdateWebhookByIdFlags.webhookId)

	req := &httpclient.Request{
		Method:      "PUT",
		Path:        httpclient.SubstitutePath("/webhooks/{webhookId}", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if webhooksUpdateWebhookByIdFlags.body != "" {
		if err := json.Unmarshal([]byte(webhooksUpdateWebhookByIdFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("description") {
		bodyMap["description"] = webhooksUpdateWebhookByIdFlags.description
	}
	if cmd.Flags().Changed("url") {
		bodyMap["url"] = webhooksUpdateWebhookByIdFlags.url
	}
	if cmd.Flags().Changed("events") {
		bodyMap["events"] = webhooksUpdateWebhookByIdFlags.events
	}
	if cmd.Flags().Changed("domain") {
		bodyMap["domain"] = webhooksUpdateWebhookByIdFlags.domain
	}
	if cmd.Flags().Changed("batched") {
		bodyMap["batched"] = webhooksUpdateWebhookByIdFlags.batched
	}
	if cmd.Flags().Changed("headers") {
		bodyMap["headers"] = webhooksUpdateWebhookByIdFlags.headers
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
