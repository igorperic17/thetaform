package provider

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type projectDataSource struct {
	client *Client
}

func ProjectDataSource() datasource.DataSource {
	return &projectDataSource{}
}

func (d *projectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "theta_projects"
}

func (d *projectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for fetching projects in a Theta organization",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the organization",
				Required:            true,
			},
			"projects": schema.ListNestedAttribute{
				MarkdownDescription: "List of projects",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The ID of the project",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the project",
							Computed:            true,
						},
						"org_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the organization",
							Computed:            true,
						},
						"tva_id": schema.StringAttribute{
							MarkdownDescription: "The TVA ID of the project",
							Computed:            true,
						},
						"gateway_id": schema.StringAttribute{
							MarkdownDescription: "The gateway ID of the project",
							Computed:            true,
						},
						"create_time": schema.StringAttribute{
							MarkdownDescription: "The creation time of the project",
							Computed:            true,
						},
						"user_join_time": schema.StringAttribute{
							MarkdownDescription: "The user join time to the project",
							Computed:            true,
						},
						"user_ids": schema.ListAttribute{
							ElementType:         types.StringType,
							MarkdownDescription: "List of user IDs associated with the project",
							Computed:            true,
						},
						"user_role": schema.StringAttribute{
							MarkdownDescription: "The user's role in the project",
							Computed:            true,
						},
						"tva_secret": schema.StringAttribute{
							MarkdownDescription: "The TVA secret of the project",
							Computed:            true,
						},
						"gateway_key": schema.StringAttribute{
							MarkdownDescription: "The gateway key of the project",
							Computed:            true,
						},
						"gateway_secret": schema.StringAttribute{
							MarkdownDescription: "The gateway secret of the project",
							Computed:            true,
						},
						"disabled": schema.BoolAttribute{
							MarkdownDescription: "Whether the project is disabled",
							Computed:            true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (d *projectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *projectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Client Error", "The client is not configured")
		log.Println("Client is not configured in Read method")
		return
	}

	var state struct {
		OrganizationID types.String `tfsdk:"organization_id"`
		Projects       []struct {
			ID            types.String   `tfsdk:"id"`
			Name          types.String   `tfsdk:"name"`
			OrgID         types.String   `tfsdk:"org_id"`
			TvaID         types.String   `tfsdk:"tva_id"`
			GatewayID     types.String   `tfsdk:"gateway_id"`
			CreateTime    types.String   `tfsdk:"create_time"`
			UserJoinTime  types.String   `tfsdk:"user_join_time"`
			UserIDs       []types.String `tfsdk:"user_ids"`
			UserRole      types.String   `tfsdk:"user_role"`
			TvaSecret     types.String   `tfsdk:"tva_secret"`
			GatewayKey    types.String   `tfsdk:"gateway_key"`
			GatewaySecret types.String   `tfsdk:"gateway_secret"`
			Disabled      types.Bool     `tfsdk:"disabled"`
		} `tfsdk:"projects"`
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		log.Println("Error getting config:", resp.Diagnostics)
		return
	}

	organizationID := state.OrganizationID.ValueString()

	projects, err := d.client.GetProjects(organizationID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read projects, got error: %s", err))
		log.Println("Unable to read projects, got error:", err)
		return
	}

	for _, project := range *projects {
		var userIDs []types.String
		for _, userID := range project.UserIDs {
			userIDs = append(userIDs, types.StringValue(userID))
		}

		state.Projects = append(state.Projects, struct {
			ID            types.String   `tfsdk:"id"`
			Name          types.String   `tfsdk:"name"`
			OrgID         types.String   `tfsdk:"org_id"`
			TvaID         types.String   `tfsdk:"tva_id"`
			GatewayID     types.String   `tfsdk:"gateway_id"`
			CreateTime    types.String   `tfsdk:"create_time"`
			UserJoinTime  types.String   `tfsdk:"user_join_time"`
			UserIDs       []types.String `tfsdk:"user_ids"`
			UserRole      types.String   `tfsdk:"user_role"`
			TvaSecret     types.String   `tfsdk:"tva_secret"`
			GatewayKey    types.String   `tfsdk:"gateway_key"`
			GatewaySecret types.String   `tfsdk:"gateway_secret"`
			Disabled      types.Bool     `tfsdk:"disabled"`
		}{
			ID:            types.StringValue(project.ID),
			Name:          types.StringValue(project.Name),
			OrgID:         types.StringValue(project.OrgID),
			TvaID:         types.StringValue(project.TvaID),
			GatewayID:     stringPtrToValue(project.GatewayID),
			CreateTime:    types.StringValue(project.CreateTime),
			UserJoinTime:  types.StringValue(project.UserJoinTime),
			UserIDs:       userIDs,
			UserRole:      types.StringValue(project.UserRole),
			TvaSecret:     types.StringValue(project.TvaSecret),
			GatewayKey:    stringPtrToValue(project.GatewayKey),
			GatewaySecret: stringPtrToValue(project.GatewaySecret),
			Disabled:      types.BoolValue(project.Disabled),
		})
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func stringPtrToValue(ptr *string) types.String {
	if ptr == nil {
		return types.StringNull()
	}
	return types.StringValue(*ptr)
}
