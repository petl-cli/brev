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

var emailCampaignsUpdateCampaignCmd = &cobra.Command{
	Use:   "update-campaign",
	Short: "Update an email campaign",
	RunE:  runEmailCampaignsUpdateCampaign,
}

var emailCampaignsUpdateCampaignFlags struct {
	campaignId            int
	tag                   string
	name                  string
	htmlContent           string
	htmlUrl               string
	scheduledAt           string
	subject               string
	previewText           string
	replyTo               string
	toField               string
	attachmentUrl         string
	inlineImageActivation bool
	mirrorActive          bool
	recurring             bool
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
	emailCampaignsUpdateCampaignCmd.Flags().IntVar(&emailCampaignsUpdateCampaignFlags.campaignId, "campaign-id", 0, "Id of the campaign")
	emailCampaignsUpdateCampaignCmd.MarkFlagRequired("campaign-id")
	emailCampaignsUpdateCampaignCmd.Flags().StringVar(&emailCampaignsUpdateCampaignFlags.tag, "tag", "", "Tag of the campaign")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().StringVar(&emailCampaignsUpdateCampaignFlags.name, "name", "", "Name of the campaign")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().StringVar(&emailCampaignsUpdateCampaignFlags.htmlContent, "html-content", "", "Body of the message (HTML version). If the campaign is designed using Drag & Drop editor via HTML content, then the design page will not have Drag & Drop editor access for that campaign. **REQUIRED if htmlUrl is empty** ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().StringVar(&emailCampaignsUpdateCampaignFlags.htmlUrl, "html-url", "", "Url which contents the body of the email message. **REQUIRED if htmlContent is empty** ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().StringVar(&emailCampaignsUpdateCampaignFlags.scheduledAt, "scheduled-at", "", "UTC date-time on which the campaign has to run (YYYY-MM-DDTHH:mm:ss.SSSZ). **Prefer to pass your timezone in date-time format for accurate result.** If sendAtBestTime is set to true, your campaign will be sent according to the date passed (ignoring the time part). ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().StringVar(&emailCampaignsUpdateCampaignFlags.subject, "subject", "", "Subject of the campaign")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().StringVar(&emailCampaignsUpdateCampaignFlags.previewText, "preview-text", "", "Preview text or preheader of the email campaign")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().StringVar(&emailCampaignsUpdateCampaignFlags.replyTo, "reply-to", "", "Email on which campaign recipients will be able to reply to")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().StringVar(&emailCampaignsUpdateCampaignFlags.toField, "to-field", "", "To personalize the **To** Field. If you want to include the first name and last name of your recipient, add **{FNAME} {LNAME}**. These contact attributes must already exist in your Brevo account. If input parameter **params** used please use **{{contact.FNAME}} {{contact.LNAME}}** for personalization ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().StringVar(&emailCampaignsUpdateCampaignFlags.attachmentUrl, "attachment-url", "", "Absolute url of the attachment (no local file). Extension allowed: #### xlsx, xls, ods, docx, docm, doc, csv, pdf, txt, gif, jpg, jpeg, png, tif, tiff, rtf, bmp, cgm, css, shtml, html, htm, zip, xml, ppt, pptx, tar, ez, ics, mobi, msg, pub and eps' ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().BoolVar(&emailCampaignsUpdateCampaignFlags.inlineImageActivation, "inline-image-activation", false, "Status of inline image. inlineImageActivation = false means image can’t be embedded, & inlineImageActivation = true means image can be embedded, in the email. You cannot send a campaign of more than **4MB** with images embedded in the email. Campaigns with the images embedded in the email _must be sent to less than 5000 contacts_. ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().BoolVar(&emailCampaignsUpdateCampaignFlags.mirrorActive, "mirror-active", false, "Status of mirror links in campaign. mirrorActive = false means mirror links are deactivated, & mirrorActive = true means mirror links are activated, in the campaign")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().BoolVar(&emailCampaignsUpdateCampaignFlags.recurring, "recurring", false, "**FOR TRIGGER ONLY !** Type of trigger campaign.recurring = false means contact can receive the same Trigger campaign only once, & recurring = true means contact can receive the same Trigger campaign several times ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().StringVar(&emailCampaignsUpdateCampaignFlags.footer, "footer", "", "Footer of the email campaign")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().StringVar(&emailCampaignsUpdateCampaignFlags.header, "header", "", "Header of the email campaign")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().StringVar(&emailCampaignsUpdateCampaignFlags.utmCampaign, "utm-campaign", "", "Customize the utm_campaign value. If this field is empty, the campaign name will be used. Only alphanumeric characters and spaces are allowed")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().BoolVar(&emailCampaignsUpdateCampaignFlags.sendAtBestTime, "send-at-best-time", false, "Set this to true if you want to send your campaign at best time. Note:- **if true, warmup ip will be disabled.** ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().BoolVar(&emailCampaignsUpdateCampaignFlags.abTesting, "ab-testing", false, "Status of A/B Test. abTesting = false means it is disabled & abTesting = true means it is enabled. **subjectA, subjectB, splitRule, winnerCriteria & winnerDelay** will be considered when abTesting is set to true. subjectA & subjectB are mandatory together & subject if passed is ignored. **Can be set to true only if sendAtBestTime is false**. You will be able to set up two subject lines for your campaign and send them to a random sample of your total recipients. Half of the test group will receive version A, and the other half will receive version B ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().StringVar(&emailCampaignsUpdateCampaignFlags.subjectA, "subject-a", "", "Subject A of the campaign. **Mandatory if abTesting = true**. subjectA & subjectB should have unique value ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().StringVar(&emailCampaignsUpdateCampaignFlags.subjectB, "subject-b", "", "Subject B of the campaign. **Mandatory if abTesting = true**. subjectA & subjectB should have unique value ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().IntVar(&emailCampaignsUpdateCampaignFlags.splitRule, "split-rule", 0, "Add the size of your test groups. **Mandatory if abTesting = true & 'recipients' is passed**. We'll send version A and B to a random sample of recipients, and then the winning version to everyone else ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().StringVar(&emailCampaignsUpdateCampaignFlags.winnerCriteria, "winner-criteria", "", "Choose the metrics that will determinate the winning version. **Mandatory if _splitRule_ >= 1 and < 50**. If splitRule = 50, `winnerCriteria` is ignored if passed ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().IntVar(&emailCampaignsUpdateCampaignFlags.winnerDelay, "winner-delay", 0, "Choose the duration of the test in hours. Maximum is 7 days, pass 24*7 = 168 hours. The winning version will be sent at the end of the test. **Mandatory if _splitRule_ >= 1 and < 50**. If splitRule = 50, `winnerDelay` is ignored if passed ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().BoolVar(&emailCampaignsUpdateCampaignFlags.ipWarmupEnable, "ip-warmup-enable", false, "**Available for dedicated ip clients**. Set this to true if you wish to warm up your ip. ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().IntVar(&emailCampaignsUpdateCampaignFlags.initialQuota, "initial-quota", 0, "Set an initial quota greater than 1 for warming up your ip. We recommend you set a value of 3000. ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().IntVar(&emailCampaignsUpdateCampaignFlags.increaseRate, "increase-rate", 0, "Set a percentage increase rate for warming up your ip. We recommend you set the increase rate to 30% per day. If you want to send the same number of emails every day, set the daily increase value to 0%. ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().StringVar(&emailCampaignsUpdateCampaignFlags.unsubscriptionPageId, "unsubscription-page-id", "", "Enter an unsubscription page id. The page id is a 24 digit alphanumeric id that can be found in the URL when editing the page. ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().StringVar(&emailCampaignsUpdateCampaignFlags.updateFormId, "update-form-id", "", "**Mandatory if templateId is used containing the {{ update_profile }} tag**. Enter an update profile form id. The form id is a 24 digit alphanumeric id that can be found in the URL when editing the form. ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	emailCampaignsUpdateCampaignCmd.Flags().StringVar(&emailCampaignsUpdateCampaignFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	emailCampaignsCmd.AddCommand(emailCampaignsUpdateCampaignCmd)
}

func runEmailCampaignsUpdateCampaign(cmd *cobra.Command, args []string) error {
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
			Name:        "campaign-id",
			Type:        "integer",
			Required:    true,
			Location:    "path",
			Description: "Id of the campaign",
		})
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
			Required:    false,
			Location:    "body",
			Description: "Sender details including id or email and name (optional). Only one of either Sender's email or Sender's ID shall be passed in one request at a time. For example: **{\"name\":\"xyz\", \"email\":\"example@abc.com\"}** **{\"name\":\"xyz\", \"id\":123}** ",
		})
		flags = append(flags, flagSchema{
			Name:        "name",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Name of the campaign",
		})
		flags = append(flags, flagSchema{
			Name:        "html-content",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Body of the message (HTML version). If the campaign is designed using Drag & Drop editor via HTML content, then the design page will not have Drag & Drop editor access for that campaign. **REQUIRED if htmlUrl is empty** ",
		})
		flags = append(flags, flagSchema{
			Name:        "html-url",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Url which contents the body of the email message. **REQUIRED if htmlContent is empty** ",
		})
		flags = append(flags, flagSchema{
			Name:        "scheduled-at",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "UTC date-time on which the campaign has to run (YYYY-MM-DDTHH:mm:ss.SSSZ). **Prefer to pass your timezone in date-time format for accurate result.** If sendAtBestTime is set to true, your campaign will be sent according to the date passed (ignoring the time part). ",
		})
		flags = append(flags, flagSchema{
			Name:        "subject",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Subject of the campaign",
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
			Description: "Absolute url of the attachment (no local file). Extension allowed: #### xlsx, xls, ods, docx, docm, doc, csv, pdf, txt, gif, jpg, jpeg, png, tif, tiff, rtf, bmp, cgm, css, shtml, html, htm, zip, xml, ppt, pptx, tar, ez, ics, mobi, msg, pub and eps' ",
		})
		flags = append(flags, flagSchema{
			Name:        "inline-image-activation",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Status of inline image. inlineImageActivation = false means image can’t be embedded, & inlineImageActivation = true means image can be embedded, in the email. You cannot send a campaign of more than **4MB** with images embedded in the email. Campaigns with the images embedded in the email _must be sent to less than 5000 contacts_. ",
		})
		flags = append(flags, flagSchema{
			Name:        "mirror-active",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Status of mirror links in campaign. mirrorActive = false means mirror links are deactivated, & mirrorActive = true means mirror links are activated, in the campaign",
		})
		flags = append(flags, flagSchema{
			Name:        "recurring",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "**FOR TRIGGER ONLY !** Type of trigger campaign.recurring = false means contact can receive the same Trigger campaign only once, & recurring = true means contact can receive the same Trigger campaign several times ",
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
			Description: "Set this to true if you want to send your campaign at best time. Note:- **if true, warmup ip will be disabled.** ",
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
			Description: "Set an initial quota greater than 1 for warming up your ip. We recommend you set a value of 3000. ",
		})
		flags = append(flags, flagSchema{
			Name:        "increase-rate",
			Type:        "integer",
			Required:    false,
			Location:    "body",
			Description: "Set a percentage increase rate for warming up your ip. We recommend you set the increase rate to 30% per day. If you want to send the same number of emails every day, set the daily increase value to 0%. ",
		})
		flags = append(flags, flagSchema{
			Name:        "unsubscription-page-id",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Enter an unsubscription page id. The page id is a 24 digit alphanumeric id that can be found in the URL when editing the page. ",
		})
		flags = append(flags, flagSchema{
			Name:        "update-form-id",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "**Mandatory if templateId is used containing the {{ update_profile }} tag**. Enter an update profile form id. The form id is a 24 digit alphanumeric id that can be found in the URL when editing the form. ",
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
			Description: "Email campaign updated",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})
		responses = append(responses, responseSchema{
			Status:      "404",
			ContentType: "application/json",
			Description: "Campaign ID not found",
		})

		schema := map[string]any{
			"command":     "update-campaign",
			"description": "Update an email campaign",
			"http": map[string]any{
				"method": "PUT",
				"path":   "/emailCampaigns/{campaignId}",
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
	pathParams["campaignId"] = fmt.Sprintf("%v", emailCampaignsUpdateCampaignFlags.campaignId)

	req := &httpclient.Request{
		Method:      "PUT",
		Path:        httpclient.SubstitutePath("/emailCampaigns/{campaignId}", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if emailCampaignsUpdateCampaignFlags.body != "" {
		if err := json.Unmarshal([]byte(emailCampaignsUpdateCampaignFlags.body), &bodyMap); err != nil {
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
		bodyMap["tag"] = emailCampaignsUpdateCampaignFlags.tag
	}
	if cmd.Flags().Changed("name") {
		bodyMap["name"] = emailCampaignsUpdateCampaignFlags.name
	}
	if cmd.Flags().Changed("html-content") {
		bodyMap["htmlContent"] = emailCampaignsUpdateCampaignFlags.htmlContent
	}
	if cmd.Flags().Changed("html-url") {
		bodyMap["htmlUrl"] = emailCampaignsUpdateCampaignFlags.htmlUrl
	}
	if cmd.Flags().Changed("scheduled-at") {
		bodyMap["scheduledAt"] = emailCampaignsUpdateCampaignFlags.scheduledAt
	}
	if cmd.Flags().Changed("subject") {
		bodyMap["subject"] = emailCampaignsUpdateCampaignFlags.subject
	}
	if cmd.Flags().Changed("preview-text") {
		bodyMap["previewText"] = emailCampaignsUpdateCampaignFlags.previewText
	}
	if cmd.Flags().Changed("reply-to") {
		bodyMap["replyTo"] = emailCampaignsUpdateCampaignFlags.replyTo
	}
	if cmd.Flags().Changed("to-field") {
		bodyMap["toField"] = emailCampaignsUpdateCampaignFlags.toField
	}
	if cmd.Flags().Changed("attachment-url") {
		bodyMap["attachmentUrl"] = emailCampaignsUpdateCampaignFlags.attachmentUrl
	}
	if cmd.Flags().Changed("inline-image-activation") {
		bodyMap["inlineImageActivation"] = emailCampaignsUpdateCampaignFlags.inlineImageActivation
	}
	if cmd.Flags().Changed("mirror-active") {
		bodyMap["mirrorActive"] = emailCampaignsUpdateCampaignFlags.mirrorActive
	}
	if cmd.Flags().Changed("recurring") {
		bodyMap["recurring"] = emailCampaignsUpdateCampaignFlags.recurring
	}
	if cmd.Flags().Changed("footer") {
		bodyMap["footer"] = emailCampaignsUpdateCampaignFlags.footer
	}
	if cmd.Flags().Changed("header") {
		bodyMap["header"] = emailCampaignsUpdateCampaignFlags.header
	}
	if cmd.Flags().Changed("utm-campaign") {
		bodyMap["utmCampaign"] = emailCampaignsUpdateCampaignFlags.utmCampaign
	}
	if cmd.Flags().Changed("send-at-best-time") {
		bodyMap["sendAtBestTime"] = emailCampaignsUpdateCampaignFlags.sendAtBestTime
	}
	if cmd.Flags().Changed("ab-testing") {
		bodyMap["abTesting"] = emailCampaignsUpdateCampaignFlags.abTesting
	}
	if cmd.Flags().Changed("subject-a") {
		bodyMap["subjectA"] = emailCampaignsUpdateCampaignFlags.subjectA
	}
	if cmd.Flags().Changed("subject-b") {
		bodyMap["subjectB"] = emailCampaignsUpdateCampaignFlags.subjectB
	}
	if cmd.Flags().Changed("split-rule") {
		bodyMap["splitRule"] = emailCampaignsUpdateCampaignFlags.splitRule
	}
	if cmd.Flags().Changed("winner-criteria") {
		bodyMap["winnerCriteria"] = emailCampaignsUpdateCampaignFlags.winnerCriteria
	}
	if cmd.Flags().Changed("winner-delay") {
		bodyMap["winnerDelay"] = emailCampaignsUpdateCampaignFlags.winnerDelay
	}
	if cmd.Flags().Changed("ip-warmup-enable") {
		bodyMap["ipWarmupEnable"] = emailCampaignsUpdateCampaignFlags.ipWarmupEnable
	}
	if cmd.Flags().Changed("initial-quota") {
		bodyMap["initialQuota"] = emailCampaignsUpdateCampaignFlags.initialQuota
	}
	if cmd.Flags().Changed("increase-rate") {
		bodyMap["increaseRate"] = emailCampaignsUpdateCampaignFlags.increaseRate
	}
	if cmd.Flags().Changed("unsubscription-page-id") {
		bodyMap["unsubscriptionPageId"] = emailCampaignsUpdateCampaignFlags.unsubscriptionPageId
	}
	if cmd.Flags().Changed("update-form-id") {
		bodyMap["updateFormId"] = emailCampaignsUpdateCampaignFlags.updateFormId
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
