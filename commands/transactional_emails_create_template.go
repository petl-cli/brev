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

var transactionalEmailsCreateTemplateCmd = &cobra.Command{
	Use:   "create-template",
	Short: "Create an email template",
	RunE:  runTransactionalEmailsCreateTemplate,
}

var transactionalEmailsCreateTemplateFlags struct {
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
	transactionalEmailsCreateTemplateCmd.Flags().StringVar(&transactionalEmailsCreateTemplateFlags.tag, "tag", "", "Tag of the template")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsCreateTemplateCmd.Flags().StringVar(&transactionalEmailsCreateTemplateFlags.templateName, "template-name", "", "Name of the template")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsCreateTemplateCmd.Flags().StringVar(&transactionalEmailsCreateTemplateFlags.htmlContent, "html-content", "", "Body of the message (HTML version). The field must have more than 10 characters. **REQUIRED if htmlUrl is empty** ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsCreateTemplateCmd.Flags().StringVar(&transactionalEmailsCreateTemplateFlags.htmlUrl, "html-url", "", "Url which contents the body of the email message. REQUIRED if htmlContent is empty")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsCreateTemplateCmd.Flags().StringVar(&transactionalEmailsCreateTemplateFlags.subject, "subject", "", "Subject of the template")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsCreateTemplateCmd.Flags().StringVar(&transactionalEmailsCreateTemplateFlags.replyTo, "reply-to", "", "Email on which campaign recipients will be able to reply to")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsCreateTemplateCmd.Flags().StringVar(&transactionalEmailsCreateTemplateFlags.toField, "to-field", "", "To personalize the **To** Field. If you want to include the first name and last name of your recipient, add **{FNAME} {LNAME}**. These contact attributes must already exist in your Brevo account. If input parameter **params** used please use **{{contact.FNAME}} {{contact.LNAME}}** for personalization ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsCreateTemplateCmd.Flags().StringVar(&transactionalEmailsCreateTemplateFlags.attachmentUrl, "attachment-url", "", "Absolute url of the attachment (**no local file**). Extension allowed: #### xlsx, xls, ods, docx, docm, doc, csv, pdf, txt, gif, jpg, jpeg, png, tif, tiff, rtf, bmp, cgm, css, shtml, html, htm, zip, xml, ppt, pptx, tar, ez, ics, mobi, msg, pub and eps' ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsCreateTemplateCmd.Flags().BoolVar(&transactionalEmailsCreateTemplateFlags.isActive, "is-active", false, "Status of template. isActive = true means template is active and isActive = false means template is inactive")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsCreateTemplateCmd.Flags().StringVar(&transactionalEmailsCreateTemplateFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	transactionalEmailsCmd.AddCommand(transactionalEmailsCreateTemplateCmd)
}

func runTransactionalEmailsCreateTemplate(cmd *cobra.Command, args []string) error {
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
			Name:        "tag",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Tag of the template",
		})
		flags = append(flags, flagSchema{
			Name:        "sender",
			Type:        "object",
			Required:    true,
			Location:    "body",
			Description: "Sender details including id or email and name (_optional_). Only one of either Sender's email or Sender's ID shall be passed in one request at a time. For example: **{\"name\":\"xyz\", \"email\":\"example@abc.com\"}** **{\"name\":\"xyz\", \"id\":123}** ",
		})
		flags = append(flags, flagSchema{
			Name:        "template-name",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Name of the template",
		})
		flags = append(flags, flagSchema{
			Name:        "html-content",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Body of the message (HTML version). The field must have more than 10 characters. **REQUIRED if htmlUrl is empty** ",
		})
		flags = append(flags, flagSchema{
			Name:        "html-url",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Url which contents the body of the email message. REQUIRED if htmlContent is empty",
		})
		flags = append(flags, flagSchema{
			Name:        "subject",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Subject of the template",
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
			Description: "Absolute url of the attachment (**no local file**). Extension allowed: #### xlsx, xls, ods, docx, docm, doc, csv, pdf, txt, gif, jpg, jpeg, png, tif, tiff, rtf, bmp, cgm, css, shtml, html, htm, zip, xml, ppt, pptx, tar, ez, ics, mobi, msg, pub and eps' ",
		})
		flags = append(flags, flagSchema{
			Name:        "is-active",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Status of template. isActive = true means template is active and isActive = false means template is inactive",
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
			Description: "successfully created",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "create-template",
			"description": "Create an email template",
			"http": map[string]any{
				"method": "POST",
				"path":   "/smtp/templates",
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
		Path:        httpclient.SubstitutePath("/smtp/templates", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if transactionalEmailsCreateTemplateFlags.body != "" {
		if err := json.Unmarshal([]byte(transactionalEmailsCreateTemplateFlags.body), &bodyMap); err != nil {
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
		bodyMap["tag"] = transactionalEmailsCreateTemplateFlags.tag
	}
	if cmd.Flags().Changed("template-name") {
		bodyMap["templateName"] = transactionalEmailsCreateTemplateFlags.templateName
	}
	if cmd.Flags().Changed("html-content") {
		bodyMap["htmlContent"] = transactionalEmailsCreateTemplateFlags.htmlContent
	}
	if cmd.Flags().Changed("html-url") {
		bodyMap["htmlUrl"] = transactionalEmailsCreateTemplateFlags.htmlUrl
	}
	if cmd.Flags().Changed("subject") {
		bodyMap["subject"] = transactionalEmailsCreateTemplateFlags.subject
	}
	if cmd.Flags().Changed("reply-to") {
		bodyMap["replyTo"] = transactionalEmailsCreateTemplateFlags.replyTo
	}
	if cmd.Flags().Changed("to-field") {
		bodyMap["toField"] = transactionalEmailsCreateTemplateFlags.toField
	}
	if cmd.Flags().Changed("attachment-url") {
		bodyMap["attachmentUrl"] = transactionalEmailsCreateTemplateFlags.attachmentUrl
	}
	if cmd.Flags().Changed("is-active") {
		bodyMap["isActive"] = transactionalEmailsCreateTemplateFlags.isActive
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
