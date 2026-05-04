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

var contactsUpdateAttributeCmd = &cobra.Command{
	Use:   "update-attribute",
	Short: "Update contact attribute",
	RunE:  runContactsUpdateAttribute,
}

var contactsUpdateAttributeFlags struct {
	attributeCategory string
	attributeName     string
	value             string
	enumeration       []string
	body              string
}

func init() {
	contactsUpdateAttributeCmd.Flags().StringVar(&contactsUpdateAttributeFlags.attributeCategory, "attribute-category", "", "Category of the attribute")
	contactsUpdateAttributeCmd.MarkFlagRequired("attribute-category")
	contactsUpdateAttributeCmd.Flags().StringVar(&contactsUpdateAttributeFlags.attributeName, "attribute-name", "", "Name of the existing attribute")
	contactsUpdateAttributeCmd.MarkFlagRequired("attribute-name")
	contactsUpdateAttributeCmd.Flags().StringVar(&contactsUpdateAttributeFlags.value, "value", "", "Value of the attribute to update. **Use only if the attribute's category is 'calculated' or 'global'** ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	contactsUpdateAttributeCmd.Flags().StringSliceVar(&contactsUpdateAttributeFlags.enumeration, "enumeration", nil, "List of the values and labels that the attribute can take. **Use only if the attribute's category is \"category\"**. For example, **[{\"value\":1, \"label\":\"male\"}, {\"value\":2, \"label\":\"female\"}]** ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	contactsUpdateAttributeCmd.Flags().StringVar(&contactsUpdateAttributeFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	contactsCmd.AddCommand(contactsUpdateAttributeCmd)
}

func runContactsUpdateAttribute(cmd *cobra.Command, args []string) error {
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
		flags = append(flags, flagSchema{
			Name:        "value",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Value of the attribute to update. **Use only if the attribute's category is 'calculated' or 'global'** ",
		})
		flags = append(flags, flagSchema{
			Name:        "enumeration",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "List of the values and labels that the attribute can take. **Use only if the attribute's category is \"category\"**. For example, **[{\"value\":1, \"label\":\"male\"}, {\"value\":2, \"label\":\"female\"}]** ",
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
			Description: "Attribute updated",
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
			"command":     "update-attribute",
			"description": "Update contact attribute",
			"http": map[string]any{
				"method": "PUT",
				"path":   "/contacts/attributes/{attributeCategory}/{attributeName}",
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
	pathParams["attributeCategory"] = fmt.Sprintf("%v", contactsUpdateAttributeFlags.attributeCategory)
	pathParams["attributeName"] = fmt.Sprintf("%v", contactsUpdateAttributeFlags.attributeName)

	req := &httpclient.Request{
		Method:      "PUT",
		Path:        httpclient.SubstitutePath("/contacts/attributes/{attributeCategory}/{attributeName}", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if contactsUpdateAttributeFlags.body != "" {
		if err := json.Unmarshal([]byte(contactsUpdateAttributeFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("value") {
		bodyMap["value"] = contactsUpdateAttributeFlags.value
	}
	if cmd.Flags().Changed("enumeration") {
		bodyMap["enumeration"] = contactsUpdateAttributeFlags.enumeration
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
