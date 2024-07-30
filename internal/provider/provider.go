package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type ThetaProvider struct {
	client *Client
}

func New() provider.Provider {
	return &ThetaProvider{}
}

func (p *ThetaProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "theta"
}

func (p *ThetaProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Theta provider",
		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				MarkdownDescription: "Email for the Theta API",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password for the Theta API",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *ThetaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	log.Println("Configure method called")

	var config struct {
		Email    string `tfsdk:"email"`
		Password string `tfsdk:"password"`
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		log.Println("Error getting config:", resp.Diagnostics)
		return
	}

	log.Println("Config Email:", config.Email)

	client := NewClient(config.Email, config.Password)
	if client == nil || client.authToken == "" {
		resp.Diagnostics.AddError("Authentication Error", "Failed to authenticate with the Theta API")
		log.Println("Failed to authenticate with the Theta API")
		return
	}

	log.Println("Client successfully created with authToken:", client.authToken)
	p.client = client

	resp.ResourceData = client
	resp.DataSourceData = client
}

func (p *ThetaProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		DeploymentResource,
		DeploymentTemplateResource,
	}
}

func (p *ThetaProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		OrganizationDataSource,
		ProjectDataSource,
		DeploymentTemplateDataSource,
	}
}
