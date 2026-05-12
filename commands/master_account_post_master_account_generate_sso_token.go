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

var masterAccountPostMasterAccountGenerateSsoTokenCmd = &cobra.Command{
	Use:   "post-master-account-generate-sso-token",
	Short: "Generate SSO token to access sub-account",
	RunE:  runMasterAccountPostMasterAccountGenerateSsoToken,
}

var masterAccountPostMasterAccountGenerateSsoTokenFlags struct {
	id     int
	email  string
	target string
	url    string
	body   string
}

func init() {
	masterAccountPostMasterAccountGenerateSsoTokenCmd.Flags().IntVar(&masterAccountPostMasterAccountGenerateSsoTokenFlags.id, "id", 0, "Id of the sub-account organization")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	masterAccountPostMasterAccountGenerateSsoTokenCmd.Flags().StringVar(&masterAccountPostMasterAccountGenerateSsoTokenFlags.email, "email", "", "User email of sub-account organization")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	masterAccountPostMasterAccountGenerateSsoTokenCmd.Flags().StringVar(&masterAccountPostMasterAccountGenerateSsoTokenFlags.target, "target", "", "**Set target after login success** * **automation** - Redirect to Automation after login * **email_campaign** - Redirect to Email Campaign after login * **contacts** - Redirect to Contacts after login * **landing_pages** - Redirect to Landing Pages after login * **email_transactional** - Redirect to Email Transactional after login * **senders** - Redirect to Senders after login * **sms_campaign** - Redirect to Sms Campaign after login * **sms_transactional** - Redirect to Sms Transactional after login ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	masterAccountPostMasterAccountGenerateSsoTokenCmd.Flags().StringVar(&masterAccountPostMasterAccountGenerateSsoTokenFlags.url, "url", "", "Set the full target URL after login success. The user will land directly on this target URL after login")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	masterAccountPostMasterAccountGenerateSsoTokenCmd.Flags().StringVar(&masterAccountPostMasterAccountGenerateSsoTokenFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	masterAccountCmd.AddCommand(masterAccountPostMasterAccountGenerateSsoTokenCmd)
}

func runMasterAccountPostMasterAccountGenerateSsoToken(cmd *cobra.Command, args []string) error {
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
			Name:        "id",
			Type:        "integer",
			Required:    true,
			Location:    "body",
			Description: "Id of the sub-account organization",
		})
		flags = append(flags, flagSchema{
			Name:        "email",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "User email of sub-account organization",
		})
		flags = append(flags, flagSchema{
			Name:        "target",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "**Set target after login success** * **automation** - Redirect to Automation after login * **email_campaign** - Redirect to Email Campaign after login * **contacts** - Redirect to Contacts after login * **landing_pages** - Redirect to Landing Pages after login * **email_transactional** - Redirect to Email Transactional after login * **senders** - Redirect to Senders after login * **sms_campaign** - Redirect to Sms Campaign after login * **sms_transactional** - Redirect to Sms Transactional after login ",
		})
		flags = append(flags, flagSchema{
			Name:        "url",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Set the full target URL after login success. The user will land directly on this target URL after login",
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
			Description: "Session token",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})
		responses = append(responses, responseSchema{
			Status:      "403",
			ContentType: "application/json",
			Description: "Current account is not a master account",
		})

		schema := map[string]any{
			"command":     "post-master-account-generate-sso-token",
			"description": "Generate SSO token to access sub-account",
			"http": map[string]any{
				"method": "POST",
				"path":   "/corporate/subAccount/ssoToken",
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
		Path:        httpclient.SubstitutePath("/corporate/subAccount/ssoToken", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if masterAccountPostMasterAccountGenerateSsoTokenFlags.body != "" {
		if err := json.Unmarshal([]byte(masterAccountPostMasterAccountGenerateSsoTokenFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("id") {
		bodyMap["id"] = masterAccountPostMasterAccountGenerateSsoTokenFlags.id
	}
	if cmd.Flags().Changed("email") {
		bodyMap["email"] = masterAccountPostMasterAccountGenerateSsoTokenFlags.email
	}
	if cmd.Flags().Changed("target") {
		bodyMap["target"] = masterAccountPostMasterAccountGenerateSsoTokenFlags.target
	}
	if cmd.Flags().Changed("url") {
		bodyMap["url"] = masterAccountPostMasterAccountGenerateSsoTokenFlags.url
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
