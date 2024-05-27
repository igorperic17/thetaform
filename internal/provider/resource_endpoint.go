package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &EndpointResource{}
var _ resource.ResourceWithImportState = &EndpointResource{}

func NewEndpoint() resource.Resource {
	return &EndpointResource{}
}

type EndpointResource struct {
	client *Client
}

type EndpointResourceModel struct {
	Name              types.String `tfsdk:"name"`
	ProjectID         types.String `tfsdk:"project_id"`
	DeploymentImageID types.String `tfsdk:"deployment_image_id"`
	ContainerImage    types.String `tfsdk:"container_image"`
	MinReplicas       types.Int64  `tfsdk:"min_replicas"`
	MaxReplicas       types.Int64  `tfsdk:"max_replicas"`
	VMID              types.String `tfsdk:"vm_id"`
	Annotations       types.Map    `tfsdk:"annotations"`
	EnvVars           types.Map    `tfsdk:"env_vars"`
	Id                types.String `tfsdk:"id"`
}

func (r *EndpointResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoint"
}

func (r *EndpointResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Theta Endpoint Deployment",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Endpoint name",
				Required:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "Project ID",
				Required:            true,
			},
			"deployment_image_id": schema.StringAttribute{
				MarkdownDescription: "Deployment Image ID",
				Required:            true,
			},
			"container_image": schema.StringAttribute{
				MarkdownDescription: "Container Image",
				Required:            true,
			},
			"min_replicas": schema.Int64Attribute{
				MarkdownDescription: "Minimum Replicas",
				Required:            true,
			},
			"max_replicas": schema.Int64Attribute{
				MarkdownDescription: "Maximum Replicas",
				Required:            true,
			},
			"vm_id": schema.StringAttribute{
				MarkdownDescription: "VM ID",
				Required:            true,
			},
			"annotations": schema.MapAttribute{
				MarkdownDescription: "Annotations",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"env_vars": schema.MapAttribute{
				MarkdownDescription: "Environment Variables",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *EndpointResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *EndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EndpointResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := &Endpoint{
		Name:              data.Name.ValueString(),
		ProjectID:         data.ProjectID.ValueString(),
		DeploymentImageID: data.DeploymentImageID.ValueString(),
		ContainerImage:    data.ContainerImage.ValueString(),
		MinReplicas:       int(data.MinReplicas.ValueInt64()),
		MaxReplicas:       int(data.MaxReplicas.ValueInt64()),
		VMID:              data.VMID.ValueString(),
		Annotations:       convertMapToStringMap(data.Annotations.Elements()),
		EnvVars:           convertMapToStringMap(data.EnvVars.Elements()),
	}

	createdEndpoint, err := r.client.CreateEndpoint(endpoint)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Endpoint",
			"Could not create endpoint, unexpected error: "+err.Error(),
		)
		return
	}

	data.Id = types.StringValue(createdEndpoint.ID)
	tflog.Trace(ctx, "created an endpoint resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func convertMapToStringMap(input map[string]attr.Value) map[string]string {
	output := make(map[string]string)
	for key, value := range input {
		output[key] = value.(types.String).ValueString()
	}
	return output
}

func convertStringMapToMap(input map[string]string) types.Map {
	elements := make(map[string]attr.Value)
	for key, value := range input {
		elements[key] = types.StringValue(value)
	}
	return types.MapValueMust(types.StringType, elements)
}

func (r *EndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data EndpointResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint, err := r.client.GetEndpoint(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Endpoint",
			"Could not read endpoint, unexpected error: "+err.Error(),
		)
		return
	}

	data.Name = types.StringValue(endpoint.Name)
	data.ProjectID = types.StringValue(endpoint.ProjectID)
	data.DeploymentImageID = types.StringValue(endpoint.DeploymentImageID)
	data.ContainerImage = types.StringValue(endpoint.ContainerImage)
	data.MinReplicas = types.Int64Value(int64(endpoint.MinReplicas))
	data.MaxReplicas = types.Int64Value(int64(endpoint.MaxReplicas))
	data.VMID = types.StringValue(endpoint.VMID)
	data.Annotations = convertStringMapToMap(endpoint.Annotations)
	data.EnvVars = convertStringMapToMap(endpoint.EnvVars)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data EndpointResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := &Endpoint{
		Name:              data.Name.ValueString(),
		ProjectID:         data.ProjectID.ValueString(),
		DeploymentImageID: data.DeploymentImageID.ValueString(),
		ContainerImage:    data.ContainerImage.ValueString(),
		MinReplicas:       int(data.MinReplicas.ValueInt64()),
		MaxReplicas:       int(data.MaxReplicas.ValueInt64()),
		VMID:              data.VMID.ValueString(),
		Annotations:       convertMapToStringMap(data.Annotations.Elements()),
		EnvVars:           convertMapToStringMap(data.EnvVars.Elements()),
	}

	_, err := r.client.UpdateEndpoint(data.Id.ValueString(), endpoint)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Endpoint",
			"Could not update endpoint, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data EndpointResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteEndpoint(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Endpoint",
			"Could not delete endpoint, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *EndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
