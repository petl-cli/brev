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

var couponsCreateCollectionCmd = &cobra.Command{
	Use:   "create-collection",
	Short: "Create а coupon collection",
	RunE:  runCouponsCreateCollection,
}

var couponsCreateCollectionFlags struct {
	name                  string
	defaultCoupon         string
	expirationDate        string
	remainingDaysAlert    int
	remainingCouponsAlert int
	body                  string
}

func init() {
	couponsCreateCollectionCmd.Flags().StringVar(&couponsCreateCollectionFlags.name, "name", "", "Name of the coupons collection")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	couponsCreateCollectionCmd.Flags().StringVar(&couponsCreateCollectionFlags.defaultCoupon, "default-coupon", "", "Default coupons collection name")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	couponsCreateCollectionCmd.Flags().StringVar(&couponsCreateCollectionFlags.expirationDate, "expiration-date", "", "Specify an expiration date for the coupon collection in RFC3339 format. Use null to remove the expiration date.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	couponsCreateCollectionCmd.Flags().IntVar(&couponsCreateCollectionFlags.remainingDaysAlert, "remaining-days-alert", 0, "Send a notification alert (email) when the remaining days until the expiration date are equal or fall bellow this number. Use null to disable alerts.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	couponsCreateCollectionCmd.Flags().IntVar(&couponsCreateCollectionFlags.remainingCouponsAlert, "remaining-coupons-alert", 0, "Send a notification alert (email) when the remaining coupons count is equal or fall bellow this number. Use null to disable alerts.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	couponsCreateCollectionCmd.Flags().StringVar(&couponsCreateCollectionFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	couponsCmd.AddCommand(couponsCreateCollectionCmd)
}

func runCouponsCreateCollection(cmd *cobra.Command, args []string) error {
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
			Description: "Name of the coupons collection",
		})
		flags = append(flags, flagSchema{
			Name:        "default-coupon",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Default coupons collection name",
		})
		flags = append(flags, flagSchema{
			Name:        "expiration-date",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Specify an expiration date for the coupon collection in RFC3339 format. Use null to remove the expiration date.",
		})
		flags = append(flags, flagSchema{
			Name:        "remaining-days-alert",
			Type:        "integer",
			Required:    false,
			Location:    "body",
			Description: "Send a notification alert (email) when the remaining days until the expiration date are equal or fall bellow this number. Use null to disable alerts.",
		})
		flags = append(flags, flagSchema{
			Name:        "remaining-coupons-alert",
			Type:        "integer",
			Required:    false,
			Location:    "body",
			Description: "Send a notification alert (email) when the remaining coupons count is equal or fall bellow this number. Use null to disable alerts.",
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
			Description: "Coupon collection created",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})
		responses = append(responses, responseSchema{
			Status:      "401",
			ContentType: "application/json",
			Description: "unauthorized",
		})

		schema := map[string]any{
			"command":     "create-collection",
			"description": "Create а coupon collection",
			"http": map[string]any{
				"method": "POST",
				"path":   "/couponCollections",
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
		Path:        httpclient.SubstitutePath("/couponCollections", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if couponsCreateCollectionFlags.body != "" {
		if err := json.Unmarshal([]byte(couponsCreateCollectionFlags.body), &bodyMap); err != nil {
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
		bodyMap["name"] = couponsCreateCollectionFlags.name
	}
	if cmd.Flags().Changed("default-coupon") {
		bodyMap["defaultCoupon"] = couponsCreateCollectionFlags.defaultCoupon
	}
	if cmd.Flags().Changed("expiration-date") {
		bodyMap["expirationDate"] = couponsCreateCollectionFlags.expirationDate
	}
	if cmd.Flags().Changed("remaining-days-alert") {
		bodyMap["remainingDaysAlert"] = couponsCreateCollectionFlags.remainingDaysAlert
	}
	if cmd.Flags().Changed("remaining-coupons-alert") {
		bodyMap["remainingCouponsAlert"] = couponsCreateCollectionFlags.remainingCouponsAlert
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
