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

var transactionalSmsGetAllEventsCmd = &cobra.Command{
	Use:   "get-all-events",
	Short: "Get all your SMS activity (unaggregated events)",
	RunE:  runTransactionalSmsGetAllEvents,
}

var transactionalSmsGetAllEventsFlags struct {
	limit       int
	startDate   string
	endDate     string
	offset      int
	days        int
	phoneNumber string
	event       string
	tags        string
	sort        string
}

func init() {
	transactionalSmsGetAllEventsCmd.Flags().IntVar(&transactionalSmsGetAllEventsFlags.limit, "limit", 0, "Number of documents per page")
	transactionalSmsGetAllEventsCmd.Flags().StringVar(&transactionalSmsGetAllEventsFlags.startDate, "start-date", "", "**Mandatory if endDate is used.** Starting date (YYYY-MM-DD) of the report ")
	transactionalSmsGetAllEventsCmd.Flags().StringVar(&transactionalSmsGetAllEventsFlags.endDate, "end-date", "", "**Mandatory if startDate is used.** Ending date (YYYY-MM-DD) of the report ")
	transactionalSmsGetAllEventsCmd.Flags().IntVar(&transactionalSmsGetAllEventsFlags.offset, "offset", 0, "Index of the first document of the page")
	transactionalSmsGetAllEventsCmd.Flags().IntVar(&transactionalSmsGetAllEventsFlags.days, "days", 0, "Number of days in the past including today (positive integer). **Not compatible with 'startDate' and 'endDate'** ")
	transactionalSmsGetAllEventsCmd.Flags().StringVar(&transactionalSmsGetAllEventsFlags.phoneNumber, "phone-number", "", "Filter the report for a specific phone number")
	transactionalSmsGetAllEventsCmd.Flags().StringVar(&transactionalSmsGetAllEventsFlags.event, "event", "", "Filter the report for specific events")
	transactionalSmsGetAllEventsCmd.Flags().StringVar(&transactionalSmsGetAllEventsFlags.tags, "tags", "", "Filter the report for specific tags passed as a serialized urlencoded array")
	transactionalSmsGetAllEventsCmd.Flags().StringVar(&transactionalSmsGetAllEventsFlags.sort, "sort", "", "Sort the results in the ascending/descending order of record creation. Default order is **descending** if `sort` is not passed")

	transactionalSmsCmd.AddCommand(transactionalSmsGetAllEventsCmd)
}

func runTransactionalSmsGetAllEvents(cmd *cobra.Command, args []string) error {
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
			Name:        "limit",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Number of documents per page",
		})
		flags = append(flags, flagSchema{
			Name:        "start-date",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "**Mandatory if endDate is used.** Starting date (YYYY-MM-DD) of the report ",
		})
		flags = append(flags, flagSchema{
			Name:        "end-date",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "**Mandatory if startDate is used.** Ending date (YYYY-MM-DD) of the report ",
		})
		flags = append(flags, flagSchema{
			Name:        "offset",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Index of the first document of the page",
		})
		flags = append(flags, flagSchema{
			Name:        "days",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Number of days in the past including today (positive integer). **Not compatible with 'startDate' and 'endDate'** ",
		})
		flags = append(flags, flagSchema{
			Name:        "phone-number",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter the report for a specific phone number",
		})
		flags = append(flags, flagSchema{
			Name:        "event",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter the report for specific events",
		})
		flags = append(flags, flagSchema{
			Name:        "tags",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter the report for specific tags passed as a serialized urlencoded array",
		})
		flags = append(flags, flagSchema{
			Name:        "sort",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Sort the results in the ascending/descending order of record creation. Default order is **descending** if `sort` is not passed",
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
			Description: "Sms events report informations",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "get-all-events",
			"description": "Get all your SMS activity (unaggregated events)",
			"http": map[string]any{
				"method": "GET",
				"path":   "/transactionalSMS/statistics/events",
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
		Path:        httpclient.SubstitutePath("/transactionalSMS/statistics/events", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("limit") {
		req.QueryParams["limit"] = fmt.Sprintf("%v", transactionalSmsGetAllEventsFlags.limit)
	}
	if cmd.Flags().Changed("start-date") {
		req.QueryParams["startDate"] = fmt.Sprintf("%v", transactionalSmsGetAllEventsFlags.startDate)
	}
	if cmd.Flags().Changed("end-date") {
		req.QueryParams["endDate"] = fmt.Sprintf("%v", transactionalSmsGetAllEventsFlags.endDate)
	}
	if cmd.Flags().Changed("offset") {
		req.QueryParams["offset"] = fmt.Sprintf("%v", transactionalSmsGetAllEventsFlags.offset)
	}
	if cmd.Flags().Changed("days") {
		req.QueryParams["days"] = fmt.Sprintf("%v", transactionalSmsGetAllEventsFlags.days)
	}
	if cmd.Flags().Changed("phone-number") {
		req.QueryParams["phoneNumber"] = fmt.Sprintf("%v", transactionalSmsGetAllEventsFlags.phoneNumber)
	}
	if cmd.Flags().Changed("event") {
		req.QueryParams["event"] = fmt.Sprintf("%v", transactionalSmsGetAllEventsFlags.event)
	}
	if cmd.Flags().Changed("tags") {
		req.QueryParams["tags"] = fmt.Sprintf("%v", transactionalSmsGetAllEventsFlags.tags)
	}
	if cmd.Flags().Changed("sort") {
		req.QueryParams["sort"] = fmt.Sprintf("%v", transactionalSmsGetAllEventsFlags.sort)
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
