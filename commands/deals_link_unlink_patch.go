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

var dealsLinkUnlinkPatchCmd = &cobra.Command{
	Use:   "link-unlink-patch",
	Short: "Link and Unlink a deal with contacts and companies",
	RunE:  runDealsLinkUnlinkPatch,
}

var dealsLinkUnlinkPatchFlags struct {
	id               string
	linkContactIds   []string
	unlinkContactIds []string
	linkCompanyIds   []string
	unlinkCompanyIds []string
	body             string
}

func init() {
	dealsLinkUnlinkPatchCmd.Flags().StringVar(&dealsLinkUnlinkPatchFlags.id, "id", "", "")
	dealsLinkUnlinkPatchCmd.MarkFlagRequired("id")
	dealsLinkUnlinkPatchCmd.Flags().StringSliceVar(&dealsLinkUnlinkPatchFlags.linkContactIds, "link-contact-ids", nil, "Contact ids for contacts to be linked with deal")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	dealsLinkUnlinkPatchCmd.Flags().StringSliceVar(&dealsLinkUnlinkPatchFlags.unlinkContactIds, "unlink-contact-ids", nil, "Contact ids for contacts to be unlinked from deal")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	dealsLinkUnlinkPatchCmd.Flags().StringSliceVar(&dealsLinkUnlinkPatchFlags.linkCompanyIds, "link-company-ids", nil, "Company ids to be linked with deal")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	dealsLinkUnlinkPatchCmd.Flags().StringSliceVar(&dealsLinkUnlinkPatchFlags.unlinkCompanyIds, "unlink-company-ids", nil, "Company ids to be unlinked from deal")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	dealsLinkUnlinkPatchCmd.Flags().StringVar(&dealsLinkUnlinkPatchFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	dealsCmd.AddCommand(dealsLinkUnlinkPatchCmd)
}

func runDealsLinkUnlinkPatch(cmd *cobra.Command, args []string) error {
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
			Name:        "link-contact-ids",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "Contact ids for contacts to be linked with deal",
		})
		flags = append(flags, flagSchema{
			Name:        "unlink-contact-ids",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "Contact ids for contacts to be unlinked from deal",
		})
		flags = append(flags, flagSchema{
			Name:        "link-company-ids",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "Company ids to be linked with deal",
		})
		flags = append(flags, flagSchema{
			Name:        "unlink-company-ids",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "Company ids to be unlinked from deal",
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
			Description: "Successfully linked/unlinked contacts/companies with the deal.",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "Returned when query params are invalid or invalid data provided in request.",
		})

		schema := map[string]any{
			"command":     "link-unlink-patch",
			"description": "Link and Unlink a deal with contacts and companies",
			"http": map[string]any{
				"method": "PATCH",
				"path":   "/crm/deals/link-unlink/{id}",
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
	pathParams["id"] = fmt.Sprintf("%v", dealsLinkUnlinkPatchFlags.id)

	req := &httpclient.Request{
		Method:      "PATCH",
		Path:        httpclient.SubstitutePath("/crm/deals/link-unlink/{id}", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if dealsLinkUnlinkPatchFlags.body != "" {
		if err := json.Unmarshal([]byte(dealsLinkUnlinkPatchFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("link-contact-ids") {
		bodyMap["linkContactIds"] = dealsLinkUnlinkPatchFlags.linkContactIds
	}
	if cmd.Flags().Changed("unlink-contact-ids") {
		bodyMap["unlinkContactIds"] = dealsLinkUnlinkPatchFlags.unlinkContactIds
	}
	if cmd.Flags().Changed("link-company-ids") {
		bodyMap["linkCompanyIds"] = dealsLinkUnlinkPatchFlags.linkCompanyIds
	}
	if cmd.Flags().Changed("unlink-company-ids") {
		bodyMap["unlinkCompanyIds"] = dealsLinkUnlinkPatchFlags.unlinkCompanyIds
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
