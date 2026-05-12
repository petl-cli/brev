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

var contactsCreateDoubleOptInContactCmd = &cobra.Command{
	Use:   "create-double-opt-in-contact",
	Short: "Create Contact via DOI (Double-Opt-In) Flow",
	RunE:  runContactsCreateDoubleOptInContact,
}

var contactsCreateDoubleOptInContactFlags struct {
	email          string
	includeListIds []string
	excludeListIds []string
	templateId     int
	redirectionUrl string
	body           string
}

func init() {
	contactsCreateDoubleOptInContactCmd.Flags().StringVar(&contactsCreateDoubleOptInContactFlags.email, "email", "", "Email address where the confirmation email will be sent. This email address will be the identifier for all other contact attributes.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsCreateDoubleOptInContactCmd.Flags().StringSliceVar(&contactsCreateDoubleOptInContactFlags.includeListIds, "include-list-ids", nil, "Lists under user account where contact should be added")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsCreateDoubleOptInContactCmd.Flags().StringSliceVar(&contactsCreateDoubleOptInContactFlags.excludeListIds, "exclude-list-ids", nil, "Lists under user account where contact should not be added")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsCreateDoubleOptInContactCmd.Flags().IntVar(&contactsCreateDoubleOptInContactFlags.templateId, "template-id", 0, "Id of the Double opt-in (DOI) template")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsCreateDoubleOptInContactCmd.Flags().StringVar(&contactsCreateDoubleOptInContactFlags.redirectionUrl, "redirection-url", "", "URL of the web page that user will be redirected to after clicking on the double opt in URL. When editing your DOI template you can reference this URL by using the tag **{{ params.DOIurl }}**. ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	contactsCreateDoubleOptInContactCmd.Flags().StringVar(&contactsCreateDoubleOptInContactFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	contactsCmd.AddCommand(contactsCreateDoubleOptInContactCmd)
}

func runContactsCreateDoubleOptInContact(cmd *cobra.Command, args []string) error {
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
			Name:        "email",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Email address where the confirmation email will be sent. This email address will be the identifier for all other contact attributes.",
		})
		flags = append(flags, flagSchema{
			Name:        "attributes",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "Pass the set of attributes and their values. **These attributes must be present in your Brevo account**. For eg. **{'FNAME':'Elly', 'LNAME':'Roger'}** ",
		})
		flags = append(flags, flagSchema{
			Name:        "include-list-ids",
			Type:        "array",
			Required:    true,
			Location:    "body",
			Description: "Lists under user account where contact should be added",
		})
		flags = append(flags, flagSchema{
			Name:        "exclude-list-ids",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "Lists under user account where contact should not be added",
		})
		flags = append(flags, flagSchema{
			Name:        "template-id",
			Type:        "integer",
			Required:    true,
			Location:    "body",
			Description: "Id of the Double opt-in (DOI) template",
		})
		flags = append(flags, flagSchema{
			Name:        "redirection-url",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "URL of the web page that user will be redirected to after clicking on the double opt in URL. When editing your DOI template you can reference this URL by using the tag **{{ params.DOIurl }}**. ",
		})

		type responseSchema struct {
			Status      string `json:"status"`
			ContentType string `json:"content_type,omitempty"`
			Description string `json:"description,omitempty"`
		}
		var responses []responseSchema
		responses = append(responses, responseSchema{
			Status:      "201",
			ContentType: "",
			Description: "DOI Contact created",
		})
		responses = append(responses, responseSchema{
			Status:      "204",
			ContentType: "",
			Description: "DOI Contact updated",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "create-double-opt-in-contact",
			"description": "Create Contact via DOI (Double-Opt-In) Flow",
			"http": map[string]any{
				"method": "POST",
				"path":   "/contacts/doubleOptinConfirmation",
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
		Path:        httpclient.SubstitutePath("/contacts/doubleOptinConfirmation", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if contactsCreateDoubleOptInContactFlags.body != "" {
		if err := json.Unmarshal([]byte(contactsCreateDoubleOptInContactFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("email") {
		bodyMap["email"] = contactsCreateDoubleOptInContactFlags.email
	}
	if cmd.Flags().Changed("include-list-ids") {
		bodyMap["includeListIds"] = contactsCreateDoubleOptInContactFlags.includeListIds
	}
	if cmd.Flags().Changed("exclude-list-ids") {
		bodyMap["excludeListIds"] = contactsCreateDoubleOptInContactFlags.excludeListIds
	}
	if cmd.Flags().Changed("template-id") {
		bodyMap["templateId"] = contactsCreateDoubleOptInContactFlags.templateId
	}
	if cmd.Flags().Changed("redirection-url") {
		bodyMap["redirectionUrl"] = contactsCreateDoubleOptInContactFlags.redirectionUrl
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
