package provider

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Resource for deployment
type deploymentResource struct {
	client *Client
}

func DeploymentResource() resource.Resource {
	return &deploymentResource{}
}

func (r *deploymentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "theta_deployment"
}
func (r *deploymentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing Theta deployments",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the deployment",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the deployment",
				Required:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the project",
				Required:            true,
			},
			"deployment_image_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the deployment image",
				Required:            true,
			},
			"container_image": schema.StringAttribute{
				MarkdownDescription: "The container image",
				Required:            true,
			},
			"min_replicas": schema.Int64Attribute{
				MarkdownDescription: "Minimum number of replicas",
				Required:            true,
			},
			"max_replicas": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of replicas",
				Required:            true,
			},
			"vm_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the VM",
				Optional:            true,
			},
			"annotations": schema.MapAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Annotations for the deployment",
				Optional:            true,
			},
			"auth_username": schema.StringAttribute{
				MarkdownDescription: "The authentication username",
				Required:            true,
			},
			"auth_password": schema.StringAttribute{
				MarkdownDescription: "The authentication password",
				Required:            true,
			},
			"deployment_url": schema.StringAttribute{
				MarkdownDescription: "URL used to access successfull deployment",
				Computed:            true,
			},
		},
	}
}

func (r *deploymentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *deploymentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DeploymentCreateRequest
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nativePlan := convertDeploymentToNativePlan(plan)

	// Call Client's CreateDeployment method
	deployment, err := r.client.CreateDeployment(nativePlan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating deployment",
			"Could not create deployment, unexpected error: "+err.Error(),
		)
		return
	}

	state := convertToDeploymentTerraformState(deployment)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
func (r *deploymentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DeploymentTerraformState

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call Client's GetDeploymentByID method
	deployment, err := r.client.GetDeploymentByID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read deployment, got error: %s", err))
		return
	}

	newState := convertToDeploymentTerraformState(deployment)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
func (r *deploymentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DeploymentCreateRequest
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state DeploymentTerraformState
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nativePlan := convertDeploymentToNativePlan(plan)

	// Call Client's UpdateDeployment method
	deployment, err := r.client.UpdateDeployment(state.ID.ValueString(), nativePlan)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update deployment, got error: %s", err))
		return
	}

	newState := convertToDeploymentTerraformState(deployment)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
func (r *deploymentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DeploymentTerraformState
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call Client's DeleteDeployment method
	_, err := r.client.DeleteDeployment(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete deployment, got error: %s", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

type DeploymentTerraformState struct {
	ID                types.String            `tfsdk:"id"`
	Name              types.String            `tfsdk:"name"`
	ProjectID         types.String            `tfsdk:"project_id"`
	DeploymentImageID types.String            `tfsdk:"deployment_image_id"`
	ContainerImage    types.String            `tfsdk:"container_image"`
	MinReplicas       types.Int64             `tfsdk:"min_replicas"`
	MaxReplicas       types.Int64             `tfsdk:"max_replicas"`
	VMID              types.String            `tfsdk:"vm_id"`
	Annotations       map[string]types.String `tfsdk:"annotations"`
	AuthUsername      types.String            `tfsdk:"auth_username"`
	AuthPassword      types.String            `tfsdk:"auth_password"`
	URL               types.String            `tfsdk:"deployment_url"`
}

func convertToNativeMap(attributes map[string]types.String) map[string]string {
	result := make(map[string]string)
	for k, v := range attributes {
		result[k] = v.ValueString()
	}
	return result
}

func convertDeploymentToNativePlan(plan DeploymentCreateRequest) DeploymentCreateRequestNative {
	return DeploymentCreateRequestNative{
		Name:              plan.Name.ValueString(),
		ProjectID:         plan.ProjectID.ValueString(),
		DeploymentImageID: plan.DeploymentImageID.ValueString(),
		ContainerImage:    plan.ContainerImage.ValueString(),
		MinReplicas:       plan.MinReplicas.ValueInt64(),
		MaxReplicas:       plan.MaxReplicas.ValueInt64(),
		VMID:              plan.VMID.ValueString(),
		Annotations:       convertToNativeMap(plan.Annotations),
		AuthUsername:      plan.AuthUsername.ValueString(),
		AuthPassword:      plan.AuthPassword.ValueString(),
		URL:               plan.URL.ValueString(),
	}
}
func convertToDeploymentTerraformState(deployment *Deployment) DeploymentTerraformState {
	return DeploymentTerraformState{
		ID:                types.StringValue(deployment.Suffix), // Use Suffix as ID
		Name:              types.StringValue(deployment.Name),
		ProjectID:         types.StringValue(deployment.ProjectID),
		DeploymentImageID: types.StringNull(), // Not used, so set to null
		ContainerImage:    types.StringValue(deployment.ImageURL),
		MinReplicas:       types.Int64Value(1), // Assuming default value for MinReplicas
		MaxReplicas:       types.Int64Value(deployment.Replicas),
		VMID:              types.StringValue(deployment.MachineType),
		Annotations:       convertToTypesStringMap(deployment.Annotations),
		AuthUsername:      types.StringValue(deployment.AuthUsername),
		AuthPassword:      types.StringValue(deployment.AuthPassword),
		URL:               types.StringValue(deployment.Endpoint),
	}
}
