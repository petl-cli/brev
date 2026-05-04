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

var tasksGetAllCmd = &cobra.Command{
	Use:   "get-all",
	Short: "Get all tasks",
	RunE:  runTasksGetAll,
}

var tasksGetAllFlags struct {
	filterType      string
	filterStatus    string
	filterDate      string
	filterAssignTo  string
	filterContacts  string
	filterDeals     string
	filterCompanies string
	dateFrom        int
	dateTo          int
	offset          int
	limit           int
	sort            string
	sortBy          string
}

func init() {
	tasksGetAllCmd.Flags().StringVar(&tasksGetAllFlags.filterType, "filter-type", "", "Filter by task type (ID)")
	tasksGetAllCmd.Flags().StringVar(&tasksGetAllFlags.filterStatus, "filter-status", "", "Filter by task status")
	tasksGetAllCmd.Flags().StringVar(&tasksGetAllFlags.filterDate, "filter-date", "", "Filter by date")
	tasksGetAllCmd.Flags().StringVar(&tasksGetAllFlags.filterAssignTo, "filter-assign-to", "", "Filter by the \"assignTo\" ID. You can utilize account emails for the \"assignTo\" attribute.")
	tasksGetAllCmd.Flags().StringVar(&tasksGetAllFlags.filterContacts, "filter-contacts", "", "Filter by contact ids")
	tasksGetAllCmd.Flags().StringVar(&tasksGetAllFlags.filterDeals, "filter-deals", "", "Filter by deals ids")
	tasksGetAllCmd.Flags().StringVar(&tasksGetAllFlags.filterCompanies, "filter-companies", "", "Filter by companies ids")
	tasksGetAllCmd.Flags().IntVar(&tasksGetAllFlags.dateFrom, "date-from", 0, "dateFrom to date range filter type (timestamp in milliseconds)")
	tasksGetAllCmd.Flags().IntVar(&tasksGetAllFlags.dateTo, "date-to", 0, "dateTo to date range filter type (timestamp in milliseconds)")
	tasksGetAllCmd.Flags().IntVar(&tasksGetAllFlags.offset, "offset", 0, "Index of the first document of the page")
	tasksGetAllCmd.Flags().IntVar(&tasksGetAllFlags.limit, "limit", 0, "Number of documents per page")
	tasksGetAllCmd.Flags().StringVar(&tasksGetAllFlags.sort, "sort", "", "Sort the results in the ascending/descending order. Default order is **descending** by creation if `sort` is not passed")
	tasksGetAllCmd.Flags().StringVar(&tasksGetAllFlags.sortBy, "sort-by", "", "The field used to sort field names.")

	tasksCmd.AddCommand(tasksGetAllCmd)
}

func runTasksGetAll(cmd *cobra.Command, args []string) error {
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
			Name:        "filter-type",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter by task type (ID)",
		})
		flags = append(flags, flagSchema{
			Name:        "filter-status",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter by task status",
		})
		flags = append(flags, flagSchema{
			Name:        "filter-date",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter by date",
		})
		flags = append(flags, flagSchema{
			Name:        "filter-assign-to",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter by the \"assignTo\" ID. You can utilize account emails for the \"assignTo\" attribute.",
		})
		flags = append(flags, flagSchema{
			Name:        "filter-contacts",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter by contact ids",
		})
		flags = append(flags, flagSchema{
			Name:        "filter-deals",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter by deals ids",
		})
		flags = append(flags, flagSchema{
			Name:        "filter-companies",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter by companies ids",
		})
		flags = append(flags, flagSchema{
			Name:        "date-from",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "dateFrom to date range filter type (timestamp in milliseconds)",
		})
		flags = append(flags, flagSchema{
			Name:        "date-to",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "dateTo to date range filter type (timestamp in milliseconds)",
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
			Description: "Returns task list with filters",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "Returned when query params are invalid",
		})

		schema := map[string]any{
			"command":     "get-all",
			"description": "Get all tasks",
			"http": map[string]any{
				"method": "GET",
				"path":   "/crm/tasks",
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
		Path:        httpclient.SubstitutePath("/crm/tasks", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("filter-type") {
		req.QueryParams["filter[type]"] = fmt.Sprintf("%v", tasksGetAllFlags.filterType)
	}
	if cmd.Flags().Changed("filter-status") {
		req.QueryParams["filter[status]"] = fmt.Sprintf("%v", tasksGetAllFlags.filterStatus)
	}
	if cmd.Flags().Changed("filter-date") {
		req.QueryParams["filter[date]"] = fmt.Sprintf("%v", tasksGetAllFlags.filterDate)
	}
	if cmd.Flags().Changed("filter-assign-to") {
		req.QueryParams["filter[assignTo]"] = fmt.Sprintf("%v", tasksGetAllFlags.filterAssignTo)
	}
	if cmd.Flags().Changed("filter-contacts") {
		req.QueryParams["filter[contacts]"] = fmt.Sprintf("%v", tasksGetAllFlags.filterContacts)
	}
	if cmd.Flags().Changed("filter-deals") {
		req.QueryParams["filter[deals]"] = fmt.Sprintf("%v", tasksGetAllFlags.filterDeals)
	}
	if cmd.Flags().Changed("filter-companies") {
		req.QueryParams["filter[companies]"] = fmt.Sprintf("%v", tasksGetAllFlags.filterCompanies)
	}
	if cmd.Flags().Changed("date-from") {
		req.QueryParams["dateFrom"] = fmt.Sprintf("%v", tasksGetAllFlags.dateFrom)
	}
	if cmd.Flags().Changed("date-to") {
		req.QueryParams["dateTo"] = fmt.Sprintf("%v", tasksGetAllFlags.dateTo)
	}
	if cmd.Flags().Changed("offset") {
		req.QueryParams["offset"] = fmt.Sprintf("%v", tasksGetAllFlags.offset)
	}
	if cmd.Flags().Changed("limit") {
		req.QueryParams["limit"] = fmt.Sprintf("%v", tasksGetAllFlags.limit)
	}
	if cmd.Flags().Changed("sort") {
		req.QueryParams["sort"] = fmt.Sprintf("%v", tasksGetAllFlags.sort)
	}
	if cmd.Flags().Changed("sort-by") {
		req.QueryParams["sortBy"] = fmt.Sprintf("%v", tasksGetAllFlags.sortBy)
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
