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

var emailCampaignsGetAllCmd = &cobra.Command{
	Use:   "get-all",
	Short: "Return all your created email campaigns",
	RunE:  runEmailCampaignsGetAll,
}

var emailCampaignsGetAllFlags struct {
	type_              string
	status             string
	statistics         string
	startDate          string
	endDate            string
	limit              int
	offset             int
	sort               string
	excludeHtmlContent bool
}

func init() {
	emailCampaignsGetAllCmd.Flags().StringVar(&emailCampaignsGetAllFlags.type_, "type", "", "Filter on the type of the campaigns")
	emailCampaignsGetAllCmd.Flags().StringVar(&emailCampaignsGetAllFlags.status, "status", "", "Filter on the status of the campaign")
	emailCampaignsGetAllCmd.Flags().StringVar(&emailCampaignsGetAllFlags.statistics, "statistics", "", "Filter on the type of statistics required. Example **globalStats** value will only fetch globalStats info of the campaign in returned response.")
	emailCampaignsGetAllCmd.Flags().StringVar(&emailCampaignsGetAllFlags.startDate, "start-date", "", "**Mandatory if endDate is used**. Starting (urlencoded) UTC date-time (YYYY-MM-DDTHH:mm:ss.SSSZ) to filter the sent email campaigns. **Prefer to pass your timezone in date-time format for accurate result** ( only available if either 'status' not passed and if passed is set to 'sent' ) ")
	emailCampaignsGetAllCmd.Flags().StringVar(&emailCampaignsGetAllFlags.endDate, "end-date", "", "**Mandatory if startDate is used**. Ending (urlencoded) UTC date-time (YYYY-MM-DDTHH:mm:ss.SSSZ) to filter the sent email campaigns. **Prefer to pass your timezone in date-time format for accurate result** ( only available if either 'status' not passed and if passed is set to 'sent' ) ")
	emailCampaignsGetAllCmd.Flags().IntVar(&emailCampaignsGetAllFlags.limit, "limit", 0, "Number of documents per page")
	emailCampaignsGetAllCmd.Flags().IntVar(&emailCampaignsGetAllFlags.offset, "offset", 0, "Index of the first document in the page")
	emailCampaignsGetAllCmd.Flags().StringVar(&emailCampaignsGetAllFlags.sort, "sort", "", "Sort the results in the ascending/descending order of record creation. Default order is **descending** if `sort` is not passed")
	emailCampaignsGetAllCmd.Flags().BoolVar(&emailCampaignsGetAllFlags.excludeHtmlContent, "exclude-html-content", false, "Use this flag to exclude htmlContent from the response body. If set to **true**, htmlContent field will be returned as empty string in the response body")

	emailCampaignsCmd.AddCommand(emailCampaignsGetAllCmd)
}

func runEmailCampaignsGetAll(cmd *cobra.Command, args []string) error {
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
			Name:        "type",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter on the type of the campaigns",
		})
		flags = append(flags, flagSchema{
			Name:        "status",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter on the status of the campaign",
		})
		flags = append(flags, flagSchema{
			Name:        "statistics",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter on the type of statistics required. Example **globalStats** value will only fetch globalStats info of the campaign in returned response.",
		})
		flags = append(flags, flagSchema{
			Name:        "start-date",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "**Mandatory if endDate is used**. Starting (urlencoded) UTC date-time (YYYY-MM-DDTHH:mm:ss.SSSZ) to filter the sent email campaigns. **Prefer to pass your timezone in date-time format for accurate result** ( only available if either 'status' not passed and if passed is set to 'sent' ) ",
		})
		flags = append(flags, flagSchema{
			Name:        "end-date",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "**Mandatory if startDate is used**. Ending (urlencoded) UTC date-time (YYYY-MM-DDTHH:mm:ss.SSSZ) to filter the sent email campaigns. **Prefer to pass your timezone in date-time format for accurate result** ( only available if either 'status' not passed and if passed is set to 'sent' ) ",
		})
		flags = append(flags, flagSchema{
			Name:        "limit",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Number of documents per page",
		})
		flags = append(flags, flagSchema{
			Name:        "offset",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Index of the first document in the page",
		})
		flags = append(flags, flagSchema{
			Name:        "sort",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Sort the results in the ascending/descending order of record creation. Default order is **descending** if `sort` is not passed",
		})
		flags = append(flags, flagSchema{
			Name:        "exclude-html-content",
			Type:        "boolean",
			Required:    false,
			Location:    "query",
			Description: "Use this flag to exclude htmlContent from the response body. If set to **true**, htmlContent field will be returned as empty string in the response body",
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
			Description: "Email campaigns informations",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "get-all",
			"description": "Return all your created email campaigns",
			"http": map[string]any{
				"method": "GET",
				"path":   "/emailCampaigns",
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
		Path:        httpclient.SubstitutePath("/emailCampaigns", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("type") {
		req.QueryParams["type"] = fmt.Sprintf("%v", emailCampaignsGetAllFlags.type_)
	}
	if cmd.Flags().Changed("status") {
		req.QueryParams["status"] = fmt.Sprintf("%v", emailCampaignsGetAllFlags.status)
	}
	if cmd.Flags().Changed("statistics") {
		req.QueryParams["statistics"] = fmt.Sprintf("%v", emailCampaignsGetAllFlags.statistics)
	}
	if cmd.Flags().Changed("start-date") {
		req.QueryParams["startDate"] = fmt.Sprintf("%v", emailCampaignsGetAllFlags.startDate)
	}
	if cmd.Flags().Changed("end-date") {
		req.QueryParams["endDate"] = fmt.Sprintf("%v", emailCampaignsGetAllFlags.endDate)
	}
	if cmd.Flags().Changed("limit") {
		req.QueryParams["limit"] = fmt.Sprintf("%v", emailCampaignsGetAllFlags.limit)
	}
	if cmd.Flags().Changed("offset") {
		req.QueryParams["offset"] = fmt.Sprintf("%v", emailCampaignsGetAllFlags.offset)
	}
	if cmd.Flags().Changed("sort") {
		req.QueryParams["sort"] = fmt.Sprintf("%v", emailCampaignsGetAllFlags.sort)
	}
	if cmd.Flags().Changed("exclude-html-content") {
		req.QueryParams["excludeHtmlContent"] = fmt.Sprintf("%v", emailCampaignsGetAllFlags.excludeHtmlContent)
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
