package provider

import (
	"context"

	"github.com/devopsarr/readarr-go/readarr"
	"github.com/devopsarr/terraform-provider-readarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const tagDataSourceName = "tag"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &TagDataSource{}

func NewTagDataSource() datasource.DataSource {
	return &TagDataSource{}
}

// TagDataSource defines the tag implementation.
type TagDataSource struct {
	client *readarr.APIClient
}

func (d *TagDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + tagDataSourceName
}

func (d *TagDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "<!-- subcategory:Tags -->Single [Tag](../resources/tag).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Tag ID.",
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Tag label.",
				Required:            true,
			},
		},
	}
}

func (d *TagDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *TagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *Tag

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get tags current value
	response, _, err := d.client.TagAPI.ListTag(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, tagDataSourceName, err))

		return
	}

	tags := make([]*readarr.TagResource, len(response))
	for i := range response {
		tags[i] = &response[i]
	}
	data.find(data.Label.ValueString(), tags, &resp.Diagnostics)
	tflog.Trace(ctx, "read "+tagDataSourceName)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (t *Tag) find(label string, tags []*readarr.TagResource, diags *diag.Diagnostics) {
	for _, tag := range tags {
		if tag.GetLabel() == label {
			t.write(tag)

			return
		}
	}

	diags.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(tagDataSourceName, "label", label))
}
