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

var transactionalEmailsDeleteScheduledEmailsCmd = &cobra.Command{
	Use:   "delete-scheduled-emails",
	Short: "Delete scheduled emails by batchId or messageId",
	RunE:  runTransactionalEmailsDeleteScheduledEmails,
}

var transactionalEmailsDeleteScheduledEmailsFlags struct {
	identifier string
}

func init() {
	transactionalEmailsDeleteScheduledEmailsCmd.Flags().StringVar(&transactionalEmailsDeleteScheduledEmailsFlags.identifier, "identifier", "", "The `batchId` of scheduled emails batch (Should be a valid UUIDv4) or the `messageId` of scheduled email.")
	transactionalEmailsDeleteScheduledEmailsCmd.MarkFlagRequired("identifier")

	transactionalEmailsCmd.AddCommand(transactionalEmailsDeleteScheduledEmailsCmd)
}

func runTransactionalEmailsDeleteScheduledEmails(cmd *cobra.Command, args []string) error {
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
			Name:        "identifier",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "The `batchId` of scheduled emails batch (Should be a valid UUIDv4) or the `messageId` of scheduled email.",
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
			Description: "Scheduled email(s) deleted",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "Invalid parameters passed",
		})
		responses = append(responses, responseSchema{
			Status:      "404",
			ContentType: "application/json",
			Description: "Record(s) for identifier not found",
		})

		schema := map[string]any{
			"command":     "delete-scheduled-emails",
			"description": "Delete scheduled emails by batchId or messageId",
			"http": map[string]any{
				"method": "DELETE",
				"path":   "/smtp/email/{identifier}",
			},
			"input": map[string]any{
				"flags":         flags,
				"body_flag":     false,
				"body_required": false,
			},
			"output": map[string]any{
				"responses": responses,
			},
			"semantics": map[string]any{
				"safe":         false,
				"idempotent":   true,
				"reversible":   false,
				"side_effects": []string{"destroys_resource"},
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
	pathParams["identifier"] = fmt.Sprintf("%v", transactionalEmailsDeleteScheduledEmailsFlags.identifier)

	req := &httpclient.Request{
		Method:      "DELETE",
		Path:        httpclient.SubstitutePath("/smtp/email/{identifier}", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

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
