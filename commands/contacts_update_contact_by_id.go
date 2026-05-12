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

var contactsUpdateContactByIdCmd = &cobra.Command{
	Use:   "update-contact-by-id",
	Short: "Update a contact",
	RunE:  runContactsUpdateContactById,
}

var contactsUpdateContactByIdFlags struct {
	identifier          string
	extId               string
	emailBlacklisted    bool
	smsBlacklisted      bool
	listIds             []string
	unlinkListIds       []string
	smtpBlacklistSender []string
	body                string
}

func init() {
	contactsUpdateContactByIdCmd.Flags().StringVar(&contactsUpdateContactByIdFlags.identifier, "identifier", "", "Email (urlencoded) OR ID of the contact")
	contactsUpdateContactByIdCmd.MarkFlagRequired("identifier")
	contactsUpdateContactByIdCmd.Flags().StringVar(&contactsUpdateContactByIdFlags.extId, "ext-id", "", "Pass your own Id to update ext_id of a contact.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsUpdateContactByIdCmd.Flags().BoolVar(&contactsUpdateContactByIdFlags.emailBlacklisted, "email-blacklisted", false, "Set/unset this field to blacklist/allow the contact for emails (emailBlacklisted = true)")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsUpdateContactByIdCmd.Flags().BoolVar(&contactsUpdateContactByIdFlags.smsBlacklisted, "sms-blacklisted", false, "Set/unset this field to blacklist/allow the contact for SMS (smsBlacklisted = true)")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsUpdateContactByIdCmd.Flags().StringSliceVar(&contactsUpdateContactByIdFlags.listIds, "list-ids", nil, "Ids of the lists to add the contact to")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsUpdateContactByIdCmd.Flags().StringSliceVar(&contactsUpdateContactByIdFlags.unlinkListIds, "unlink-list-ids", nil, "Ids of the lists to remove the contact from")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsUpdateContactByIdCmd.Flags().StringSliceVar(&contactsUpdateContactByIdFlags.smtpBlacklistSender, "smtp-blacklist-sender", nil, "transactional email forbidden sender for contact. Use only for email Contact")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsUpdateContactByIdCmd.Flags().StringVar(&contactsUpdateContactByIdFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	contactsCmd.AddCommand(contactsUpdateContactByIdCmd)
}

func runContactsUpdateContactById(cmd *cobra.Command, args []string) error {
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
			Name:        "identifier",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "Email (urlencoded) OR ID of the contact",
		})
		flags = append(flags, flagSchema{
			Name:        "attributes",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "Pass the set of attributes to be updated. **These attributes must be present in your account**. To update existing email address of a contact with the new one please pass EMAIL in attributes. For example, **{ \"EMAIL\":\"newemail@domain.com\", \"FNAME\":\"Ellie\", \"LNAME\":\"Roger\"}**. The attribute's parameter should be passed in capital letter while updating a contact. Values that don't match the attribute type (e.g. text or string in a date attribute) will be ignored. Keep in mind transactional attributes can be updated the same way as normal attributes. Mobile Number in **SMS** field should be passed with proper country code. For example: **{\"SMS\":\"+91xxxxxxxxxx\"} or {\"SMS\":\"0091xxxxxxxxxx\"}** ",
		})
		flags = append(flags, flagSchema{
			Name:        "ext-id",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Pass your own Id to update ext_id of a contact.",
		})
		flags = append(flags, flagSchema{
			Name:        "email-blacklisted",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Set/unset this field to blacklist/allow the contact for emails (emailBlacklisted = true)",
		})
		flags = append(flags, flagSchema{
			Name:        "sms-blacklisted",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Set/unset this field to blacklist/allow the contact for SMS (smsBlacklisted = true)",
		})
		flags = append(flags, flagSchema{
			Name:        "list-ids",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "Ids of the lists to add the contact to",
		})
		flags = append(flags, flagSchema{
			Name:        "unlink-list-ids",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "Ids of the lists to remove the contact from",
		})
		flags = append(flags, flagSchema{
			Name:        "smtp-blacklist-sender",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "transactional email forbidden sender for contact. Use only for email Contact",
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
			Description: "Contact updated",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})
		responses = append(responses, responseSchema{
			Status:      "404",
			ContentType: "application/json",
			Description: "Contact's email not found",
		})

		schema := map[string]any{
			"command":     "update-contact-by-id",
			"description": "Update a contact",
			"http": map[string]any{
				"method": "PUT",
				"path":   "/contacts/{identifier}",
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
	pathParams["identifier"] = fmt.Sprintf("%v", contactsUpdateContactByIdFlags.identifier)

	req := &httpclient.Request{
		Method:      "PUT",
		Path:        httpclient.SubstitutePath("/contacts/{identifier}", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if contactsUpdateContactByIdFlags.body != "" {
		if err := json.Unmarshal([]byte(contactsUpdateContactByIdFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("ext-id") {
		bodyMap["ext_id"] = contactsUpdateContactByIdFlags.extId
	}
	if cmd.Flags().Changed("email-blacklisted") {
		bodyMap["emailBlacklisted"] = contactsUpdateContactByIdFlags.emailBlacklisted
	}
	if cmd.Flags().Changed("sms-blacklisted") {
		bodyMap["smsBlacklisted"] = contactsUpdateContactByIdFlags.smsBlacklisted
	}
	if cmd.Flags().Changed("list-ids") {
		bodyMap["listIds"] = contactsUpdateContactByIdFlags.listIds
	}
	if cmd.Flags().Changed("unlink-list-ids") {
		bodyMap["unlinkListIds"] = contactsUpdateContactByIdFlags.unlinkListIds
	}
	if cmd.Flags().Changed("smtp-blacklist-sender") {
		bodyMap["smtpBlacklistSender"] = contactsUpdateContactByIdFlags.smtpBlacklistSender
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
