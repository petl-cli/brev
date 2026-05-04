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

var transactionalEmailsListBlockedContactsCmd = &cobra.Command{
	Use:   "list-blocked-contacts",
	Short: "Get the list of blocked or unsubscribed transactional contacts",
	RunE:  runTransactionalEmailsListBlockedContacts,
}

var transactionalEmailsListBlockedContactsFlags struct {
	startDate string
	endDate   string
	limit     int
	offset    int
	senders   []string
	sort      string
}

func init() {
	transactionalEmailsListBlockedContactsCmd.Flags().StringVar(&transactionalEmailsListBlockedContactsFlags.startDate, "start-date", "", "**Mandatory if endDate is used.** Starting date (YYYY-MM-DD) from which you want to fetch the blocked or unsubscribed contacts ")
	transactionalEmailsListBlockedContactsCmd.Flags().StringVar(&transactionalEmailsListBlockedContactsFlags.endDate, "end-date", "", "**Mandatory if startDate is used.** Ending date (YYYY-MM-DD) till which you want to fetch the blocked or unsubscribed contacts ")
	transactionalEmailsListBlockedContactsCmd.Flags().IntVar(&transactionalEmailsListBlockedContactsFlags.limit, "limit", 0, "Number of documents returned per page")
	transactionalEmailsListBlockedContactsCmd.Flags().IntVar(&transactionalEmailsListBlockedContactsFlags.offset, "offset", 0, "Index of the first document on the page")
	transactionalEmailsListBlockedContactsCmd.Flags().StringSliceVar(&transactionalEmailsListBlockedContactsFlags.senders, "senders", nil, "Comma separated list of emails of the senders from which contacts are blocked or unsubscribed")
	transactionalEmailsListBlockedContactsCmd.Flags().StringVar(&transactionalEmailsListBlockedContactsFlags.sort, "sort", "", "Sort the results in the ascending/descending order of record creation. Default order is **descending** if `sort` is not passed")

	transactionalEmailsCmd.AddCommand(transactionalEmailsListBlockedContactsCmd)
}

func runTransactionalEmailsListBlockedContacts(cmd *cobra.Command, args []string) error {
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
			Description: "**Mandatory if endDate is used.** Starting date (YYYY-MM-DD) from which you want to fetch the blocked or unsubscribed contacts ",
		})
		flags = append(flags, flagSchema{
			Name:        "end-date",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "**Mandatory if startDate is used.** Ending date (YYYY-MM-DD) till which you want to fetch the blocked or unsubscribed contacts ",
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
			Description: "Index of the first document on the page",
		})
		flags = append(flags, flagSchema{
			Name:        "senders",
			Type:        "array",
			Required:    false,
			Location:    "query",
			Description: "Comma separated list of emails of the senders from which contacts are blocked or unsubscribed",
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
			Description: "List of blocked or unsubscribed transactional contacts",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "list-blocked-contacts",
			"description": "Get the list of blocked or unsubscribed transactional contacts",
			"http": map[string]any{
				"method": "GET",
				"path":   "/smtp/blockedContacts",
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
		Path:        httpclient.SubstitutePath("/smtp/blockedContacts", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("start-date") {
		req.QueryParams["startDate"] = fmt.Sprintf("%v", transactionalEmailsListBlockedContactsFlags.startDate)
	}
	if cmd.Flags().Changed("end-date") {
		req.QueryParams["endDate"] = fmt.Sprintf("%v", transactionalEmailsListBlockedContactsFlags.endDate)
	}
	if cmd.Flags().Changed("limit") {
		req.QueryParams["limit"] = fmt.Sprintf("%v", transactionalEmailsListBlockedContactsFlags.limit)
	}
	if cmd.Flags().Changed("offset") {
		req.QueryParams["offset"] = fmt.Sprintf("%v", transactionalEmailsListBlockedContactsFlags.offset)
	}
	if cmd.Flags().Changed("senders") {
		req.ArrayParams["senders"] = transactionalEmailsListBlockedContactsFlags.senders
	}
	if cmd.Flags().Changed("sort") {
		req.QueryParams["sort"] = fmt.Sprintf("%v", transactionalEmailsListBlockedContactsFlags.sort)
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
