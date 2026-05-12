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

var ecommerceCreateProductCmd = &cobra.Command{
	Use:   "create-product",
	Short: "Create/Update a product",
	RunE:  runEcommerceCreateProduct,
}

var ecommerceCreateProductFlags struct {
	id            string
	name          string
	url           string
	imageUrl      string
	sku           string
	price         float64
	categories    []string
	parentId      string
	updateEnabled bool
	deletedAt     string
	body          string
}

func init() {
	ecommerceCreateProductCmd.Flags().StringVar(&ecommerceCreateProductFlags.id, "id", "", "Product ID for which you requested the details")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	ecommerceCreateProductCmd.Flags().StringVar(&ecommerceCreateProductFlags.name, "name", "", "Mandatory in case of creation**. Name of the product for which you requested the details")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	ecommerceCreateProductCmd.Flags().StringVar(&ecommerceCreateProductFlags.url, "url", "", "URL to the product")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	ecommerceCreateProductCmd.Flags().StringVar(&ecommerceCreateProductFlags.imageUrl, "image-url", "", "Absolute URL to the cover image of the product")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	ecommerceCreateProductCmd.Flags().StringVar(&ecommerceCreateProductFlags.sku, "sku", "", "Product identifier from the shop")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	ecommerceCreateProductCmd.Flags().Float64Var(&ecommerceCreateProductFlags.price, "price", 0, "Price of the product")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	ecommerceCreateProductCmd.Flags().StringSliceVar(&ecommerceCreateProductFlags.categories, "categories", nil, "Category ID-s of the product")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	ecommerceCreateProductCmd.Flags().StringVar(&ecommerceCreateProductFlags.parentId, "parent-id", "", "Parent product id of the product")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	ecommerceCreateProductCmd.Flags().BoolVar(&ecommerceCreateProductFlags.updateEnabled, "update-enabled", false, "Facilitate to update the existing category in the same request (updateEnabled = true)")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	ecommerceCreateProductCmd.Flags().StringVar(&ecommerceCreateProductFlags.deletedAt, "deleted-at", "", "UTC date-time (YYYY-MM-DDTHH:mm:ss.SSSZ) of the product deleted from the shop's database")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	ecommerceCreateProductCmd.Flags().StringVar(&ecommerceCreateProductFlags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	ecommerceCmd.AddCommand(ecommerceCreateProductCmd)
}

func runEcommerceCreateProduct(cmd *cobra.Command, args []string) error {
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
			Location:    "body",
			Description: "Product ID for which you requested the details",
		})
		flags = append(flags, flagSchema{
			Name:        "name",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Mandatory in case of creation**. Name of the product for which you requested the details",
		})
		flags = append(flags, flagSchema{
			Name:        "url",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "URL to the product",
		})
		flags = append(flags, flagSchema{
			Name:        "image-url",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Absolute URL to the cover image of the product",
		})
		flags = append(flags, flagSchema{
			Name:        "sku",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Product identifier from the shop",
		})
		flags = append(flags, flagSchema{
			Name:        "price",
			Type:        "number",
			Required:    false,
			Location:    "body",
			Description: "Price of the product",
		})
		flags = append(flags, flagSchema{
			Name:        "categories",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "Category ID-s of the product",
		})
		flags = append(flags, flagSchema{
			Name:        "parent-id",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Parent product id of the product",
		})
		flags = append(flags, flagSchema{
			Name:        "meta-info",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "Meta data of product such as description, vendor, producer, stock level. The size of cumulative metaInfo shall not exceed **1000 KB**. Maximum length of metaInfo object can be 10.",
		})
		flags = append(flags, flagSchema{
			Name:        "update-enabled",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Facilitate to update the existing category in the same request (updateEnabled = true)",
		})
		flags = append(flags, flagSchema{
			Name:        "deleted-at",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "UTC date-time (YYYY-MM-DDTHH:mm:ss.SSSZ) of the product deleted from the shop's database",
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
			Description: "Product created",
		})
		responses = append(responses, responseSchema{
			Status:      "204",
			ContentType: "",
			Description: "Product updated",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "create-product",
			"description": "Create/Update a product",
			"http": map[string]any{
				"method": "POST",
				"path":   "/products",
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
		Path:        httpclient.SubstitutePath("/products", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if ecommerceCreateProductFlags.body != "" {
		if err := json.Unmarshal([]byte(ecommerceCreateProductFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("id") {
		bodyMap["id"] = ecommerceCreateProductFlags.id
	}
	if cmd.Flags().Changed("name") {
		bodyMap["name"] = ecommerceCreateProductFlags.name
	}
	if cmd.Flags().Changed("url") {
		bodyMap["url"] = ecommerceCreateProductFlags.url
	}
	if cmd.Flags().Changed("image-url") {
		bodyMap["imageUrl"] = ecommerceCreateProductFlags.imageUrl
	}
	if cmd.Flags().Changed("sku") {
		bodyMap["sku"] = ecommerceCreateProductFlags.sku
	}
	if cmd.Flags().Changed("price") {
		bodyMap["price"] = ecommerceCreateProductFlags.price
	}
	if cmd.Flags().Changed("categories") {
		bodyMap["categories"] = ecommerceCreateProductFlags.categories
	}
	if cmd.Flags().Changed("parent-id") {
		bodyMap["parentId"] = ecommerceCreateProductFlags.parentId
	}
	if cmd.Flags().Changed("update-enabled") {
		bodyMap["updateEnabled"] = ecommerceCreateProductFlags.updateEnabled
	}
	if cmd.Flags().Changed("deleted-at") {
		bodyMap["deletedAt"] = ecommerceCreateProductFlags.deletedAt
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
