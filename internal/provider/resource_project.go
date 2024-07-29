package provider

// import (
// 	"context"
// 	"fmt"

// 	"github.com/hashicorp/terraform-plugin-framework/resource"
// 	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
// 	"github.com/hashicorp/terraform-plugin-framework/types"
// )

// type projectResource struct {
// 	client *Client
// }

// func NewProjectResource() resource.Resource {
// 	return &projectResource{}
// }

// func (r *projectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
// 	resp.TypeName = "theta_project"
// }

// func (r *projectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
// 	resp.Schema = schema.Schema{
// 		Attributes: map[string]schema.Attribute{
// 			"id": schema.StringAttribute{
// 				Computed: true,
// 			},
// 			"name": schema.StringAttribute{
// 				Required: true,
// 			},
// 			"description": schema.StringAttribute{
// 				Optional: true,
// 			},
// 		},
// 	}
// }

// func (r *projectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
// 	if req.ProviderData == nil {
// 		return
// 	}

// 	client, ok := req.ProviderData.(*Client)
// 	if !ok {
// 		resp.Diagnostics.AddError(
// 			"Unexpected Resource Configure Type",
// 			fmt.Sprintf("Expected *Client, got: %T", req.ProviderData),
// 		)
// 		return
// 	}

// 	r.client = client
// }

// func (r *projectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
// 	var data Project

// 	diags := req.Plan.Get(ctx, &data)
// 	resp.Diagnostics.Append(diags...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	projectID, err := r.client.CreateProject(data.Name.ValueString())
// 	if err != nil {
// 		resp.Diagnostics.AddError(
// 			"Error Creating Project",
// 			"Could not create project, unexpected error: "+err.Error(),
// 		)
// 		return
// 	}

// 	data.ID = types.StringValue(projectID)
// 	diags = resp.State.Set(ctx, &data)
// 	resp.Diagnostics.Append(diags...)
// }

// func (r *projectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
// 	var data Project

// 	diags := req.State.Get(ctx, &data)
// 	resp.Diagnostics.Append(diags...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	project, err := r.client.GetProject(data.ID.ValueString())
// 	if err != nil {
// 		resp.Diagnostics.AddError(
// 			"Error Reading Project",
// 			"Could not read project, unexpected error: "+err.Error(),
// 		)
// 		return
// 	}

// 	data = *project
// 	diags = resp.State.Set(ctx, &data)
// 	resp.Diagnostics.Append(diags...)
// }

// func (r *projectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
// 	var data Project

// 	diags := req.Plan.Get(ctx, &data)
// 	resp.Diagnostics.Append(diags...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	project, err := r.client.UpdateProject(data.ID.ValueString(), &data)
// 	if err != nil {
// 		resp.Diagnostics.AddError(
// 			"Error Updating Project",
// 			"Could not update project, unexpected error: "+err.Error(),
// 		)
// 		return
// 	}

// 	data = *project
// 	diags = resp.State.Set(ctx, &data)
// 	resp.Diagnostics.Append(diags...)
// }

// func (r *projectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
// 	var data Project

// 	diags := req.State.Get(ctx, &data)
// 	resp.Diagnostics.Append(diags...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	err := r.client.DeleteProject(data.ID.ValueString())
// 	if err != nil {
// 		resp.Diagnostics.AddError(
// 			"Error Deleting Project",
// 			"Could not delete project, unexpected error: "+err.Error(),
// 		)
// 		return
// 	}
// }
