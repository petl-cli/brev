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

var ecommerceListAllProductsCmd = &cobra.Command{
	Use:   "list-all-products",
	Short: "Return all your products",
	RunE:  runEcommerceListAllProducts,
}

var ecommerceListAllProductsFlags struct {
	limit         int
	offset        int
	sort          string
	ids           []string
	name          string
	priceLte      float64
	priceGte      float64
	priceLt       float64
	priceGt       float64
	priceEq       float64
	priceNe       float64
	categories    []string
	modifiedSince string
	createdSince  string
}

func init() {
	ecommerceListAllProductsCmd.Flags().IntVar(&ecommerceListAllProductsFlags.limit, "limit", 0, "Number of documents per page")
	ecommerceListAllProductsCmd.Flags().IntVar(&ecommerceListAllProductsFlags.offset, "offset", 0, "Index of the first document in the page")
	ecommerceListAllProductsCmd.Flags().StringVar(&ecommerceListAllProductsFlags.sort, "sort", "", "Sort the results in the ascending/descending order of record creation. Default order is **descending** if `sort` is not passed")
	ecommerceListAllProductsCmd.Flags().StringSliceVar(&ecommerceListAllProductsFlags.ids, "ids", nil, "Filter by product ids")
	ecommerceListAllProductsCmd.Flags().StringVar(&ecommerceListAllProductsFlags.name, "name", "", "Filter by product name, minimum 3 characters should be present for search")
	ecommerceListAllProductsCmd.Flags().Float64Var(&ecommerceListAllProductsFlags.priceLte, "price-lte", 0, "Price filter for products less than and equals to particular amount")
	ecommerceListAllProductsCmd.Flags().Float64Var(&ecommerceListAllProductsFlags.priceGte, "price-gte", 0, "Price filter for products greater than and equals to particular amount")
	ecommerceListAllProductsCmd.Flags().Float64Var(&ecommerceListAllProductsFlags.priceLt, "price-lt", 0, "Price filter for products less than particular amount")
	ecommerceListAllProductsCmd.Flags().Float64Var(&ecommerceListAllProductsFlags.priceGt, "price-gt", 0, "Price filter for products greater than particular amount")
	ecommerceListAllProductsCmd.Flags().Float64Var(&ecommerceListAllProductsFlags.priceEq, "price-eq", 0, "Price filter for products equals to particular amount")
	ecommerceListAllProductsCmd.Flags().Float64Var(&ecommerceListAllProductsFlags.priceNe, "price-ne", 0, "Price filter for products not equals to particular amount")
	ecommerceListAllProductsCmd.Flags().StringSliceVar(&ecommerceListAllProductsFlags.categories, "categories", nil, "Filter by product categories")
	ecommerceListAllProductsCmd.Flags().StringVar(&ecommerceListAllProductsFlags.modifiedSince, "modified-since", "", "Filter (urlencoded) the orders modified after a given UTC date-time (YYYY-MM-DDTHH:mm:ss.SSSZ). **Prefer to pass your timezone in date-time format for accurate result.** ")
	ecommerceListAllProductsCmd.Flags().StringVar(&ecommerceListAllProductsFlags.createdSince, "created-since", "", "Filter (urlencoded) the orders created after a given UTC date-time (YYYY-MM-DDTHH:mm:ss.SSSZ). **Prefer to pass your timezone in date-time format for accurate result.** ")

	ecommerceCmd.AddCommand(ecommerceListAllProductsCmd)
}

func runEcommerceListAllProducts(cmd *cobra.Command, args []string) error {
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
			Name:        "limit",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Number of documents per page",
		})
		flags = append(flags, flagSchema{
			Name:        "offset",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Index of the first document in the page",
		})
		flags = append(flags, flagSchema{
			Name:        "sort",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Sort the results in the ascending/descending order of record creation. Default order is **descending** if `sort` is not passed",
		})
		flags = append(flags, flagSchema{
			Name:        "ids",
			Type:        "array",
			Required:    false,
			Location:    "query",
			Description: "Filter by product ids",
		})
		flags = append(flags, flagSchema{
			Name:        "name",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter by product name, minimum 3 characters should be present for search",
		})
		flags = append(flags, flagSchema{
			Name:        "price-lte",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "Price filter for products less than and equals to particular amount",
		})
		flags = append(flags, flagSchema{
			Name:        "price-gte",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "Price filter for products greater than and equals to particular amount",
		})
		flags = append(flags, flagSchema{
			Name:        "price-lt",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "Price filter for products less than particular amount",
		})
		flags = append(flags, flagSchema{
			Name:        "price-gt",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "Price filter for products greater than particular amount",
		})
		flags = append(flags, flagSchema{
			Name:        "price-eq",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "Price filter for products equals to particular amount",
		})
		flags = append(flags, flagSchema{
			Name:        "price-ne",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "Price filter for products not equals to particular amount",
		})
		flags = append(flags, flagSchema{
			Name:        "categories",
			Type:        "array",
			Required:    false,
			Location:    "query",
			Description: "Filter by product categories",
		})
		flags = append(flags, flagSchema{
			Name:        "modified-since",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter (urlencoded) the orders modified after a given UTC date-time (YYYY-MM-DDTHH:mm:ss.SSSZ). **Prefer to pass your timezone in date-time format for accurate result.** ",
		})
		flags = append(flags, flagSchema{
			Name:        "created-since",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter (urlencoded) the orders created after a given UTC date-time (YYYY-MM-DDTHH:mm:ss.SSSZ). **Prefer to pass your timezone in date-time format for accurate result.** ",
		})

		type responseSchema struct {
			Status      string `json:"status"`
			ContentType string `json:"content_type,omitempty"`
			Description string `json:"description,omitempty"`
		}
		var responses []responseSchema
		responses = append(responses, responseSchema{
			Status:      "200",
			ContentType: "application/json",
			Description: "All products listed",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "application/json",
			Description: "bad request",
		})

		schema := map[string]any{
			"command":     "list-all-products",
			"description": "Return all your products",
			"http": map[string]any{
				"method": "GET",
				"path":   "/products",
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
				"safe":         true,
				"idempotent":   true,
				"reversible":   true,
				"side_effects": []string{},
				"impact":       "low",
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
		Method:      "GET",
		Path:        httpclient.SubstitutePath("/products", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("limit") {
		req.QueryParams["limit"] = fmt.Sprintf("%v", ecommerceListAllProductsFlags.limit)
	}
	if cmd.Flags().Changed("offset") {
		req.QueryParams["offset"] = fmt.Sprintf("%v", ecommerceListAllProductsFlags.offset)
	}
	if cmd.Flags().Changed("sort") {
		req.QueryParams["sort"] = fmt.Sprintf("%v", ecommerceListAllProductsFlags.sort)
	}
	if cmd.Flags().Changed("ids") {
		req.ArrayParams["ids"] = ecommerceListAllProductsFlags.ids
	}
	if cmd.Flags().Changed("name") {
		req.QueryParams["name"] = fmt.Sprintf("%v", ecommerceListAllProductsFlags.name)
	}
	if cmd.Flags().Changed("price-lte") {
		req.QueryParams["price[lte]"] = fmt.Sprintf("%v", ecommerceListAllProductsFlags.priceLte)
	}
	if cmd.Flags().Changed("price-gte") {
		req.QueryParams["price[gte]"] = fmt.Sprintf("%v", ecommerceListAllProductsFlags.priceGte)
	}
	if cmd.Flags().Changed("price-lt") {
		req.QueryParams["price[lt]"] = fmt.Sprintf("%v", ecommerceListAllProductsFlags.priceLt)
	}
	if cmd.Flags().Changed("price-gt") {
		req.QueryParams["price[gt]"] = fmt.Sprintf("%v", ecommerceListAllProductsFlags.priceGt)
	}
	if cmd.Flags().Changed("price-eq") {
		req.QueryParams["price[eq]"] = fmt.Sprintf("%v", ecommerceListAllProductsFlags.priceEq)
	}
	if cmd.Flags().Changed("price-ne") {
		req.QueryParams["price[ne]"] = fmt.Sprintf("%v", ecommerceListAllProductsFlags.priceNe)
	}
	if cmd.Flags().Changed("categories") {
		req.ArrayParams["categories"] = ecommerceListAllProductsFlags.categories
	}
	if cmd.Flags().Changed("modified-since") {
		req.QueryParams["modifiedSince"] = fmt.Sprintf("%v", ecommerceListAllProductsFlags.modifiedSince)
	}
	if cmd.Flags().Changed("created-since") {
		req.QueryParams["createdSince"] = fmt.Sprintf("%v", ecommerceListAllProductsFlags.createdSince)
	}

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
