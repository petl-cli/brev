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

var resellerRemoveCreditsFromChildCmd = &cobra.Command{
	Use:   "remove-credits-from-child",
	Short: "Remove Email and/or SMS credits from a specific child account",
	RunE:  runResellerRemoveCreditsFromChild,
}

var resellerRemoveCreditsFromChildFlags struct {
	childIdentifier string
	sms             int
	email           int
	body            string
}

func init() {
	resellerRemoveCreditsFromChildCmd.Flags().StringVar(&resellerRemoveCreditsFromChildFlags.childIdentifier, "child-identifier", "", "Either auth key or id of reseller's child")
	resellerRemoveCreditsFromChildCmd.MarkFlagRequired("child-identifier")
	resellerRemoveCreditsFromChildCmd.Flags().IntVar(&resellerRemoveCreditsFromChildFlags.sms, "sms", 0, "**Required if email credits are empty.** SMS credits to be removed from the child account ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	resellerRemoveCreditsFromChildCmd.Flags().IntVar(&resellerRemoveCreditsFromChildFlags.email, "email", 0, "**Required if sms credits are empty.** Email credits to be removed from the child account ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	resellerRemoveCreditsFromChildCmd.Flags().StringVar(&resellerRemoveCreditsFromChildFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	resellerCmd.AddCommand(resellerRemoveCreditsFromChildCmd)
}

func runResellerRemoveCreditsFromChild(cmd *cobra.Command, args []string) error {
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
			Name:        "child-identifier",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "Either auth key or id of reseller's child",
		})
		flags = append(flags, flagSchema{
			Name:        "sms",
			Type:        "integer",
			Required:    false,
			Location:    "body",
			Description: "**Required if email credits are empty.** SMS credits to be removed from the child account ",
		})
		flags = append(flags, flagSchema{
			Name:        "email",
			Type:        "integer",
			Required:    false,
			Location:    "body",
			Description: "**Required if sms credits are empty.** Email credits to be removed from the child account ",
		})

		type responseSchema struct {
			Status      string `json:"status"`
			ContentType string `json:"content_type,omitempty"`
			Description string `json:"description,omitempty"`
		}
		var responses []responseSchema
		responses = append(responses, responseSchema{
			Status:      "200",
			ContentType: "application/json",
			Description: "Credits removed",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})
		responses = append(responses, responseSchema{
			Status:      "403",
			ContentType: "application/json",
			Description: "Current account is not a reseller",
		})
		responses = append(responses, responseSchema{
			Status:      "404",
			ContentType: "application/json",
			Description: "Child auth key or child id not found",
		})

		schema := map[string]any{
			"command":     "remove-credits-from-child",
			"description": "Remove Email and/or SMS credits from a specific child account",
			"http": map[string]any{
				"method": "POST",
				"path":   "/reseller/children/{childIdentifier}/credits/remove",
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
	pathParams["childIdentifier"] = fmt.Sprintf("%v", resellerRemoveCreditsFromChildFlags.childIdentifier)

	req := &httpclient.Request{
		Method:      "POST",
		Path:        httpclient.SubstitutePath("/reseller/children/{childIdentifier}/credits/remove", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if resellerRemoveCreditsFromChildFlags.body != "" {
		if err := json.Unmarshal([]byte(resellerRemoveCreditsFromChildFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("sms") {
		bodyMap["sms"] = resellerRemoveCreditsFromChildFlags.sms
	}
	if cmd.Flags().Changed("email") {
		bodyMap["email"] = resellerRemoveCreditsFromChildFlags.email
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
