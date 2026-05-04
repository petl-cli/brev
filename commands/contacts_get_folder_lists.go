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

var contactsGetFolderListsCmd = &cobra.Command{
	Use:   "get-folder-lists",
	Short: "Get lists in a folder",
	RunE:  runContactsGetFolderLists,
}

var contactsGetFolderListsFlags struct {
	folderId int
	limit    int
	offset   int
	sort     string
}

func init() {
	contactsGetFolderListsCmd.Flags().IntVar(&contactsGetFolderListsFlags.folderId, "folder-id", 0, "Id of the folder")
	contactsGetFolderListsCmd.MarkFlagRequired("folder-id")
	contactsGetFolderListsCmd.Flags().IntVar(&contactsGetFolderListsFlags.limit, "limit", 0, "Number of documents per page")
	contactsGetFolderListsCmd.Flags().IntVar(&contactsGetFolderListsFlags.offset, "offset", 0, "Index of the first document of the page")
	contactsGetFolderListsCmd.Flags().StringVar(&contactsGetFolderListsFlags.sort, "sort", "", "Sort the results in the ascending/descending order of record creation. Default order is **descending** if `sort` is not passed")

	contactsCmd.AddCommand(contactsGetFolderListsCmd)
}

func runContactsGetFolderLists(cmd *cobra.Command, args []string) error {
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
			Name:        "folder-id",
			Type:        "integer",
			Required:    true,
			Location:    "path",
			Description: "Id of the folder",
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
			Description: "Index of the first document of the page",
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
			Description: "Folder's Lists details",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})
		responses = append(responses, responseSchema{
			Status:      "404",
			ContentType: "application/json",
			Description: "Folder ID not found",
		})

		schema := map[string]any{
			"command":     "get-folder-lists",
			"description": "Get lists in a folder",
			"http": map[string]any{
				"method": "GET",
				"path":   "/contacts/folders/{folderId}/lists",
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
	pathParams["folderId"] = fmt.Sprintf("%v", contactsGetFolderListsFlags.folderId)

	req := &httpclient.Request{
		Method:      "GET",
		Path:        httpclient.SubstitutePath("/contacts/folders/{folderId}/lists", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("limit") {
		req.QueryParams["limit"] = fmt.Sprintf("%v", contactsGetFolderListsFlags.limit)
	}
	if cmd.Flags().Changed("offset") {
		req.QueryParams["offset"] = fmt.Sprintf("%v", contactsGetFolderListsFlags.offset)
	}
	if cmd.Flags().Changed("sort") {
		req.QueryParams["sort"] = fmt.Sprintf("%v", contactsGetFolderListsFlags.sort)
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
