package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/readarr-go/readarr"
	"github.com/devopsarr/terraform-provider-readarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	notificationCustomScriptResourceName   = "notification_custom_script"
	notificationCustomScriptImplementation = "CustomScript"
	notificationCustomScriptConfigContract = "CustomScriptSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationCustomScriptResource{}
	_ resource.ResourceWithImportState = &NotificationCustomScriptResource{}
)

func NewNotificationCustomScriptResource() resource.Resource {
	return &NotificationCustomScriptResource{}
}

// NotificationCustomScriptResource defines the notification implementation.
type NotificationCustomScriptResource struct {
	client *readarr.APIClient
}

// NotificationCustomScript describes the notification data model.
type NotificationCustomScript struct {
	Tags                       types.Set    `tfsdk:"tags"`
	Arguments                  types.String `tfsdk:"arguments"`
	Path                       types.String `tfsdk:"path"`
	Name                       types.String `tfsdk:"name"`
	ID                         types.Int64  `tfsdk:"id"`
	OnGrab                     types.Bool   `tfsdk:"on_grab"`
	IncludeHealthWarnings      types.Bool   `tfsdk:"include_health_warnings"`
	OnHealthIssue              types.Bool   `tfsdk:"on_health_issue"`
	OnRename                   types.Bool   `tfsdk:"on_rename"`
	OnUpgrade                  types.Bool   `tfsdk:"on_upgrade"`
	OnReleaseImport            types.Bool   `tfsdk:"on_release_import"`
	OnAuthorDelete             types.Bool   `tfsdk:"on_author_delete"`
	OnBookDelete               types.Bool   `tfsdk:"on_book_delete"`
	OnBookFileDelete           types.Bool   `tfsdk:"on_book_file_delete"`
	OnBookFileDeleteForUpgrade types.Bool   `tfsdk:"on_book_file_delete_for_upgrade"`
	OnDownloadFailure          types.Bool   `tfsdk:"on_download_failure"`
	OnImportFailure            types.Bool   `tfsdk:"on_import_failure"`
	OnBookRetag                types.Bool   `tfsdk:"on_book_retag"`
}

func (n NotificationCustomScript) toNotification() *Notification {
	return &Notification{
		Tags:                       n.Tags,
		Path:                       n.Path,
		Arguments:                  n.Arguments,
		Name:                       n.Name,
		ID:                         n.ID,
		OnGrab:                     n.OnGrab,
		OnReleaseImport:            n.OnReleaseImport,
		OnAuthorDelete:             n.OnAuthorDelete,
		IncludeHealthWarnings:      n.IncludeHealthWarnings,
		OnHealthIssue:              n.OnHealthIssue,
		OnBookDelete:               n.OnBookDelete,
		OnBookFileDelete:           n.OnBookFileDelete,
		OnRename:                   n.OnRename,
		OnUpgrade:                  n.OnUpgrade,
		OnBookFileDeleteForUpgrade: n.OnBookFileDeleteForUpgrade,
		OnDownloadFailure:          n.OnDownloadFailure,
		OnImportFailure:            n.OnImportFailure,
		OnBookRetag:                n.OnBookRetag,
	}
}

func (n *NotificationCustomScript) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.Path = notification.Path
	n.Arguments = notification.Arguments
	n.Name = notification.Name
	n.ID = notification.ID
	n.OnGrab = notification.OnGrab
	n.OnBookFileDeleteForUpgrade = notification.OnBookFileDeleteForUpgrade
	n.OnBookFileDelete = notification.OnBookFileDelete
	n.IncludeHealthWarnings = notification.IncludeHealthWarnings
	n.OnHealthIssue = notification.OnHealthIssue
	n.OnAuthorDelete = notification.OnAuthorDelete
	n.OnBookDelete = notification.OnBookDelete
	n.OnRename = notification.OnRename
	n.OnUpgrade = notification.OnUpgrade
	n.OnDownloadFailure = notification.OnDownloadFailure
	n.OnBookRetag = notification.OnBookRetag
	n.OnImportFailure = notification.OnImportFailure
	n.OnReleaseImport = notification.OnReleaseImport
}

func (r *NotificationCustomScriptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationCustomScriptResourceName
}

func (r *NotificationCustomScriptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Custom Script resource.\nFor more information refer to [Notification](https://wiki.servarr.com/readarr/settings#connect) and [Custom Script](https://wiki.servarr.com/readarr/supported#customscript).",
		Attributes: map[string]schema.Attribute{
			"on_grab": schema.BoolAttribute{
				MarkdownDescription: "On grab flag.",
				Required:            true,
			},
			"on_download_failure": schema.BoolAttribute{
				MarkdownDescription: "On download failure flag.",
				Required:            true,
			},
			"on_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On upgrade flag.",
				Required:            true,
			},
			"on_rename": schema.BoolAttribute{
				MarkdownDescription: "On rename flag.",
				Required:            true,
			},
			"on_author_delete": schema.BoolAttribute{
				MarkdownDescription: "On author deleted flag.",
				Required:            true,
			},
			"on_book_delete": schema.BoolAttribute{
				MarkdownDescription: "On book delete flag.",
				Required:            true,
			},
			"on_book_file_delete": schema.BoolAttribute{
				MarkdownDescription: "On book file delete flag.",
				Required:            true,
			},
			"on_book_file_delete_for_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On book file delete for upgrade flag.",
				Required:            true,
			},
			"on_health_issue": schema.BoolAttribute{
				MarkdownDescription: "On health issue flag.",
				Required:            true,
			},
			"on_import_failure": schema.BoolAttribute{
				MarkdownDescription: "On import failure flag.",
				Required:            true,
			},
			"on_book_retag": schema.BoolAttribute{
				MarkdownDescription: "On book retag flag.",
				Required:            true,
			},
			"on_release_import": schema.BoolAttribute{
				MarkdownDescription: "On release import flag.",
				Required:            true,
			},
			"include_health_warnings": schema.BoolAttribute{
				MarkdownDescription: "Include health warnings.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "NotificationCustomScript name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Notification ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"arguments": schema.StringAttribute{
				MarkdownDescription: "Arguments.",
				Optional:            true,
				Computed:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "Path.",
				Required:            true,
			},
		},
	}
}

func (r *NotificationCustomScriptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*readarr.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *readarr.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *NotificationCustomScriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationCustomScript

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationCustomScript
	request := notification.read(ctx)

	response, _, err := r.client.NotificationApi.CreateNotification(ctx).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", notificationCustomScriptResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationCustomScriptResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationCustomScriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationCustomScript

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationCustomScript current value
	response, _, err := r.client.NotificationApi.GetNotificationById(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationCustomScriptResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationCustomScriptResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationCustomScriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationCustomScript

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationCustomScript
	request := notification.read(ctx)

	response, _, err := r.client.NotificationApi.UpdateNotification(ctx, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", notificationCustomScriptResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationCustomScriptResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationCustomScriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var notification *NotificationCustomScript

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationCustomScript current value
	_, err := r.client.NotificationApi.DeleteNotification(ctx, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationCustomScriptResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationCustomScriptResourceName+": "+strconv.Itoa(int(notification.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationCustomScriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+notificationCustomScriptResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (n *NotificationCustomScript) write(ctx context.Context, notification *readarr.NotificationResource) {
	genericNotification := Notification{
		OnGrab:          types.BoolValue(notification.GetOnGrab()),
		OnReleaseImport: types.BoolValue(notification.GetOnReleaseImport()),
		OnUpgrade:       types.BoolValue(notification.GetOnUpgrade()),
		OnRename:        types.BoolValue(notification.GetOnRename()),
		// OnAuthorDelete:             types.BoolValue(notification.GetOnAuthorDelete()),
		// OnBookDelete:               types.BoolValue(notification.GetOnBookDelete()),
		// OnBookFileDelete:           types.BoolValue(notification.GetOnBookFileDelete()),
		// OnBookFileDeleteForUpgrade: types.BoolValue(notification.GetOnBookFileDeleteForUpgrade()),
		OnHealthIssue:         types.BoolValue(notification.GetOnHealthIssue()),
		OnDownloadFailure:     types.BoolValue(notification.GetOnDownloadFailure()),
		OnImportFailure:       types.BoolValue(notification.GetOnImportFailure()),
		OnBookRetag:           types.BoolValue(notification.GetOnBookRetag()),
		IncludeHealthWarnings: types.BoolValue(notification.GetIncludeHealthWarnings()),
		ID:                    types.Int64Value(int64(notification.GetId())),
		Name:                  types.StringValue(notification.GetName()),
		Tags:                  types.SetValueMust(types.Int64Type, nil),
	}
	tfsdk.ValueFrom(ctx, notification.Tags, genericNotification.Tags.Type(ctx), &genericNotification.Tags)
	genericNotification.writeFields(ctx, notification.Fields)
	n.fromNotification(&genericNotification)
}

func (n *NotificationCustomScript) read(ctx context.Context) *readarr.NotificationResource {
	var tags []*int32

	tfsdk.ValueAs(ctx, n.Tags, &tags)

	notification := readarr.NewNotificationResource()
	notification.SetOnGrab(n.OnGrab.ValueBool())
	notification.SetOnReleaseImport(n.OnReleaseImport.ValueBool())
	notification.SetOnUpgrade(n.OnUpgrade.ValueBool())
	notification.SetOnRename(n.OnRename.ValueBool())
	// notification.SetOnAuthorDelete(n.OnAuthorDelete.ValueBool())
	// notification.SetOnBookDelete(n.OnBookDelete.ValueBool())
	// notification.SetOnBookFileDelete(n.OnBookFileDelete.ValueBool())
	// notification.SetOnBookFileDeleteForUpgrade(n.OnBookFileDeleteForUpgrade.ValueBool())
	notification.SetOnHealthIssue(n.OnHealthIssue.ValueBool())
	notification.SetOnDownloadFailure(n.OnDownloadFailure.ValueBool())
	notification.SetOnImportFailure(n.OnImportFailure.ValueBool())
	notification.SetOnBookRetag(n.OnBookRetag.ValueBool())
	notification.SetIncludeHealthWarnings(n.IncludeHealthWarnings.ValueBool())
	notification.SetId(int32(n.ID.ValueInt64()))
	notification.SetName(n.Name.ValueString())
	notification.SetConfigContract(notificationCustomScriptConfigContract)
	notification.SetImplementation(notificationCustomScriptImplementation)
	notification.SetTags(tags)
	notification.SetFields(n.toNotification().readFields(ctx))

	return notification
}
