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

var couponsCreateCouponCollectionCmd = &cobra.Command{
	Use:   "create-coupon-collection",
	Short: "Create coupons for a coupon collection",
	RunE:  runCouponsCreateCouponCollection,
}

var couponsCreateCouponCollectionFlags struct {
	collectionId string
	coupons      []string
	body         string
}

func init() {
	couponsCreateCouponCollectionCmd.Flags().StringVar(&couponsCreateCouponCollectionFlags.collectionId, "collection-id", "", "The id of the coupon collection for which the coupons will be created")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	couponsCreateCouponCollectionCmd.Flags().StringSliceVar(&couponsCreateCouponCollectionFlags.coupons, "coupons", nil, "")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	couponsCreateCouponCollectionCmd.Flags().StringVar(&couponsCreateCouponCollectionFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	couponsCmd.AddCommand(couponsCreateCouponCollectionCmd)
}

func runCouponsCreateCouponCollection(cmd *cobra.Command, args []string) error {
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
			Name:        "collection-id",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "The id of the coupon collection for which the coupons will be created",
		})
		flags = append(flags, flagSchema{
			Name:        "coupons",
			Type:        "array",
			Required:    true,
			Location:    "body",
			Description: "",
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
			Description: "Coupons creation in progress",
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
		responses = append(responses, responseSchema{
			Status:      "404",
			ContentType: "application/json",
			Description: "Coupon collection not found",
		})

		schema := map[string]any{
			"command":     "create-coupon-collection",
			"description": "Create coupons for a coupon collection",
			"http": map[string]any{
				"method": "POST",
				"path":   "/coupons",
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
		Path:        httpclient.SubstitutePath("/coupons", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if couponsCreateCouponCollectionFlags.body != "" {
		if err := json.Unmarshal([]byte(couponsCreateCouponCollectionFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("collection-id") {
		bodyMap["collectionId"] = couponsCreateCouponCollectionFlags.collectionId
	}
	if cmd.Flags().Changed("coupons") {
		bodyMap["coupons"] = couponsCreateCouponCollectionFlags.coupons
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
