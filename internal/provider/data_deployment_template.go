package provider

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DataSource for deployment templates
type deploymentTemplateDataSource struct {
	client *Client
}

func DeploymentTemplateDataSource() datasource.DataSource {
	return &deploymentTemplateDataSource{}
}

func (d *deploymentTemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "theta_deployment_templates"
}

func (d *deploymentTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for fetching Theta deployment templates",
		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the project",
				Required:            true,
			},
			"deployment_templates": schema.ListNestedAttribute{
				MarkdownDescription: "List of deployment templates",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The ID of the deployment template",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the deployment template",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description of the deployment template",
							Computed:            true,
						},
						"tags": schema.ListAttribute{
							ElementType:         types.StringType,
							MarkdownDescription: "The tags of the deployment template",
							Computed:            true,
						},
						"category": schema.StringAttribute{
							MarkdownDescription: "The category of the deployment template",
							Computed:            true,
						},
						"project_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the project",
							Computed:            true,
						},
						"container_images": schema.ListAttribute{
							ElementType:         types.StringType,
							MarkdownDescription: "The container image of the deployment template",
							Computed:            true,
						},
						"container_port": schema.Int64Attribute{
							MarkdownDescription: "The container port of the deployment template",
							Computed:            true,
						},
						"container_args": schema.ListAttribute{
							ElementType:         types.StringType,
							MarkdownDescription: "The container arguments of the deployment template",
							Computed:            true,
						},
						"env_vars": schema.MapAttribute{
							ElementType:         types.StringType,
							MarkdownDescription: "The environment variables of the deployment template",
							Computed:            true,
						},
						"require_env_vars": schema.BoolAttribute{
							MarkdownDescription: "Whether the deployment template requires environment variables",
							Computed:            true,
						},
						"rank": schema.Int64Attribute{
							MarkdownDescription: "The rank of the deployment template",
							Computed:            true,
						},
						"icon_url": schema.StringAttribute{
							MarkdownDescription: "The icon URL of the deployment template",
							Computed:            true,
						},
						"create_time": schema.StringAttribute{
							MarkdownDescription: "The creation time of the deployment template",
							Computed:            true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (d *deploymentTemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *deploymentTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Client Error", "The client is not configured")
		log.Println("Client is not configured in Read method")
		return
	}

	var state struct {
		ProjectID           types.String `tfsdk:"project_id"`
		DeploymentTemplates []struct {
			ID              types.String            `tfsdk:"id"`
			Name            types.String            `tfsdk:"name"`
			Description     types.String            `tfsdk:"description"`
			Tags            []types.String          `tfsdk:"tags"`
			Category        types.String            `tfsdk:"category"`
			ProjectID       types.String            `tfsdk:"project_id"`
			ContainerImages []types.String          `tfsdk:"container_images"`
			ContainerPort   types.Int64             `tfsdk:"container_port"`
			ContainerArgs   []types.String          `tfsdk:"container_args"`
			EnvVars         map[string]types.String `tfsdk:"env_vars"`
			RequireEnvVars  types.Bool              `tfsdk:"require_env_vars"`
			Rank            types.Int64             `tfsdk:"rank"`
			IconURL         types.String            `tfsdk:"icon_url"`
			CreateTime      types.String            `tfsdk:"create_time"`
		} `tfsdk:"deployment_templates"`
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		log.Println("Error getting config:", resp.Diagnostics)
		return
	}

	projectID := state.ProjectID.ValueString()

	templates, err := d.client.GetDeploymentTemplates(projectID, 0, 100)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read deployment templates, got error: %s", err))
		log.Println("Unable to read deployment templates, got error:", err)
		return
	}

	for _, template := range templates {
		envVars := make(map[string]types.String)
		for k, v := range template.EnvVars {
			envVars[k] = types.StringValue(v)
		}

		tags := make([]types.String, len(template.Tags))
		for i, tag := range template.Tags {
			tags[i] = types.StringValue(tag)
		}

		containerImages := make([]types.String, len(template.ContainerImages))
		for i, img := range template.ContainerImages {
			containerImages[i] = types.StringValue(img)
		}

		containerArgs := make([]types.String, len(template.ContainerArgs))
		for i, arg := range template.ContainerArgs {
			containerArgs[i] = types.StringValue(arg)
		}

		var requireEnvVars types.Bool
		if template.RequireEnvVars != nil {
			requireEnvVars = types.BoolValue(*template.RequireEnvVars)
		} else {
			requireEnvVars = types.BoolNull()
		}

		var rank types.Int64
		if template.Rank != nil {
			rank = types.Int64Value(*template.Rank)
		} else {
			rank = types.Int64Null()
		}

		state.DeploymentTemplates = append(state.DeploymentTemplates, struct {
			ID              types.String            `tfsdk:"id"`
			Name            types.String            `tfsdk:"name"`
			Description     types.String            `tfsdk:"description"`
			Tags            []types.String          `tfsdk:"tags"`
			Category        types.String            `tfsdk:"category"`
			ProjectID       types.String            `tfsdk:"project_id"`
			ContainerImages []types.String          `tfsdk:"container_images"`
			ContainerPort   types.Int64             `tfsdk:"container_port"`
			ContainerArgs   []types.String          `tfsdk:"container_args"`
			EnvVars         map[string]types.String `tfsdk:"env_vars"`
			RequireEnvVars  types.Bool              `tfsdk:"require_env_vars"`
			Rank            types.Int64             `tfsdk:"rank"`
			IconURL         types.String            `tfsdk:"icon_url"`
			CreateTime      types.String            `tfsdk:"create_time"`
		}{
			ID:              types.StringValue(template.ID),
			Name:            types.StringValue(template.Name),
			Description:     types.StringValue(template.Description),
			Tags:            tags,
			Category:        types.StringValue(template.Category),
			ProjectID:       types.StringValue(template.ProjectID),
			ContainerImages: containerImages,
			ContainerPort:   types.Int64Value(template.ContainerPort),
			ContainerArgs:   containerArgs,
			EnvVars:         envVars,
			RequireEnvVars:  requireEnvVars,
			Rank:            rank,
			IconURL:         types.StringValue(template.IconURL),
			CreateTime:      types.StringValue(template.CreateTime.Format(time.RFC3339)),
		})
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
