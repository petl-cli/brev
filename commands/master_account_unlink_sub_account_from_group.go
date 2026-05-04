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

var masterAccountUnlinkSubAccountFromGroupCmd = &cobra.Command{
	Use:   "unlink-sub-account-from-group",
	Short: "Delete sub-account from group",
	RunE:  runMasterAccountUnlinkSubAccountFromGroup,
}

var masterAccountUnlinkSubAccountFromGroupFlags struct {
	groupId       string
	subAccountIds []string
	body          string
}

func init() {
	masterAccountUnlinkSubAccountFromGroupCmd.Flags().StringVar(&masterAccountUnlinkSubAccountFromGroupFlags.groupId, "group-id", "", "Group id")
	masterAccountUnlinkSubAccountFromGroupCmd.MarkFlagRequired("group-id")
	masterAccountUnlinkSubAccountFromGroupCmd.Flags().StringSliceVar(&masterAccountUnlinkSubAccountFromGroupFlags.subAccountIds, "sub-account-ids", nil, "List of sub-account ids")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	masterAccountUnlinkSubAccountFromGroupCmd.Flags().StringVar(&masterAccountUnlinkSubAccountFromGroupFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	masterAccountCmd.AddCommand(masterAccountUnlinkSubAccountFromGroupCmd)
}

func runMasterAccountUnlinkSubAccountFromGroup(cmd *cobra.Command, args []string) error {
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
			Name:        "group-id",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "Group id",
		})
		flags = append(flags, flagSchema{
			Name:        "sub-account-ids",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "List of sub-account ids",
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
			Description: "SubAccounts removed from the group",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "unlink-sub-account-from-group",
			"description": "Delete sub-account from group",
			"http": map[string]any{
				"method": "PUT",
				"path":   "/corporate/group/unlink/{groupId}/subAccounts",
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
	pathParams["groupId"] = fmt.Sprintf("%v", masterAccountUnlinkSubAccountFromGroupFlags.groupId)

	req := &httpclient.Request{
		Method:      "PUT",
		Path:        httpclient.SubstitutePath("/corporate/group/unlink/{groupId}/subAccounts", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if masterAccountUnlinkSubAccountFromGroupFlags.body != "" {
		if err := json.Unmarshal([]byte(masterAccountUnlinkSubAccountFromGroupFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("sub-account-ids") {
		bodyMap["subAccountIds"] = masterAccountUnlinkSubAccountFromGroupFlags.subAccountIds
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
