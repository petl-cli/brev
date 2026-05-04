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

var transactionalEmailsGetListCmd = &cobra.Command{
	Use:   "get-list",
	Short: "Get the list of transactional emails on the basis of allowed filters",
	RunE:  runTransactionalEmailsGetList,
}

var transactionalEmailsGetListFlags struct {
	email      string
	templateId int
	messageId  string
	startDate  string
	endDate    string
	sort       string
	limit      int
	offset     int
}

func init() {
	transactionalEmailsGetListCmd.Flags().StringVar(&transactionalEmailsGetListFlags.email, "email", "", "**Mandatory if templateId and messageId are not passed in query filters.** Email address to which transactional email has been sent. ")
	transactionalEmailsGetListCmd.Flags().IntVar(&transactionalEmailsGetListFlags.templateId, "template-id", 0, "**Mandatory if email and messageId are not passed in query filters.** Id of the template that was used to compose transactional email. ")
	transactionalEmailsGetListCmd.Flags().StringVar(&transactionalEmailsGetListFlags.messageId, "message-id", "", "**Mandatory if templateId and email are not passed in query filters.** Message ID of the transactional email sent. ")
	transactionalEmailsGetListCmd.Flags().StringVar(&transactionalEmailsGetListFlags.startDate, "start-date", "", "**Mandatory if endDate is used.** Starting date (YYYY-MM-DD) from which you want to fetch the list. **Maximum time period that can be selected is one month**. ")
	transactionalEmailsGetListCmd.Flags().StringVar(&transactionalEmailsGetListFlags.endDate, "end-date", "", "**Mandatory if startDate is used.** Ending date (YYYY-MM-DD) till which you want to fetch the list. **Maximum time period that can be selected is one month.** ")
	transactionalEmailsGetListCmd.Flags().StringVar(&transactionalEmailsGetListFlags.sort, "sort", "", "Sort the results in the ascending/descending order of record creation. Default order is **descending** if `sort` is not passed")
	transactionalEmailsGetListCmd.Flags().IntVar(&transactionalEmailsGetListFlags.limit, "limit", 0, "Number of documents returned per page")
	transactionalEmailsGetListCmd.Flags().IntVar(&transactionalEmailsGetListFlags.offset, "offset", 0, "Index of the first document in the page")

	transactionalEmailsCmd.AddCommand(transactionalEmailsGetListCmd)
}

func runTransactionalEmailsGetList(cmd *cobra.Command, args []string) error {
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
			Name:        "email",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "**Mandatory if templateId and messageId are not passed in query filters.** Email address to which transactional email has been sent. ",
		})
		flags = append(flags, flagSchema{
			Name:        "template-id",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "**Mandatory if email and messageId are not passed in query filters.** Id of the template that was used to compose transactional email. ",
		})
		flags = append(flags, flagSchema{
			Name:        "message-id",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "**Mandatory if templateId and email are not passed in query filters.** Message ID of the transactional email sent. ",
		})
		flags = append(flags, flagSchema{
			Name:        "start-date",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "**Mandatory if endDate is used.** Starting date (YYYY-MM-DD) from which you want to fetch the list. **Maximum time period that can be selected is one month**. ",
		})
		flags = append(flags, flagSchema{
			Name:        "end-date",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "**Mandatory if startDate is used.** Ending date (YYYY-MM-DD) till which you want to fetch the list. **Maximum time period that can be selected is one month.** ",
		})
		flags = append(flags, flagSchema{
			Name:        "sort",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Sort the results in the ascending/descending order of record creation. Default order is **descending** if `sort` is not passed",
		})
		flags = append(flags, flagSchema{
			Name:        "limit",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Number of documents returned per page",
		})
		flags = append(flags, flagSchema{
			Name:        "offset",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Index of the first document in the page",
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
			Description: "List of transactional emails",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "get-list",
			"description": "Get the list of transactional emails on the basis of allowed filters",
			"http": map[string]any{
				"method": "GET",
				"path":   "/smtp/emails",
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
		Path:        httpclient.SubstitutePath("/smtp/emails", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("email") {
		req.QueryParams["email"] = fmt.Sprintf("%v", transactionalEmailsGetListFlags.email)
	}
	if cmd.Flags().Changed("template-id") {
		req.QueryParams["templateId"] = fmt.Sprintf("%v", transactionalEmailsGetListFlags.templateId)
	}
	if cmd.Flags().Changed("message-id") {
		req.QueryParams["messageId"] = fmt.Sprintf("%v", transactionalEmailsGetListFlags.messageId)
	}
	if cmd.Flags().Changed("start-date") {
		req.QueryParams["startDate"] = fmt.Sprintf("%v", transactionalEmailsGetListFlags.startDate)
	}
	if cmd.Flags().Changed("end-date") {
		req.QueryParams["endDate"] = fmt.Sprintf("%v", transactionalEmailsGetListFlags.endDate)
	}
	if cmd.Flags().Changed("sort") {
		req.QueryParams["sort"] = fmt.Sprintf("%v", transactionalEmailsGetListFlags.sort)
	}
	if cmd.Flags().Changed("limit") {
		req.QueryParams["limit"] = fmt.Sprintf("%v", transactionalEmailsGetListFlags.limit)
	}
	if cmd.Flags().Changed("offset") {
		req.QueryParams["offset"] = fmt.Sprintf("%v", transactionalEmailsGetListFlags.offset)
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
