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

var filesGetAllFilesCmd = &cobra.Command{
	Use:   "get-all-files",
	Short: "Get all files",
	RunE:  runFilesGetAllFiles,
}

var filesGetAllFilesFlags struct {
	entity    string
	entityIds string
	dateFrom  int
	dateTo    int
	offset    int
	limit     int
	sort      string
}

func init() {
	filesGetAllFilesCmd.Flags().StringVar(&filesGetAllFilesFlags.entity, "entity", "", "Filter by file entity type")
	filesGetAllFilesCmd.Flags().StringVar(&filesGetAllFilesFlags.entityIds, "entity-ids", "", "Filter by file entity IDs")
	filesGetAllFilesCmd.Flags().IntVar(&filesGetAllFilesFlags.dateFrom, "date-from", 0, "dateFrom to date range filter type (timestamp in milliseconds)")
	filesGetAllFilesCmd.Flags().IntVar(&filesGetAllFilesFlags.dateTo, "date-to", 0, "dateTo to date range filter type (timestamp in milliseconds)")
	filesGetAllFilesCmd.Flags().IntVar(&filesGetAllFilesFlags.offset, "offset", 0, "Index of the first document of the page")
	filesGetAllFilesCmd.Flags().IntVar(&filesGetAllFilesFlags.limit, "limit", 0, "Number of documents per page")
	filesGetAllFilesCmd.Flags().StringVar(&filesGetAllFilesFlags.sort, "sort", "", "Sort the results in the ascending/descending order. Default order is **descending** by creation if `sort` is not passed")

	filesCmd.AddCommand(filesGetAllFilesCmd)
}

func runFilesGetAllFiles(cmd *cobra.Command, args []string) error {
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
			Name:        "entity",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter by file entity type",
		})
		flags = append(flags, flagSchema{
			Name:        "entity-ids",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter by file entity IDs",
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

		type responseSchema struct {
			Status      string `json:"status"`
			ContentType string `json:"content_type,omitempty"`
			Description string `json:"description,omitempty"`
		}
		var responses []responseSchema
		responses = append(responses, responseSchema{
			Status:      "200",
			ContentType: "application/json",
			Description: "Returns files list with filters",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "Returned when query params are invalid",
		})

		schema := map[string]any{
			"command":     "get-all-files",
			"description": "Get all files",
			"http": map[string]any{
				"method": "GET",
				"path":   "/crm/files",
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
		Path:        httpclient.SubstitutePath("/crm/files", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("entity") {
		req.QueryParams["entity"] = fmt.Sprintf("%v", filesGetAllFilesFlags.entity)
	}
	if cmd.Flags().Changed("entity-ids") {
		req.QueryParams["entityIds"] = fmt.Sprintf("%v", filesGetAllFilesFlags.entityIds)
	}
	if cmd.Flags().Changed("date-from") {
		req.QueryParams["dateFrom"] = fmt.Sprintf("%v", filesGetAllFilesFlags.dateFrom)
	}
	if cmd.Flags().Changed("date-to") {
		req.QueryParams["dateTo"] = fmt.Sprintf("%v", filesGetAllFilesFlags.dateTo)
	}
	if cmd.Flags().Changed("offset") {
		req.QueryParams["offset"] = fmt.Sprintf("%v", filesGetAllFilesFlags.offset)
	}
	if cmd.Flags().Changed("limit") {
		req.QueryParams["limit"] = fmt.Sprintf("%v", filesGetAllFilesFlags.limit)
	}
	if cmd.Flags().Changed("sort") {
		req.QueryParams["sort"] = fmt.Sprintf("%v", filesGetAllFilesFlags.sort)
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
