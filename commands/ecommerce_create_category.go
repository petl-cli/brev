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

var ecommerceCreateCategoryCmd = &cobra.Command{
	Use:   "create-category",
	Short: "Create/Update a category",
	RunE:  runEcommerceCreateCategory,
}

var ecommerceCreateCategoryFlags struct {
	id            string
	name          string
	url           string
	updateEnabled bool
	deletedAt     string
	body          string
}

func init() {
	ecommerceCreateCategoryCmd.Flags().StringVar(&ecommerceCreateCategoryFlags.id, "id", "", "Unique Category ID as saved in the shop ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	ecommerceCreateCategoryCmd.Flags().StringVar(&ecommerceCreateCategoryFlags.name, "name", "", "**Mandatory in case of creation**. Name of the Category, as displayed in the shop ")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	ecommerceCreateCategoryCmd.Flags().StringVar(&ecommerceCreateCategoryFlags.url, "url", "", "URL to the category")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	ecommerceCreateCategoryCmd.Flags().BoolVar(&ecommerceCreateCategoryFlags.updateEnabled, "update-enabled", false, "Facilitate to update the existing category in the same request (updateEnabled = true)")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	ecommerceCreateCategoryCmd.Flags().StringVar(&ecommerceCreateCategoryFlags.deletedAt, "deleted-at", "", "UTC date-time (YYYY-MM-DDTHH:mm:ss.SSSZ) of the category deleted from the shop's database")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	ecommerceCreateCategoryCmd.Flags().StringVar(&ecommerceCreateCategoryFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	ecommerceCmd.AddCommand(ecommerceCreateCategoryCmd)
}

func runEcommerceCreateCategory(cmd *cobra.Command, args []string) error {
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
			Required:    false,
			Location:    "body",
			Description: "Unique Category ID as saved in the shop ",
		})
		flags = append(flags, flagSchema{
			Name:        "name",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "**Mandatory in case of creation**. Name of the Category, as displayed in the shop ",
		})
		flags = append(flags, flagSchema{
			Name:        "url",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "URL to the category",
		})
		flags = append(flags, flagSchema{
			Name:        "update-enabled",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Facilitate to update the existing category in the same request (updateEnabled = true)",
		})
		flags = append(flags, flagSchema{
			Name:        "deleted-at",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "UTC date-time (YYYY-MM-DDTHH:mm:ss.SSSZ) of the category deleted from the shop's database",
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
			Description: "Category created",
		})
		responses = append(responses, responseSchema{
			Status:      "204",
			ContentType: "",
			Description: "Category updated",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "create-category",
			"description": "Create/Update a category",
			"http": map[string]any{
				"method": "POST",
				"path":   "/categories",
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
		Path:        httpclient.SubstitutePath("/categories", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if ecommerceCreateCategoryFlags.body != "" {
		if err := json.Unmarshal([]byte(ecommerceCreateCategoryFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("id") {
		bodyMap["id"] = ecommerceCreateCategoryFlags.id
	}
	if cmd.Flags().Changed("name") {
		bodyMap["name"] = ecommerceCreateCategoryFlags.name
	}
	if cmd.Flags().Changed("url") {
		bodyMap["url"] = ecommerceCreateCategoryFlags.url
	}
	if cmd.Flags().Changed("update-enabled") {
		bodyMap["updateEnabled"] = ecommerceCreateCategoryFlags.updateEnabled
	}
	if cmd.Flags().Changed("deleted-at") {
		bodyMap["deletedAt"] = ecommerceCreateCategoryFlags.deletedAt
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
