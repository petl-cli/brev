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

var transactionalEmailsUpdateTemplateCmd = &cobra.Command{
	Use:   "update-template",
	Short: "Update an email template",
	RunE:  runTransactionalEmailsUpdateTemplate,
}

var transactionalEmailsUpdateTemplateFlags struct {
	templateId    int
	tag           string
	templateName  string
	htmlContent   string
	htmlUrl       string
	subject       string
	replyTo       string
	toField       string
	attachmentUrl string
	isActive      bool
	body          string
}

func init() {
	transactionalEmailsUpdateTemplateCmd.Flags().IntVar(&transactionalEmailsUpdateTemplateFlags.templateId, "template-id", 0, "id of the template")
	transactionalEmailsUpdateTemplateCmd.MarkFlagRequired("template-id")
	transactionalEmailsUpdateTemplateCmd.Flags().StringVar(&transactionalEmailsUpdateTemplateFlags.tag, "tag", "", "Tag of the template")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsUpdateTemplateCmd.Flags().StringVar(&transactionalEmailsUpdateTemplateFlags.templateName, "template-name", "", "Name of the template")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsUpdateTemplateCmd.Flags().StringVar(&transactionalEmailsUpdateTemplateFlags.htmlContent, "html-content", "", "**Required if htmlUrl is empty**. If the template is designed using Drag & Drop editor via HTML content, then the design page will not have Drag & Drop editor access for that template. Body of the message (HTML must have more than 10 characters) ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsUpdateTemplateCmd.Flags().StringVar(&transactionalEmailsUpdateTemplateFlags.htmlUrl, "html-url", "", "**Required if htmlContent is empty**. URL to the body of the email (HTML) ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsUpdateTemplateCmd.Flags().StringVar(&transactionalEmailsUpdateTemplateFlags.subject, "subject", "", "Subject of the email")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsUpdateTemplateCmd.Flags().StringVar(&transactionalEmailsUpdateTemplateFlags.replyTo, "reply-to", "", "Email on which campaign recipients will be able to reply to")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsUpdateTemplateCmd.Flags().StringVar(&transactionalEmailsUpdateTemplateFlags.toField, "to-field", "", "To personalize the **To** Field. If you want to include the first name and last name of your recipient, add **{FNAME} {LNAME}**. These contact attributes must already exist in your Brevo account. If input parameter **params** used please use **{{contact.FNAME}} {{contact.LNAME}}** for personalization ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsUpdateTemplateCmd.Flags().StringVar(&transactionalEmailsUpdateTemplateFlags.attachmentUrl, "attachment-url", "", "Absolute url of the attachment (**no local file**). Extensions allowed: #### xlsx, xls, ods, docx, docm, doc, csv, pdf, txt, gif, jpg, jpeg, png, tif, tiff, rtf, bmp, cgm, css, shtml, html, htm, zip, xml, ppt, pptx, tar, ez, ics, mobi, msg, pub and eps ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsUpdateTemplateCmd.Flags().BoolVar(&transactionalEmailsUpdateTemplateFlags.isActive, "is-active", false, "Status of the template. isActive = false means template is inactive, isActive = true means template is active")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsUpdateTemplateCmd.Flags().StringVar(&transactionalEmailsUpdateTemplateFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	transactionalEmailsCmd.AddCommand(transactionalEmailsUpdateTemplateCmd)
}

func runTransactionalEmailsUpdateTemplate(cmd *cobra.Command, args []string) error {
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
			Name:        "template-id",
			Type:        "integer",
			Required:    true,
			Location:    "path",
			Description: "id of the template",
		})
		flags = append(flags, flagSchema{
			Name:        "tag",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Tag of the template",
		})
		flags = append(flags, flagSchema{
			Name:        "sender",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "Sender details including id or email and name (_optional_). Only one of either Sender's email or Sender's ID shall be passed in one request at a time. For example: **{\"name\":\"xyz\", \"email\":\"example@abc.com\"}** **{\"name\":\"xyz\", \"id\":123}** ",
		})
		flags = append(flags, flagSchema{
			Name:        "template-name",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Name of the template",
		})
		flags = append(flags, flagSchema{
			Name:        "html-content",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "**Required if htmlUrl is empty**. If the template is designed using Drag & Drop editor via HTML content, then the design page will not have Drag & Drop editor access for that template. Body of the message (HTML must have more than 10 characters) ",
		})
		flags = append(flags, flagSchema{
			Name:        "html-url",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "**Required if htmlContent is empty**. URL to the body of the email (HTML) ",
		})
		flags = append(flags, flagSchema{
			Name:        "subject",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Subject of the email",
		})
		flags = append(flags, flagSchema{
			Name:        "reply-to",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Email on which campaign recipients will be able to reply to",
		})
		flags = append(flags, flagSchema{
			Name:        "to-field",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "To personalize the **To** Field. If you want to include the first name and last name of your recipient, add **{FNAME} {LNAME}**. These contact attributes must already exist in your Brevo account. If input parameter **params** used please use **{{contact.FNAME}} {{contact.LNAME}}** for personalization ",
		})
		flags = append(flags, flagSchema{
			Name:        "attachment-url",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Absolute url of the attachment (**no local file**). Extensions allowed: #### xlsx, xls, ods, docx, docm, doc, csv, pdf, txt, gif, jpg, jpeg, png, tif, tiff, rtf, bmp, cgm, css, shtml, html, htm, zip, xml, ppt, pptx, tar, ez, ics, mobi, msg, pub and eps ",
		})
		flags = append(flags, flagSchema{
			Name:        "is-active",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Status of the template. isActive = false means template is inactive, isActive = true means template is active",
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
			Description: "transactional email template updated",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})
		responses = append(responses, responseSchema{
			Status:      "404",
			ContentType: "application/json",
			Description: "Template ID not found",
		})

		schema := map[string]any{
			"command":     "update-template",
			"description": "Update an email template",
			"http": map[string]any{
				"method": "PUT",
				"path":   "/smtp/templates/{templateId}",
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
	pathParams["templateId"] = fmt.Sprintf("%v", transactionalEmailsUpdateTemplateFlags.templateId)

	req := &httpclient.Request{
		Method:      "PUT",
		Path:        httpclient.SubstitutePath("/smtp/templates/{templateId}", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if transactionalEmailsUpdateTemplateFlags.body != "" {
		if err := json.Unmarshal([]byte(transactionalEmailsUpdateTemplateFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("tag") {
		bodyMap["tag"] = transactionalEmailsUpdateTemplateFlags.tag
	}
	if cmd.Flags().Changed("template-name") {
		bodyMap["templateName"] = transactionalEmailsUpdateTemplateFlags.templateName
	}
	if cmd.Flags().Changed("html-content") {
		bodyMap["htmlContent"] = transactionalEmailsUpdateTemplateFlags.htmlContent
	}
	if cmd.Flags().Changed("html-url") {
		bodyMap["htmlUrl"] = transactionalEmailsUpdateTemplateFlags.htmlUrl
	}
	if cmd.Flags().Changed("subject") {
		bodyMap["subject"] = transactionalEmailsUpdateTemplateFlags.subject
	}
	if cmd.Flags().Changed("reply-to") {
		bodyMap["replyTo"] = transactionalEmailsUpdateTemplateFlags.replyTo
	}
	if cmd.Flags().Changed("to-field") {
		bodyMap["toField"] = transactionalEmailsUpdateTemplateFlags.toField
	}
	if cmd.Flags().Changed("attachment-url") {
		bodyMap["attachmentUrl"] = transactionalEmailsUpdateTemplateFlags.attachmentUrl
	}
	if cmd.Flags().Changed("is-active") {
		bodyMap["isActive"] = transactionalEmailsUpdateTemplateFlags.isActive
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
