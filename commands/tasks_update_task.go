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

var tasksUpdateTaskCmd = &cobra.Command{
	Use:   "update-task",
	Short: "Update a task",
	RunE:  runTasksUpdateTask,
}

var tasksUpdateTaskFlags struct {
	id           string
	name         string
	duration     int
	taskTypeId   string
	date         string
	notes        string
	done         bool
	assignToId   string
	contactsIds  []string
	dealsIds     []string
	companiesIds []string
	body         string
}

func init() {
	tasksUpdateTaskCmd.Flags().StringVar(&tasksUpdateTaskFlags.id, "id", "", "")
	tasksUpdateTaskCmd.MarkFlagRequired("id")
	tasksUpdateTaskCmd.Flags().StringVar(&tasksUpdateTaskFlags.name, "name", "", "Name of task")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	tasksUpdateTaskCmd.Flags().IntVar(&tasksUpdateTaskFlags.duration, "duration", 0, "Duration of task in milliseconds [1 minute = 60000 ms]")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	tasksUpdateTaskCmd.Flags().StringVar(&tasksUpdateTaskFlags.taskTypeId, "task-type-id", "", "Id for type of task e.g Call / Email / Meeting etc.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	tasksUpdateTaskCmd.Flags().StringVar(&tasksUpdateTaskFlags.date, "date", "", "Task date/time")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	tasksUpdateTaskCmd.Flags().StringVar(&tasksUpdateTaskFlags.notes, "notes", "", "Notes added to a task")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	tasksUpdateTaskCmd.Flags().BoolVar(&tasksUpdateTaskFlags.done, "done", false, "Task marked as done")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	tasksUpdateTaskCmd.Flags().StringVar(&tasksUpdateTaskFlags.assignToId, "assign-to-id", "", "To assign a task to a user you can use either the account email or ID.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	tasksUpdateTaskCmd.Flags().StringSliceVar(&tasksUpdateTaskFlags.contactsIds, "contacts-ids", nil, "Contact ids for contacts linked to this task")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	tasksUpdateTaskCmd.Flags().StringSliceVar(&tasksUpdateTaskFlags.dealsIds, "deals-ids", nil, "Deal ids for deals a task is linked to")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	tasksUpdateTaskCmd.Flags().StringSliceVar(&tasksUpdateTaskFlags.companiesIds, "companies-ids", nil, "Companies ids for companies a task is linked to")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	tasksUpdateTaskCmd.Flags().StringVar(&tasksUpdateTaskFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	tasksCmd.AddCommand(tasksUpdateTaskCmd)
}

func runTasksUpdateTask(cmd *cobra.Command, args []string) error {
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
			Name:        "id",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "",
		})
		flags = append(flags, flagSchema{
			Name:        "name",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Name of task",
		})
		flags = append(flags, flagSchema{
			Name:        "duration",
			Type:        "integer",
			Required:    false,
			Location:    "body",
			Description: "Duration of task in milliseconds [1 minute = 60000 ms]",
		})
		flags = append(flags, flagSchema{
			Name:        "task-type-id",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Id for type of task e.g Call / Email / Meeting etc.",
		})
		flags = append(flags, flagSchema{
			Name:        "date",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Task date/time",
		})
		flags = append(flags, flagSchema{
			Name:        "notes",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Notes added to a task",
		})
		flags = append(flags, flagSchema{
			Name:        "done",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Task marked as done",
		})
		flags = append(flags, flagSchema{
			Name:        "assign-to-id",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "To assign a task to a user you can use either the account email or ID.",
		})
		flags = append(flags, flagSchema{
			Name:        "contacts-ids",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "Contact ids for contacts linked to this task",
		})
		flags = append(flags, flagSchema{
			Name:        "deals-ids",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "Deal ids for deals a task is linked to",
		})
		flags = append(flags, flagSchema{
			Name:        "companies-ids",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "Companies ids for companies a task is linked to",
		})
		flags = append(flags, flagSchema{
			Name:        "reminder",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "Task reminder date/time for a task",
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
			Description: "Task updated succussfully",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "Returned when invalid data posted",
		})
		responses = append(responses, responseSchema{
			Status:      "404",
			ContentType: "application/json",
			Description: "Returned when task id is not found",
		})

		schema := map[string]any{
			"command":     "update-task",
			"description": "Update a task",
			"http": map[string]any{
				"method": "PATCH",
				"path":   "/crm/tasks/{id}",
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
	pathParams["id"] = fmt.Sprintf("%v", tasksUpdateTaskFlags.id)

	req := &httpclient.Request{
		Method:      "PATCH",
		Path:        httpclient.SubstitutePath("/crm/tasks/{id}", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if tasksUpdateTaskFlags.body != "" {
		if err := json.Unmarshal([]byte(tasksUpdateTaskFlags.body), &bodyMap); err != nil {
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
		bodyMap["name"] = tasksUpdateTaskFlags.name
	}
	if cmd.Flags().Changed("duration") {
		bodyMap["duration"] = tasksUpdateTaskFlags.duration
	}
	if cmd.Flags().Changed("task-type-id") {
		bodyMap["taskTypeId"] = tasksUpdateTaskFlags.taskTypeId
	}
	if cmd.Flags().Changed("date") {
		bodyMap["date"] = tasksUpdateTaskFlags.date
	}
	if cmd.Flags().Changed("notes") {
		bodyMap["notes"] = tasksUpdateTaskFlags.notes
	}
	if cmd.Flags().Changed("done") {
		bodyMap["done"] = tasksUpdateTaskFlags.done
	}
	if cmd.Flags().Changed("assign-to-id") {
		bodyMap["assignToId"] = tasksUpdateTaskFlags.assignToId
	}
	if cmd.Flags().Changed("contacts-ids") {
		bodyMap["contactsIds"] = tasksUpdateTaskFlags.contactsIds
	}
	if cmd.Flags().Changed("deals-ids") {
		bodyMap["dealsIds"] = tasksUpdateTaskFlags.dealsIds
	}
	if cmd.Flags().Changed("companies-ids") {
		bodyMap["companiesIds"] = tasksUpdateTaskFlags.companiesIds
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
