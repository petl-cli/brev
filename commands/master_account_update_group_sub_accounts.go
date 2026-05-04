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

var masterAccountUpdateGroupSubAccountsCmd = &cobra.Command{
	Use:   "update-group-sub-accounts",
	Short: "Update a group of sub-accounts",
	RunE:  runMasterAccountUpdateGroupSubAccounts,
}

var masterAccountUpdateGroupSubAccountsFlags struct {
	id            string
	groupName     string
	subAccountIds []string
	body          string
}

func init() {
	masterAccountUpdateGroupSubAccountsCmd.Flags().StringVar(&masterAccountUpdateGroupSubAccountsFlags.id, "id", "", "Id of the group")
	masterAccountUpdateGroupSubAccountsCmd.MarkFlagRequired("id")
	masterAccountUpdateGroupSubAccountsCmd.Flags().StringVar(&masterAccountUpdateGroupSubAccountsFlags.groupName, "group-name", "", "The name of the group of sub-accounts")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	masterAccountUpdateGroupSubAccountsCmd.Flags().StringSliceVar(&masterAccountUpdateGroupSubAccountsFlags.subAccountIds, "sub-account-ids", nil, "Pass the list of sub-account Ids to be included in the group")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	masterAccountUpdateGroupSubAccountsCmd.Flags().StringVar(&masterAccountUpdateGroupSubAccountsFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	masterAccountCmd.AddCommand(masterAccountUpdateGroupSubAccountsCmd)
}

func runMasterAccountUpdateGroupSubAccounts(cmd *cobra.Command, args []string) error {
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
			Name:        "id",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "Id of the group",
		})
		flags = append(flags, flagSchema{
			Name:        "group-name",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "The name of the group of sub-accounts",
		})
		flags = append(flags, flagSchema{
			Name:        "sub-account-ids",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "Pass the list of sub-account Ids to be included in the group",
		})

		type responseSchema struct {
			Status      string `json:"status"`
			ContentType string `json:"content_type,omitempty"`
			Description string `json:"description,omitempty"`
		}
		var responses []responseSchema
		responses = append(responses, responseSchema{
			Status:      "204",
			ContentType: "",
			Description: "Group details updated",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "update-group-sub-accounts",
			"description": "Update a group of sub-accounts",
			"http": map[string]any{
				"method": "PUT",
				"path":   "/corporate/group/{id}",
			},
			"input": map[string]any{
				"flags":         flags,
				"body_flag":     true,
				"body_required": true,
			},
			"output": map[string]any{
				"responses": responses,
			},
			"semantics": map[string]any{
				"safe":         false,
				"idempotent":   true,
				"reversible":   true,
				"side_effects": []string{"mutates_resource"},
				"impact":       "medium",
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
	pathParams["id"] = fmt.Sprintf("%v", masterAccountUpdateGroupSubAccountsFlags.id)

	req := &httpclient.Request{
		Method:      "PUT",
		Path:        httpclient.SubstitutePath("/corporate/group/{id}", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if masterAccountUpdateGroupSubAccountsFlags.body != "" {
		if err := json.Unmarshal([]byte(masterAccountUpdateGroupSubAccountsFlags.body), &bodyMap); err != nil {
			_invState.errorType = "parse_error"
			cliErr := &output.CLIError{
				Error:    true,
				Code:     "validation_error",
				Message:  fmt.Sprintf("invalid JSON in --body: %v", err),
				ExitCode: output.ExitValidation,
			}
			cliErr.Write(os.Stderr)
			return output.NewExitError(cliErr)
		}
	}
	// Individual flags overlay onto body (flags take precedence over --body JSON)
	if cmd.Flags().Changed("group-name") {
		bodyMap["groupName"] = masterAccountUpdateGroupSubAccountsFlags.groupName
	}
	if cmd.Flags().Changed("sub-account-ids") {
		bodyMap["subAccountIds"] = masterAccountUpdateGroupSubAccountsFlags.subAccountIds
	}
	req.Body = bodyMap

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
