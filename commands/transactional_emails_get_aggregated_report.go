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

var transactionalEmailsGetAggregatedReportCmd = &cobra.Command{
	Use:   "get-aggregated-report",
	Short: "Get your transactional email activity aggregated over a period of time",
	RunE:  runTransactionalEmailsGetAggregatedReport,
}

var transactionalEmailsGetAggregatedReportFlags struct {
	startDate string
	endDate   string
	days      int
	tag       string
}

func init() {
	transactionalEmailsGetAggregatedReportCmd.Flags().StringVar(&transactionalEmailsGetAggregatedReportFlags.startDate, "start-date", "", "**Mandatory if endDate is used.** Starting date of the report (YYYY-MM-DD). Must be lower than equal to endDate ")
	transactionalEmailsGetAggregatedReportCmd.Flags().StringVar(&transactionalEmailsGetAggregatedReportFlags.endDate, "end-date", "", "**Mandatory if startDate is used.** Ending date of the report (YYYY-MM-DD). Must be greater than equal to startDate ")
	transactionalEmailsGetAggregatedReportCmd.Flags().IntVar(&transactionalEmailsGetAggregatedReportFlags.days, "days", 0, "Number of days in the past including today (positive integer). _Not compatible with 'startDate' and 'endDate'_ ")
	transactionalEmailsGetAggregatedReportCmd.Flags().StringVar(&transactionalEmailsGetAggregatedReportFlags.tag, "tag", "", "Tag of the emails")

	transactionalEmailsCmd.AddCommand(transactionalEmailsGetAggregatedReportCmd)
}

func runTransactionalEmailsGetAggregatedReport(cmd *cobra.Command, args []string) error {
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
			Name:        "tag",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Tag of the emails",
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
			Description: "Aggregated report informations",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "get-aggregated-report",
			"description": "Get your transactional email activity aggregated over a period of time",
			"http": map[string]any{
				"method": "GET",
				"path":   "/smtp/statistics/aggregatedReport",
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
		Path:        httpclient.SubstitutePath("/smtp/statistics/aggregatedReport", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("start-date") {
		req.QueryParams["startDate"] = fmt.Sprintf("%v", transactionalEmailsGetAggregatedReportFlags.startDate)
	}
	if cmd.Flags().Changed("end-date") {
		req.QueryParams["endDate"] = fmt.Sprintf("%v", transactionalEmailsGetAggregatedReportFlags.endDate)
	}
	if cmd.Flags().Changed("days") {
		req.QueryParams["days"] = fmt.Sprintf("%v", transactionalEmailsGetAggregatedReportFlags.days)
	}
	if cmd.Flags().Changed("tag") {
		req.QueryParams["tag"] = fmt.Sprintf("%v", transactionalEmailsGetAggregatedReportFlags.tag)
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
