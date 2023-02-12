package provider

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"time"

	"github.com/TGNThump/terraform-provider-vyos/internal/vyos"
	"github.com/foltik/vyos-client-go/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure VyOSProvider satisfies various provider interfaces.
var _ provider.Provider = &VyOSProvider{}

// VyOSProvider defines the provider implementation.
type VyOSProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// VyOSProviderModel describes the provider data model.
type VyOSProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	ApiKey   types.String `tfsdk:"api_key"`
}

func (p *VyOSProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "vyos"
	resp.Version = p.version
}

func (p *VyOSProvider) Schema(ctx context.Context, request provider.SchemaRequest, response *provider.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Endpoint of the VyOS HTTP API",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API Key for the VyOS HTTP API",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *VyOSProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config VyOSProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown VyOS API Endpoint",
			"The provider cannot create the VyOS API client as there is an unknown configuration value for the VyOS API Endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VYOS_ENDPOINT environment variable.",
		)
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown VyOS API key",
			"The provider cannot create the VyOS API client as there is an unknown configuration value for the VyOS API ApiKey. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VYOS_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := os.Getenv("VYOS_ENDPOINT")
	api_key := os.Getenv("VYOS_API_KEY")

	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	if !config.ApiKey.IsNull() {
		api_key = config.ApiKey.ValueString()
	}

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing VyOS API Endpoint",
			"The provider cannot create the VyOS API client as there is an missing or empty value for the VyOS API Endpoint. "+
				"Set the username value in the configuration or use the VYOS_ENDPOINT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if api_key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing VyOS API key",
			"The provider cannot create the VyOS API client as there is an missing or empty value for the VyOS API Key. "+
				"Set the username value in the configuration or use the VYOS_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	httpClient := &http.Client{Transport: transport, Timeout: 10 * time.Minute}

	apiClient := client.NewWithClient(httpClient, endpoint, api_key)
	vyosConfig := vyos.New(apiClient)

	resp.DataSourceData = vyosConfig
	resp.ResourceData = vyosConfig
}

func (p *VyOSProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewConfigResource,
	}
}

func (p *VyOSProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &VyOSProvider{
			version: version,
		}
	}
}
