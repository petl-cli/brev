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

var masterAccountResendCancelAdminUserInvitationCmd = &cobra.Command{
	Use:   "resend-cancel-admin-user-invitation",
	Short: "Resend / cancel admin user invitation",
	RunE:  runMasterAccountResendCancelAdminUserInvitation,
}

var masterAccountResendCancelAdminUserInvitationFlags struct {
	action string
	email  string
}

func init() {
	masterAccountResendCancelAdminUserInvitationCmd.Flags().StringVar(&masterAccountResendCancelAdminUserInvitationFlags.action, "action", "", "Action to be performed (cancel / resend)")
	masterAccountResendCancelAdminUserInvitationCmd.MarkFlagRequired("action")
	masterAccountResendCancelAdminUserInvitationCmd.Flags().StringVar(&masterAccountResendCancelAdminUserInvitationFlags.email, "email", "", "Email address of the recipient")
	masterAccountResendCancelAdminUserInvitationCmd.MarkFlagRequired("email")

	masterAccountCmd.AddCommand(masterAccountResendCancelAdminUserInvitationCmd)
}

func runMasterAccountResendCancelAdminUserInvitation(cmd *cobra.Command, args []string) error {
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
			Name:        "action",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "Action to be performed (cancel / resend)",
		})
		flags = append(flags, flagSchema{
			Name:        "email",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "Email address of the recipient",
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
			Description: "Response of the action performed",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "resend-cancel-admin-user-invitation",
			"description": "Resend / cancel admin user invitation",
			"http": map[string]any{
				"method": "PUT",
				"path":   "/corporate/user/invitation/{action}/{email}",
			},
			"input": map[string]any{
				"flags":         flags,
				"body_flag":     false,
				"body_required": false,
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
	pathParams["action"] = fmt.Sprintf("%v", masterAccountResendCancelAdminUserInvitationFlags.action)
	pathParams["email"] = fmt.Sprintf("%v", masterAccountResendCancelAdminUserInvitationFlags.email)

	req := &httpclient.Request{
		Method:      "PUT",
		Path:        httpclient.SubstitutePath("/corporate/user/invitation/{action}/{email}", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

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
