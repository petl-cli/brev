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

var contactsUpdateListCmd = &cobra.Command{
	Use:   "update-list",
	Short: "Update a list",
	RunE:  runContactsUpdateList,
}

var contactsUpdateListFlags struct {
	listId   int
	name     string
	folderId int
	body     string
}

func init() {
	contactsUpdateListCmd.Flags().IntVar(&contactsUpdateListFlags.listId, "list-id", 0, "Id of the list")
	contactsUpdateListCmd.MarkFlagRequired("list-id")
	contactsUpdateListCmd.Flags().StringVar(&contactsUpdateListFlags.name, "name", "", "Name of the list. Either of the two parameters (name, folderId) can be updated at a time.")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	contactsUpdateListCmd.Flags().IntVar(&contactsUpdateListFlags.folderId, "folder-id", 0, "Id of the folder in which the list is to be moved. Either of the two parameters (name, folderId) can be updated at a time.")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	contactsUpdateListCmd.Flags().StringVar(&contactsUpdateListFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	contactsCmd.AddCommand(contactsUpdateListCmd)
}

func runContactsUpdateList(cmd *cobra.Command, args []string) error {
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
			Name:        "list-id",
			Type:        "integer",
			Required:    true,
			Location:    "path",
			Description: "Id of the list",
		})
		flags = append(flags, flagSchema{
			Name:        "name",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Name of the list. Either of the two parameters (name, folderId) can be updated at a time.",
		})
		flags = append(flags, flagSchema{
			Name:        "folder-id",
			Type:        "integer",
			Required:    false,
			Location:    "body",
			Description: "Id of the folder in which the list is to be moved. Either of the two parameters (name, folderId) can be updated at a time.",
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
			Description: "List updated",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})
		responses = append(responses, responseSchema{
			Status:      "404",
			ContentType: "application/json",
			Description: "List ID not found",
		})

		schema := map[string]any{
			"command":     "update-list",
			"description": "Update a list",
			"http": map[string]any{
				"method": "PUT",
				"path":   "/contacts/lists/{listId}",
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
	pathParams["listId"] = fmt.Sprintf("%v", contactsUpdateListFlags.listId)

	req := &httpclient.Request{
		Method:      "PUT",
		Path:        httpclient.SubstitutePath("/contacts/lists/{listId}", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if contactsUpdateListFlags.body != "" {
		if err := json.Unmarshal([]byte(contactsUpdateListFlags.body), &bodyMap); err != nil {
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
		bodyMap["name"] = contactsUpdateListFlags.name
	}
	if cmd.Flags().Changed("folder-id") {
		bodyMap["folderId"] = contactsUpdateListFlags.folderId
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
