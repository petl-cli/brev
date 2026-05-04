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

var filesUploadFileCmd = &cobra.Command{
	Use:   "upload-file",
	Short: "Upload a file",
	RunE:  runFilesUploadFile,
}

var filesUploadFileFlags struct {
	file      string
	dealId    string
	contactId int
	companyId string
	body      string
}

func init() {
	filesUploadFileCmd.Flags().StringVar(&filesUploadFileFlags.file, "file", "", "File data to create a file.")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	filesUploadFileCmd.Flags().StringVar(&filesUploadFileFlags.dealId, "deal-id", "", "")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	filesUploadFileCmd.Flags().IntVar(&filesUploadFileFlags.contactId, "contact-id", 0, "")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	filesUploadFileCmd.Flags().StringVar(&filesUploadFileFlags.companyId, "company-id", "", "")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	filesUploadFileCmd.Flags().StringVar(&filesUploadFileFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	filesCmd.AddCommand(filesUploadFileCmd)
}

func runFilesUploadFile(cmd *cobra.Command, args []string) error {
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
			Name:        "file",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "File data to create a file.",
		})
		flags = append(flags, flagSchema{
			Name:        "deal-id",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "contact-id",
			Type:        "integer",
			Required:    false,
			Location:    "body",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "company-id",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "",
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
			Description: "Returns the created File with additional details",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "Returned when invalid data posted",
		})

		schema := map[string]any{
			"command":     "upload-file",
			"description": "Upload a file",
			"http": map[string]any{
				"method": "POST",
				"path":   "/crm/files",
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
		Path:        httpclient.SubstitutePath("/crm/files", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if filesUploadFileFlags.body != "" {
		if err := json.Unmarshal([]byte(filesUploadFileFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("file") {
		bodyMap["file"] = filesUploadFileFlags.file
	}
	if cmd.Flags().Changed("deal-id") {
		bodyMap["dealId"] = filesUploadFileFlags.dealId
	}
	if cmd.Flags().Changed("contact-id") {
		bodyMap["contactId"] = filesUploadFileFlags.contactId
	}
	if cmd.Flags().Changed("company-id") {
		bodyMap["companyId"] = filesUploadFileFlags.companyId
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
