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

var transactionalSmsGetAggregatedReportCmd = &cobra.Command{
	Use:   "get-aggregated-report",
	Short: "Get your SMS activity aggregated over a period of time",
	RunE:  runTransactionalSmsGetAggregatedReport,
}

var transactionalSmsGetAggregatedReportFlags struct {
	startDate string
	endDate   string
	days      int
	tag       string
}

func init() {
	transactionalSmsGetAggregatedReportCmd.Flags().StringVar(&transactionalSmsGetAggregatedReportFlags.startDate, "start-date", "", "**Mandatory if endDate is used.** Starting date (YYYY-MM-DD) of the report ")
	transactionalSmsGetAggregatedReportCmd.Flags().StringVar(&transactionalSmsGetAggregatedReportFlags.endDate, "end-date", "", "**Mandatory if startDate is used.** Ending date (YYYY-MM-DD) of the report ")
	transactionalSmsGetAggregatedReportCmd.Flags().IntVar(&transactionalSmsGetAggregatedReportFlags.days, "days", 0, "Number of days in the past including today (positive integer). **Not compatible with startDate and endDate** ")
	transactionalSmsGetAggregatedReportCmd.Flags().StringVar(&transactionalSmsGetAggregatedReportFlags.tag, "tag", "", "Filter on a tag")

	transactionalSmsCmd.AddCommand(transactionalSmsGetAggregatedReportCmd)
}

func runTransactionalSmsGetAggregatedReport(cmd *cobra.Command, args []string) error {
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
			Name:        "days",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Number of days in the past including today (positive integer). **Not compatible with startDate and endDate** ",
		})
		flags = append(flags, flagSchema{
			Name:        "tag",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter on a tag",
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
			Description: "Aggregated SMS report informations",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "get-aggregated-report",
			"description": "Get your SMS activity aggregated over a period of time",
			"http": map[string]any{
				"method": "GET",
				"path":   "/transactionalSMS/statistics/aggregatedReport",
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
		Path:        httpclient.SubstitutePath("/transactionalSMS/statistics/aggregatedReport", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("start-date") {
		req.QueryParams["startDate"] = fmt.Sprintf("%v", transactionalSmsGetAggregatedReportFlags.startDate)
	}
	if cmd.Flags().Changed("end-date") {
		req.QueryParams["endDate"] = fmt.Sprintf("%v", transactionalSmsGetAggregatedReportFlags.endDate)
	}
	if cmd.Flags().Changed("days") {
		req.QueryParams["days"] = fmt.Sprintf("%v", transactionalSmsGetAggregatedReportFlags.days)
	}
	if cmd.Flags().Changed("tag") {
		req.QueryParams["tag"] = fmt.Sprintf("%v", transactionalSmsGetAggregatedReportFlags.tag)
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
