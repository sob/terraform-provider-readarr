package provider

import (
	"context"

	"github.com/devopsarr/readarr-go/readarr"
	"github.com/devopsarr/terraform-provider-readarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const indexerDataSourceName = "indexer"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &IndexerDataSource{}

func NewIndexerDataSource() datasource.DataSource {
	return &IndexerDataSource{}
}

// IndexerDataSource defines the indexer implementation.
type IndexerDataSource struct {
	client *readarr.APIClient
}

func (d *IndexerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerDataSourceName
}

func (d *IndexerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Indexers -->Single [Indexer](../resources/indexer).",
		Attributes: map[string]schema.Attribute{
			"enable_automatic_search": schema.BoolAttribute{
				MarkdownDescription: "Enable automatic search flag.",
				Computed:            true,
			},
			"enable_interactive_search": schema.BoolAttribute{
				MarkdownDescription: "Enable interactive search flag.",
				Computed:            true,
			},
			"enable_rss": schema.BoolAttribute{
				MarkdownDescription: "Enable RSS flag.",
				Computed:            true,
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Priority.",
				Computed:            true,
			},
			"config_contract": schema.StringAttribute{
				MarkdownDescription: "Indexer configuration template.",
				Computed:            true,
			},
			"implementation": schema.StringAttribute{
				MarkdownDescription: "Indexer implementation name.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Indexer name.",
				Required:            true,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "Protocol. Valid values are 'usenet' and 'torrent'.",
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Indexer ID.",
				Computed:            true,
			},
			// Field values
			"allow_zero_size": schema.BoolAttribute{
				MarkdownDescription: "Allow zero size files.",
				Computed:            true,
			},
			"ranked_only": schema.BoolAttribute{
				MarkdownDescription: "Allow ranked only.",
				Computed:            true,
			},
			"delay": schema.Int64Attribute{
				MarkdownDescription: "Delay before grabbing.",
				Computed:            true,
			},
			"minimum_seeders": schema.Int64Attribute{
				MarkdownDescription: "Minimum seeders.",
				Computed:            true,
			},
			"early_release_limit": schema.Int64Attribute{
				MarkdownDescription: "Early release limit.",
				Computed:            true,
			},
			"seed_time": schema.Int64Attribute{
				MarkdownDescription: "Seed time.",
				Computed:            true,
			},
			"author_seed_time": schema.Int64Attribute{
				MarkdownDescription: "Author seed time.",
				Computed:            true,
			},
			"seed_ratio": schema.Float64Attribute{
				MarkdownDescription: "Seed ratio.",
				Computed:            true,
			},
			"additional_parameters": schema.StringAttribute{
				MarkdownDescription: "Additional parameters.",
				Computed:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Computed:            true,
			},
			"api_user": schema.StringAttribute{
				MarkdownDescription: "API User.",
				Computed:            true,
			},
			"api_path": schema.StringAttribute{
				MarkdownDescription: "API path.",
				Computed:            true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Computed:            true,
			},
			"captcha_token": schema.StringAttribute{
				MarkdownDescription: "Captcha token.",
				Computed:            true,
			},
			"cookie": schema.StringAttribute{
				MarkdownDescription: "Cookie.",
				Computed:            true,
			},
			"passkey": schema.StringAttribute{
				MarkdownDescription: "Passkey.",
				Computed:            true,
				Sensitive:           true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password.",
				Computed:            true,
				Sensitive:           true,
			},
			"categories": schema.SetAttribute{
				MarkdownDescription: "Series list.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (d *IndexerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *IndexerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *Indexer

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get indexer current value
	response, _, err := d.client.IndexerAPI.ListIndexer(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerDataSourceName, err))

		return
	}

	indexers := make([]*readarr.IndexerResource, len(response))
	for i := range response {
		indexers[i] = &response[i]
	}
	data.find(ctx, data.Name.ValueString(), indexers, &resp.Diagnostics)
	tflog.Trace(ctx, "read "+indexerDataSourceName)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (i *Indexer) find(ctx context.Context, name string, indexers []*readarr.IndexerResource, diags *diag.Diagnostics) {
	for _, indexer := range indexers {
		if indexer.GetName() == name {
			i.write(ctx, indexer, diags)

			return
		}
	}

	diags.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(indexerDataSourceName, "name", name))
}
