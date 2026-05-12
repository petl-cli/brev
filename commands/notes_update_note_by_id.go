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

var notesUpdateNoteByIdCmd = &cobra.Command{
	Use:   "update-note-by-id",
	Short: "Update a note",
	RunE:  runNotesUpdateNoteById,
}

var notesUpdateNoteByIdFlags struct {
	id         string
	text       string
	contactIds []string
	dealIds    []string
	companyIds []string
	body       string
}

func init() {
	notesUpdateNoteByIdCmd.Flags().StringVar(&notesUpdateNoteByIdFlags.id, "id", "", "Note ID to update")
	notesUpdateNoteByIdCmd.MarkFlagRequired("id")
	notesUpdateNoteByIdCmd.Flags().StringVar(&notesUpdateNoteByIdFlags.text, "text", "", "Text content of a note")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	notesUpdateNoteByIdCmd.Flags().StringSliceVar(&notesUpdateNoteByIdFlags.contactIds, "contact-ids", nil, "Contact Ids linked to a note")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	notesUpdateNoteByIdCmd.Flags().StringSliceVar(&notesUpdateNoteByIdFlags.dealIds, "deal-ids", nil, "Deal Ids linked to a note")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	notesUpdateNoteByIdCmd.Flags().StringSliceVar(&notesUpdateNoteByIdFlags.companyIds, "company-ids", nil, "Company Ids linked to a note")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	notesUpdateNoteByIdCmd.Flags().StringVar(&notesUpdateNoteByIdFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	notesCmd.AddCommand(notesUpdateNoteByIdCmd)
}

func runNotesUpdateNoteById(cmd *cobra.Command, args []string) error {
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
			Name:        "id",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "Note ID to update",
		})
		flags = append(flags, flagSchema{
			Name:        "text",
			Type:        "string",
			Required:    true,
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
			Status:      "204",
			ContentType: "",
			Description: "Note updated successfully",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "Returned when invalid data posted",
		})
		responses = append(responses, responseSchema{
			Status:      "404",
			ContentType: "application/json",
			Description: "Returned when item not found",
		})
		responses = append(responses, responseSchema{
			Status:      "415",
			ContentType: "application/json",
			Description: "Format is not supported",
		})

		schema := map[string]any{
			"command":     "update-note-by-id",
			"description": "Update a note",
			"http": map[string]any{
				"method": "PATCH",
				"path":   "/crm/notes/{id}",
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
	pathParams["id"] = fmt.Sprintf("%v", notesUpdateNoteByIdFlags.id)

	req := &httpclient.Request{
		Method:      "PATCH",
		Path:        httpclient.SubstitutePath("/crm/notes/{id}", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if notesUpdateNoteByIdFlags.body != "" {
		if err := json.Unmarshal([]byte(notesUpdateNoteByIdFlags.body), &bodyMap); err != nil {
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
		bodyMap["text"] = notesUpdateNoteByIdFlags.text
	}
	if cmd.Flags().Changed("contact-ids") {
		bodyMap["contactIds"] = notesUpdateNoteByIdFlags.contactIds
	}
	if cmd.Flags().Changed("deal-ids") {
		bodyMap["dealIds"] = notesUpdateNoteByIdFlags.dealIds
	}
	if cmd.Flags().Changed("company-ids") {
		bodyMap["companyIds"] = notesUpdateNoteByIdFlags.companyIds
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
