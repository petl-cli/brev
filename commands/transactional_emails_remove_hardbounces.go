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

var transactionalEmailsRemoveHardbouncesCmd = &cobra.Command{
	Use:   "remove-hardbounces",
	Short: "Delete hardbounces",
	RunE:  runTransactionalEmailsRemoveHardbounces,
}

var transactionalEmailsRemoveHardbouncesFlags struct {
	startDate    string
	endDate      string
	contactEmail string
	body         string
}

func init() {
	transactionalEmailsRemoveHardbouncesCmd.Flags().StringVar(&transactionalEmailsRemoveHardbouncesFlags.startDate, "start-date", "", "Starting date (YYYY-MM-DD) of the time period for deletion. The hardbounces occurred after this date will be deleted. Must be less than or equal to the endDate")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsRemoveHardbouncesCmd.Flags().StringVar(&transactionalEmailsRemoveHardbouncesFlags.endDate, "end-date", "", "Ending date (YYYY-MM-DD) of the time period for deletion. The hardbounces until this date will be deleted. Must be greater than or equal to the startDate")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsRemoveHardbouncesCmd.Flags().StringVar(&transactionalEmailsRemoveHardbouncesFlags.contactEmail, "contact-email", "", "Target a specific email address")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsRemoveHardbouncesCmd.Flags().StringVar(&transactionalEmailsRemoveHardbouncesFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	transactionalEmailsCmd.AddCommand(transactionalEmailsRemoveHardbouncesCmd)
}

func runTransactionalEmailsRemoveHardbounces(cmd *cobra.Command, args []string) error {
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
			Name:        "start-date",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Starting date (YYYY-MM-DD) of the time period for deletion. The hardbounces occurred after this date will be deleted. Must be less than or equal to the endDate",
		})
		flags = append(flags, flagSchema{
			Name:        "end-date",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Ending date (YYYY-MM-DD) of the time period for deletion. The hardbounces until this date will be deleted. Must be greater than or equal to the startDate",
		})
		flags = append(flags, flagSchema{
			Name:        "contact-email",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Target a specific email address",
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
			Description: "Hardbounces deleted",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "remove-hardbounces",
			"description": "Delete hardbounces",
			"http": map[string]any{
				"method": "POST",
				"path":   "/smtp/deleteHardbounces",
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
		Path:        httpclient.SubstitutePath("/smtp/deleteHardbounces", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if transactionalEmailsRemoveHardbouncesFlags.body != "" {
		if err := json.Unmarshal([]byte(transactionalEmailsRemoveHardbouncesFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("start-date") {
		bodyMap["startDate"] = transactionalEmailsRemoveHardbouncesFlags.startDate
	}
	if cmd.Flags().Changed("end-date") {
		bodyMap["endDate"] = transactionalEmailsRemoveHardbouncesFlags.endDate
	}
	if cmd.Flags().Changed("contact-email") {
		bodyMap["contactEmail"] = transactionalEmailsRemoveHardbouncesFlags.contactEmail
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
