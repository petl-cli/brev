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

var sendersCreateNewSenderCmd = &cobra.Command{
	Use:   "create-new-sender",
	Short: "Create a new sender",
	RunE:  runSendersCreateNewSender,
}

var sendersCreateNewSenderFlags struct {
	name  string
	email string
	ips   []string
	body  string
}

func init() {
	sendersCreateNewSenderCmd.Flags().StringVar(&sendersCreateNewSenderFlags.name, "name", "", "From Name to use for the sender")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	sendersCreateNewSenderCmd.Flags().StringVar(&sendersCreateNewSenderFlags.email, "email", "", "From email to use for the sender. A verification email will be sent to this address.")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	sendersCreateNewSenderCmd.Flags().StringSliceVar(&sendersCreateNewSenderFlags.ips, "ips", nil, "**Mandatory in case of dedicated IP**. IPs to associate to the sender ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	sendersCreateNewSenderCmd.Flags().StringVar(&sendersCreateNewSenderFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	sendersCmd.AddCommand(sendersCreateNewSenderCmd)
}

func runSendersCreateNewSender(cmd *cobra.Command, args []string) error {
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
			Required:    false,
			Location:    "body",
			Description: "From Name to use for the sender",
		})
		flags = append(flags, flagSchema{
			Name:        "email",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "From email to use for the sender. A verification email will be sent to this address.",
		})
		flags = append(flags, flagSchema{
			Name:        "ips",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "**Mandatory in case of dedicated IP**. IPs to associate to the sender ",
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
			Description: "sender created",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "create-new-sender",
			"description": "Create a new sender",
			"http": map[string]any{
				"method": "POST",
				"path":   "/senders",
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

	req := &httpclient.Request{
		Method:      "POST",
		Path:        httpclient.SubstitutePath("/senders", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if sendersCreateNewSenderFlags.body != "" {
		if err := json.Unmarshal([]byte(sendersCreateNewSenderFlags.body), &bodyMap); err != nil {
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
		bodyMap["name"] = sendersCreateNewSenderFlags.name
	}
	if cmd.Flags().Changed("email") {
		bodyMap["email"] = sendersCreateNewSenderFlags.email
	}
	if cmd.Flags().Changed("ips") {
		bodyMap["ips"] = sendersCreateNewSenderFlags.ips
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
