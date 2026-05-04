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

var notesCreateNewNoteCmd = &cobra.Command{
	Use:   "create-new-note",
	Short: "Create a note",
	RunE:  runNotesCreateNewNote,
}

var notesCreateNewNoteFlags struct {
	text       string
	contactIds []string
	dealIds    []string
	companyIds []string
	body       string
}

func init() {
	notesCreateNewNoteCmd.Flags().StringVar(&notesCreateNewNoteFlags.text, "text", "", "Text content of a note")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	notesCreateNewNoteCmd.Flags().StringSliceVar(&notesCreateNewNoteFlags.contactIds, "contact-ids", nil, "Contact Ids linked to a note")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	notesCreateNewNoteCmd.Flags().StringSliceVar(&notesCreateNewNoteFlags.dealIds, "deal-ids", nil, "Deal Ids linked to a note")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	notesCreateNewNoteCmd.Flags().StringSliceVar(&notesCreateNewNoteFlags.companyIds, "company-ids", nil, "Company Ids linked to a note")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	notesCreateNewNoteCmd.Flags().StringVar(&notesCreateNewNoteFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	notesCmd.AddCommand(notesCreateNewNoteCmd)
}

func runNotesCreateNewNote(cmd *cobra.Command, args []string) error {
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
			Name:        "text",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Text content of a note",
		})
		flags = append(flags, flagSchema{
			Name:        "contact-ids",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "Contact Ids linked to a note",
		})
		flags = append(flags, flagSchema{
			Name:        "deal-ids",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "Deal Ids linked to a note",
		})
		flags = append(flags, flagSchema{
			Name:        "company-ids",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "Company Ids linked to a note",
		})

		type responseSchema struct {
			Status      string `json:"status"`
			ContentType string `json:"content_type,omitempty"`
			Description string `json:"description,omitempty"`
		}
		var responses []responseSchema
		responses = append(responses, responseSchema{
			Status:      "200",
			ContentType: "application/json",
			Description: "Created new note",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "Returned when invalid data posted",
		})
		responses = append(responses, responseSchema{
			Status:      "415",
			ContentType: "application/json",
			Description: "Format is not supported",
		})

		schema := map[string]any{
			"command":     "create-new-note",
			"description": "Create a note",
			"http": map[string]any{
				"method": "POST",
				"path":   "/crm/notes",
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
		Path:        httpclient.SubstitutePath("/crm/notes", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if notesCreateNewNoteFlags.body != "" {
		if err := json.Unmarshal([]byte(notesCreateNewNoteFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("text") {
		bodyMap["text"] = notesCreateNewNoteFlags.text
	}
	if cmd.Flags().Changed("contact-ids") {
		bodyMap["contactIds"] = notesCreateNewNoteFlags.contactIds
	}
	if cmd.Flags().Changed("deal-ids") {
		bodyMap["dealIds"] = notesCreateNewNoteFlags.dealIds
	}
	if cmd.Flags().Changed("company-ids") {
		bodyMap["companyIds"] = notesCreateNewNoteFlags.companyIds
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
