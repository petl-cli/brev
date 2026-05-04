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

var resellerDeleteSenderDomainByChildIdentifierAndDomainNameCmd = &cobra.Command{
	Use:   "delete-sender-domain-by-child-identifier-and-domain-name",
	Short: "Delete the sender domain of the reseller child based on the childIdentifier and domainName passed",
	RunE:  runResellerDeleteSenderDomainByChildIdentifierAndDomainName,
}

var resellerDeleteSenderDomainByChildIdentifierAndDomainNameFlags struct {
	childIdentifier string
	domainName      string
}

func init() {
	resellerDeleteSenderDomainByChildIdentifierAndDomainNameCmd.Flags().StringVar(&resellerDeleteSenderDomainByChildIdentifierAndDomainNameFlags.childIdentifier, "child-identifier", "", "Either auth key or id of reseller's child")
	resellerDeleteSenderDomainByChildIdentifierAndDomainNameCmd.MarkFlagRequired("child-identifier")
	resellerDeleteSenderDomainByChildIdentifierAndDomainNameCmd.Flags().StringVar(&resellerDeleteSenderDomainByChildIdentifierAndDomainNameFlags.domainName, "domain-name", "", "Pass the existing domain that needs to be deleted")
	resellerDeleteSenderDomainByChildIdentifierAndDomainNameCmd.MarkFlagRequired("domain-name")

	resellerCmd.AddCommand(resellerDeleteSenderDomainByChildIdentifierAndDomainNameCmd)
}

func runResellerDeleteSenderDomainByChildIdentifierAndDomainName(cmd *cobra.Command, args []string) error {
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
			Name:        "child-identifier",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "Either auth key or id of reseller's child",
		})
		flags = append(flags, flagSchema{
			Name:        "domain-name",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "Pass the existing domain that needs to be deleted",
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
			Description: "child domain deleted",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})
		responses = append(responses, responseSchema{
			Status:      "403",
			ContentType: "application/json",
			Description: "Current account is not a reseller",
		})
		responses = append(responses, responseSchema{
			Status:      "404",
			ContentType: "application/json",
			Description: "Child auth key or child id not found",
		})

		schema := map[string]any{
			"command":     "delete-sender-domain-by-child-identifier-and-domain-name",
			"description": "Delete the sender domain of the reseller child based on the childIdentifier and domainName passed",
			"http": map[string]any{
				"method": "DELETE",
				"path":   "/reseller/children/{childIdentifier}/domains/{domainName}",
			},
			"input": map[string]any{
				"flags":         flags,
				"body_flag":     false,
				"body_required": false,
			},
			"output": map[string]any{
				"responses": responses,
			},
			"semantics": map[string]any{
				"safe":         false,
				"idempotent":   true,
				"reversible":   false,
				"side_effects": []string{"destroys_resource"},
				"impact":       "high",
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
	pathParams["childIdentifier"] = fmt.Sprintf("%v", resellerDeleteSenderDomainByChildIdentifierAndDomainNameFlags.childIdentifier)
	pathParams["domainName"] = fmt.Sprintf("%v", resellerDeleteSenderDomainByChildIdentifierAndDomainNameFlags.domainName)

	req := &httpclient.Request{
		Method:      "DELETE",
		Path:        httpclient.SubstitutePath("/reseller/children/{childIdentifier}/domains/{domainName}", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

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
