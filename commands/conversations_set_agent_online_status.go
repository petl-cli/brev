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

var conversationsSetAgentOnlineStatusCmd = &cobra.Command{
	Use:   "set-agent-online-status",
	Short: "Sets agent’s status to online for 2-3 minutes",
	RunE:  runConversationsSetAgentOnlineStatus,
}

var conversationsSetAgentOnlineStatusFlags struct {
	agentId      string
	receivedFrom string
	agentEmail   string
	agentName    string
	body         string
}

func init() {
	conversationsSetAgentOnlineStatusCmd.Flags().StringVar(&conversationsSetAgentOnlineStatusFlags.agentId, "agent-id", "", "agent ID. It can be found on agent’s page or received <a href=\"https://developers.brevo.com/docs/conversations-webhooks\">from a webhook</a>. Alternatively, you can use `agentEmail` + `agentName` + `receivedFrom` instead (all 3 fields required).")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	conversationsSetAgentOnlineStatusCmd.Flags().StringVar(&conversationsSetAgentOnlineStatusFlags.receivedFrom, "received-from", "", "mark your messages to distinguish messages created by you from the others.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	conversationsSetAgentOnlineStatusCmd.Flags().StringVar(&conversationsSetAgentOnlineStatusFlags.agentEmail, "agent-email", "", "agent email. When sending online pings from a standalone system, it’s hard to maintain a 1-to-1 relationship between the users of both systems. In this case, an agent can be specified by their email address. If there’s no agent with the specified email address in your Brevo organization, a dummy agent will be created automatically.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	conversationsSetAgentOnlineStatusCmd.Flags().StringVar(&conversationsSetAgentOnlineStatusFlags.agentName, "agent-name", "", "agent name")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	conversationsSetAgentOnlineStatusCmd.Flags().StringVar(&conversationsSetAgentOnlineStatusFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	conversationsCmd.AddCommand(conversationsSetAgentOnlineStatusCmd)
}

func runConversationsSetAgentOnlineStatus(cmd *cobra.Command, args []string) error {
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
			Name:        "agent-id",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "agent ID. It can be found on agent’s page or received <a href=\"https://developers.brevo.com/docs/conversations-webhooks\">from a webhook</a>. Alternatively, you can use `agentEmail` + `agentName` + `receivedFrom` instead (all 3 fields required).",
		})
		flags = append(flags, flagSchema{
			Name:        "received-from",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "mark your messages to distinguish messages created by you from the others.",
		})
		flags = append(flags, flagSchema{
			Name:        "agent-email",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "agent email. When sending online pings from a standalone system, it’s hard to maintain a 1-to-1 relationship between the users of both systems. In this case, an agent can be specified by their email address. If there’s no agent with the specified email address in your Brevo organization, a dummy agent will be created automatically.",
		})
		flags = append(flags, flagSchema{
			Name:        "agent-name",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "agent name",
		})

		type responseSchema struct {
			Status      string `json:"status"`
			ContentType string `json:"content_type,omitempty"`
			Description string `json:"description,omitempty"`
		}
		var responses []responseSchema
		responses = append(responses, responseSchema{
			Status:      "201",
			ContentType: "",
			Description: "Status of the agent was set successfully. Response body will be empty.",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "Returned when invalid data posted",
		})

		schema := map[string]any{
			"command":     "set-agent-online-status",
			"description": "Sets agent’s status to online for 2-3 minutes",
			"http": map[string]any{
				"method": "POST",
				"path":   "/conversations/agentOnlinePing",
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

	req := &httpclient.Request{
		Method:      "POST",
		Path:        httpclient.SubstitutePath("/conversations/agentOnlinePing", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if conversationsSetAgentOnlineStatusFlags.body != "" {
		if err := json.Unmarshal([]byte(conversationsSetAgentOnlineStatusFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("agent-id") {
		bodyMap["agentId"] = conversationsSetAgentOnlineStatusFlags.agentId
	}
	if cmd.Flags().Changed("received-from") {
		bodyMap["receivedFrom"] = conversationsSetAgentOnlineStatusFlags.receivedFrom
	}
	if cmd.Flags().Changed("agent-email") {
		bodyMap["agentEmail"] = conversationsSetAgentOnlineStatusFlags.agentEmail
	}
	if cmd.Flags().Changed("agent-name") {
		bodyMap["agentName"] = conversationsSetAgentOnlineStatusFlags.agentName
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
