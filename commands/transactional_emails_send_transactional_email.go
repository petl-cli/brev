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

var transactionalEmailsSendTransactionalEmailCmd = &cobra.Command{
	Use:   "send-transactional-email",
	Short: "Send a transactional email",
	RunE:  runTransactionalEmailsSendTransactionalEmail,
}

var transactionalEmailsSendTransactionalEmailFlags struct {
	tags            []string
	to              []string
	bcc             []string
	cc              []string
	htmlContent     string
	textContent     string
	subject         string
	attachment      []string
	templateId      int
	messageVersions []string
	scheduledAt     string
	batchId         string
	body            string
}

func init() {
	transactionalEmailsSendTransactionalEmailCmd.Flags().StringSliceVar(&transactionalEmailsSendTransactionalEmailFlags.tags, "tags", nil, "Tag your emails to find them more easily")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsSendTransactionalEmailCmd.Flags().StringSliceVar(&transactionalEmailsSendTransactionalEmailFlags.to, "to", nil, "**Mandatory if messageVersions are not passed, ignored if messageVersions are passed** List of email addresses and names (_optional_) of the recipients. For example, **[{\"name\":\"Jimmy\", \"email\":\"jimmy98@example.com\"}, {\"name\":\"Joe\", \"email\":\"joe@example.com\"}]** ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsSendTransactionalEmailCmd.Flags().StringSliceVar(&transactionalEmailsSendTransactionalEmailFlags.bcc, "bcc", nil, "List of email addresses and names (_optional_) of the recipients in bcc ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsSendTransactionalEmailCmd.Flags().StringSliceVar(&transactionalEmailsSendTransactionalEmailFlags.cc, "cc", nil, "List of email addresses and names (_optional_) of the recipients in cc ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsSendTransactionalEmailCmd.Flags().StringVar(&transactionalEmailsSendTransactionalEmailFlags.htmlContent, "html-content", "", "HTML body of the message. **Mandatory if 'templateId' is not passed, ignored if 'templateId' is passed** ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsSendTransactionalEmailCmd.Flags().StringVar(&transactionalEmailsSendTransactionalEmailFlags.textContent, "text-content", "", "Plain Text body of the message. **Ignored if 'templateId' is passed** ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsSendTransactionalEmailCmd.Flags().StringVar(&transactionalEmailsSendTransactionalEmailFlags.subject, "subject", "", "Subject of the message. **Mandatory if 'templateId' is not passed** ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsSendTransactionalEmailCmd.Flags().StringSliceVar(&transactionalEmailsSendTransactionalEmailFlags.attachment, "attachment", nil, "Pass the _absolute URL_ (**no local file**) or the _base64 content_ of the attachment along with the attachment name. **Mandatory if attachment content is passed**. For example, **[{\"url\":\"https://attachment.domain.com/myAttachmentFromUrl.jpg\", \"name\":\"myAttachmentFromUrl.jpg\"}, {\"content\":\"base64 example content\", \"name\":\"myAttachmentFromBase64.jpg\"}]**. Allowed extensions for attachment file: #### xlsx, xls, ods, docx, docm, doc, csv, pdf, txt, gif, jpg, jpeg, png, tif, tiff, rtf, bmp, cgm, css, shtml, html, htm, zip, xml, ppt, pptx, tar, ez, ics, mobi, msg, pub, eps, odt, mp3, m4a, m4v, wma, ogg, flac, wav, aif, aifc, aiff, mp4, mov, avi, mkv, mpeg, mpg, wmv, pkpass and xlsm. If `templateId` is passed and is in New Template Language format then both attachment url and content are accepted. If template is in Old template Language format, then `attachment` is ignored ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsSendTransactionalEmailCmd.Flags().IntVar(&transactionalEmailsSendTransactionalEmailFlags.templateId, "template-id", 0, "Id of the template.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsSendTransactionalEmailCmd.Flags().StringSliceVar(&transactionalEmailsSendTransactionalEmailFlags.messageVersions, "message-versions", nil, "You can customize and send out multiple versions of a mail. **templateId** can be customized only if global parameter contains templateId. **htmlContent and textContent** can be customized only if any of the two, htmlContent or textContent, is present in global parameters. Some global parameters such as **to(mandatory), bcc, cc, replyTo, subject** can also be customized specific to each version. Total number of recipients in one API request must not exceed 2000. However, you can still pass upto 99 recipients maximum in one message version. The size of individual params in all the messageVersions shall not exceed **100 KB** limit and that of cumulative params shall not exceed **1000 KB**. You can follow this **step-by-step guide** on how to use **messageVersions** to batch send emails - **https://developers.brevo.com/docs/batch-send-transactional-emails** ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsSendTransactionalEmailCmd.Flags().StringVar(&transactionalEmailsSendTransactionalEmailFlags.scheduledAt, "scheduled-at", "", "UTC date-time on which the email has to schedule (YYYY-MM-DDTHH:mm:ss.SSSZ). Prefer to pass your timezone in date-time format for scheduling. There can be an expected delay of +5 minutes in scheduled email delivery. **Please note this feature is currently a public beta**.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsSendTransactionalEmailCmd.Flags().StringVar(&transactionalEmailsSendTransactionalEmailFlags.batchId, "batch-id", "", "Valid UUIDv4 batch id to identify the scheduled batches transactional email. If not passed we will create a valid UUIDv4 batch id at our end.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	transactionalEmailsSendTransactionalEmailCmd.Flags().StringVar(&transactionalEmailsSendTransactionalEmailFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	transactionalEmailsCmd.AddCommand(transactionalEmailsSendTransactionalEmailCmd)
}

func runTransactionalEmailsSendTransactionalEmail(cmd *cobra.Command, args []string) error {
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
			Name:        "tags",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "Tag your emails to find them more easily",
		})
		flags = append(flags, flagSchema{
			Name:        "sender",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "**Mandatory if `templateId` is not passed**. Pass name (_optional_) and email or id of sender from which emails will be sent. **`name` will be ignored if passed along with sender `id`**. For example, **{\"name\":\"Mary from MyShop\", \"email\":\"no-reply@myshop.com\"}** **{\"id\":2}** ",
		})
		flags = append(flags, flagSchema{
			Name:        "to",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "**Mandatory if messageVersions are not passed, ignored if messageVersions are passed** List of email addresses and names (_optional_) of the recipients. For example, **[{\"name\":\"Jimmy\", \"email\":\"jimmy98@example.com\"}, {\"name\":\"Joe\", \"email\":\"joe@example.com\"}]** ",
		})
		flags = append(flags, flagSchema{
			Name:        "bcc",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "List of email addresses and names (_optional_) of the recipients in bcc ",
		})
		flags = append(flags, flagSchema{
			Name:        "cc",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "List of email addresses and names (_optional_) of the recipients in cc ",
		})
		flags = append(flags, flagSchema{
			Name:        "html-content",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "HTML body of the message. **Mandatory if 'templateId' is not passed, ignored if 'templateId' is passed** ",
		})
		flags = append(flags, flagSchema{
			Name:        "text-content",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Plain Text body of the message. **Ignored if 'templateId' is passed** ",
		})
		flags = append(flags, flagSchema{
			Name:        "subject",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Subject of the message. **Mandatory if 'templateId' is not passed** ",
		})
		flags = append(flags, flagSchema{
			Name:        "reply-to",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "Email (**required**), along with name (_optional_), on which transactional mail recipients will be able to reply back. For example, **{\"email\":\"ann6533@example.com\", \"name\":\"Ann\"}** ",
		})
		flags = append(flags, flagSchema{
			Name:        "attachment",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "Pass the _absolute URL_ (**no local file**) or the _base64 content_ of the attachment along with the attachment name. **Mandatory if attachment content is passed**. For example, **[{\"url\":\"https://attachment.domain.com/myAttachmentFromUrl.jpg\", \"name\":\"myAttachmentFromUrl.jpg\"}, {\"content\":\"base64 example content\", \"name\":\"myAttachmentFromBase64.jpg\"}]**. Allowed extensions for attachment file: #### xlsx, xls, ods, docx, docm, doc, csv, pdf, txt, gif, jpg, jpeg, png, tif, tiff, rtf, bmp, cgm, css, shtml, html, htm, zip, xml, ppt, pptx, tar, ez, ics, mobi, msg, pub, eps, odt, mp3, m4a, m4v, wma, ogg, flac, wav, aif, aifc, aiff, mp4, mov, avi, mkv, mpeg, mpg, wmv, pkpass and xlsm. If `templateId` is passed and is in New Template Language format then both attachment url and content are accepted. If template is in Old template Language format, then `attachment` is ignored ",
		})
		flags = append(flags, flagSchema{
			Name:        "headers",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "Pass the set of custom headers (_not the standard headers_) that shall be sent along the mail headers in the original email. **'sender.ip'** header can be set (**only for dedicated ip users**) to mention the IP to be used for sending transactional emails. Headers are allowed in `This-Case-Only` (i.e. words separated by hyphen with first letter of each word in capital letter), they will be converted to such case styling if not in this format in the request payload. For example, **{\"sender.ip\":\"1.2.3.4\", \"X-Mailin-custom\":\"some_custom_header\", \"idempotencyKey\":\"abc-123\"}**. ",
		})
		flags = append(flags, flagSchema{
			Name:        "template-id",
			Type:        "integer",
			Required:    false,
			Location:    "body",
			Description: "Id of the template.",
		})
		flags = append(flags, flagSchema{
			Name:        "params",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "Pass the set of attributes to customize the template. For example, **{\"FNAME\":\"Joe\", \"LNAME\":\"Doe\"}**. It's **considered only if template is in New Template Language format**. ",
		})
		flags = append(flags, flagSchema{
			Name:        "message-versions",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "You can customize and send out multiple versions of a mail. **templateId** can be customized only if global parameter contains templateId. **htmlContent and textContent** can be customized only if any of the two, htmlContent or textContent, is present in global parameters. Some global parameters such as **to(mandatory), bcc, cc, replyTo, subject** can also be customized specific to each version. Total number of recipients in one API request must not exceed 2000. However, you can still pass upto 99 recipients maximum in one message version. The size of individual params in all the messageVersions shall not exceed **100 KB** limit and that of cumulative params shall not exceed **1000 KB**. You can follow this **step-by-step guide** on how to use **messageVersions** to batch send emails - **https://developers.brevo.com/docs/batch-send-transactional-emails** ",
		})
		flags = append(flags, flagSchema{
			Name:        "scheduled-at",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "UTC date-time on which the email has to schedule (YYYY-MM-DDTHH:mm:ss.SSSZ). Prefer to pass your timezone in date-time format for scheduling. There can be an expected delay of +5 minutes in scheduled email delivery. **Please note this feature is currently a public beta**.",
		})
		flags = append(flags, flagSchema{
			Name:        "batch-id",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Valid UUIDv4 batch id to identify the scheduled batches transactional email. If not passed we will create a valid UUIDv4 batch id at our end.",
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
			Description: "transactional email sent",
		})
		responses = append(responses, responseSchema{
			Status:      "202",
			ContentType: "application/json",
			Description: "transactional email scheduled",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "send-transactional-email",
			"description": "Send a transactional email",
			"http": map[string]any{
				"method": "POST",
				"path":   "/smtp/email",
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
		Path:        httpclient.SubstitutePath("/smtp/email", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if transactionalEmailsSendTransactionalEmailFlags.body != "" {
		if err := json.Unmarshal([]byte(transactionalEmailsSendTransactionalEmailFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("tags") {
		bodyMap["tags"] = transactionalEmailsSendTransactionalEmailFlags.tags
	}
	if cmd.Flags().Changed("to") {
		bodyMap["to"] = transactionalEmailsSendTransactionalEmailFlags.to
	}
	if cmd.Flags().Changed("bcc") {
		bodyMap["bcc"] = transactionalEmailsSendTransactionalEmailFlags.bcc
	}
	if cmd.Flags().Changed("cc") {
		bodyMap["cc"] = transactionalEmailsSendTransactionalEmailFlags.cc
	}
	if cmd.Flags().Changed("html-content") {
		bodyMap["htmlContent"] = transactionalEmailsSendTransactionalEmailFlags.htmlContent
	}
	if cmd.Flags().Changed("text-content") {
		bodyMap["textContent"] = transactionalEmailsSendTransactionalEmailFlags.textContent
	}
	if cmd.Flags().Changed("subject") {
		bodyMap["subject"] = transactionalEmailsSendTransactionalEmailFlags.subject
	}
	if cmd.Flags().Changed("attachment") {
		bodyMap["attachment"] = transactionalEmailsSendTransactionalEmailFlags.attachment
	}
	if cmd.Flags().Changed("template-id") {
		bodyMap["templateId"] = transactionalEmailsSendTransactionalEmailFlags.templateId
	}
	if cmd.Flags().Changed("message-versions") {
		bodyMap["messageVersions"] = transactionalEmailsSendTransactionalEmailFlags.messageVersions
	}
	if cmd.Flags().Changed("scheduled-at") {
		bodyMap["scheduledAt"] = transactionalEmailsSendTransactionalEmailFlags.scheduledAt
	}
	if cmd.Flags().Changed("batch-id") {
		bodyMap["batchId"] = transactionalEmailsSendTransactionalEmailFlags.batchId
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
