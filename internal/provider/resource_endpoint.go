package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type EndpointResource struct {
	client *Client
}

func NewEndpoint() resource.Resource {
	return &EndpointResource{}
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
				MarkdownDescription: "Minimum number of replicas",
				Required:            true,
			},
			"max_replicas": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of replicas",
				Required:            true,
			},
			"vm_id": schema.StringAttribute{
				MarkdownDescription: "VM ID",
				Required:            true,
			},
			"annotations": schema.MapAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Annotations",
				Optional:            true,
			},
			"env_vars": schema.MapAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Environment Variables",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"suffix": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Suffix",
			},
			"url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Deployment URL",
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

type EndpointResourceModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	ProjectID         types.String `tfsdk:"project_id"`
	DeploymentImageID types.String `tfsdk:"deployment_image_id"`
	ContainerImage    types.String `tfsdk:"container_image"`
	MinReplicas       types.Int64  `tfsdk:"min_replicas"`
	MaxReplicas       types.Int64  `tfsdk:"max_replicas"`
	VMID              types.String `tfsdk:"vm_id"`
	Annotations       types.Map    `tfsdk:"annotations"`
	EnvVars           types.Map    `tfsdk:"env_vars"`
	Suffix            types.String `tfsdk:"suffix"`
	URL               types.String `tfsdk:"url"`
}

func convertMapToStringMap(input types.Map) (map[string]string, error) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	output := make(map[string]string)
	for key, value := range input.Elements() {
		if valueStr, ok := value.(types.String); ok {
			output[key] = valueStr.ValueString()
		} else {
			return nil, fmt.Errorf("expected types.String, got %T", value)
		}
	}
	return output, nil
}

func convertStringMapToMap(input map[string]string) (types.Map, diag.Diagnostics) {
	elements := make(map[string]attr.Value)
	for key, value := range input {
		elements[key] = types.StringValue(value)
	}
	return types.MapValue(types.StringType, elements)
}

func (r *EndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EndpointResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	annotations, err := convertMapToStringMap(data.Annotations)
	if err != nil {
		resp.Diagnostics.AddError("Attribute Conversion Error", fmt.Sprintf("Error converting annotations: %s", err))
		return
	}
	envVars, err := convertMapToStringMap(data.EnvVars)
	if err != nil {
		resp.Diagnostics.AddError("Attribute Conversion Error", fmt.Sprintf("Error converting environment variables: %s", err))
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
		Annotations:       annotations,
		EnvVars:           envVars,
	}

	createdEndpoint, err := r.client.CreateEndpoint(endpoint)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create endpoint, got error: %s", err))
		return
	}

	data.ID = types.StringValue(createdEndpoint.ID)
	data.Suffix = types.StringValue(createdEndpoint.Suffix)
	data.URL = types.StringValue(createdEndpoint.URL)

	tflog.Trace(ctx, "created a resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data EndpointResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint, err := r.client.GetEndpoint(data.Suffix.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read endpoint, got error: %s", err))
		return
	}

	annotations, diags := convertStringMapToMap(endpoint.Annotations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	envVars, diags := convertStringMapToMap(endpoint.EnvVars)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue(endpoint.ID)
	data.Name = types.StringValue(endpoint.Name)
	data.ProjectID = types.StringValue(endpoint.ProjectID)
	data.DeploymentImageID = types.StringValue(endpoint.DeploymentImageID)
	data.ContainerImage = types.StringValue(endpoint.ContainerImage)
	data.MinReplicas = types.Int64Value(int64(endpoint.MinReplicas))
	data.MaxReplicas = types.Int64Value(int64(endpoint.MaxReplicas))
	data.VMID = types.StringValue(endpoint.VMID)
	data.Annotations = annotations
	data.EnvVars = envVars
	data.Suffix = types.StringValue(endpoint.Suffix)
	data.URL = types.StringValue(endpoint.URL)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data EndpointResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	annotations, err := convertMapToStringMap(data.Annotations)
	if err != nil {
		resp.Diagnostics.AddError("Attribute Conversion Error", fmt.Sprintf("Error converting annotations: %s", err))
		return
	}
	envVars, err := convertMapToStringMap(data.EnvVars)
	if err != nil {
		resp.Diagnostics.AddError("Attribute Conversion Error", fmt.Sprintf("Error converting environment variables: %s", err))
		return
	}

	endpoint := &Endpoint{
		ID:                data.ID.ValueString(),
		Name:              data.Name.ValueString(),
		ProjectID:         data.ProjectID.ValueString(),
		DeploymentImageID: data.DeploymentImageID.ValueString(),
		ContainerImage:    data.ContainerImage.ValueString(),
		MinReplicas:       int(data.MinReplicas.ValueInt64()),
		MaxReplicas:       int(data.MaxReplicas.ValueInt64()),
		VMID:              data.VMID.ValueString(),
		Annotations:       annotations,
		EnvVars:           envVars,
	}

	_, err = r.client.UpdateEndpoint(data.Suffix.ValueString(), endpoint)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update endpoint, got error: %s", err))
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

	err := r.client.DeleteEndpoint(data.Suffix.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete endpoint, got error: %s", err))
		return
	}
}

func (r *EndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
