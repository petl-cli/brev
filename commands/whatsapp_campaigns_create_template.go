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

var whatsappCampaignsCreateTemplateCmd = &cobra.Command{
	Use:   "create-template",
	Short: "Create a WhatsApp template",
	RunE:  runWhatsappCampaignsCreateTemplate,
}

var whatsappCampaignsCreateTemplateFlags struct {
	name       string
	language   string
	category   string
	mediaUrl   string
	bodyText   string
	headerText string
	source     string
	body       string
}

func init() {
	whatsappCampaignsCreateTemplateCmd.Flags().StringVar(&whatsappCampaignsCreateTemplateFlags.name, "name", "", "Name of the template")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	whatsappCampaignsCreateTemplateCmd.Flags().StringVar(&whatsappCampaignsCreateTemplateFlags.language, "language", "", "Language of the template. For Example : **en** for English ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	whatsappCampaignsCreateTemplateCmd.Flags().StringVar(&whatsappCampaignsCreateTemplateFlags.category, "category", "", "Category of the template")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	whatsappCampaignsCreateTemplateCmd.Flags().StringVar(&whatsappCampaignsCreateTemplateFlags.mediaUrl, "media-url", "", "Absolute url of the media file **(no local file)** for the header. **Use this field in you want to add media in Template header and headerText is empty**. Allowed extensions for media files are: #### jpeg | png | mp4 | pdf ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	whatsappCampaignsCreateTemplateCmd.Flags().StringVar(&whatsappCampaignsCreateTemplateFlags.bodyText, "body-text", "", "Body of the template. **Maximum allowed characters are 1024**")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	whatsappCampaignsCreateTemplateCmd.Flags().StringVar(&whatsappCampaignsCreateTemplateFlags.headerText, "header-text", "", "Text content of the header in the template. **Maximum allowed characters are 45** **Use this field to add text content in template header and if mediaUrl is empty** ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	whatsappCampaignsCreateTemplateCmd.Flags().StringVar(&whatsappCampaignsCreateTemplateFlags.source, "source", "", "source of the template")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	whatsappCampaignsCreateTemplateCmd.Flags().StringVar(&whatsappCampaignsCreateTemplateFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	whatsappCampaignsCmd.AddCommand(whatsappCampaignsCreateTemplateCmd)
}

func runWhatsappCampaignsCreateTemplate(cmd *cobra.Command, args []string) error {
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
			Name:        "name",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Name of the template",
		})
		flags = append(flags, flagSchema{
			Name:        "language",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Language of the template. For Example : **en** for English ",
		})
		flags = append(flags, flagSchema{
			Name:        "category",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Category of the template",
		})
		flags = append(flags, flagSchema{
			Name:        "media-url",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Absolute url of the media file **(no local file)** for the header. **Use this field in you want to add media in Template header and headerText is empty**. Allowed extensions for media files are: #### jpeg | png | mp4 | pdf ",
		})
		flags = append(flags, flagSchema{
			Name:        "body-text",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Body of the template. **Maximum allowed characters are 1024**",
		})
		flags = append(flags, flagSchema{
			Name:        "header-text",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Text content of the header in the template. **Maximum allowed characters are 45** **Use this field to add text content in template header and if mediaUrl is empty** ",
		})
		flags = append(flags, flagSchema{
			Name:        "source",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "source of the template",
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
			"description": "Create a WhatsApp template",
			"http": map[string]any{
				"method": "POST",
				"path":   "/whatsappCampaigns/template",
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
		Path:        httpclient.SubstitutePath("/whatsappCampaigns/template", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if whatsappCampaignsCreateTemplateFlags.body != "" {
		if err := json.Unmarshal([]byte(whatsappCampaignsCreateTemplateFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("name") {
		bodyMap["name"] = whatsappCampaignsCreateTemplateFlags.name
	}
	if cmd.Flags().Changed("language") {
		bodyMap["language"] = whatsappCampaignsCreateTemplateFlags.language
	}
	if cmd.Flags().Changed("category") {
		bodyMap["category"] = whatsappCampaignsCreateTemplateFlags.category
	}
	if cmd.Flags().Changed("media-url") {
		bodyMap["mediaUrl"] = whatsappCampaignsCreateTemplateFlags.mediaUrl
	}
	if cmd.Flags().Changed("body-text") {
		bodyMap["bodyText"] = whatsappCampaignsCreateTemplateFlags.bodyText
	}
	if cmd.Flags().Changed("header-text") {
		bodyMap["headerText"] = whatsappCampaignsCreateTemplateFlags.headerText
	}
	if cmd.Flags().Changed("source") {
		bodyMap["source"] = whatsappCampaignsCreateTemplateFlags.source
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
