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

var companiesGetAllCmd = &cobra.Command{
	Use:   "get-all",
	Short: "Get all Companies",
	RunE:  runCompaniesGetAll,
}

var companiesGetAllFlags struct {
	filters           string
	linkedContactsIds int
	linkedDealsIds    string
	page              int
	limit             int
	sort              string
	sortBy            string
}

func init() {
	companiesGetAllCmd.Flags().StringVar(&companiesGetAllFlags.filters, "filters", "", "Filter by attrbutes. If you have filter for owner on your side please send it as {\"attributes.owner\":\"6299dcf3874a14eacbc65c46\"}")
	companiesGetAllCmd.Flags().IntVar(&companiesGetAllFlags.linkedContactsIds, "linked-contacts-ids", 0, "Filter by linked contacts ids")
	companiesGetAllCmd.Flags().StringVar(&companiesGetAllFlags.linkedDealsIds, "linked-deals-ids", "", "Filter by linked Deals ids")
	companiesGetAllCmd.Flags().IntVar(&companiesGetAllFlags.page, "page", 0, "Index of the first document of the page")
	companiesGetAllCmd.Flags().IntVar(&companiesGetAllFlags.limit, "limit", 0, "Number of documents per page")
	companiesGetAllCmd.Flags().StringVar(&companiesGetAllFlags.sort, "sort", "", "Sort the results in the ascending/descending order. Default order is **descending** by creation if `sort` is not passed")
	companiesGetAllCmd.Flags().StringVar(&companiesGetAllFlags.sortBy, "sort-by", "", "The field used to sort field names.")

	companiesCmd.AddCommand(companiesGetAllCmd)
}

func runCompaniesGetAll(cmd *cobra.Command, args []string) error {
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
			Name:        "filters",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter by attrbutes. If you have filter for owner on your side please send it as {\"attributes.owner\":\"6299dcf3874a14eacbc65c46\"}",
		})
		flags = append(flags, flagSchema{
			Name:        "linked-contacts-ids",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Filter by linked contacts ids",
		})
		flags = append(flags, flagSchema{
			Name:        "linked-deals-ids",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter by linked Deals ids",
		})
		flags = append(flags, flagSchema{
			Name:        "page",
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
		flags = append(flags, flagSchema{
			Name:        "sort-by",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "The field used to sort field names.",
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
			Description: "Returns companies list with filters",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "Returned when query params are invalid",
		})

		schema := map[string]any{
			"command":     "get-all",
			"description": "Get all Companies",
			"http": map[string]any{
				"method": "GET",
				"path":   "/companies",
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
		Path:        httpclient.SubstitutePath("/companies", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("filters") {
		req.QueryParams["filters"] = fmt.Sprintf("%v", companiesGetAllFlags.filters)
	}
	if cmd.Flags().Changed("linked-contacts-ids") {
		req.QueryParams["linkedContactsIds"] = fmt.Sprintf("%v", companiesGetAllFlags.linkedContactsIds)
	}
	if cmd.Flags().Changed("linked-deals-ids") {
		req.QueryParams["linkedDealsIds"] = fmt.Sprintf("%v", companiesGetAllFlags.linkedDealsIds)
	}
	if cmd.Flags().Changed("page") {
		req.QueryParams["page"] = fmt.Sprintf("%v", companiesGetAllFlags.page)
	}
	if cmd.Flags().Changed("limit") {
		req.QueryParams["limit"] = fmt.Sprintf("%v", companiesGetAllFlags.limit)
	}
	if cmd.Flags().Changed("sort") {
		req.QueryParams["sort"] = fmt.Sprintf("%v", companiesGetAllFlags.sort)
	}
	if cmd.Flags().Changed("sort-by") {
		req.QueryParams["sortBy"] = fmt.Sprintf("%v", companiesGetAllFlags.sortBy)
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
