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

var contactsGetDetailsCmd = &cobra.Command{
	Use:   "get-details",
	Short: "Get a contact's details",
	RunE:  runContactsGetDetails,
}

var contactsGetDetailsFlags struct {
	identifier string
	startDate  string
	endDate    string
}

func init() {
	contactsGetDetailsCmd.Flags().StringVar(&contactsGetDetailsFlags.identifier, "identifier", "", "Email (urlencoded) OR ID of the contact OR its SMS attribute value")
	contactsGetDetailsCmd.MarkFlagRequired("identifier")
	contactsGetDetailsCmd.Flags().StringVar(&contactsGetDetailsFlags.startDate, "start-date", "", "**Mandatory if endDate is used.** Starting date (YYYY-MM-DD) of the statistic events specific to campaigns. Must be lower than equal to endDate ")
	contactsGetDetailsCmd.Flags().StringVar(&contactsGetDetailsFlags.endDate, "end-date", "", "**Mandatory if startDate is used.** Ending date (YYYY-MM-DD) of the statistic events specific to campaigns. Must be greater than equal to startDate. ")

	contactsCmd.AddCommand(contactsGetDetailsCmd)
}

func runContactsGetDetails(cmd *cobra.Command, args []string) error {
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
			Description: "Email (urlencoded) OR ID of the contact OR its SMS attribute value",
		})
		flags = append(flags, flagSchema{
			Name:        "start-date",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "**Mandatory if endDate is used.** Starting date (YYYY-MM-DD) of the statistic events specific to campaigns. Must be lower than equal to endDate ",
		})
		flags = append(flags, flagSchema{
			Name:        "end-date",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "**Mandatory if startDate is used.** Ending date (YYYY-MM-DD) of the statistic events specific to campaigns. Must be greater than equal to startDate. ",
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
			Description: "Contact informations",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})
		responses = append(responses, responseSchema{
			Status:      "404",
			ContentType: "application/json",
			Description: "Contact's email not found",
		})

		schema := map[string]any{
			"command":     "get-details",
			"description": "Get a contact's details",
			"http": map[string]any{
				"method": "GET",
				"path":   "/contacts/{identifier}",
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
				"safe":         true,
				"idempotent":   true,
				"reversible":   true,
				"side_effects": []string{},
				"impact":       "low",
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
	pathParams["identifier"] = fmt.Sprintf("%v", contactsGetDetailsFlags.identifier)

	req := &httpclient.Request{
		Method:      "GET",
		Path:        httpclient.SubstitutePath("/contacts/{identifier}", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("start-date") {
		req.QueryParams["startDate"] = fmt.Sprintf("%v", contactsGetDetailsFlags.startDate)
	}
	if cmd.Flags().Changed("end-date") {
		req.QueryParams["endDate"] = fmt.Sprintf("%v", contactsGetDetailsFlags.endDate)
	}

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
