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

var emailCampaignsCreateCampaignCmd = &cobra.Command{
	Use:   "create-campaign",
	Short: "Create an email campaign",
	RunE:  runEmailCampaignsCreateCampaign,
}

var emailCampaignsCreateCampaignFlags struct {
	tag                   string
	name                  string
	htmlContent           string
	htmlUrl               string
	templateId            int
	scheduledAt           string
	subject               string
	previewText           string
	replyTo               string
	toField               string
	attachmentUrl         string
	inlineImageActivation bool
	mirrorActive          bool
	footer                string
	header                string
	utmCampaign           string
	sendAtBestTime        bool
	abTesting             bool
	subjectA              string
	subjectB              string
	splitRule             int
	winnerCriteria        string
	winnerDelay           int
	ipWarmupEnable        bool
	initialQuota          int
	increaseRate          int
	unsubscriptionPageId  string
	updateFormId          string
	body                  string
}

func init() {
	emailCampaignsCreateCampaignCmd.Flags().StringVar(&emailCampaignsCreateCampaignFlags.tag, "tag", "", "Tag of the campaign")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().StringVar(&emailCampaignsCreateCampaignFlags.name, "name", "", "Name of the campaign")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().StringVar(&emailCampaignsCreateCampaignFlags.htmlContent, "html-content", "", "Mandatory if htmlUrl and templateId are empty. Body of the message (HTML). ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().StringVar(&emailCampaignsCreateCampaignFlags.htmlUrl, "html-url", "", "**Mandatory if htmlContent and templateId are empty**. Url to the message (HTML). For example: **https://html.domain.com** ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().IntVar(&emailCampaignsCreateCampaignFlags.templateId, "template-id", 0, "**Mandatory if htmlContent and htmlUrl are empty**. Id of the transactional email template with status _active_. Used to copy only its content fetched from htmlContent/htmlUrl to an email campaign for RSS feature. ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().StringVar(&emailCampaignsCreateCampaignFlags.scheduledAt, "scheduled-at", "", "Sending UTC date-time (YYYY-MM-DDTHH:mm:ss.SSSZ). **Prefer to pass your timezone in date-time format for accurate result**. If sendAtBestTime is set to true, your campaign will be sent according to the date passed (ignoring the time part). For example: **2017-06-01T12:30:00+02:00** ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().StringVar(&emailCampaignsCreateCampaignFlags.subject, "subject", "", "Subject of the campaign. **Mandatory if abTesting is false**. Ignored if abTesting is true. ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().StringVar(&emailCampaignsCreateCampaignFlags.previewText, "preview-text", "", "Preview text or preheader of the email campaign")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().StringVar(&emailCampaignsCreateCampaignFlags.replyTo, "reply-to", "", "Email on which the campaign recipients will be able to reply to")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().StringVar(&emailCampaignsCreateCampaignFlags.toField, "to-field", "", "To personalize the **To** Field. If you want to include the first name and last name of your recipient, add **{FNAME} {LNAME}**. These contact attributes must already exist in your Brevo account. If input parameter **params** used please use **{{contact.FNAME}} {{contact.LNAME}}** for personalization ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().StringVar(&emailCampaignsCreateCampaignFlags.attachmentUrl, "attachment-url", "", "Absolute url of the attachment (no local file). Extension allowed: #### xlsx, xls, ods, docx, docm, doc, csv, pdf, txt, gif, jpg, jpeg, png, tif, tiff, rtf, bmp, cgm, css, shtml, html, htm, zip, xml, ppt, pptx, tar, ez, ics, mobi, msg, pub and eps ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().BoolVar(&emailCampaignsCreateCampaignFlags.inlineImageActivation, "inline-image-activation", false, "Use true to embedded the images in your email. Final size of the email should be less than **4MB**. Campaigns with embedded images can _not be sent to more than 5000 contacts_ ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().BoolVar(&emailCampaignsCreateCampaignFlags.mirrorActive, "mirror-active", false, "Use true to enable the mirror link")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().StringVar(&emailCampaignsCreateCampaignFlags.footer, "footer", "", "Footer of the email campaign")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().StringVar(&emailCampaignsCreateCampaignFlags.header, "header", "", "Header of the email campaign")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().StringVar(&emailCampaignsCreateCampaignFlags.utmCampaign, "utm-campaign", "", "Customize the utm_campaign value. If this field is empty, the campaign name will be used. Only alphanumeric characters and spaces are allowed")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().BoolVar(&emailCampaignsCreateCampaignFlags.sendAtBestTime, "send-at-best-time", false, "Set this to true if you want to send your campaign at best time.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().BoolVar(&emailCampaignsCreateCampaignFlags.abTesting, "ab-testing", false, "Status of A/B Test. abTesting = false means it is disabled & abTesting = true means it is enabled. **subjectA, subjectB, splitRule, winnerCriteria & winnerDelay** will be considered when abTesting is set to true. subjectA & subjectB are mandatory together & subject if passed is ignored. **Can be set to true only if sendAtBestTime is false**. You will be able to set up two subject lines for your campaign and send them to a random sample of your total recipients. Half of the test group will receive version A, and the other half will receive version B ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().StringVar(&emailCampaignsCreateCampaignFlags.subjectA, "subject-a", "", "Subject A of the campaign. **Mandatory if abTesting = true**. subjectA & subjectB should have unique value ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().StringVar(&emailCampaignsCreateCampaignFlags.subjectB, "subject-b", "", "Subject B of the campaign. **Mandatory if abTesting = true**. subjectA & subjectB should have unique value ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().IntVar(&emailCampaignsCreateCampaignFlags.splitRule, "split-rule", 0, "Add the size of your test groups. **Mandatory if abTesting = true & 'recipients' is passed**. We'll send version A and B to a random sample of recipients, and then the winning version to everyone else ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().StringVar(&emailCampaignsCreateCampaignFlags.winnerCriteria, "winner-criteria", "", "Choose the metrics that will determinate the winning version. **Mandatory if _splitRule_ >= 1 and < 50**. If splitRule = 50, `winnerCriteria` is ignored if passed ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().IntVar(&emailCampaignsCreateCampaignFlags.winnerDelay, "winner-delay", 0, "Choose the duration of the test in hours. Maximum is 7 days, pass 24*7 = 168 hours. The winning version will be sent at the end of the test. **Mandatory if _splitRule_ >= 1 and < 50**. If splitRule = 50, `winnerDelay` is ignored if passed ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().BoolVar(&emailCampaignsCreateCampaignFlags.ipWarmupEnable, "ip-warmup-enable", false, "**Available for dedicated ip clients**. Set this to true if you wish to warm up your ip. ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().IntVar(&emailCampaignsCreateCampaignFlags.initialQuota, "initial-quota", 0, "**Mandatory if ipWarmupEnable is set to true**. Set an initial quota greater than 1 for warming up your ip. We recommend you set a value of 3000. ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().IntVar(&emailCampaignsCreateCampaignFlags.increaseRate, "increase-rate", 0, "**Mandatory if ipWarmupEnable is set to true**. Set a percentage increase rate for warming up your ip. We recommend you set the increase rate to 30% per day. If you want to send the same number of emails every day, set the daily increase value to 0%. ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().StringVar(&emailCampaignsCreateCampaignFlags.unsubscriptionPageId, "unsubscription-page-id", "", "Enter an unsubscription page id. The page id is a 24 digit alphanumeric id that can be found in the URL when editing the page. If not entered, then the default unsubscription page will be used. ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().StringVar(&emailCampaignsCreateCampaignFlags.updateFormId, "update-form-id", "", "**Mandatory if templateId is used containing the {{ update_profile }} tag**. Enter an update profile form id. The form id is a 24 digit alphanumeric id that can be found in the URL when editing the form. If not entered, then the default update profile form will be used. ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsCreateCampaignCmd.Flags().StringVar(&emailCampaignsCreateCampaignFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	emailCampaignsCmd.AddCommand(emailCampaignsCreateCampaignCmd)
}

func runEmailCampaignsCreateCampaign(cmd *cobra.Command, args []string) error {
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
			Description: "Tag of the campaign",
		})
		flags = append(flags, flagSchema{
			Name:        "sender",
			Type:        "object",
			Required:    true,
			Location:    "body",
			Description: "Sender details including id or email and name (_optional_). Only one of either Sender's email or Sender's ID shall be passed in one request at a time. For example: **{\"name\":\"xyz\", \"email\":\"example@abc.com\"}** **{\"name\":\"xyz\", \"id\":123}** ",
		})
		flags = append(flags, flagSchema{
			Name:        "name",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Name of the campaign",
		})
		flags = append(flags, flagSchema{
			Name:        "html-content",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Mandatory if htmlUrl and templateId are empty. Body of the message (HTML). ",
		})
		flags = append(flags, flagSchema{
			Name:        "html-url",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "**Mandatory if htmlContent and templateId are empty**. Url to the message (HTML). For example: **https://html.domain.com** ",
		})
		flags = append(flags, flagSchema{
			Name:        "template-id",
			Type:        "integer",
			Required:    false,
			Location:    "body",
			Description: "**Mandatory if htmlContent and htmlUrl are empty**. Id of the transactional email template with status _active_. Used to copy only its content fetched from htmlContent/htmlUrl to an email campaign for RSS feature. ",
		})
		flags = append(flags, flagSchema{
			Name:        "scheduled-at",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Sending UTC date-time (YYYY-MM-DDTHH:mm:ss.SSSZ). **Prefer to pass your timezone in date-time format for accurate result**. If sendAtBestTime is set to true, your campaign will be sent according to the date passed (ignoring the time part). For example: **2017-06-01T12:30:00+02:00** ",
		})
		flags = append(flags, flagSchema{
			Name:        "subject",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Subject of the campaign. **Mandatory if abTesting is false**. Ignored if abTesting is true. ",
		})
		flags = append(flags, flagSchema{
			Name:        "preview-text",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Preview text or preheader of the email campaign",
		})
		flags = append(flags, flagSchema{
			Name:        "reply-to",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Email on which the campaign recipients will be able to reply to",
		})
		flags = append(flags, flagSchema{
			Name:        "to-field",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "To personalize the **To** Field. If you want to include the first name and last name of your recipient, add **{FNAME} {LNAME}**. These contact attributes must already exist in your Brevo account. If input parameter **params** used please use **{{contact.FNAME}} {{contact.LNAME}}** for personalization ",
		})
		flags = append(flags, flagSchema{
			Name:        "recipients",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "Segment ids and List ids to include/exclude from campaign",
		})
		flags = append(flags, flagSchema{
			Name:        "attachment-url",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Absolute url of the attachment (no local file). Extension allowed: #### xlsx, xls, ods, docx, docm, doc, csv, pdf, txt, gif, jpg, jpeg, png, tif, tiff, rtf, bmp, cgm, css, shtml, html, htm, zip, xml, ppt, pptx, tar, ez, ics, mobi, msg, pub and eps ",
		})
		flags = append(flags, flagSchema{
			Name:        "inline-image-activation",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Use true to embedded the images in your email. Final size of the email should be less than **4MB**. Campaigns with embedded images can _not be sent to more than 5000 contacts_ ",
		})
		flags = append(flags, flagSchema{
			Name:        "mirror-active",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Use true to enable the mirror link",
		})
		flags = append(flags, flagSchema{
			Name:        "footer",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Footer of the email campaign",
		})
		flags = append(flags, flagSchema{
			Name:        "header",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Header of the email campaign",
		})
		flags = append(flags, flagSchema{
			Name:        "utm-campaign",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Customize the utm_campaign value. If this field is empty, the campaign name will be used. Only alphanumeric characters and spaces are allowed",
		})
		flags = append(flags, flagSchema{
			Name:        "params",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "Pass the set of attributes to customize the type classic campaign. For example: **{\"FNAME\":\"Joe\", \"LNAME\":\"Doe\"}**. Only available if **type** is **classic**. It's considered only if campaign is in _New Template Language format_. The New Template Language is dependent on the values of **subject, htmlContent/htmlUrl, sender.name & toField** ",
		})
		flags = append(flags, flagSchema{
			Name:        "send-at-best-time",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Set this to true if you want to send your campaign at best time.",
		})
		flags = append(flags, flagSchema{
			Name:        "ab-testing",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Status of A/B Test. abTesting = false means it is disabled & abTesting = true means it is enabled. **subjectA, subjectB, splitRule, winnerCriteria & winnerDelay** will be considered when abTesting is set to true. subjectA & subjectB are mandatory together & subject if passed is ignored. **Can be set to true only if sendAtBestTime is false**. You will be able to set up two subject lines for your campaign and send them to a random sample of your total recipients. Half of the test group will receive version A, and the other half will receive version B ",
		})
		flags = append(flags, flagSchema{
			Name:        "subject-a",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Subject A of the campaign. **Mandatory if abTesting = true**. subjectA & subjectB should have unique value ",
		})
		flags = append(flags, flagSchema{
			Name:        "subject-b",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Subject B of the campaign. **Mandatory if abTesting = true**. subjectA & subjectB should have unique value ",
		})
		flags = append(flags, flagSchema{
			Name:        "split-rule",
			Type:        "integer",
			Required:    false,
			Location:    "body",
			Description: "Add the size of your test groups. **Mandatory if abTesting = true & 'recipients' is passed**. We'll send version A and B to a random sample of recipients, and then the winning version to everyone else ",
		})
		flags = append(flags, flagSchema{
			Name:        "winner-criteria",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Choose the metrics that will determinate the winning version. **Mandatory if _splitRule_ >= 1 and < 50**. If splitRule = 50, `winnerCriteria` is ignored if passed ",
		})
		flags = append(flags, flagSchema{
			Name:        "winner-delay",
			Type:        "integer",
			Required:    false,
			Location:    "body",
			Description: "Choose the duration of the test in hours. Maximum is 7 days, pass 24*7 = 168 hours. The winning version will be sent at the end of the test. **Mandatory if _splitRule_ >= 1 and < 50**. If splitRule = 50, `winnerDelay` is ignored if passed ",
		})
		flags = append(flags, flagSchema{
			Name:        "ip-warmup-enable",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "**Available for dedicated ip clients**. Set this to true if you wish to warm up your ip. ",
		})
		flags = append(flags, flagSchema{
			Name:        "initial-quota",
			Type:        "integer",
			Required:    false,
			Location:    "body",
			Description: "**Mandatory if ipWarmupEnable is set to true**. Set an initial quota greater than 1 for warming up your ip. We recommend you set a value of 3000. ",
		})
		flags = append(flags, flagSchema{
			Name:        "increase-rate",
			Type:        "integer",
			Required:    false,
			Location:    "body",
			Description: "**Mandatory if ipWarmupEnable is set to true**. Set a percentage increase rate for warming up your ip. We recommend you set the increase rate to 30% per day. If you want to send the same number of emails every day, set the daily increase value to 0%. ",
		})
		flags = append(flags, flagSchema{
			Name:        "unsubscription-page-id",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Enter an unsubscription page id. The page id is a 24 digit alphanumeric id that can be found in the URL when editing the page. If not entered, then the default unsubscription page will be used. ",
		})
		flags = append(flags, flagSchema{
			Name:        "update-form-id",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "**Mandatory if templateId is used containing the {{ update_profile }} tag**. Enter an update profile form id. The form id is a 24 digit alphanumeric id that can be found in the URL when editing the form. If not entered, then the default update profile form will be used. ",
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
			"command":     "create-campaign",
			"description": "Create an email campaign",
			"http": map[string]any{
				"method": "POST",
				"path":   "/emailCampaigns",
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
		Path:        httpclient.SubstitutePath("/emailCampaigns", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if emailCampaignsCreateCampaignFlags.body != "" {
		if err := json.Unmarshal([]byte(emailCampaignsCreateCampaignFlags.body), &bodyMap); err != nil {
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
		bodyMap["tag"] = emailCampaignsCreateCampaignFlags.tag
	}
	if cmd.Flags().Changed("name") {
		bodyMap["name"] = emailCampaignsCreateCampaignFlags.name
	}
	if cmd.Flags().Changed("html-content") {
		bodyMap["htmlContent"] = emailCampaignsCreateCampaignFlags.htmlContent
	}
	if cmd.Flags().Changed("html-url") {
		bodyMap["htmlUrl"] = emailCampaignsCreateCampaignFlags.htmlUrl
	}
	if cmd.Flags().Changed("template-id") {
		bodyMap["templateId"] = emailCampaignsCreateCampaignFlags.templateId
	}
	if cmd.Flags().Changed("scheduled-at") {
		bodyMap["scheduledAt"] = emailCampaignsCreateCampaignFlags.scheduledAt
	}
	if cmd.Flags().Changed("subject") {
		bodyMap["subject"] = emailCampaignsCreateCampaignFlags.subject
	}
	if cmd.Flags().Changed("preview-text") {
		bodyMap["previewText"] = emailCampaignsCreateCampaignFlags.previewText
	}
	if cmd.Flags().Changed("reply-to") {
		bodyMap["replyTo"] = emailCampaignsCreateCampaignFlags.replyTo
	}
	if cmd.Flags().Changed("to-field") {
		bodyMap["toField"] = emailCampaignsCreateCampaignFlags.toField
	}
	if cmd.Flags().Changed("attachment-url") {
		bodyMap["attachmentUrl"] = emailCampaignsCreateCampaignFlags.attachmentUrl
	}
	if cmd.Flags().Changed("inline-image-activation") {
		bodyMap["inlineImageActivation"] = emailCampaignsCreateCampaignFlags.inlineImageActivation
	}
	if cmd.Flags().Changed("mirror-active") {
		bodyMap["mirrorActive"] = emailCampaignsCreateCampaignFlags.mirrorActive
	}
	if cmd.Flags().Changed("footer") {
		bodyMap["footer"] = emailCampaignsCreateCampaignFlags.footer
	}
	if cmd.Flags().Changed("header") {
		bodyMap["header"] = emailCampaignsCreateCampaignFlags.header
	}
	if cmd.Flags().Changed("utm-campaign") {
		bodyMap["utmCampaign"] = emailCampaignsCreateCampaignFlags.utmCampaign
	}
	if cmd.Flags().Changed("send-at-best-time") {
		bodyMap["sendAtBestTime"] = emailCampaignsCreateCampaignFlags.sendAtBestTime
	}
	if cmd.Flags().Changed("ab-testing") {
		bodyMap["abTesting"] = emailCampaignsCreateCampaignFlags.abTesting
	}
	if cmd.Flags().Changed("subject-a") {
		bodyMap["subjectA"] = emailCampaignsCreateCampaignFlags.subjectA
	}
	if cmd.Flags().Changed("subject-b") {
		bodyMap["subjectB"] = emailCampaignsCreateCampaignFlags.subjectB
	}
	if cmd.Flags().Changed("split-rule") {
		bodyMap["splitRule"] = emailCampaignsCreateCampaignFlags.splitRule
	}
	if cmd.Flags().Changed("winner-criteria") {
		bodyMap["winnerCriteria"] = emailCampaignsCreateCampaignFlags.winnerCriteria
	}
	if cmd.Flags().Changed("winner-delay") {
		bodyMap["winnerDelay"] = emailCampaignsCreateCampaignFlags.winnerDelay
	}
	if cmd.Flags().Changed("ip-warmup-enable") {
		bodyMap["ipWarmupEnable"] = emailCampaignsCreateCampaignFlags.ipWarmupEnable
	}
	if cmd.Flags().Changed("initial-quota") {
		bodyMap["initialQuota"] = emailCampaignsCreateCampaignFlags.initialQuota
	}
	if cmd.Flags().Changed("increase-rate") {
		bodyMap["increaseRate"] = emailCampaignsCreateCampaignFlags.increaseRate
	}
	if cmd.Flags().Changed("unsubscription-page-id") {
		bodyMap["unsubscriptionPageId"] = emailCampaignsCreateCampaignFlags.unsubscriptionPageId
	}
	if cmd.Flags().Changed("update-form-id") {
		bodyMap["updateFormId"] = emailCampaignsCreateCampaignFlags.updateFormId
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
