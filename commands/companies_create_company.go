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

var companiesCreateCompanyCmd = &cobra.Command{
	Use:   "create-company",
	Short: "Create a company",
	RunE:  runCompaniesCreateCompany,
}

var companiesCreateCompanyFlags struct {
	name        string
	countryCode int
	body        string
}

func init() {
	companiesCreateCompanyCmd.Flags().StringVar(&companiesCreateCompanyFlags.name, "name", "", "Name of company")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	companiesCreateCompanyCmd.Flags().IntVar(&companiesCreateCompanyFlags.countryCode, "country-code", 0, "Country code if phone_number is passed in attributes.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	companiesCreateCompanyCmd.Flags().StringVar(&companiesCreateCompanyFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	companiesCmd.AddCommand(companiesCreateCompanyCmd)
}

func runCompaniesCreateCompany(cmd *cobra.Command, args []string) error {
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
			Name:        "name",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Name of company",
		})
		flags = append(flags, flagSchema{
			Name:        "attributes",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "Attributes for company creation",
		})
		flags = append(flags, flagSchema{
			Name:        "country-code",
			Type:        "integer",
			Required:    false,
			Location:    "body",
			Description: "Country code if phone_number is passed in attributes.",
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
			Description: "Created new Company",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "Returned when invalid data posted",
		})

		schema := map[string]any{
			"command":     "create-company",
			"description": "Create a company",
			"http": map[string]any{
				"method": "POST",
				"path":   "/companies",
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
		Path:        httpclient.SubstitutePath("/companies", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if companiesCreateCompanyFlags.body != "" {
		if err := json.Unmarshal([]byte(companiesCreateCompanyFlags.body), &bodyMap); err != nil {
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
		bodyMap["name"] = companiesCreateCompanyFlags.name
	}
	if cmd.Flags().Changed("country-code") {
		bodyMap["countryCode"] = companiesCreateCompanyFlags.countryCode
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
