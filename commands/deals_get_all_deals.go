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

var dealsGetAllDealsCmd = &cobra.Command{
	Use:   "get-all-deals",
	Short: "Get all deals",
	RunE:  runDealsGetAllDeals,
}

var dealsGetAllDealsFlags struct {
	filtersAttributesDealName string
	filtersLinkedCompaniesIds string
	filtersLinkedContactsIds  string
	offset                    int
	limit                     int
	sort                      string
}

func init() {
	dealsGetAllDealsCmd.Flags().StringVar(&dealsGetAllDealsFlags.filtersAttributesDealName, "filters-attributes-deal-name", "", "Filter by attributes. If you have a filter for the owner on your end, please send it as filters[attributes.deal_owner] and utilize the account email for the filtering.")
	dealsGetAllDealsCmd.Flags().StringVar(&dealsGetAllDealsFlags.filtersLinkedCompaniesIds, "filters-linked-companies-ids", "", "Filter by linked companies ids")
	dealsGetAllDealsCmd.Flags().StringVar(&dealsGetAllDealsFlags.filtersLinkedContactsIds, "filters-linked-contacts-ids", "", "Filter by linked companies ids")
	dealsGetAllDealsCmd.Flags().IntVar(&dealsGetAllDealsFlags.offset, "offset", 0, "Index of the first document of the page")
	dealsGetAllDealsCmd.Flags().IntVar(&dealsGetAllDealsFlags.limit, "limit", 0, "Number of documents per page")
	dealsGetAllDealsCmd.Flags().StringVar(&dealsGetAllDealsFlags.sort, "sort", "", "Sort the results in the ascending/descending order. Default order is **descending** by creation if `sort` is not passed")

	dealsCmd.AddCommand(dealsGetAllDealsCmd)
}

func runDealsGetAllDeals(cmd *cobra.Command, args []string) error {
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
			Name:        "filters-attributes-deal-name",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter by attributes. If you have a filter for the owner on your end, please send it as filters[attributes.deal_owner] and utilize the account email for the filtering.",
		})
		flags = append(flags, flagSchema{
			Name:        "filters-linked-companies-ids",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter by linked companies ids",
		})
		flags = append(flags, flagSchema{
			Name:        "filters-linked-contacts-ids",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter by linked companies ids",
		})
		flags = append(flags, flagSchema{
			Name:        "offset",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Index of the first document of the page",
		})
		flags = append(flags, flagSchema{
			Name:        "limit",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Number of documents per page",
		})
		flags = append(flags, flagSchema{
			Name:        "sort",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Sort the results in the ascending/descending order. Default order is **descending** by creation if `sort` is not passed",
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
			Description: "Returns deals list with filters",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "Returned when query params are invalid",
		})

		schema := map[string]any{
			"command":     "get-all-deals",
			"description": "Get all deals",
			"http": map[string]any{
				"method": "GET",
				"path":   "/crm/deals",
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

	req := &httpclient.Request{
		Method:      "GET",
		Path:        httpclient.SubstitutePath("/crm/deals", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("filters-attributes-deal-name") {
		req.QueryParams["filters[attributes.deal_name]"] = fmt.Sprintf("%v", dealsGetAllDealsFlags.filtersAttributesDealName)
	}
	if cmd.Flags().Changed("filters-linked-companies-ids") {
		req.QueryParams["filters[linkedCompaniesIds]"] = fmt.Sprintf("%v", dealsGetAllDealsFlags.filtersLinkedCompaniesIds)
	}
	if cmd.Flags().Changed("filters-linked-contacts-ids") {
		req.QueryParams["filters[linkedContactsIds]"] = fmt.Sprintf("%v", dealsGetAllDealsFlags.filtersLinkedContactsIds)
	}
	if cmd.Flags().Changed("offset") {
		req.QueryParams["offset"] = fmt.Sprintf("%v", dealsGetAllDealsFlags.offset)
	}
	if cmd.Flags().Changed("limit") {
		req.QueryParams["limit"] = fmt.Sprintf("%v", dealsGetAllDealsFlags.limit)
	}
	if cmd.Flags().Changed("sort") {
		req.QueryParams["sort"] = fmt.Sprintf("%v", dealsGetAllDealsFlags.sort)
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
