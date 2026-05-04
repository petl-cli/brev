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

var resellerAssociateDedicatedIpToChildCmd = &cobra.Command{
	Use:   "associate-dedicated-ip-to-child",
	Short: "Associate a dedicated IP to the child",
	RunE:  runResellerAssociateDedicatedIpToChild,
}

var resellerAssociateDedicatedIpToChildFlags struct {
	childIdentifier string
	ip              string
	body            string
}

func init() {
	resellerAssociateDedicatedIpToChildCmd.Flags().StringVar(&resellerAssociateDedicatedIpToChildFlags.childIdentifier, "child-identifier", "", "Either auth key or id of reseller's child")
	resellerAssociateDedicatedIpToChildCmd.MarkFlagRequired("child-identifier")
	resellerAssociateDedicatedIpToChildCmd.Flags().StringVar(&resellerAssociateDedicatedIpToChildFlags.ip, "ip", "", "Dedicated ID")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	resellerAssociateDedicatedIpToChildCmd.Flags().StringVar(&resellerAssociateDedicatedIpToChildFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	resellerCmd.AddCommand(resellerAssociateDedicatedIpToChildCmd)
}

func runResellerAssociateDedicatedIpToChild(cmd *cobra.Command, args []string) error {
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
			Name:        "child-identifier",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "Either auth key or id of reseller's child",
		})
		flags = append(flags, flagSchema{
			Name:        "ip",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Dedicated ID",
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
			Description: "Dedicated IP is associated to the child",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})
		responses = append(responses, responseSchema{
			Status:      "404",
			ContentType: "application/json",
			Description: "Child auth key or child id not found",
		})

		schema := map[string]any{
			"command":     "associate-dedicated-ip-to-child",
			"description": "Associate a dedicated IP to the child",
			"http": map[string]any{
				"method": "POST",
				"path":   "/reseller/children/{childIdentifier}/ips/associate",
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
	pathParams["childIdentifier"] = fmt.Sprintf("%v", resellerAssociateDedicatedIpToChildFlags.childIdentifier)

	req := &httpclient.Request{
		Method:      "POST",
		Path:        httpclient.SubstitutePath("/reseller/children/{childIdentifier}/ips/associate", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if resellerAssociateDedicatedIpToChildFlags.body != "" {
		if err := json.Unmarshal([]byte(resellerAssociateDedicatedIpToChildFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("ip") {
		bodyMap["ip"] = resellerAssociateDedicatedIpToChildFlags.ip
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
