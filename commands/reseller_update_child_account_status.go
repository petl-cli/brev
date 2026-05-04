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

var resellerUpdateChildAccountStatusCmd = &cobra.Command{
	Use:   "update-child-account-status",
	Short: "Update info of reseller's child account status based on the identifier supplied",
	RunE:  runResellerUpdateChildAccountStatus,
}

var resellerUpdateChildAccountStatusFlags struct {
	childIdentifier     string
	transactionalEmail  bool
	transactionalSms    bool
	marketingAutomation bool
	smsCampaign         bool
	body                string
}

func init() {
	resellerUpdateChildAccountStatusCmd.Flags().StringVar(&resellerUpdateChildAccountStatusFlags.childIdentifier, "child-identifier", "", "Either auth key or id of reseller's child")
	resellerUpdateChildAccountStatusCmd.MarkFlagRequired("child-identifier")
	resellerUpdateChildAccountStatusCmd.Flags().BoolVar(&resellerUpdateChildAccountStatusFlags.transactionalEmail, "transactional-email", false, "Status of Transactional Email Platform activation for your account (true=enabled, false=disabled)")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	resellerUpdateChildAccountStatusCmd.Flags().BoolVar(&resellerUpdateChildAccountStatusFlags.transactionalSms, "transactional-sms", false, "Status of Transactional SMS Platform activation for your account (true=enabled, false=disabled)")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	resellerUpdateChildAccountStatusCmd.Flags().BoolVar(&resellerUpdateChildAccountStatusFlags.marketingAutomation, "marketing-automation", false, "Status of Marketing Automation Platform activation for your account (true=enabled, false=disabled)")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	resellerUpdateChildAccountStatusCmd.Flags().BoolVar(&resellerUpdateChildAccountStatusFlags.smsCampaign, "sms-campaign", false, "Status of SMS Campaign Platform activation for your account (true=enabled, false=disabled)")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	resellerUpdateChildAccountStatusCmd.Flags().StringVar(&resellerUpdateChildAccountStatusFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	resellerCmd.AddCommand(resellerUpdateChildAccountStatusCmd)
}

func runResellerUpdateChildAccountStatus(cmd *cobra.Command, args []string) error {
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
			Name:        "transactional-email",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Status of Transactional Email Platform activation for your account (true=enabled, false=disabled)",
		})
		flags = append(flags, flagSchema{
			Name:        "transactional-sms",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Status of Transactional SMS Platform activation for your account (true=enabled, false=disabled)",
		})
		flags = append(flags, flagSchema{
			Name:        "marketing-automation",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Status of Marketing Automation Platform activation for your account (true=enabled, false=disabled)",
		})
		flags = append(flags, flagSchema{
			Name:        "sms-campaign",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Status of SMS Campaign Platform activation for your account (true=enabled, false=disabled)",
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
			Description: "reseller's child account status updated",
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
			"command":     "update-child-account-status",
			"description": "Update info of reseller's child account status based on the identifier supplied",
			"http": map[string]any{
				"method": "PUT",
				"path":   "/reseller/children/{childIdentifier}/accountStatus",
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
	pathParams["childIdentifier"] = fmt.Sprintf("%v", resellerUpdateChildAccountStatusFlags.childIdentifier)

	req := &httpclient.Request{
		Method:      "PUT",
		Path:        httpclient.SubstitutePath("/reseller/children/{childIdentifier}/accountStatus", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if resellerUpdateChildAccountStatusFlags.body != "" {
		if err := json.Unmarshal([]byte(resellerUpdateChildAccountStatusFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("transactional-email") {
		bodyMap["transactionalEmail"] = resellerUpdateChildAccountStatusFlags.transactionalEmail
	}
	if cmd.Flags().Changed("transactional-sms") {
		bodyMap["transactionalSms"] = resellerUpdateChildAccountStatusFlags.transactionalSms
	}
	if cmd.Flags().Changed("marketing-automation") {
		bodyMap["marketingAutomation"] = resellerUpdateChildAccountStatusFlags.marketingAutomation
	}
	if cmd.Flags().Changed("sms-campaign") {
		bodyMap["smsCampaign"] = resellerUpdateChildAccountStatusFlags.smsCampaign
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
