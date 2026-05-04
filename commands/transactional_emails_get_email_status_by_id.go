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

var transactionalEmailsGetEmailStatusByIdCmd = &cobra.Command{
	Use:   "get-email-status-by-id",
	Short: "Fetch scheduled emails by batchId or messageId",
	RunE:  runTransactionalEmailsGetEmailStatusById,
}

var transactionalEmailsGetEmailStatusByIdFlags struct {
	identifier string
	startDate  string
	endDate    string
	sort       string
	status     string
	limit      int
	offset     int
}

func init() {
	transactionalEmailsGetEmailStatusByIdCmd.Flags().StringVar(&transactionalEmailsGetEmailStatusByIdFlags.identifier, "identifier", "", "The `batchId` of scheduled emails batch (Should be a valid UUIDv4) or the `messageId` of scheduled email.")
	transactionalEmailsGetEmailStatusByIdCmd.MarkFlagRequired("identifier")
	transactionalEmailsGetEmailStatusByIdCmd.Flags().StringVar(&transactionalEmailsGetEmailStatusByIdFlags.startDate, "start-date", "", "Mandatory if `endDate` is used. Starting date (YYYY-MM-DD) from which you want to fetch the list. Can be maximum 30 days older tha current date.")
	transactionalEmailsGetEmailStatusByIdCmd.Flags().StringVar(&transactionalEmailsGetEmailStatusByIdFlags.endDate, "end-date", "", "Mandatory if `startDate` is used. Ending date (YYYY-MM-DD) till which you want to fetch the list. Maximum time period that can be selected is one month.")
	transactionalEmailsGetEmailStatusByIdCmd.Flags().StringVar(&transactionalEmailsGetEmailStatusByIdFlags.sort, "sort", "", "Sort the results in the ascending/descending order of record creation. Default order is **descending** if `sort` is not passed. Not valid when identifier is `messageId`.")
	transactionalEmailsGetEmailStatusByIdCmd.Flags().StringVar(&transactionalEmailsGetEmailStatusByIdFlags.status, "status", "", "Filter the records by `status` of the scheduled email batch or message. Not valid when identifier is `messageId`.")
	transactionalEmailsGetEmailStatusByIdCmd.Flags().IntVar(&transactionalEmailsGetEmailStatusByIdFlags.limit, "limit", 0, "Number of documents returned per page. Not valid when identifier is `messageId`.")
	transactionalEmailsGetEmailStatusByIdCmd.Flags().IntVar(&transactionalEmailsGetEmailStatusByIdFlags.offset, "offset", 0, "Index of the first document on the page.  Not valid when identifier is `messageId`.")

	transactionalEmailsCmd.AddCommand(transactionalEmailsGetEmailStatusByIdCmd)
}

func runTransactionalEmailsGetEmailStatusById(cmd *cobra.Command, args []string) error {
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
			Description: "The `batchId` of scheduled emails batch (Should be a valid UUIDv4) or the `messageId` of scheduled email.",
		})
		flags = append(flags, flagSchema{
			Name:        "start-date",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Mandatory if `endDate` is used. Starting date (YYYY-MM-DD) from which you want to fetch the list. Can be maximum 30 days older tha current date.",
		})
		flags = append(flags, flagSchema{
			Name:        "end-date",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Mandatory if `startDate` is used. Ending date (YYYY-MM-DD) till which you want to fetch the list. Maximum time period that can be selected is one month.",
		})
		flags = append(flags, flagSchema{
			Name:        "sort",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Sort the results in the ascending/descending order of record creation. Default order is **descending** if `sort` is not passed. Not valid when identifier is `messageId`.",
		})
		flags = append(flags, flagSchema{
			Name:        "status",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter the records by `status` of the scheduled email batch or message. Not valid when identifier is `messageId`.",
		})
		flags = append(flags, flagSchema{
			Name:        "limit",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Number of documents returned per page. Not valid when identifier is `messageId`.",
		})
		flags = append(flags, flagSchema{
			Name:        "offset",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Index of the first document on the page.  Not valid when identifier is `messageId`.",
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
			Description: "Scheduled email batches",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "Invalid parameters passed",
		})
		responses = append(responses, responseSchema{
			Status:      "404",
			ContentType: "application/json",
			Description: "Record(s) for identifier not found",
		})

		schema := map[string]any{
			"command":     "get-email-status-by-id",
			"description": "Fetch scheduled emails by batchId or messageId",
			"http": map[string]any{
				"method": "GET",
				"path":   "/smtp/emailStatus/{identifier}",
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
	pathParams["identifier"] = fmt.Sprintf("%v", transactionalEmailsGetEmailStatusByIdFlags.identifier)

	req := &httpclient.Request{
		Method:      "GET",
		Path:        httpclient.SubstitutePath("/smtp/emailStatus/{identifier}", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("start-date") {
		req.QueryParams["startDate"] = fmt.Sprintf("%v", transactionalEmailsGetEmailStatusByIdFlags.startDate)
	}
	if cmd.Flags().Changed("end-date") {
		req.QueryParams["endDate"] = fmt.Sprintf("%v", transactionalEmailsGetEmailStatusByIdFlags.endDate)
	}
	if cmd.Flags().Changed("sort") {
		req.QueryParams["sort"] = fmt.Sprintf("%v", transactionalEmailsGetEmailStatusByIdFlags.sort)
	}
	if cmd.Flags().Changed("status") {
		req.QueryParams["status"] = fmt.Sprintf("%v", transactionalEmailsGetEmailStatusByIdFlags.status)
	}
	if cmd.Flags().Changed("limit") {
		req.QueryParams["limit"] = fmt.Sprintf("%v", transactionalEmailsGetEmailStatusByIdFlags.limit)
	}
	if cmd.Flags().Changed("offset") {
		req.QueryParams["offset"] = fmt.Sprintf("%v", transactionalEmailsGetEmailStatusByIdFlags.offset)
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
