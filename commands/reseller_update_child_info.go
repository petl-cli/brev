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

var resellerUpdateChildInfoCmd = &cobra.Command{
	Use:   "update-child-info",
	Short: "Update info of reseller's child based on the child identifier supplied",
	RunE:  runResellerUpdateChildInfo,
}

var resellerUpdateChildInfoFlags struct {
	childIdentifier string
	email           string
	firstName       string
	lastName        string
	companyName     string
	password        string
	body            string
}

func init() {
	resellerUpdateChildInfoCmd.Flags().StringVar(&resellerUpdateChildInfoFlags.childIdentifier, "child-identifier", "", "Either auth key or id of reseller's child")
	resellerUpdateChildInfoCmd.MarkFlagRequired("child-identifier")
	resellerUpdateChildInfoCmd.Flags().StringVar(&resellerUpdateChildInfoFlags.email, "email", "", "New Email address to update the child account")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	resellerUpdateChildInfoCmd.Flags().StringVar(&resellerUpdateChildInfoFlags.firstName, "first-name", "", "New First name to use to update the child account")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	resellerUpdateChildInfoCmd.Flags().StringVar(&resellerUpdateChildInfoFlags.lastName, "last-name", "", "New Last name to use to update the child account")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	resellerUpdateChildInfoCmd.Flags().StringVar(&resellerUpdateChildInfoFlags.companyName, "company-name", "", "New Company name to use to update the child account")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	resellerUpdateChildInfoCmd.Flags().StringVar(&resellerUpdateChildInfoFlags.password, "password", "", "New password for the child account to login")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	resellerUpdateChildInfoCmd.Flags().StringVar(&resellerUpdateChildInfoFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	resellerCmd.AddCommand(resellerUpdateChildInfoCmd)
}

func runResellerUpdateChildInfo(cmd *cobra.Command, args []string) error {
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
			Name:        "email",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "New Email address to update the child account",
		})
		flags = append(flags, flagSchema{
			Name:        "first-name",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "New First name to use to update the child account",
		})
		flags = append(flags, flagSchema{
			Name:        "last-name",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "New Last name to use to update the child account",
		})
		flags = append(flags, flagSchema{
			Name:        "company-name",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "New Company name to use to update the child account",
		})
		flags = append(flags, flagSchema{
			Name:        "password",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "New password for the child account to login",
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
			Description: "reseller's child updated",
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
			"command":     "update-child-info",
			"description": "Update info of reseller's child based on the child identifier supplied",
			"http": map[string]any{
				"method": "PUT",
				"path":   "/reseller/children/{childIdentifier}",
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
	pathParams["childIdentifier"] = fmt.Sprintf("%v", resellerUpdateChildInfoFlags.childIdentifier)

	req := &httpclient.Request{
		Method:      "PUT",
		Path:        httpclient.SubstitutePath("/reseller/children/{childIdentifier}", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if resellerUpdateChildInfoFlags.body != "" {
		if err := json.Unmarshal([]byte(resellerUpdateChildInfoFlags.body), &bodyMap); err != nil {
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
		bodyMap["email"] = resellerUpdateChildInfoFlags.email
	}
	if cmd.Flags().Changed("first-name") {
		bodyMap["firstName"] = resellerUpdateChildInfoFlags.firstName
	}
	if cmd.Flags().Changed("last-name") {
		bodyMap["lastName"] = resellerUpdateChildInfoFlags.lastName
	}
	if cmd.Flags().Changed("company-name") {
		bodyMap["companyName"] = resellerUpdateChildInfoFlags.companyName
	}
	if cmd.Flags().Changed("password") {
		bodyMap["password"] = resellerUpdateChildInfoFlags.password
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
