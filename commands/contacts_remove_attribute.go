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

var contactsRemoveAttributeCmd = &cobra.Command{
	Use:   "remove-attribute",
	Short: "Delete an attribute",
	RunE:  runContactsRemoveAttribute,
}

var contactsRemoveAttributeFlags struct {
	attributeCategory string
	attributeName     string
}

func init() {
	contactsRemoveAttributeCmd.Flags().StringVar(&contactsRemoveAttributeFlags.attributeCategory, "attribute-category", "", "Category of the attribute")
	contactsRemoveAttributeCmd.MarkFlagRequired("attribute-category")
	contactsRemoveAttributeCmd.Flags().StringVar(&contactsRemoveAttributeFlags.attributeName, "attribute-name", "", "Name of the existing attribute")
	contactsRemoveAttributeCmd.MarkFlagRequired("attribute-name")

	contactsCmd.AddCommand(contactsRemoveAttributeCmd)
}

func runContactsRemoveAttribute(cmd *cobra.Command, args []string) error {
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
			Name:        "attribute-category",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "Category of the attribute",
		})
		flags = append(flags, flagSchema{
			Name:        "attribute-name",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "Name of the existing attribute",
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
			Description: "Attribute deleted",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})
		responses = append(responses, responseSchema{
			Status:      "404",
			ContentType: "application/json",
			Description: "Attribute not found",
		})

		schema := map[string]any{
			"command":     "remove-attribute",
			"description": "Delete an attribute",
			"http": map[string]any{
				"method": "DELETE",
				"path":   "/contacts/attributes/{attributeCategory}/{attributeName}",
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
	pathParams["attributeCategory"] = fmt.Sprintf("%v", contactsRemoveAttributeFlags.attributeCategory)
	pathParams["attributeName"] = fmt.Sprintf("%v", contactsRemoveAttributeFlags.attributeName)

	req := &httpclient.Request{
		Method:      "DELETE",
		Path:        httpclient.SubstitutePath("/contacts/attributes/{attributeCategory}/{attributeName}", pathParams),
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
