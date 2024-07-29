package provider

import (
	"context"
	"encoding/json"
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

type DeploymentTemplateResource struct {
	client *Client
}

func NewDeploymentTemplateResource() resource.Resource {
	return &DeploymentTemplateResource{}
}

func (r *DeploymentTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deployment_template"
}

func (r *DeploymentTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Theta Deployment Template",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Template name",
				Required:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "Project ID",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Template description",
				Optional:            true,
			},
			"container_image": schema.StringAttribute{
				MarkdownDescription: "Container Image",
				Required:            true,
			},
			"tags": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Tags",
				Optional:            true,
			},
			"container_port": schema.StringAttribute{
				MarkdownDescription: "Container Port",
				Required:            true,
			},
			"container_args": schema.StringAttribute{
				MarkdownDescription: "Container Arguments",
				Optional:            true,
			},
			"env_vars": schema.StringAttribute{
				MarkdownDescription: "Environment Variables",
				Optional:            true,
			},
			"icon_url": schema.StringAttribute{
				MarkdownDescription: "Icon URL",
				Optional:            true,
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

func (r *DeploymentTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

type DeploymentTemplateResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	ProjectID      types.String `tfsdk:"project_id"`
	Description    types.String `tfsdk:"description"`
	ContainerImage types.String `tfsdk:"container_image"`
	Tags           types.List   `tfsdk:"tags"`
	ContainerPort  types.String `tfsdk:"container_port"`
	ContainerArgs  types.String `tfsdk:"container_args"`
	EnvVars        types.String `tfsdk:"env_vars"`
	IconURL        types.String `tfsdk:"icon_url"`
}

func (r *DeploymentTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DeploymentTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	template := DeploymentTemplateRequest{
		Name:           data.Name.ValueString(),
		ProjectID:      data.ProjectID.ValueString(),
		Description:    data.Description.ValueString(),
		ContainerImage: data.ContainerImage.ValueString(),
		ContainerPort:  data.ContainerPort.ValueString(),
		ContainerArgs:  data.ContainerArgs.ValueString(),
		IconURL:        data.IconURL.ValueString(),
	}

	if !data.EnvVars.IsNull() {
		var envVars map[string]string
		err := json.Unmarshal([]byte(data.EnvVars.ValueString()), &envVars)
		if err != nil {
			resp.Diagnostics.AddError("EnvVars Unmarshal Error", fmt.Sprintf("Error unmarshaling env_vars: %s", err))
			return
		}
		template.EnvVars = envVars
	}

	if !data.Tags.IsNull() {
		var tags []string
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		template.Tags = tags
	}

	// Log the request payload
	requestPayload, err := json.Marshal(template)
	if err == nil {
		tflog.Debug(ctx, fmt.Sprintf("Create request payload: %s", requestPayload))
	}

	createdTemplate, err := r.client.CreateDeploymentTemplate(template)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create deployment template, got error: %s", err))
		return
	}

	data.ID = types.StringValue(createdTemplate.ID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeploymentTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DeploymentTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	template, err := r.client.GetDeploymentTemplateByID(data.ProjectID.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read deployment template, got error: %s", err))
		return
	}

	data.Name = types.StringValue(template.Name)
	data.ProjectID = types.StringValue(template.ProjectID)
	data.Description = types.StringValue(template.Description)
	data.ContainerImage = types.StringValue(template.ContainerImage)
	data.ContainerPort = types.StringValue(template.ContainerPort)
	data.ContainerArgs = types.StringValue(template.ContainerArgs)
	data.IconURL = types.StringValue(template.IconURL)

	if template.EnvVars != nil {
		envVarsJson, err := json.Marshal(template.EnvVars)
		if err != nil {
			resp.Diagnostics.AddError("EnvVars Marshal Error", fmt.Sprintf("Error marshaling env_vars: %s", err))
			return
		}
		data.EnvVars = types.StringValue(string(envVarsJson))
	} else {
		data.EnvVars = types.StringNull()
	}

	if len(template.Tags) > 0 {
		var tags []attr.Value
		for _, tag := range template.Tags {
			tags = append(tags, types.StringValue(tag))
		}
		data.Tags = types.ListValueMust(types.StringType, tags)
	} else {
		data.Tags = types.ListNull(types.StringType)
	}

	// Debugging statements
	fmt.Printf("Read method debug: Name: %s, ProjectID: %s, Description: %s, ContainerImage: %s, ContainerPort: %s\n",
		data.Name.ValueString(), data.ProjectID.ValueString(), data.Description.ValueString(), data.ContainerImage.ValueString(), data.ContainerPort.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeploymentTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DeploymentTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	template := DeploymentTemplateRequest{
		Name:           data.Name.ValueString(),
		ProjectID:      data.ProjectID.ValueString(),
		Description:    data.Description.ValueString(),
		ContainerImage: data.ContainerImage.ValueString(),
		ContainerPort:  data.ContainerPort.ValueString(),
		ContainerArgs:  data.ContainerArgs.ValueString(),
		IconURL:        data.IconURL.ValueString(),
	}

	if !data.EnvVars.IsNull() {
		var envVars map[string]string
		err := json.Unmarshal([]byte(data.EnvVars.ValueString()), &envVars)
		if err != nil {
			resp.Diagnostics.AddError("EnvVars Unmarshal Error", fmt.Sprintf("Error unmarshaling env_vars: %s", err))
			return
		}
		template.EnvVars = envVars
	} else {
		template.EnvVars = nil
	}

	if !data.Tags.IsNull() {
		var tags []string
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		template.Tags = tags
	}

	// Log the request payload
	requestPayload, err := json.Marshal(template)
	if err == nil {
		tflog.Debug(ctx, fmt.Sprintf("Update request payload: %s", requestPayload))
	}

	_, err = r.client.UpdateDeploymentTemplate(data.ID.ValueString(), template)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update deployment template, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeploymentTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DeploymentTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteDeploymentTemplate(data.ID.ValueString(), data.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete deployment template, got error: %s", err))
		return
	}
}

func (r *DeploymentTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
