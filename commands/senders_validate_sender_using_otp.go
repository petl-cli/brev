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

var sendersValidateSenderUsingOtpCmd = &cobra.Command{
	Use:   "validate-sender-using-otp",
	Short: "Validate Sender using OTP",
	RunE:  runSendersValidateSenderUsingOtp,
}

var sendersValidateSenderUsingOtpFlags struct {
	senderId int
	otp      int
	body     string
}

func init() {
	sendersValidateSenderUsingOtpCmd.Flags().IntVar(&sendersValidateSenderUsingOtpFlags.senderId, "sender-id", 0, "Id of the sender")
	sendersValidateSenderUsingOtpCmd.MarkFlagRequired("sender-id")
	sendersValidateSenderUsingOtpCmd.Flags().IntVar(&sendersValidateSenderUsingOtpFlags.otp, "otp", 0, "6 digit OTP received on email")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	sendersValidateSenderUsingOtpCmd.Flags().StringVar(&sendersValidateSenderUsingOtpFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	sendersCmd.AddCommand(sendersValidateSenderUsingOtpCmd)
}

func runSendersValidateSenderUsingOtp(cmd *cobra.Command, args []string) error {
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
			Name:        "sender-id",
			Type:        "integer",
			Required:    true,
			Location:    "path",
			Description: "Id of the sender",
		})
		flags = append(flags, flagSchema{
			Name:        "otp",
			Type:        "integer",
			Required:    false,
			Location:    "body",
			Description: "6 digit OTP received on email",
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
			Description: "Sender verified",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})
		responses = append(responses, responseSchema{
			Status:      "404",
			ContentType: "application/json",
			Description: "Sender ID not found",
		})

		schema := map[string]any{
			"command":     "validate-sender-using-otp",
			"description": "Validate Sender using OTP",
			"http": map[string]any{
				"method": "PUT",
				"path":   "/senders/{senderId}/validate",
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
	pathParams["senderId"] = fmt.Sprintf("%v", sendersValidateSenderUsingOtpFlags.senderId)

	req := &httpclient.Request{
		Method:      "PUT",
		Path:        httpclient.SubstitutePath("/senders/{senderId}/validate", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if sendersValidateSenderUsingOtpFlags.body != "" {
		if err := json.Unmarshal([]byte(sendersValidateSenderUsingOtpFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("otp") {
		bodyMap["otp"] = sendersValidateSenderUsingOtpFlags.otp
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
