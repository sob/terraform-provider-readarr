package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/readarr-go/readarr"
	"github.com/devopsarr/terraform-provider-readarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	indexerNewznabResourceName   = "indexer_newznab"
	indexerNewznabImplementation = "Newznab"
	indexerNewznabConfigContract = "NewznabSettings"
	indexerNewznabProtocol       = "usenet"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &IndexerNewznabResource{}
	_ resource.ResourceWithImportState = &IndexerNewznabResource{}
)

func NewIndexerNewznabResource() resource.Resource {
	return &IndexerNewznabResource{}
}

// IndexerNewznabResource defines the Newznab indexer implementation.
type IndexerNewznabResource struct {
	client *readarr.APIClient
}

// IndexerNewznab describes the Newznab indexer data model.
type IndexerNewznab struct {
	Tags                    types.Set    `tfsdk:"tags"`
	Categories              types.Set    `tfsdk:"categories"`
	AdditionalParameters    types.String `tfsdk:"additional_parameters"`
	BaseURL                 types.String `tfsdk:"base_url"`
	APIPath                 types.String `tfsdk:"api_path"`
	APIKey                  types.String `tfsdk:"api_key"`
	Name                    types.String `tfsdk:"name"`
	EarlyReleaseLimit       types.Int64  `tfsdk:"early_release_limit"`
	ID                      types.Int64  `tfsdk:"id"`
	Priority                types.Int64  `tfsdk:"priority"`
	EnableRss               types.Bool   `tfsdk:"enable_rss"`
	EnableInteractiveSearch types.Bool   `tfsdk:"enable_interactive_search"`
	EnableAutomaticSearch   types.Bool   `tfsdk:"enable_automatic_search"`
}

func (i IndexerNewznab) toIndexer() *Indexer {
	return &Indexer{
		EnableAutomaticSearch:   i.EnableAutomaticSearch,
		EnableInteractiveSearch: i.EnableInteractiveSearch,
		EnableRss:               i.EnableRss,
		EarlyReleaseLimit:       i.EarlyReleaseLimit,
		Priority:                i.Priority,
		ID:                      i.ID,
		Name:                    i.Name,
		AdditionalParameters:    i.AdditionalParameters,
		APIKey:                  i.APIKey,
		APIPath:                 i.APIKey,
		BaseURL:                 i.BaseURL,
		Categories:              i.Categories,
		Tags:                    i.Tags,
		Implementation:          types.StringValue(indexerNewznabImplementation),
		ConfigContract:          types.StringValue(indexerNewznabConfigContract),
		Protocol:                types.StringValue(indexerNewznabProtocol),
	}
}

func (i *IndexerNewznab) fromIndexer(indexer *Indexer) {
	i.EnableAutomaticSearch = indexer.EnableAutomaticSearch
	i.EnableInteractiveSearch = indexer.EnableInteractiveSearch
	i.EnableRss = indexer.EnableRss
	i.EarlyReleaseLimit = indexer.EarlyReleaseLimit
	i.Priority = indexer.Priority
	i.ID = indexer.ID
	i.Name = indexer.Name
	i.AdditionalParameters = indexer.AdditionalParameters
	i.APIKey = indexer.APIKey
	i.APIPath = indexer.APIPath
	i.BaseURL = indexer.BaseURL
	i.Categories = indexer.Categories
	i.Tags = indexer.Tags
}

func (r *IndexerNewznabResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerNewznabResourceName
}

func (r *IndexerNewznabResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Indexers -->Indexer Newznab resource.\nFor more information refer to [Indexer](https://wiki.servarr.com/readarr/settings#indexers) and [Newznab](https://wiki.servarr.com/readarr/supported#newznab).",
		Attributes: map[string]schema.Attribute{
			"enable_automatic_search": schema.BoolAttribute{
				MarkdownDescription: "Enable automatic search flag.",
				Optional:            true,
				Computed:            true,
			},
			"enable_interactive_search": schema.BoolAttribute{
				MarkdownDescription: "Enable interactive search flag.",
				Optional:            true,
				Computed:            true,
			},
			"enable_rss": schema.BoolAttribute{
				MarkdownDescription: "Enable RSS flag.",
				Optional:            true,
				Computed:            true,
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Priority.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "IndexerNewznab name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "IndexerNewznab ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"early_release_limit": schema.Int64Attribute{
				MarkdownDescription: "Early release limit.",
				Optional:            true,
				Computed:            true,
			},
			"additional_parameters": schema.StringAttribute{
				MarkdownDescription: "Additional parameters.",
				Optional:            true,
				Computed:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Optional:            true,
				Computed:            true,
			},
			"api_path": schema.StringAttribute{
				MarkdownDescription: "API path.",
				Optional:            true,
				Computed:            true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Optional:            true,
				Computed:            true,
			},
			"categories": schema.SetAttribute{
				MarkdownDescription: "Series list.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (r *IndexerNewznabResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *IndexerNewznabResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var indexer *IndexerNewznab

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new IndexerNewznab
	request := indexer.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.IndexerAPI.CreateIndexer(ctx).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, indexerNewznabResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+indexerNewznabResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	indexer.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerNewznabResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var indexer *IndexerNewznab

	resp.Diagnostics.Append(req.State.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get IndexerNewznab current value
	response, _, err := r.client.IndexerAPI.GetIndexerById(ctx, int32(indexer.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerNewznabResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerNewznabResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	indexer.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerNewznabResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var indexer *IndexerNewznab

	resp.Diagnostics.Append(req.Plan.Get(ctx, &indexer)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update IndexerNewznab
	request := indexer.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.IndexerAPI.UpdateIndexer(ctx, strconv.Itoa(int(request.GetId()))).IndexerResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, indexerNewznabResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+indexerNewznabResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	indexer.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &indexer)...)
}

func (r *IndexerNewznabResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete IndexerNewznab current value
	_, err := r.client.IndexerAPI.DeleteIndexer(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, indexerNewznabResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+indexerNewznabResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *IndexerNewznabResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+indexerNewznabResourceName+": "+req.ID)
}

func (i *IndexerNewznab) write(ctx context.Context, indexer *readarr.IndexerResource, diags *diag.Diagnostics) {
	genericIndexer := i.toIndexer()
	genericIndexer.write(ctx, indexer, diags)
	i.fromIndexer(genericIndexer)
}

func (i *IndexerNewznab) read(ctx context.Context, diags *diag.Diagnostics) *readarr.IndexerResource {
	return i.toIndexer().read(ctx, diags)
}
