package provider

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type organizationDataSource struct {
	client *Client
}

func OrganizationDataSource() datasource.DataSource {
	return &organizationDataSource{}
}

func (d *organizationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "theta_organizations"
}

func (d *organizationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for fetching Theta organizations",
		Attributes: map[string]schema.Attribute{
			"organizations": schema.ListNestedAttribute{
				MarkdownDescription: "List of organizations",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The ID of the organization",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the organization",
							Computed:            true,
						},
						"logo_url": schema.StringAttribute{
							MarkdownDescription: "The logo URL of the organization",
							Computed:            true,
						},
						"create_time": schema.StringAttribute{
							MarkdownDescription: "The creation time of the organization",
							Computed:            true,
						},
						"user_join_time": schema.StringAttribute{
							MarkdownDescription: "The user join time to the organization",
							Computed:            true,
						},
						"user_role": schema.StringAttribute{
							MarkdownDescription: "The user's role in the organization",
							Computed:            true,
						},
						"disabled": schema.BoolAttribute{
							MarkdownDescription: "Whether the organization is disabled",
							Computed:            true,
						},
						"suspended": schema.BoolAttribute{
							MarkdownDescription: "Whether the organization is suspended",
							Computed:            true,
						},
						"email": schema.StringAttribute{
							MarkdownDescription: "The email associated with the organization",
							Computed:            true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (d *organizationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	log.Println("Data source Configure method called")

	if req.ProviderData == nil {
		log.Println("Provider data is nil")
		return
	}

	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", "Expected *Client")
		log.Println("Unexpected Data Source Configure Type")
		return
	}

	d.client = client
	log.Println("Client configured in data source")
}

func (d *organizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Client Error", "The client is not configured")
		log.Println("Client is not configured in Read method")
		return
	}

	var state struct {
		Organizations []struct {
			ID           types.String `tfsdk:"id"`
			Name         types.String `tfsdk:"name"`
			LogoURL      types.String `tfsdk:"logo_url"`
			CreateTime   types.String `tfsdk:"create_time"`
			UserJoinTime types.String `tfsdk:"user_join_time"`
			UserRole     types.String `tfsdk:"user_role"`
			Disabled     types.Bool   `tfsdk:"disabled"`
			Suspended    types.Bool   `tfsdk:"suspended"`
			Email        types.String `tfsdk:"email"`
		} `tfsdk:"organizations"`
	}

	organizations, err := d.client.GetOrganizations()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read organizations, got error: %s", err))
		log.Println("Unable to read organizations, got error:", err)
		return
	}

	for _, organization := range organizations {
		state.Organizations = append(state.Organizations, struct {
			ID           types.String `tfsdk:"id"`
			Name         types.String `tfsdk:"name"`
			LogoURL      types.String `tfsdk:"logo_url"`
			CreateTime   types.String `tfsdk:"create_time"`
			UserJoinTime types.String `tfsdk:"user_join_time"`
			UserRole     types.String `tfsdk:"user_role"`
			Disabled     types.Bool   `tfsdk:"disabled"`
			Suspended    types.Bool   `tfsdk:"suspended"`
			Email        types.String `tfsdk:"email"`
		}{
			ID:           types.StringValue(organization.ID),
			Name:         types.StringValue(organization.Name),
			LogoURL:      types.StringValue(organization.LogoURL),
			CreateTime:   types.StringValue(organization.CreateTime),
			UserJoinTime: types.StringValue(organization.UserJoinTime),
			UserRole:     types.StringValue(organization.UserRole),
			Disabled:     types.BoolValue(organization.Disabled),
			Suspended:    types.BoolValue(organization.Suspended),
			Email:        types.StringValue(organization.Email),
		})
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
