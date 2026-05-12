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

var eventTrackInteractionCmd = &cobra.Command{
	Use:   "track-interaction",
	Short: "Create an event",
	RunE:  runEventTrackInteraction,
}

var eventTrackInteractionFlags struct {
	eventName string
	eventDate string
	body      string
}

func init() {
	eventTrackInteractionCmd.Flags().StringVar(&eventTrackInteractionFlags.eventName, "event-name", "", "The name of the event that occurred. This is how you will find your event in Brevo. Limited to 255 characters, alphanumerical characters and - _ only.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	eventTrackInteractionCmd.Flags().StringVar(&eventTrackInteractionFlags.eventDate, "event-date", "", "Timestamp of when the event occurred (e.g. \"2024-01-24T17:39:57+01:00\"). If no value is passed, the timestamp of the event creation is used.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	eventTrackInteractionCmd.Flags().StringVar(&eventTrackInteractionFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	eventCmd.AddCommand(eventTrackInteractionCmd)
}

func runEventTrackInteraction(cmd *cobra.Command, args []string) error {
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
			Name:        "event-name",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "The name of the event that occurred. This is how you will find your event in Brevo. Limited to 255 characters, alphanumerical characters and - _ only.",
		})
		flags = append(flags, flagSchema{
			Name:        "event-date",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Timestamp of when the event occurred (e.g. \"2024-01-24T17:39:57+01:00\"). If no value is passed, the timestamp of the event creation is used.",
		})
		flags = append(flags, flagSchema{
			Name:        "identifiers",
			Type:        "object",
			Required:    true,
			Location:    "body",
			Description: "Identifies the contact associated with the event. At least one identifier is required.",
		})
		flags = append(flags, flagSchema{
			Name:        "contact-properties",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "Properties defining the state of the contact associated to this event. Useful to update contact attributes defined in your contacts database while passing the event. For example: **\"FIRSTNAME\": \"Jane\" , \"AGE\": 37**",
		})
		flags = append(flags, flagSchema{
			Name:        "event-properties",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "Properties of the event. Top level properties and nested properties can be used to better segment contacts and personalise workflow conditions. The following field type are supported: string, number, boolean (true/false), date (Timestamp e.g. \"2024-01-24T17:39:57+01:00\"). Keys are limited to 255 characters, alphanumerical characters and - _ only. Size is limited to 50Kb.",
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
			Description: "An event posted",
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
			"command":     "track-interaction",
			"description": "Create an event",
			"http": map[string]any{
				"method": "POST",
				"path":   "/events",
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
		Path:        httpclient.SubstitutePath("/events", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if eventTrackInteractionFlags.body != "" {
		if err := json.Unmarshal([]byte(eventTrackInteractionFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("event-name") {
		bodyMap["event_name"] = eventTrackInteractionFlags.eventName
	}
	if cmd.Flags().Changed("event-date") {
		bodyMap["event_date"] = eventTrackInteractionFlags.eventDate
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
