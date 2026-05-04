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

var transactionalEmailsGetAllActivityCmd = &cobra.Command{
	Use:   "get-all-activity",
	Short: "Get all your transactional email activity (unaggregated events)",
	RunE:  runTransactionalEmailsGetAllActivity,
}

var transactionalEmailsGetAllActivityFlags struct {
	limit      int
	offset     int
	startDate  string
	endDate    string
	days       int
	email      string
	event      string
	tags       string
	messageId  string
	templateId int
	sort       string
}

func init() {
	transactionalEmailsGetAllActivityCmd.Flags().IntVar(&transactionalEmailsGetAllActivityFlags.limit, "limit", 0, "Number limitation for the result returned")
	transactionalEmailsGetAllActivityCmd.Flags().IntVar(&transactionalEmailsGetAllActivityFlags.offset, "offset", 0, "Beginning point in the list to retrieve from.")
	transactionalEmailsGetAllActivityCmd.Flags().StringVar(&transactionalEmailsGetAllActivityFlags.startDate, "start-date", "", "**Mandatory if endDate is used.** Starting date of the report (YYYY-MM-DD). Must be lower than equal to endDate ")
	transactionalEmailsGetAllActivityCmd.Flags().StringVar(&transactionalEmailsGetAllActivityFlags.endDate, "end-date", "", "**Mandatory if startDate is used.** Ending date of the report (YYYY-MM-DD). Must be greater than equal to startDate ")
	transactionalEmailsGetAllActivityCmd.Flags().IntVar(&transactionalEmailsGetAllActivityFlags.days, "days", 0, "Number of days in the past including today (positive integer). _Not compatible with 'startDate' and 'endDate'_ ")
	transactionalEmailsGetAllActivityCmd.Flags().StringVar(&transactionalEmailsGetAllActivityFlags.email, "email", "", "Filter the report for a specific email addresses")
	transactionalEmailsGetAllActivityCmd.Flags().StringVar(&transactionalEmailsGetAllActivityFlags.event, "event", "", "Filter the report for a specific event type")
	transactionalEmailsGetAllActivityCmd.Flags().StringVar(&transactionalEmailsGetAllActivityFlags.tags, "tags", "", "Filter the report for tags (serialized and urlencoded array)")
	transactionalEmailsGetAllActivityCmd.Flags().StringVar(&transactionalEmailsGetAllActivityFlags.messageId, "message-id", "", "Filter on a specific message id")
	transactionalEmailsGetAllActivityCmd.Flags().IntVar(&transactionalEmailsGetAllActivityFlags.templateId, "template-id", 0, "Filter on a specific template id")
	transactionalEmailsGetAllActivityCmd.Flags().StringVar(&transactionalEmailsGetAllActivityFlags.sort, "sort", "", "Sort the results in the ascending/descending order of record creation. Default order is **descending** if `sort` is not passed")

	transactionalEmailsCmd.AddCommand(transactionalEmailsGetAllActivityCmd)
}

func runTransactionalEmailsGetAllActivity(cmd *cobra.Command, args []string) error {
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
			Description: "Number limitation for the result returned",
		})
		flags = append(flags, flagSchema{
			Name:        "offset",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Beginning point in the list to retrieve from.",
		})
		flags = append(flags, flagSchema{
			Name:        "start-date",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "**Mandatory if endDate is used.** Starting date of the report (YYYY-MM-DD). Must be lower than equal to endDate ",
		})
		flags = append(flags, flagSchema{
			Name:        "end-date",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "**Mandatory if startDate is used.** Ending date of the report (YYYY-MM-DD). Must be greater than equal to startDate ",
		})
		flags = append(flags, flagSchema{
			Name:        "days",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Number of days in the past including today (positive integer). _Not compatible with 'startDate' and 'endDate'_ ",
		})
		flags = append(flags, flagSchema{
			Name:        "email",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter the report for a specific email addresses",
		})
		flags = append(flags, flagSchema{
			Name:        "event",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter the report for a specific event type",
		})
		flags = append(flags, flagSchema{
			Name:        "tags",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter the report for tags (serialized and urlencoded array)",
		})
		flags = append(flags, flagSchema{
			Name:        "message-id",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter on a specific message id",
		})
		flags = append(flags, flagSchema{
			Name:        "template-id",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Filter on a specific template id",
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
			Description: "Email events report informations",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "get-all-activity",
			"description": "Get all your transactional email activity (unaggregated events)",
			"http": map[string]any{
				"method": "GET",
				"path":   "/smtp/statistics/events",
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
		Path:        httpclient.SubstitutePath("/smtp/statistics/events", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("limit") {
		req.QueryParams["limit"] = fmt.Sprintf("%v", transactionalEmailsGetAllActivityFlags.limit)
	}
	if cmd.Flags().Changed("offset") {
		req.QueryParams["offset"] = fmt.Sprintf("%v", transactionalEmailsGetAllActivityFlags.offset)
	}
	if cmd.Flags().Changed("start-date") {
		req.QueryParams["startDate"] = fmt.Sprintf("%v", transactionalEmailsGetAllActivityFlags.startDate)
	}
	if cmd.Flags().Changed("end-date") {
		req.QueryParams["endDate"] = fmt.Sprintf("%v", transactionalEmailsGetAllActivityFlags.endDate)
	}
	if cmd.Flags().Changed("days") {
		req.QueryParams["days"] = fmt.Sprintf("%v", transactionalEmailsGetAllActivityFlags.days)
	}
	if cmd.Flags().Changed("email") {
		req.QueryParams["email"] = fmt.Sprintf("%v", transactionalEmailsGetAllActivityFlags.email)
	}
	if cmd.Flags().Changed("event") {
		req.QueryParams["event"] = fmt.Sprintf("%v", transactionalEmailsGetAllActivityFlags.event)
	}
	if cmd.Flags().Changed("tags") {
		req.QueryParams["tags"] = fmt.Sprintf("%v", transactionalEmailsGetAllActivityFlags.tags)
	}
	if cmd.Flags().Changed("message-id") {
		req.QueryParams["messageId"] = fmt.Sprintf("%v", transactionalEmailsGetAllActivityFlags.messageId)
	}
	if cmd.Flags().Changed("template-id") {
		req.QueryParams["templateId"] = fmt.Sprintf("%v", transactionalEmailsGetAllActivityFlags.templateId)
	}
	if cmd.Flags().Changed("sort") {
		req.QueryParams["sort"] = fmt.Sprintf("%v", transactionalEmailsGetAllActivityFlags.sort)
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
