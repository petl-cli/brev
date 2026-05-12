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

var contactsExportProcessIdCmd = &cobra.Command{
	Use:   "export-process-id",
	Short: "Export contacts",
	RunE:  runContactsExportProcessId,
}

var contactsExportProcessIdFlags struct {
	exportAttributes []string
	notifyUrl        string
	body             string
}

func init() {
	contactsExportProcessIdCmd.Flags().StringSliceVar(&contactsExportProcessIdFlags.exportAttributes, "export-attributes", nil, "List of all the attributes that you want to export. **These attributes must be present in your contact database.** For example: **['fname', 'lname', 'email']** ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsExportProcessIdCmd.Flags().StringVar(&contactsExportProcessIdFlags.notifyUrl, "notify-url", "", "Webhook that will be called once the export process is finished. For reference, https://help.brevo.com/hc/en-us/articles/360007666479")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsExportProcessIdCmd.Flags().StringVar(&contactsExportProcessIdFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	contactsCmd.AddCommand(contactsExportProcessIdCmd)
}

func runContactsExportProcessId(cmd *cobra.Command, args []string) error {
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
			Name:        "export-attributes",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "List of all the attributes that you want to export. **These attributes must be present in your contact database.** For example: **['fname', 'lname', 'email']** ",
		})
		flags = append(flags, flagSchema{
			Name:        "custom-contact-filter",
			Type:        "object",
			Required:    true,
			Location:    "body",
			Description: "Set the filter for the contacts to be exported. ",
		})
		flags = append(flags, flagSchema{
			Name:        "notify-url",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Webhook that will be called once the export process is finished. For reference, https://help.brevo.com/hc/en-us/articles/360007666479",
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
			"command":     "export-process-id",
			"description": "Export contacts",
			"http": map[string]any{
				"method": "POST",
				"path":   "/contacts/export",
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
		Path:        httpclient.SubstitutePath("/contacts/export", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if contactsExportProcessIdFlags.body != "" {
		if err := json.Unmarshal([]byte(contactsExportProcessIdFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("export-attributes") {
		bodyMap["exportAttributes"] = contactsExportProcessIdFlags.exportAttributes
	}
	if cmd.Flags().Changed("notify-url") {
		bodyMap["notifyUrl"] = contactsExportProcessIdFlags.notifyUrl
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
