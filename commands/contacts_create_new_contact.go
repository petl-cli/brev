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

var contactsCreateNewContactCmd = &cobra.Command{
	Use:   "create-new-contact",
	Short: "Create a contact",
	RunE:  runContactsCreateNewContact,
}

var contactsCreateNewContactFlags struct {
	email               string
	extId               string
	emailBlacklisted    bool
	smsBlacklisted      bool
	listIds             []string
	updateEnabled       bool
	smtpBlacklistSender []string
	body                string
}

func init() {
	contactsCreateNewContactCmd.Flags().StringVar(&contactsCreateNewContactFlags.email, "email", "", "Email address of the user. **Mandatory if \"SMS\" field is not passed in \"attributes\" parameter**. Mobile Number in **SMS** field should be passed with proper country code. For example: **{\"SMS\":\"+91xxxxxxxxxx\"}** or **{\"SMS\":\"0091xxxxxxxxxx\"}** ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsCreateNewContactCmd.Flags().StringVar(&contactsCreateNewContactFlags.extId, "ext-id", "", "Pass your own Id to create a contact.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsCreateNewContactCmd.Flags().BoolVar(&contactsCreateNewContactFlags.emailBlacklisted, "email-blacklisted", false, "Set this field to blacklist the contact for emails (emailBlacklisted = true)")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsCreateNewContactCmd.Flags().BoolVar(&contactsCreateNewContactFlags.smsBlacklisted, "sms-blacklisted", false, "Set this field to blacklist the contact for SMS (smsBlacklisted = true)")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsCreateNewContactCmd.Flags().StringSliceVar(&contactsCreateNewContactFlags.listIds, "list-ids", nil, "Ids of the lists to add the contact to")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsCreateNewContactCmd.Flags().BoolVar(&contactsCreateNewContactFlags.updateEnabled, "update-enabled", false, "Facilitate to update the existing contact in the same request (updateEnabled = true)")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsCreateNewContactCmd.Flags().StringSliceVar(&contactsCreateNewContactFlags.smtpBlacklistSender, "smtp-blacklist-sender", nil, "transactional email forbidden sender for contact. Use only for email Contact ( only available if updateEnabled = true )")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsCreateNewContactCmd.Flags().StringVar(&contactsCreateNewContactFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	contactsCmd.AddCommand(contactsCreateNewContactCmd)
}

func runContactsCreateNewContact(cmd *cobra.Command, args []string) error {
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
			Description: "Email address of the user. **Mandatory if \"SMS\" field is not passed in \"attributes\" parameter**. Mobile Number in **SMS** field should be passed with proper country code. For example: **{\"SMS\":\"+91xxxxxxxxxx\"}** or **{\"SMS\":\"0091xxxxxxxxxx\"}** ",
		})
		flags = append(flags, flagSchema{
			Name:        "ext-id",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Pass your own Id to create a contact.",
		})
		flags = append(flags, flagSchema{
			Name:        "attributes",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "Pass the set of attributes and their values. The attribute's parameter should be passed in capital letter while creating a contact. Values that don't match the attribute type (e.g. text or string in a date attribute) will be ignored. **These attributes must be present in your Brevo account.**. For eg: **{\"FNAME\":\"Elly\", \"LNAME\":\"Roger\"}** ",
		})
		flags = append(flags, flagSchema{
			Name:        "email-blacklisted",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Set this field to blacklist the contact for emails (emailBlacklisted = true)",
		})
		flags = append(flags, flagSchema{
			Name:        "sms-blacklisted",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Set this field to blacklist the contact for SMS (smsBlacklisted = true)",
		})
		flags = append(flags, flagSchema{
			Name:        "list-ids",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "Ids of the lists to add the contact to",
		})
		flags = append(flags, flagSchema{
			Name:        "update-enabled",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Facilitate to update the existing contact in the same request (updateEnabled = true)",
		})
		flags = append(flags, flagSchema{
			Name:        "smtp-blacklist-sender",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "transactional email forbidden sender for contact. Use only for email Contact ( only available if updateEnabled = true )",
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
			Description: "Contact created",
		})
		responses = append(responses, responseSchema{
			Status:      "204",
			ContentType: "",
			Description: "Contact updated",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "create-new-contact",
			"description": "Create a contact",
			"http": map[string]any{
				"method": "POST",
				"path":   "/contacts",
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
		Path:        httpclient.SubstitutePath("/contacts", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if contactsCreateNewContactFlags.body != "" {
		if err := json.Unmarshal([]byte(contactsCreateNewContactFlags.body), &bodyMap); err != nil {
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
		bodyMap["email"] = contactsCreateNewContactFlags.email
	}
	if cmd.Flags().Changed("ext-id") {
		bodyMap["ext_id"] = contactsCreateNewContactFlags.extId
	}
	if cmd.Flags().Changed("email-blacklisted") {
		bodyMap["emailBlacklisted"] = contactsCreateNewContactFlags.emailBlacklisted
	}
	if cmd.Flags().Changed("sms-blacklisted") {
		bodyMap["smsBlacklisted"] = contactsCreateNewContactFlags.smsBlacklisted
	}
	if cmd.Flags().Changed("list-ids") {
		bodyMap["listIds"] = contactsCreateNewContactFlags.listIds
	}
	if cmd.Flags().Changed("update-enabled") {
		bodyMap["updateEnabled"] = contactsCreateNewContactFlags.updateEnabled
	}
	if cmd.Flags().Changed("smtp-blacklist-sender") {
		bodyMap["smtpBlacklistSender"] = contactsCreateNewContactFlags.smtpBlacklistSender
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
