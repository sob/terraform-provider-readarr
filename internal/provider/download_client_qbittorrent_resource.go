package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/readarr-go/readarr"
	"github.com/devopsarr/terraform-provider-readarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	downloadClientQbittorrentResourceName   = "download_client_qbittorrent"
	downloadClientQbittorrentImplementation = "QBittorrent"
	downloadClientQbittorrentConfigContract = "QBittorrentSettings"
	downloadClientQbittorrentProtocol       = "torrent"
)

var downloadClientQbittorrentInitialStates = []int64{0, 1, 2}

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DownloadClientQbittorrentResource{}
	_ resource.ResourceWithImportState = &DownloadClientQbittorrentResource{}
)

func NewDownloadClientQbittorrentResource() resource.Resource {
	return &DownloadClientQbittorrentResource{}
}

// DownloadClientQbittorrentResource defines the download client implementation.
type DownloadClientQbittorrentResource struct {
	client *readarr.APIClient
}

// DownloadClientQbittorrent describes the download client data model.
type DownloadClientQbittorrent struct {
	Tags                     types.Set    `tfsdk:"tags"`
	MusicImportedCategory    types.String `tfsdk:"book_imported_category"`
	Name                     types.String `tfsdk:"name"`
	Host                     types.String `tfsdk:"host"`
	URLBase                  types.String `tfsdk:"url_base"`
	Username                 types.String `tfsdk:"username"`
	Password                 types.String `tfsdk:"password"`
	MusicCategory            types.String `tfsdk:"book_category"`
	RecentTVPriority         types.Int64  `tfsdk:"recent_book_priority"`
	Priority                 types.Int64  `tfsdk:"priority"`
	Port                     types.Int64  `tfsdk:"port"`
	ID                       types.Int64  `tfsdk:"id"`
	OlderTVPriority          types.Int64  `tfsdk:"older_book_priority"`
	InitialState             types.Int64  `tfsdk:"initial_state"`
	UseSsl                   types.Bool   `tfsdk:"use_ssl"`
	Enable                   types.Bool   `tfsdk:"enable"`
	FirstAndLast             types.Bool   `tfsdk:"first_and_last"`
	SequentialOrder          types.Bool   `tfsdk:"sequential_order"`
	RemoveFailedDownloads    types.Bool   `tfsdk:"remove_failed_downloads"`
	RemoveCompletedDownloads types.Bool   `tfsdk:"remove_completed_downloads"`
}

func (d DownloadClientQbittorrent) toDownloadClient() *DownloadClient {
	return &DownloadClient{
		Tags:                     d.Tags,
		Name:                     d.Name,
		Host:                     d.Host,
		URLBase:                  d.URLBase,
		Username:                 d.Username,
		Password:                 d.Password,
		MusicCategory:            d.MusicCategory,
		RecentTVPriority:         d.RecentTVPriority,
		OlderTVPriority:          d.OlderTVPriority,
		Priority:                 d.Priority,
		Port:                     d.Port,
		ID:                       d.ID,
		MusicImportedCategory:    d.MusicImportedCategory,
		InitialState:             d.InitialState,
		UseSsl:                   d.UseSsl,
		Enable:                   d.Enable,
		FirstAndLast:             d.FirstAndLast,
		SequentialOrder:          d.SequentialOrder,
		RemoveFailedDownloads:    d.RemoveFailedDownloads,
		RemoveCompletedDownloads: d.RemoveCompletedDownloads,
		Implementation:           types.StringValue(downloadClientQbittorrentImplementation),
		ConfigContract:           types.StringValue(downloadClientQbittorrentConfigContract),
		Protocol:                 types.StringValue(downloadClientQbittorrentProtocol),
	}
}

func (d *DownloadClientQbittorrent) fromDownloadClient(client *DownloadClient) {
	d.Tags = client.Tags
	d.Name = client.Name
	d.Host = client.Host
	d.URLBase = client.URLBase
	d.Username = client.Username
	d.Password = client.Password
	d.MusicCategory = client.MusicCategory
	d.RecentTVPriority = client.RecentTVPriority
	d.OlderTVPriority = client.OlderTVPriority
	d.Priority = client.Priority
	d.Port = client.Port
	d.ID = client.ID
	d.MusicImportedCategory = client.MusicImportedCategory
	d.InitialState = client.InitialState
	d.UseSsl = client.UseSsl
	d.Enable = client.Enable
	d.FirstAndLast = client.FirstAndLast
	d.SequentialOrder = client.SequentialOrder
	d.RemoveFailedDownloads = client.RemoveFailedDownloads
	d.RemoveCompletedDownloads = client.RemoveCompletedDownloads
}

func (r *DownloadClientQbittorrentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientQbittorrentResourceName
}

func (r *DownloadClientQbittorrentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->Download Client qBittorrent resource.\nFor more information refer to [Download Client](https://wiki.servarr.com/readarr/settings#download-clients) and [qBittorrent](https://wiki.servarr.com/readarr/supported#qbittorrent).",
		Attributes: map[string]schema.Attribute{
			"enable": schema.BoolAttribute{
				MarkdownDescription: "Enable flag.",
				Optional:            true,
				Computed:            true,
			},
			"remove_completed_downloads": schema.BoolAttribute{
				MarkdownDescription: "Remove completed downloads flag.",
				Optional:            true,
				Computed:            true,
			},
			"remove_failed_downloads": schema.BoolAttribute{
				MarkdownDescription: "Remove failed downloads flag.",
				Optional:            true,
				Computed:            true,
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Priority.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Download Client name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Download Client ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"use_ssl": schema.BoolAttribute{
				MarkdownDescription: "Use SSL flag.",
				Optional:            true,
				Computed:            true,
			},
			"sequential_order": schema.BoolAttribute{
				MarkdownDescription: "Sequential order flag.",
				Optional:            true,
				Computed:            true,
			},
			"first_and_last": schema.BoolAttribute{
				MarkdownDescription: "First and last flag.",
				Optional:            true,
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port.",
				Optional:            true,
				Computed:            true,
			},
			"recent_book_priority": schema.Int64Attribute{
				MarkdownDescription: "Recent Music priority. `0` Last, `1` First.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
				},
			},
			"older_book_priority": schema.Int64Attribute{
				MarkdownDescription: "Older Music priority. `0` Last, `1` First.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(0, 1),
				},
			},
			"initial_state": schema.Int64Attribute{
				MarkdownDescription: "Initial state, with Stop support. `0` Start, `1` ForceStart, `2` Pause.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.OneOf(downloadClientQbittorrentInitialStates...),
				},
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "host.",
				Optional:            true,
				Computed:            true,
			},
			"url_base": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Optional:            true,
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Optional:            true,
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
			"book_category": schema.StringAttribute{
				MarkdownDescription: "Book category.",
				Optional:            true,
				Computed:            true,
			},
			"book_imported_category": schema.StringAttribute{
				MarkdownDescription: "Book imported category.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *DownloadClientQbittorrentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if client := helpers.ResourceConfigure(ctx, req, resp); client != nil {
		r.client = client
	}
}

func (r *DownloadClientQbittorrentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var client *DownloadClientQbittorrent

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new DownloadClientQbittorrent
	request := client.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.DownloadClientAPI.CreateDownloadClient(ctx).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, downloadClientQbittorrentResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+downloadClientQbittorrentResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientQbittorrentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var client DownloadClientQbittorrent

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get DownloadClientQbittorrent current value
	response, _, err := r.client.DownloadClientAPI.GetDownloadClientById(ctx, int32(client.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, downloadClientQbittorrentResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientQbittorrentResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	client.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientQbittorrentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var client *DownloadClientQbittorrent

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update DownloadClientQbittorrent
	request := client.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.DownloadClientAPI.UpdateDownloadClient(ctx, strconv.Itoa(int(request.GetId()))).DownloadClientResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, downloadClientQbittorrentResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+downloadClientQbittorrentResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	client.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientQbittorrentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete DownloadClientQbittorrent current value
	_, err := r.client.DownloadClientAPI.DeleteDownloadClient(ctx, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, downloadClientQbittorrentResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+downloadClientQbittorrentResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientQbittorrentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+downloadClientQbittorrentResourceName+": "+req.ID)
}

func (d *DownloadClientQbittorrent) write(ctx context.Context, downloadClient *readarr.DownloadClientResource, diags *diag.Diagnostics) {
	genericDownloadClient := d.toDownloadClient()
	genericDownloadClient.write(ctx, downloadClient, diags)
	d.fromDownloadClient(genericDownloadClient)
}

func (d *DownloadClientQbittorrent) read(ctx context.Context, diags *diag.Diagnostics) *readarr.DownloadClientResource {
	return d.toDownloadClient().read(ctx, diags)
}
