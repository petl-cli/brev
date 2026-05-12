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

var masterAccountCreateSubAccountCmd = &cobra.Command{
	Use:   "create-sub-account",
	Short: "Create a new sub-account under a master account.",
	RunE:  runMasterAccountCreateSubAccount,
}

var masterAccountCreateSubAccountFlags struct {
	companyName string
	email       string
	language    string
	timezone    string
	body        string
}

func init() {
	masterAccountCreateSubAccountCmd.Flags().StringVar(&masterAccountCreateSubAccountFlags.companyName, "company-name", "", "Set the name of the sub-account company")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	masterAccountCreateSubAccountCmd.Flags().StringVar(&masterAccountCreateSubAccountFlags.email, "email", "", "Email address for the organization")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	masterAccountCreateSubAccountCmd.Flags().StringVar(&masterAccountCreateSubAccountFlags.language, "language", "", "Set the language of the sub-account")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	masterAccountCreateSubAccountCmd.Flags().StringVar(&masterAccountCreateSubAccountFlags.timezone, "timezone", "", "Set the timezone of the sub-account")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	masterAccountCreateSubAccountCmd.Flags().StringVar(&masterAccountCreateSubAccountFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	masterAccountCmd.AddCommand(masterAccountCreateSubAccountCmd)
}

func runMasterAccountCreateSubAccount(cmd *cobra.Command, args []string) error {
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
			Name:        "company-name",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Set the name of the sub-account company",
		})
		flags = append(flags, flagSchema{
			Name:        "email",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Email address for the organization",
		})
		flags = append(flags, flagSchema{
			Name:        "language",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Set the language of the sub-account",
		})
		flags = append(flags, flagSchema{
			Name:        "timezone",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Set the timezone of the sub-account",
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
			Description: "Created sub-account ID",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "Bad request",
		})

		schema := map[string]any{
			"command":     "create-sub-account",
			"description": "Create a new sub-account under a master account.",
			"http": map[string]any{
				"method": "POST",
				"path":   "/corporate/subAccount",
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
		Path:        httpclient.SubstitutePath("/corporate/subAccount", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if masterAccountCreateSubAccountFlags.body != "" {
		if err := json.Unmarshal([]byte(masterAccountCreateSubAccountFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("company-name") {
		bodyMap["companyName"] = masterAccountCreateSubAccountFlags.companyName
	}
	if cmd.Flags().Changed("email") {
		bodyMap["email"] = masterAccountCreateSubAccountFlags.email
	}
	if cmd.Flags().Changed("language") {
		bodyMap["language"] = masterAccountCreateSubAccountFlags.language
	}
	if cmd.Flags().Changed("timezone") {
		bodyMap["timezone"] = masterAccountCreateSubAccountFlags.timezone
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
