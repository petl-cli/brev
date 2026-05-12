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

var masterAccountEnableDisableCmd = &cobra.Command{
	Use:   "enable-disable",
	Short: "Enable/disable sub-account application(s)",
	RunE:  runMasterAccountEnableDisable,
}

var masterAccountEnableDisableFlags struct {
	id                  int
	inbox               bool
	whatsapp            bool
	automation          bool
	emailCampaigns      bool
	smsCampaigns        bool
	landingPages        bool
	transactionalEmails bool
	transactionalSms    bool
	facebookAds         bool
	webPush             bool
	meetings            bool
	conversations       bool
	crm                 bool
	body                string
}

func init() {
	masterAccountEnableDisableCmd.Flags().IntVar(&masterAccountEnableDisableFlags.id, "id", 0, "Id of the sub-account organization (mandatory)")
	masterAccountEnableDisableCmd.MarkFlagRequired("id")
	masterAccountEnableDisableCmd.Flags().BoolVar(&masterAccountEnableDisableFlags.inbox, "inbox", false, "Set this field to enable or disable Inbox on the sub-account / Not applicable on ENTv2")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	masterAccountEnableDisableCmd.Flags().BoolVar(&masterAccountEnableDisableFlags.whatsapp, "whatsapp", false, "Set this field to enable or disable Whatsapp campaigns on the sub-account")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	masterAccountEnableDisableCmd.Flags().BoolVar(&masterAccountEnableDisableFlags.automation, "automation", false, "Set this field to enable or disable Automation on the sub-account")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	masterAccountEnableDisableCmd.Flags().BoolVar(&masterAccountEnableDisableFlags.emailCampaigns, "email-campaigns", false, "Set this field to enable or disable Email Campaigns on the sub-account")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	masterAccountEnableDisableCmd.Flags().BoolVar(&masterAccountEnableDisableFlags.smsCampaigns, "sms-campaigns", false, "Set this field to enable or disable SMS Marketing on the sub-account")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	masterAccountEnableDisableCmd.Flags().BoolVar(&masterAccountEnableDisableFlags.landingPages, "landing-pages", false, "Set this field to enable or disable Landing pages on the sub-account")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	masterAccountEnableDisableCmd.Flags().BoolVar(&masterAccountEnableDisableFlags.transactionalEmails, "transactional-emails", false, "Set this field to enable or disable Transactional Email on the sub-account")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	masterAccountEnableDisableCmd.Flags().BoolVar(&masterAccountEnableDisableFlags.transactionalSms, "transactional-sms", false, "Set this field to enable or disable Transactional SMS on the sub-account")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	masterAccountEnableDisableCmd.Flags().BoolVar(&masterAccountEnableDisableFlags.facebookAds, "facebook-ads", false, "Set this field to enable or disable Facebook ads on the sub-account")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	masterAccountEnableDisableCmd.Flags().BoolVar(&masterAccountEnableDisableFlags.webPush, "web-push", false, "Set this field to enable or disable Web Push on the sub-account")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	masterAccountEnableDisableCmd.Flags().BoolVar(&masterAccountEnableDisableFlags.meetings, "meetings", false, "Set this field to enable or disable Meetings on the sub-account")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	masterAccountEnableDisableCmd.Flags().BoolVar(&masterAccountEnableDisableFlags.conversations, "conversations", false, "Set this field to enable or disable Conversations on the sub-account")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	masterAccountEnableDisableCmd.Flags().BoolVar(&masterAccountEnableDisableFlags.crm, "crm", false, "Set this field to enable or disable Sales CRM on the sub-account")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	masterAccountEnableDisableCmd.Flags().StringVar(&masterAccountEnableDisableFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	masterAccountCmd.AddCommand(masterAccountEnableDisableCmd)
}

func runMasterAccountEnableDisable(cmd *cobra.Command, args []string) error {
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
			Location:    "path",
			Description: "Id of the sub-account organization (mandatory)",
		})
		flags = append(flags, flagSchema{
			Name:        "inbox",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Set this field to enable or disable Inbox on the sub-account / Not applicable on ENTv2",
		})
		flags = append(flags, flagSchema{
			Name:        "whatsapp",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Set this field to enable or disable Whatsapp campaigns on the sub-account",
		})
		flags = append(flags, flagSchema{
			Name:        "automation",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Set this field to enable or disable Automation on the sub-account",
		})
		flags = append(flags, flagSchema{
			Name:        "email-campaigns",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Set this field to enable or disable Email Campaigns on the sub-account",
		})
		flags = append(flags, flagSchema{
			Name:        "sms-campaigns",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Set this field to enable or disable SMS Marketing on the sub-account",
		})
		flags = append(flags, flagSchema{
			Name:        "landing-pages",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Set this field to enable or disable Landing pages on the sub-account",
		})
		flags = append(flags, flagSchema{
			Name:        "transactional-emails",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Set this field to enable or disable Transactional Email on the sub-account",
		})
		flags = append(flags, flagSchema{
			Name:        "transactional-sms",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Set this field to enable or disable Transactional SMS on the sub-account",
		})
		flags = append(flags, flagSchema{
			Name:        "facebook-ads",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Set this field to enable or disable Facebook ads on the sub-account",
		})
		flags = append(flags, flagSchema{
			Name:        "web-push",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Set this field to enable or disable Web Push on the sub-account",
		})
		flags = append(flags, flagSchema{
			Name:        "meetings",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Set this field to enable or disable Meetings on the sub-account",
		})
		flags = append(flags, flagSchema{
			Name:        "conversations",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Set this field to enable or disable Conversations on the sub-account",
		})
		flags = append(flags, flagSchema{
			Name:        "crm",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Set this field to enable or disable Sales CRM on the sub-account",
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
			Description: "Sub-account application(s) enabled/disabled",
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
			"command":     "enable-disable",
			"description": "Enable/disable sub-account application(s)",
			"http": map[string]any{
				"method": "PUT",
				"path":   "/corporate/subAccount/{id}/applications/toggle",
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
	pathParams["id"] = fmt.Sprintf("%v", masterAccountEnableDisableFlags.id)

	req := &httpclient.Request{
		Method:      "PUT",
		Path:        httpclient.SubstitutePath("/corporate/subAccount/{id}/applications/toggle", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if masterAccountEnableDisableFlags.body != "" {
		if err := json.Unmarshal([]byte(masterAccountEnableDisableFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("inbox") {
		bodyMap["inbox"] = masterAccountEnableDisableFlags.inbox
	}
	if cmd.Flags().Changed("whatsapp") {
		bodyMap["whatsapp"] = masterAccountEnableDisableFlags.whatsapp
	}
	if cmd.Flags().Changed("automation") {
		bodyMap["automation"] = masterAccountEnableDisableFlags.automation
	}
	if cmd.Flags().Changed("email-campaigns") {
		bodyMap["email-campaigns"] = masterAccountEnableDisableFlags.emailCampaigns
	}
	if cmd.Flags().Changed("sms-campaigns") {
		bodyMap["sms-campaigns"] = masterAccountEnableDisableFlags.smsCampaigns
	}
	if cmd.Flags().Changed("landing-pages") {
		bodyMap["landing-pages"] = masterAccountEnableDisableFlags.landingPages
	}
	if cmd.Flags().Changed("transactional-emails") {
		bodyMap["transactional-emails"] = masterAccountEnableDisableFlags.transactionalEmails
	}
	if cmd.Flags().Changed("transactional-sms") {
		bodyMap["transactional-sms"] = masterAccountEnableDisableFlags.transactionalSms
	}
	if cmd.Flags().Changed("facebook-ads") {
		bodyMap["facebook-ads"] = masterAccountEnableDisableFlags.facebookAds
	}
	if cmd.Flags().Changed("web-push") {
		bodyMap["web-push"] = masterAccountEnableDisableFlags.webPush
	}
	if cmd.Flags().Changed("meetings") {
		bodyMap["meetings"] = masterAccountEnableDisableFlags.meetings
	}
	if cmd.Flags().Changed("conversations") {
		bodyMap["conversations"] = masterAccountEnableDisableFlags.conversations
	}
	if cmd.Flags().Changed("crm") {
		bodyMap["crm"] = masterAccountEnableDisableFlags.crm
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
