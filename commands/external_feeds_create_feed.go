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

var externalFeedsCreateFeedCmd = &cobra.Command{
	Use:   "create-feed",
	Short: "Create an external feed",
	RunE:  runExternalFeedsCreateFeed,
}

var externalFeedsCreateFeedFlags struct {
	name       string
	url        string
	authType   string
	username   string
	password   string
	token      string
	headers    []string
	maxRetries int
	cache      bool
	body       string
}

func init() {
	externalFeedsCreateFeedCmd.Flags().StringVar(&externalFeedsCreateFeedFlags.name, "name", "", "Name of the feed")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	externalFeedsCreateFeedCmd.Flags().StringVar(&externalFeedsCreateFeedFlags.url, "url", "", "URL of the feed")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	externalFeedsCreateFeedCmd.Flags().StringVar(&externalFeedsCreateFeedFlags.authType, "auth-type", "", "Auth type of the feed:  * `basic`  * `token`  * `noAuth` ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	externalFeedsCreateFeedCmd.Flags().StringVar(&externalFeedsCreateFeedFlags.username, "username", "", "Username for authType `basic`")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	externalFeedsCreateFeedCmd.Flags().StringVar(&externalFeedsCreateFeedFlags.password, "password", "", "Password for authType `basic`")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	externalFeedsCreateFeedCmd.Flags().StringVar(&externalFeedsCreateFeedFlags.token, "token", "", "Token for authType `token`")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	externalFeedsCreateFeedCmd.Flags().StringSliceVar(&externalFeedsCreateFeedFlags.headers, "headers", nil, "Custom headers for the feed")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	externalFeedsCreateFeedCmd.Flags().IntVar(&externalFeedsCreateFeedFlags.maxRetries, "max-retries", 0, "Maximum number of retries on the feed url")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	externalFeedsCreateFeedCmd.Flags().BoolVar(&externalFeedsCreateFeedFlags.cache, "cache", false, "Toggle caching of feed url response")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	externalFeedsCreateFeedCmd.Flags().StringVar(&externalFeedsCreateFeedFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	externalFeedsCmd.AddCommand(externalFeedsCreateFeedCmd)
}

func runExternalFeedsCreateFeed(cmd *cobra.Command, args []string) error {
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
			Name:        "name",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Name of the feed",
		})
		flags = append(flags, flagSchema{
			Name:        "url",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "URL of the feed",
		})
		flags = append(flags, flagSchema{
			Name:        "auth-type",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Auth type of the feed:  * `basic`  * `token`  * `noAuth` ",
		})
		flags = append(flags, flagSchema{
			Name:        "username",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Username for authType `basic`",
		})
		flags = append(flags, flagSchema{
			Name:        "password",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Password for authType `basic`",
		})
		flags = append(flags, flagSchema{
			Name:        "token",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Token for authType `token`",
		})
		flags = append(flags, flagSchema{
			Name:        "headers",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "Custom headers for the feed",
		})
		flags = append(flags, flagSchema{
			Name:        "max-retries",
			Type:        "integer",
			Required:    false,
			Location:    "body",
			Description: "Maximum number of retries on the feed url",
		})
		flags = append(flags, flagSchema{
			Name:        "cache",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Toggle caching of feed url response",
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
			Description: "successfully created",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "create-feed",
			"description": "Create an external feed",
			"http": map[string]any{
				"method": "POST",
				"path":   "/feeds",
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
		Path:        httpclient.SubstitutePath("/feeds", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if externalFeedsCreateFeedFlags.body != "" {
		if err := json.Unmarshal([]byte(externalFeedsCreateFeedFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("name") {
		bodyMap["name"] = externalFeedsCreateFeedFlags.name
	}
	if cmd.Flags().Changed("url") {
		bodyMap["url"] = externalFeedsCreateFeedFlags.url
	}
	if cmd.Flags().Changed("auth-type") {
		bodyMap["authType"] = externalFeedsCreateFeedFlags.authType
	}
	if cmd.Flags().Changed("username") {
		bodyMap["username"] = externalFeedsCreateFeedFlags.username
	}
	if cmd.Flags().Changed("password") {
		bodyMap["password"] = externalFeedsCreateFeedFlags.password
	}
	if cmd.Flags().Changed("token") {
		bodyMap["token"] = externalFeedsCreateFeedFlags.token
	}
	if cmd.Flags().Changed("headers") {
		bodyMap["headers"] = externalFeedsCreateFeedFlags.headers
	}
	if cmd.Flags().Changed("max-retries") {
		bodyMap["maxRetries"] = externalFeedsCreateFeedFlags.maxRetries
	}
	if cmd.Flags().Changed("cache") {
		bodyMap["cache"] = externalFeedsCreateFeedFlags.cache
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
