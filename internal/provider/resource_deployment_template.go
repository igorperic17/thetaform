package provider

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Resource for deployment templates
type deploymentTemplateResource struct {
	client *Client
}

func DeploymentTemplateResource() resource.Resource {
	return &deploymentTemplateResource{}
}

func (r *deploymentTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "theta_deployment_template"
}

func (r *deploymentTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing Theta deployment templates",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the deployment template",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the deployment template",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the deployment template",
				Optional:            true,
			},
			"tags": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "The tags of the deployment template",
				Optional:            true,
			},
			"category": schema.StringAttribute{
				MarkdownDescription: "The category of the deployment template",
				Computed:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the project",
				Required:            true,
			},
			"container_image": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "The container image of the deployment template",
				Required:            true,
			},
			"container_port": schema.Int64Attribute{
				MarkdownDescription: "The container port of the deployment template",
				Optional:            true,
			},
			"container_args": schema.StringAttribute{
				MarkdownDescription: "The container arguments of the deployment template",
				Optional:            true,
			},
			"env_vars": schema.MapAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "The environment variables of the deployment template",
				Optional:            true,
			},
			"require_env_vars": schema.BoolAttribute{
				MarkdownDescription: "Whether the deployment template requires environment variables",
				Optional:            true,
			},
			"rank": schema.Int64Attribute{
				MarkdownDescription: "The rank of the deployment template",
				Optional:            true,
			},
			"icon_url": schema.StringAttribute{
				MarkdownDescription: "The icon URL of the deployment template",
				Optional:            true,
			},
			"create_time": schema.StringAttribute{
				MarkdownDescription: "The creation time of the deployment template",
				Computed:            true,
			},
		},
	}
}

func (r *deploymentTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	log.Println("DEBUG: Resource Configure method called")

	if req.ProviderData == nil {
		log.Println("DEBUG: Provider data is nil")
		return
	}

	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", "Expected *Client")
		log.Println("DEBUG: Unexpected Resource Configure Type")
		return
	}

	r.client = client
	log.Println("DEBUG: Client configured in resource")
}

func (r *deploymentTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	log.Println("DEBUG: Entering Create method")

	// Extract the plan
	var plan DeploymentTemplateRequest
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		log.Println("DEBUG: Error getting plan:", resp.Diagnostics)
		return
	}
	log.Printf("DEBUG: Plan received: %+v\n", plan)

	// Convert plan to native request format
	nativePlan := convertToNativePlan(plan)
	log.Printf("DEBUG: Native plan: %+v\n", nativePlan)

	// Call the API to create the deployment template
	template, err := r.client.CreateDeploymentTemplate(nativePlan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating deployment template",
			"Could not create deployment template, unexpected error: "+err.Error(),
		)
		return
	}
	log.Printf("DEBUG: Template created: %+v\n", template)

	// Set the state
	state := convertToTerraformState(template)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *deploymentTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state struct {
		ID        types.String `tfsdk:"id"`
		ProjectID types.String `tfsdk:"project_id"`
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	template, err := r.client.GetDeploymentTemplateByID(state.ProjectID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read deployment template, got error: %s", err))
		return
	}

	newState := convertToTerraformState(template)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *deploymentTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DeploymentTemplateRequest
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state struct {
		ID        types.String `tfsdk:"id"`
		ProjectID types.String `tfsdk:"project_id"`
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nativePlan := convertToNativePlan(plan)
	template, err := r.client.UpdateDeploymentTemplate(state.ID.ValueString(), nativePlan)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update deployment template, got error: %s", err))
		return
	}

	newState := convertToTerraformState(template)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *deploymentTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state struct {
		ID        types.String `tfsdk:"id"`
		ProjectID types.String `tfsdk:"project_id"`
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	success, err := r.client.DeleteDeploymentTemplate(state.ID.ValueString(), state.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete deployment template, got error: %s", err))
		return
	}

	if !success {
		resp.Diagnostics.AddError("Client Error", "Failed to delete deployment template")
		return
	}

	resp.State.RemoveResource(ctx)
}

func convertToTerraformState(template *DeploymentTemplate) struct {
	ID             types.String            `tfsdk:"id"`
	Name           types.String            `tfsdk:"name"`
	Description    types.String            `tfsdk:"description"`
	Tags           []types.String          `tfsdk:"tags"`
	Category       types.String            `tfsdk:"category"`
	ProjectID      types.String            `tfsdk:"project_id"`
	ContainerImage []types.String          `tfsdk:"container_image"`
	ContainerPort  types.Int64             `tfsdk:"container_port"`
	ContainerArgs  types.String            `tfsdk:"container_args"`
	EnvVars        map[string]types.String `tfsdk:"env_vars"`
	RequireEnvVars types.Bool              `tfsdk:"require_env_vars"`
	Rank           types.Int64             `tfsdk:"rank"`
	IconURL        types.String            `tfsdk:"icon_url"`
	CreateTime     types.String            `tfsdk:"create_time"`
} {
	return struct {
		ID             types.String            `tfsdk:"id"`
		Name           types.String            `tfsdk:"name"`
		Description    types.String            `tfsdk:"description"`
		Tags           []types.String          `tfsdk:"tags"`
		Category       types.String            `tfsdk:"category"`
		ProjectID      types.String            `tfsdk:"project_id"`
		ContainerImage []types.String          `tfsdk:"container_image"`
		ContainerPort  types.Int64             `tfsdk:"container_port"`
		ContainerArgs  types.String            `tfsdk:"container_args"`
		EnvVars        map[string]types.String `tfsdk:"env_vars"`
		RequireEnvVars types.Bool              `tfsdk:"require_env_vars"`
		Rank           types.Int64             `tfsdk:"rank"`
		IconURL        types.String            `tfsdk:"icon_url"`
		CreateTime     types.String            `tfsdk:"create_time"`
	}{
		ID:             types.StringValue(template.ID),
		Name:           types.StringValue(template.Name),
		Description:    types.StringValue(template.Description),
		Tags:           convertToTypesStringSlice(template.Tags),
		Category:       types.StringValue(template.Category),
		ProjectID:      types.StringValue(template.ProjectID),
		ContainerImage: convertToTypesStringSlice(template.ContainerImage),
		ContainerPort:  types.Int64Value(template.ContainerPort),
		ContainerArgs:  types.StringValue(template.ContainerArgs),
		EnvVars:        convertToTypesStringMap(template.EnvVars),
		RequireEnvVars: types.BoolValue(template.RequireEnvVars),
		Rank:           types.Int64Value(template.Rank),
		IconURL:        types.StringValue(template.IconURL),
		CreateTime:     types.StringValue(template.CreateTime.Format(time.RFC3339)),
	}
}

func convertToTypesStringSlice(input []string) []types.String {
	result := make([]types.String, len(input))
	for i, v := range input {
		result[i] = types.StringValue(v)
	}
	return result
}

func convertToTypesStringMap(input map[string]string) map[string]types.String {
	result := make(map[string]types.String, len(input))
	for k, v := range input {
		result[k] = types.StringValue(v)
	}
	return result
}

func convertToNativePlan(plan DeploymentTemplateRequest) DeploymentTemplateRequestNative {
	native := DeploymentTemplateRequestNative{
		Name:           plan.Name,
		ProjectID:      plan.ProjectID,
		Description:    plan.Description.ValueString(),
		ContainerImage: plan.ContainerImage,
		ContainerPort:  plan.ContainerPort.ValueInt64(),
		ContainerArgs:  plan.ContainerArgs.ValueString(),
		EnvVars:        plan.EnvVars,
		Tags:           plan.Tags,
		IconURL:        plan.IconURL.ValueString(),
	}
	if plan.RequireEnvVars != nil {
		native.RequireEnvVars = *plan.RequireEnvVars
	}
	if plan.Rank != nil {
		native.Rank = *plan.Rank
	}
	return native
}
