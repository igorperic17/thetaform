package provider

import (
	"context"
	"fmt"
	"net/http"
	"time"

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

const (
	maxRetries = 60
	retryDelay = 10 * time.Second
)

type DeploymentResource struct {
	client *Client
}

func NewDeployment() resource.Resource {
	return &DeploymentResource{}
}

func (r *DeploymentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deployment"
}

func (r *DeploymentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Theta Deployment",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Deployment name",
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

func (r *DeploymentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

type DeploymentResourceModel struct {
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

func (r *DeploymentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DeploymentResourceModel

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

	deployment := &Deployment{
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

	createdDeployment, err := r.client.CreateDeployment(deployment)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create deployment, got error: %s", err))
		return
	}

	data.ID = types.StringValue(createdDeployment.ID)
	data.Suffix = types.StringValue(createdDeployment.Suffix)
	data.URL = types.StringValue(createdDeployment.URL)

	tflog.Trace(ctx, "created a resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Polling mechanism to wait for the deployment URL to become available
	for i := 0; i < maxRetries; i++ {
		resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}

		resp, err := http.Get(data.URL.ValueString())
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}
		time.Sleep(retryDelay)
	}
}

func (r *DeploymentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DeploymentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deployment, err := r.client.GetDeployment(data.Suffix.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read deployment, got error: %s", err))
		return
	}

	annotations, diags := convertStringMapToMap(deployment.Annotations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	envVars, diags := convertStringMapToMap(deployment.EnvVars)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue(deployment.ID)
	data.Name = types.StringValue(deployment.Name)
	data.ProjectID = types.StringValue(deployment.ProjectID)
	data.DeploymentImageID = types.StringValue(deployment.DeploymentImageID)
	data.ContainerImage = types.StringValue(deployment.ContainerImage)
	data.MinReplicas = types.Int64Value(int64(deployment.MinReplicas))
	data.MaxReplicas = types.Int64Value(int64(deployment.MaxReplicas))
	data.VMID = types.StringValue(deployment.VMID)
	data.Annotations = annotations
	data.EnvVars = envVars
	data.Suffix = types.StringValue(deployment.Suffix)
	data.URL = types.StringValue(deployment.URL)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeploymentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DeploymentResourceModel

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

	deployment := &Deployment{
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

	_, err = r.client.UpdateDeployment(data.Suffix.ValueString(), deployment)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update deployment, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeploymentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DeploymentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDeployment(data.Suffix.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete deployment, got error: %s", err))
		return
	}
}

func (r *DeploymentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
