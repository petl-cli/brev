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

var masterAccountCreateGroupOfSubAccountsCmd = &cobra.Command{
	Use:   "create-group-of-sub-accounts",
	Short: "Create a group of sub-accounts",
	RunE:  runMasterAccountCreateGroupOfSubAccounts,
}

var masterAccountCreateGroupOfSubAccountsFlags struct {
	groupName     string
	subAccountIds []string
	body          string
}

func init() {
	masterAccountCreateGroupOfSubAccountsCmd.Flags().StringVar(&masterAccountCreateGroupOfSubAccountsFlags.groupName, "group-name", "", "The name of the group of sub-accounts")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	masterAccountCreateGroupOfSubAccountsCmd.Flags().StringSliceVar(&masterAccountCreateGroupOfSubAccountsFlags.subAccountIds, "sub-account-ids", nil, "Pass the list of sub-account Ids to be included in the group")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	masterAccountCreateGroupOfSubAccountsCmd.Flags().StringVar(&masterAccountCreateGroupOfSubAccountsFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	masterAccountCmd.AddCommand(masterAccountCreateGroupOfSubAccountsCmd)
}

func runMasterAccountCreateGroupOfSubAccounts(cmd *cobra.Command, args []string) error {
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
			Status:      "201",
			ContentType: "application/json",
			Description: "Group ID",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "Bad request",
		})

		schema := map[string]any{
			"command":     "create-group-of-sub-accounts",
			"description": "Create a group of sub-accounts",
			"http": map[string]any{
				"method": "POST",
				"path":   "/corporate/group",
			},
			"input": map[string]any{
				"flags":         flags,
				"body_flag":     true,
				"body_required": false,
			},
			"output": map[string]any{
				"responses": responses,
			},
			"semantics": map[string]any{
				"safe":         false,
				"idempotent":   false,
				"reversible":   true,
				"side_effects": []string{"creates_resource"},
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

	req := &httpclient.Request{
		Method:      "POST",
		Path:        httpclient.SubstitutePath("/corporate/group", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if masterAccountCreateGroupOfSubAccountsFlags.body != "" {
		if err := json.Unmarshal([]byte(masterAccountCreateGroupOfSubAccountsFlags.body), &bodyMap); err != nil {
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
		bodyMap["groupName"] = masterAccountCreateGroupOfSubAccountsFlags.groupName
	}
	if cmd.Flags().Changed("sub-account-ids") {
		bodyMap["subAccountIds"] = masterAccountCreateGroupOfSubAccountsFlags.subAccountIds
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
