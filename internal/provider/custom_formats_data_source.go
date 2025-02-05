package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/readarr-go/readarr"
	"github.com/devopsarr/terraform-provider-readarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const customFormatsDataSourceName = "custom_formats"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CustomFormatsDataSource{}

func NewCustomFormatsDataSource() datasource.DataSource {
	return &CustomFormatsDataSource{}
}

// CustomFormatsDataSource defines the custom formats implementation.
type CustomFormatsDataSource struct {
	client *readarr.APIClient
}

// CustomFormats describes the custom formats data model.
type CustomFormats struct {
	CustomFormats types.Set    `tfsdk:"custom_formats"`
	ID            types.String `tfsdk:"id"`
}

func (d *CustomFormatsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + customFormatsDataSourceName
}

func (d *CustomFormatsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Profiles -->List all available [Custom Formats](../resources/custom_format).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"custom_formats": schema.SetNestedAttribute{
				MarkdownDescription: "Download Client list..",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"include_custom_format_when_renaming": schema.BoolAttribute{
							MarkdownDescription: "Include custom format when renaming flag.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Custom Format name.",
							Computed:            true,
						},
						"id": schema.Int64Attribute{
							MarkdownDescription: "Custom Format ID.",
							Computed:            true,
						},
						"specifications": schema.SetNestedAttribute{
							MarkdownDescription: "Specifications.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"negate": schema.BoolAttribute{
										MarkdownDescription: "Negate flag.",
										Computed:            true,
									},
									"required": schema.BoolAttribute{
										MarkdownDescription: "Computed flag.",
										Computed:            true,
									},
									"name": schema.StringAttribute{
										MarkdownDescription: "Specification name.",
										Computed:            true,
									},
									"implementation": schema.StringAttribute{
										MarkdownDescription: "Implementation.",
										Computed:            true,
									},
									// Field values
									"value": schema.StringAttribute{
										MarkdownDescription: "Value.",
										Computed:            true,
									},
									"min": schema.Int64Attribute{
										MarkdownDescription: "Min.",
										Computed:            true,
									},
									"max": schema.Int64Attribute{
										MarkdownDescription: "Max.",
										Computed:            true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *CustomFormatsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *CustomFormatsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get custom formats current value
	response, _, err := d.client.CustomFormatAPI.ListCustomFormat(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.List, customFormatsDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+customFormatsDataSourceName)

	pointerResponse := make([]*readarr.CustomFormatResource, len(response))
	for i := range response {
		pointerResponse[i] = &response[i]
	}

	formats := make([]CustomFormat, len(pointerResponse))
	for i, p := range pointerResponse {
		formats[i].write(ctx, p, &resp.Diagnostics)
	}

	formatList, diags := types.SetValueFrom(ctx, CustomFormat{}.getType(), formats)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, CustomFormats{CustomFormats: formatList, ID: types.StringValue(strconv.Itoa(len(response)))})...)
}
