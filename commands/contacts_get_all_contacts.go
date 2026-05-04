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

var contactsGetAllContactsCmd = &cobra.Command{
	Use:   "get-all-contacts",
	Short: "Get all the contacts",
	RunE:  runContactsGetAllContacts,
}

var contactsGetAllContactsFlags struct {
	limit         int
	offset        int
	modifiedSince string
	createdSince  string
	sort          string
	segmentId     int
	listIds       []string
}

func init() {
	contactsGetAllContactsCmd.Flags().IntVar(&contactsGetAllContactsFlags.limit, "limit", 0, "Number of documents per page")
	contactsGetAllContactsCmd.Flags().IntVar(&contactsGetAllContactsFlags.offset, "offset", 0, "Index of the first document of the page")
	contactsGetAllContactsCmd.Flags().StringVar(&contactsGetAllContactsFlags.modifiedSince, "modified-since", "", "Filter (urlencoded) the contacts modified after a given UTC date-time (YYYY-MM-DDTHH:mm:ss.SSSZ). **Prefer to pass your timezone in date-time format for accurate result.** ")
	contactsGetAllContactsCmd.Flags().StringVar(&contactsGetAllContactsFlags.createdSince, "created-since", "", "Filter (urlencoded) the contacts created after a given UTC date-time (YYYY-MM-DDTHH:mm:ss.SSSZ). **Prefer to pass your timezone in date-time format for accurate result.** ")
	contactsGetAllContactsCmd.Flags().StringVar(&contactsGetAllContactsFlags.sort, "sort", "", "Sort the results in the ascending/descending order of record creation. Default order is **descending** if `sort` is not passed")
	contactsGetAllContactsCmd.Flags().IntVar(&contactsGetAllContactsFlags.segmentId, "segment-id", 0, "Id of the segment. **Either listIds or segmentId can be passed.**")
	contactsGetAllContactsCmd.Flags().StringSliceVar(&contactsGetAllContactsFlags.listIds, "list-ids", nil, "Ids of the list. **Either listIds or segmentId can be passed.**")

	contactsCmd.AddCommand(contactsGetAllContactsCmd)
}

func runContactsGetAllContacts(cmd *cobra.Command, args []string) error {
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
			Name:        "offset",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Index of the first document of the page",
		})
		flags = append(flags, flagSchema{
			Name:        "modified-since",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter (urlencoded) the contacts modified after a given UTC date-time (YYYY-MM-DDTHH:mm:ss.SSSZ). **Prefer to pass your timezone in date-time format for accurate result.** ",
		})
		flags = append(flags, flagSchema{
			Name:        "created-since",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter (urlencoded) the contacts created after a given UTC date-time (YYYY-MM-DDTHH:mm:ss.SSSZ). **Prefer to pass your timezone in date-time format for accurate result.** ",
		})
		flags = append(flags, flagSchema{
			Name:        "sort",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Sort the results in the ascending/descending order of record creation. Default order is **descending** if `sort` is not passed",
		})
		flags = append(flags, flagSchema{
			Name:        "segment-id",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Id of the segment. **Either listIds or segmentId can be passed.**",
		})
		flags = append(flags, flagSchema{
			Name:        "list-ids",
			Type:        "array",
			Required:    false,
			Location:    "query",
			Description: "Ids of the list. **Either listIds or segmentId can be passed.**",
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
			Description: "All contacts listed",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "get-all-contacts",
			"description": "Get all the contacts",
			"http": map[string]any{
				"method": "GET",
				"path":   "/contacts",
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
		Path:        httpclient.SubstitutePath("/contacts", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("limit") {
		req.QueryParams["limit"] = fmt.Sprintf("%v", contactsGetAllContactsFlags.limit)
	}
	if cmd.Flags().Changed("offset") {
		req.QueryParams["offset"] = fmt.Sprintf("%v", contactsGetAllContactsFlags.offset)
	}
	if cmd.Flags().Changed("modified-since") {
		req.QueryParams["modifiedSince"] = fmt.Sprintf("%v", contactsGetAllContactsFlags.modifiedSince)
	}
	if cmd.Flags().Changed("created-since") {
		req.QueryParams["createdSince"] = fmt.Sprintf("%v", contactsGetAllContactsFlags.createdSince)
	}
	if cmd.Flags().Changed("sort") {
		req.QueryParams["sort"] = fmt.Sprintf("%v", contactsGetAllContactsFlags.sort)
	}
	if cmd.Flags().Changed("segment-id") {
		req.QueryParams["segmentId"] = fmt.Sprintf("%v", contactsGetAllContactsFlags.segmentId)
	}
	if cmd.Flags().Changed("list-ids") {
		req.ArrayParams["listIds"] = contactsGetAllContactsFlags.listIds
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
