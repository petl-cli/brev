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

var contactsImportContactsProcessCmd = &cobra.Command{
	Use:   "import-contacts-process",
	Short: "Import contacts",
	RunE:  runContactsImportContactsProcess,
}

var contactsImportContactsProcessFlags struct {
	fileUrl                 string
	fileBody                string
	jsonBody                []string
	listIds                 []string
	notifyUrl               string
	emailBlacklist          bool
	disableNotification     bool
	smsBlacklist            bool
	updateExistingContacts  bool
	emptyContactsAttributes bool
	body                    string
}

func init() {
	contactsImportContactsProcessCmd.Flags().StringVar(&contactsImportContactsProcessFlags.fileUrl, "file-url", "", "**Mandatory if fileBody and jsonBody is not defined.** URL of the file to be imported (**no local file**). Possible file formats: #### .txt, .csv, .json ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsImportContactsProcessCmd.Flags().StringVar(&contactsImportContactsProcessFlags.fileBody, "file-body", "", "**Mandatory if fileUrl and jsonBody is not defined.** CSV content to be imported. Use semicolon to separate multiple attributes. **Maximum allowed file body size is 10MB** . However we recommend a safe limit of around 8 MB to avoid the issues caused due to increase of file body size while parsing. Please use fileUrl instead to import bigger files. ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsImportContactsProcessCmd.Flags().StringSliceVar(&contactsImportContactsProcessFlags.jsonBody, "json-body", nil, "**Mandatory if fileUrl and fileBody is not defined.** JSON content to be imported. **Maximum allowed json body size is 10MB** . However we recommend a safe limit of around 8 MB to avoid the issues caused due to increase of json body size while parsing. Please use fileUrl instead to import bigger files. ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsImportContactsProcessCmd.Flags().StringSliceVar(&contactsImportContactsProcessFlags.listIds, "list-ids", nil, "**Mandatory if newList is not defined.** Ids of the lists in which the contacts shall be imported. For example, **[2, 4, 7]**. ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsImportContactsProcessCmd.Flags().StringVar(&contactsImportContactsProcessFlags.notifyUrl, "notify-url", "", "URL that will be called once the import process is finished. For reference, https://help.brevo.com/hc/en-us/articles/360007666479")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsImportContactsProcessCmd.Flags().BoolVar(&contactsImportContactsProcessFlags.emailBlacklist, "email-blacklist", false, "To blacklist all the contacts for email")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsImportContactsProcessCmd.Flags().BoolVar(&contactsImportContactsProcessFlags.disableNotification, "disable-notification", false, "To disable email notification")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsImportContactsProcessCmd.Flags().BoolVar(&contactsImportContactsProcessFlags.smsBlacklist, "sms-blacklist", false, "To blacklist all the contacts for sms")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsImportContactsProcessCmd.Flags().BoolVar(&contactsImportContactsProcessFlags.updateExistingContacts, "update-existing-contacts", false, "To facilitate the choice to update the existing contacts")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsImportContactsProcessCmd.Flags().BoolVar(&contactsImportContactsProcessFlags.emptyContactsAttributes, "empty-contacts-attributes", false, "To facilitate the choice to erase any attribute of the existing contacts with empty value. emptyContactsAttributes = true means the empty fields in your import will erase any attribute that currently contain data in Brevo, & emptyContactsAttributes = false means the empty fields will not affect your existing data ( **only available if `updateExistingContacts` set to true **) ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsImportContactsProcessCmd.Flags().StringVar(&contactsImportContactsProcessFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	contactsCmd.AddCommand(contactsImportContactsProcessCmd)
}

func runContactsImportContactsProcess(cmd *cobra.Command, args []string) error {
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
			Name:        "file-url",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "**Mandatory if fileBody and jsonBody is not defined.** URL of the file to be imported (**no local file**). Possible file formats: #### .txt, .csv, .json ",
		})
		flags = append(flags, flagSchema{
			Name:        "file-body",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "**Mandatory if fileUrl and jsonBody is not defined.** CSV content to be imported. Use semicolon to separate multiple attributes. **Maximum allowed file body size is 10MB** . However we recommend a safe limit of around 8 MB to avoid the issues caused due to increase of file body size while parsing. Please use fileUrl instead to import bigger files. ",
		})
		flags = append(flags, flagSchema{
			Name:        "json-body",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "**Mandatory if fileUrl and fileBody is not defined.** JSON content to be imported. **Maximum allowed json body size is 10MB** . However we recommend a safe limit of around 8 MB to avoid the issues caused due to increase of json body size while parsing. Please use fileUrl instead to import bigger files. ",
		})
		flags = append(flags, flagSchema{
			Name:        "list-ids",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "**Mandatory if newList is not defined.** Ids of the lists in which the contacts shall be imported. For example, **[2, 4, 7]**. ",
		})
		flags = append(flags, flagSchema{
			Name:        "notify-url",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "URL that will be called once the import process is finished. For reference, https://help.brevo.com/hc/en-us/articles/360007666479",
		})
		flags = append(flags, flagSchema{
			Name:        "new-list",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "To create a new list and import the contacts into it, pass the listName and an optional folderId.",
		})
		flags = append(flags, flagSchema{
			Name:        "email-blacklist",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "To blacklist all the contacts for email",
		})
		flags = append(flags, flagSchema{
			Name:        "disable-notification",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "To disable email notification",
		})
		flags = append(flags, flagSchema{
			Name:        "sms-blacklist",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "To blacklist all the contacts for sms",
		})
		flags = append(flags, flagSchema{
			Name:        "update-existing-contacts",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "To facilitate the choice to update the existing contacts",
		})
		flags = append(flags, flagSchema{
			Name:        "empty-contacts-attributes",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "To facilitate the choice to erase any attribute of the existing contacts with empty value. emptyContactsAttributes = true means the empty fields in your import will erase any attribute that currently contain data in Brevo, & emptyContactsAttributes = false means the empty fields will not affect your existing data ( **only available if `updateExistingContacts` set to true **) ",
		})

		type responseSchema struct {
			Status      string `json:"status"`
			ContentType string `json:"content_type,omitempty"`
			Description string `json:"description,omitempty"`
		}
		var responses []responseSchema
		responses = append(responses, responseSchema{
			Status:      "202",
			ContentType: "application/json",
			Description: "process id created",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "import-contacts-process",
			"description": "Import contacts",
			"http": map[string]any{
				"method": "POST",
				"path":   "/contacts/import",
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
		Path:        httpclient.SubstitutePath("/contacts/import", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if contactsImportContactsProcessFlags.body != "" {
		if err := json.Unmarshal([]byte(contactsImportContactsProcessFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("file-url") {
		bodyMap["fileUrl"] = contactsImportContactsProcessFlags.fileUrl
	}
	if cmd.Flags().Changed("file-body") {
		bodyMap["fileBody"] = contactsImportContactsProcessFlags.fileBody
	}
	if cmd.Flags().Changed("json-body") {
		bodyMap["jsonBody"] = contactsImportContactsProcessFlags.jsonBody
	}
	if cmd.Flags().Changed("list-ids") {
		bodyMap["listIds"] = contactsImportContactsProcessFlags.listIds
	}
	if cmd.Flags().Changed("notify-url") {
		bodyMap["notifyUrl"] = contactsImportContactsProcessFlags.notifyUrl
	}
	if cmd.Flags().Changed("email-blacklist") {
		bodyMap["emailBlacklist"] = contactsImportContactsProcessFlags.emailBlacklist
	}
	if cmd.Flags().Changed("disable-notification") {
		bodyMap["disableNotification"] = contactsImportContactsProcessFlags.disableNotification
	}
	if cmd.Flags().Changed("sms-blacklist") {
		bodyMap["smsBlacklist"] = contactsImportContactsProcessFlags.smsBlacklist
	}
	if cmd.Flags().Changed("update-existing-contacts") {
		bodyMap["updateExistingContacts"] = contactsImportContactsProcessFlags.updateExistingContacts
	}
	if cmd.Flags().Changed("empty-contacts-attributes") {
		bodyMap["emptyContactsAttributes"] = contactsImportContactsProcessFlags.emptyContactsAttributes
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
