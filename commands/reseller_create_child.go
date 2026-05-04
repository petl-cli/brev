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

var resellerCreateChildCmd = &cobra.Command{
	Use:   "create-child",
	Short: "Creates a reseller child",
	RunE:  runResellerCreateChild,
}

var resellerCreateChildFlags struct {
	email       string
	firstName   string
	lastName    string
	companyName string
	password    string
	language    string
	body        string
}

func init() {
	resellerCreateChildCmd.Flags().StringVar(&resellerCreateChildFlags.email, "email", "", "Email address to create the child account")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	resellerCreateChildCmd.Flags().StringVar(&resellerCreateChildFlags.firstName, "first-name", "", "First name to use to create the child account")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	resellerCreateChildCmd.Flags().StringVar(&resellerCreateChildFlags.lastName, "last-name", "", "Last name to use to create the child account")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	resellerCreateChildCmd.Flags().StringVar(&resellerCreateChildFlags.companyName, "company-name", "", "Company name to use to create the child account")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	resellerCreateChildCmd.Flags().StringVar(&resellerCreateChildFlags.password, "password", "", "Password for the child account to login")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	resellerCreateChildCmd.Flags().StringVar(&resellerCreateChildFlags.language, "language", "", "Language of the child account")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	resellerCreateChildCmd.Flags().StringVar(&resellerCreateChildFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	resellerCmd.AddCommand(resellerCreateChildCmd)
}

func runResellerCreateChild(cmd *cobra.Command, args []string) error {
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
			Name:        "email",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Email address to create the child account",
		})
		flags = append(flags, flagSchema{
			Name:        "first-name",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "First name to use to create the child account",
		})
		flags = append(flags, flagSchema{
			Name:        "last-name",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Last name to use to create the child account",
		})
		flags = append(flags, flagSchema{
			Name:        "company-name",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Company name to use to create the child account",
		})
		flags = append(flags, flagSchema{
			Name:        "password",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Password for the child account to login",
		})
		flags = append(flags, flagSchema{
			Name:        "language",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Language of the child account",
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
			Description: "child created",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "Bad request",
		})
		responses = append(responses, responseSchema{
			Status:      "403",
			ContentType: "application/json",
			Description: "Current account is not a reseller",
		})

		schema := map[string]any{
			"command":     "create-child",
			"description": "Creates a reseller child",
			"http": map[string]any{
				"method": "POST",
				"path":   "/reseller/children",
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
		Path:        httpclient.SubstitutePath("/reseller/children", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if resellerCreateChildFlags.body != "" {
		if err := json.Unmarshal([]byte(resellerCreateChildFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("email") {
		bodyMap["email"] = resellerCreateChildFlags.email
	}
	if cmd.Flags().Changed("first-name") {
		bodyMap["firstName"] = resellerCreateChildFlags.firstName
	}
	if cmd.Flags().Changed("last-name") {
		bodyMap["lastName"] = resellerCreateChildFlags.lastName
	}
	if cmd.Flags().Changed("company-name") {
		bodyMap["companyName"] = resellerCreateChildFlags.companyName
	}
	if cmd.Flags().Changed("password") {
		bodyMap["password"] = resellerCreateChildFlags.password
	}
	if cmd.Flags().Changed("language") {
		bodyMap["language"] = resellerCreateChildFlags.language
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
