package dog

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ledongthuc/goterators"
	api "github.com/relaypro-open/dog_api_golang/api"
)

type (
	profileDataSource struct {
		p dogProvider
	}

	ProfileList []Profile

	Profile struct {
		ID      types.String `tfsdk:"id"`
		Name    types.String `tfsdk:"name"`
		Version types.String `tfsdk:"version"`
	}
)

var (
	_ datasource.DataSource = (*profileDataSource)(nil)
)

func NewProfileDataSource() datasource.DataSource {
	return &profileDataSource{}
}

func (*profileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_profile"
}

func (*profileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Profile data source",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Profile name",
				Optional:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Profile version",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Profile identifier",
				Computed:            true,
			},
		},
	}
}

func (d *profileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *dog.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.p.dog = client
}

type profileDataSourceData struct {
	ApiToken types.String `tfsdk:"api_token"`
	Id       types.String `tfsdk:"id"`
}

//type profileDataSource struct {
//	provider provider
//}

func (d *profileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state Profile
	var profileName string

	req.Config.GetAttribute(ctx, path.Root("name"), &profileName)

	res, statusCode, err := d.p.dog.GetProfiles(nil)
	if (statusCode < 200 || statusCode > 299) && statusCode != 404 {
		resp.Diagnostics.AddError("Client Unsuccesful", fmt.Sprintf("Status Code: %d", statusCode))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read profiles, got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}

	var filteredProfilesName []api.Profile
	if profileName != "" {
		filteredProfilesName = goterators.Filter(res, func(profile api.Profile) bool {
			//tflog.Debug(ctx, fmt.Sprintf("ZZZprofile.Name: '%s', profileName: '%s'", profile.Name,profileName))
			return profile.Name == profileName
		})
	} else {
		filteredProfilesName = res
	}

	filteredProfiles := filteredProfilesName

	if filteredProfiles == nil {
		resp.Diagnostics.AddError("Data Error", fmt.Sprintf("dog_profile data source returned no results."))
	}
	if len(filteredProfiles) > 1 {
		resp.Diagnostics.AddError("Data Error", fmt.Sprintf("dog_profile data source returned more than one result."))
	}
	if resp.Diagnostics.HasError() {
		return
	}

	profile := filteredProfiles[0]
	// Set state
	state = ApiToProfile(profile)
	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
