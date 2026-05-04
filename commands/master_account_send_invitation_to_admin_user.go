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

var masterAccountSendInvitationToAdminUserCmd = &cobra.Command{
	Use:   "send-invitation-to-admin-user",
	Short: "Send invitation to an admin user",
	RunE:  runMasterAccountSendInvitationToAdminUser,
}

var masterAccountSendInvitationToAdminUserFlags struct {
	email             string
	allFeaturesAccess bool
	groupIds          []string
	privileges        []string
	body              string
}

func init() {
	masterAccountSendInvitationToAdminUserCmd.Flags().StringVar(&masterAccountSendInvitationToAdminUserFlags.email, "email", "", "Email address for the organization")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	masterAccountSendInvitationToAdminUserCmd.Flags().BoolVar(&masterAccountSendInvitationToAdminUserFlags.allFeaturesAccess, "all-features-access", false, "All access to the features")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	masterAccountSendInvitationToAdminUserCmd.Flags().StringSliceVar(&masterAccountSendInvitationToAdminUserFlags.groupIds, "group-ids", nil, "Ids of Group")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	masterAccountSendInvitationToAdminUserCmd.Flags().StringSliceVar(&masterAccountSendInvitationToAdminUserFlags.privileges, "privileges", nil, "")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	masterAccountSendInvitationToAdminUserCmd.Flags().StringVar(&masterAccountSendInvitationToAdminUserFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	masterAccountCmd.AddCommand(masterAccountSendInvitationToAdminUserCmd)
}

func runMasterAccountSendInvitationToAdminUser(cmd *cobra.Command, args []string) error {
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
			Description: "Email address for the organization",
		})
		flags = append(flags, flagSchema{
			Name:        "all-features-access",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "All access to the features",
		})
		flags = append(flags, flagSchema{
			Name:        "group-ids",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "Ids of Group",
		})
		flags = append(flags, flagSchema{
			Name:        "privileges",
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
			Status:      "201",
			ContentType: "application/json",
			Description: "Invitation ID",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "Bad request",
		})

		schema := map[string]any{
			"command":     "send-invitation-to-admin-user",
			"description": "Send invitation to an admin user",
			"http": map[string]any{
				"method": "POST",
				"path":   "/corporate/user/invitation/send",
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
		Path:        httpclient.SubstitutePath("/corporate/user/invitation/send", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if masterAccountSendInvitationToAdminUserFlags.body != "" {
		if err := json.Unmarshal([]byte(masterAccountSendInvitationToAdminUserFlags.body), &bodyMap); err != nil {
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
		bodyMap["email"] = masterAccountSendInvitationToAdminUserFlags.email
	}
	if cmd.Flags().Changed("all-features-access") {
		bodyMap["all_features_access"] = masterAccountSendInvitationToAdminUserFlags.allFeaturesAccess
	}
	if cmd.Flags().Changed("group-ids") {
		bodyMap["groupIds"] = masterAccountSendInvitationToAdminUserFlags.groupIds
	}
	if cmd.Flags().Changed("privileges") {
		bodyMap["privileges"] = masterAccountSendInvitationToAdminUserFlags.privileges
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
