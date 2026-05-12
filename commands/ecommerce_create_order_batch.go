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

var ecommerceCreateOrderBatchCmd = &cobra.Command{
	Use:   "create-order-batch",
	Short: "Create orders in batch",
	RunE:  runEcommerceCreateOrderBatch,
}

var ecommerceCreateOrderBatchFlags struct {
	orders     []string
	notifyUrl  string
	historical bool
	body       string
}

func init() {
	ecommerceCreateOrderBatchCmd.Flags().StringSliceVar(&ecommerceCreateOrderBatchFlags.orders, "orders", nil, "array of order objects")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	ecommerceCreateOrderBatchCmd.Flags().StringVar(&ecommerceCreateOrderBatchFlags.notifyUrl, "notify-url", "", "Notify Url provided by client to get the status of batch request")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	ecommerceCreateOrderBatchCmd.Flags().BoolVar(&ecommerceCreateOrderBatchFlags.historical, "historical", false, "Defines wether you want your orders to be considered as live data or as historical data (import of past data, synchronising data). True: orders will not trigger any automation workflows. False: orders will trigger workflows as usual.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	ecommerceCreateOrderBatchCmd.Flags().StringVar(&ecommerceCreateOrderBatchFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	ecommerceCmd.AddCommand(ecommerceCreateOrderBatchCmd)
}

func runEcommerceCreateOrderBatch(cmd *cobra.Command, args []string) error {
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
			Name:        "orders",
			Type:        "array",
			Required:    true,
			Location:    "body",
			Description: "array of order objects",
		})
		flags = append(flags, flagSchema{
			Name:        "notify-url",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Notify Url provided by client to get the status of batch request",
		})
		flags = append(flags, flagSchema{
			Name:        "historical",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Defines wether you want your orders to be considered as live data or as historical data (import of past data, synchronising data). True: orders will not trigger any automation workflows. False: orders will trigger workflows as usual.",
		})

		type responseSchema struct {
			Status      string `json:"status"`
			ContentType string `json:"content_type,omitempty"`
			Description string `json:"description,omitempty"`
		}
		var responses []responseSchema
		responses = append(responses, responseSchema{
			Status:      "202",
			ContentType: "application/json",
			Description: "batch id created",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "create-order-batch",
			"description": "Create orders in batch",
			"http": map[string]any{
				"method": "POST",
				"path":   "/orders/status/batch",
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
		Path:        httpclient.SubstitutePath("/orders/status/batch", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if ecommerceCreateOrderBatchFlags.body != "" {
		if err := json.Unmarshal([]byte(ecommerceCreateOrderBatchFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("orders") {
		bodyMap["orders"] = ecommerceCreateOrderBatchFlags.orders
	}
	if cmd.Flags().Changed("notify-url") {
		bodyMap["notifyUrl"] = ecommerceCreateOrderBatchFlags.notifyUrl
	}
	if cmd.Flags().Changed("historical") {
		bodyMap["historical"] = ecommerceCreateOrderBatchFlags.historical
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
