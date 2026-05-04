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

var transactionalEmailsSendTestTemplateCmd = &cobra.Command{
	Use:   "send-test-template",
	Short: "Send a template to your test list",
	RunE:  runTransactionalEmailsSendTestTemplate,
}

var transactionalEmailsSendTestTemplateFlags struct {
	templateId int
	emailTo    []string
	body       string
}

func init() {
	transactionalEmailsSendTestTemplateCmd.Flags().IntVar(&transactionalEmailsSendTestTemplateFlags.templateId, "template-id", 0, "Id of the template")
	transactionalEmailsSendTestTemplateCmd.MarkFlagRequired("template-id")
	transactionalEmailsSendTestTemplateCmd.Flags().StringSliceVar(&transactionalEmailsSendTestTemplateFlags.emailTo, "email-to", nil, "List of the email addresses of the recipients whom you wish to send the test mail. _If left empty, the test mail will be sent to your entire test list. You can not send more than 50 test emails per day_. ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	transactionalEmailsSendTestTemplateCmd.Flags().StringVar(&transactionalEmailsSendTestTemplateFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	transactionalEmailsCmd.AddCommand(transactionalEmailsSendTestTemplateCmd)
}

func runTransactionalEmailsSendTestTemplate(cmd *cobra.Command, args []string) error {
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
			Description: "Id of the template",
		})
		flags = append(flags, flagSchema{
			Name:        "email-to",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "List of the email addresses of the recipients whom you wish to send the test mail. _If left empty, the test mail will be sent to your entire test list. You can not send more than 50 test emails per day_. ",
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
			Description: "Test email has been sent successfully to all recipients",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "Test email could not be sent to the following email addresses",
		})
		responses = append(responses, responseSchema{
			Status:      "404",
			ContentType: "application/json",
			Description: "Template ID not found",
		})

		schema := map[string]any{
			"command":     "send-test-template",
			"description": "Send a template to your test list",
			"http": map[string]any{
				"method": "POST",
				"path":   "/smtp/templates/{templateId}/sendTest",
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
	pathParams["templateId"] = fmt.Sprintf("%v", transactionalEmailsSendTestTemplateFlags.templateId)

	req := &httpclient.Request{
		Method:      "POST",
		Path:        httpclient.SubstitutePath("/smtp/templates/{templateId}/sendTest", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if transactionalEmailsSendTestTemplateFlags.body != "" {
		if err := json.Unmarshal([]byte(transactionalEmailsSendTestTemplateFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("email-to") {
		bodyMap["emailTo"] = transactionalEmailsSendTestTemplateFlags.emailTo
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
