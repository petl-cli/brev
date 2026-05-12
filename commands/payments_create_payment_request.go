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

var paymentsCreatePaymentRequestCmd = &cobra.Command{
	Use:   "create-payment-request",
	Short: "Create a payment request",
	RunE:  runPaymentsCreatePaymentRequest,
}

var paymentsCreatePaymentRequestFlags struct {
	reference     string
	contactId     int
	configuration string
	body          string
}

func init() {
	paymentsCreatePaymentRequestCmd.Flags().StringVar(&paymentsCreatePaymentRequestFlags.reference, "reference", "", "Reference of the payment request, it will appear on the payment page. ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	paymentsCreatePaymentRequestCmd.Flags().IntVar(&paymentsCreatePaymentRequestFlags.contactId, "contact-id", 0, "Brevo ID of the contact requested to pay. ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	paymentsCreatePaymentRequestCmd.Flags().StringVar(&paymentsCreatePaymentRequestFlags.configuration, "configuration", "", "Optional. Redirect contact to a custom success page once payment is successful. If empty the default Brevo page will be displayed once a payment is validated ")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	paymentsCreatePaymentRequestCmd.Flags().StringVar(&paymentsCreatePaymentRequestFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	paymentsCmd.AddCommand(paymentsCreatePaymentRequestCmd)
}

func runPaymentsCreatePaymentRequest(cmd *cobra.Command, args []string) error {
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
			Name:        "reference",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Reference of the payment request, it will appear on the payment page. ",
		})
		flags = append(flags, flagSchema{
			Name:        "cart",
			Type:        "object",
			Required:    true,
			Location:    "body",
			Description: "Specify the payment currency and amount. ",
		})
		flags = append(flags, flagSchema{
			Name:        "contact-id",
			Type:        "integer",
			Required:    true,
			Location:    "body",
			Description: "Brevo ID of the contact requested to pay. ",
		})
		flags = append(flags, flagSchema{
			Name:        "notification",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "Optional. Use this object if you want to let Brevo send an email to the contact, with the payment request URL. If empty, no notifications (message and reminders) will be sent. ",
		})
		flags = append(flags, flagSchema{
			Name:        "configuration",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Optional. Redirect contact to a custom success page once payment is successful. If empty the default Brevo page will be displayed once a payment is validated ",
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
			Description: "Payment request created.",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "Bad request.",
		})
		responses = append(responses, responseSchema{
			Status:      "401",
			ContentType: "application/json",
			Description: "Unauthorized.",
		})
		responses = append(responses, responseSchema{
			Status:      "403",
			ContentType: "application/json",
			Description: "Permission denied. Either you don't have access to Brevo Payments or your Brevo Payments account is not validated.",
		})

		schema := map[string]any{
			"command":     "create-payment-request",
			"description": "Create a payment request",
			"http": map[string]any{
				"method": "POST",
				"path":   "/payments/requests",
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
		Path:        httpclient.SubstitutePath("/payments/requests", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if paymentsCreatePaymentRequestFlags.body != "" {
		if err := json.Unmarshal([]byte(paymentsCreatePaymentRequestFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("reference") {
		bodyMap["reference"] = paymentsCreatePaymentRequestFlags.reference
	}
	if cmd.Flags().Changed("contact-id") {
		bodyMap["contactId"] = paymentsCreatePaymentRequestFlags.contactId
	}
	if cmd.Flags().Changed("configuration") {
		bodyMap["configuration"] = paymentsCreatePaymentRequestFlags.configuration
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
