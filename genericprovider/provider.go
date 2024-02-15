package genericprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	manifestprovider "github.com/hashicorp/terraform-provider-kubernetes/manifest/provider"
)


type Provider struct {
	// this provider is used to deal with Schema configuration and static answers,
	// when we need an instantiated object, but not configured.
	// This provider is NOT used for actual resource management.
	configOnlySingleClusterProvider *manifestprovider.RawProviderServer
	singleClusterProviders map[string]*manifestprovider.RawProviderServer
	logLevel string
	logger hclog.Logger
}

func NewProvider() *Provider {
	var logLevel string
	var ok bool = false
	for _, ev := range []string{"TF_LOG_PROVIDER_KUBERNETES", "TF_LOG_PROVIDER", "TF_LOG"} {
		logLevel, ok = os.LookupEnv(ev)
		if ok {
			break
		}
	}
	if !ok {
		logLevel = "off"
	}

	configOnlySingleClusterProvider := newsingleClusterProvider(logLevel)

	return &Provider{
		singleClusterProviders: map[string]*manifestprovider.RawProviderServer{},
		configOnlySingleClusterProvider: configOnlySingleClusterProvider,
		logLevel: logLevel,
		logger: hclog.New(&hclog.LoggerOptions{
			Level:  hclog.LevelFromString(logLevel),
			Output: os.Stderr,
		}),
	}
}

func newsingleClusterProvider(logLevel string) *manifestprovider.RawProviderServer {
	s := &(manifestprovider.RawProviderServer{})
	s.SetLogger(hclog.New(&hclog.LoggerOptions{
		Level:  hclog.LevelFromString(logLevel),
		Output: os.Stderr,
	}))
	return s
}

func (p *Provider) TFProvider() tfprotov5.ProviderServer {
	return p
}


// Main Provider methods (ProviderServer)
// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-go@v0.21.0/tfprotov5#ProviderServer

func (p *Provider) GetMetadata(ctx context.Context, req *tfprotov5.GetMetadataRequest) (*tfprotov5.GetMetadataResponse, error){
	// This method is not using the existing implementation of the method because we can plug directly into the public Schema
	// functions that get us the data we need, GetProviderResourceSchema and GetProviderDataSourceSchema.
	// That way, we have our own GetProviderResourceSchema and just use that wherever we need it (and especially also in GetProviderSchema).
	p.logger.Trace("[GetMetadata][Request]\n%s\n", dump(*req))

	sch := GetProviderResourceSchema()
	rs := make([]tfprotov5.ResourceMetadata, 0, len(sch))
	for k := range sch {
		rs = append(rs, tfprotov5.ResourceMetadata{TypeName: k})
	}

	sch = GetProviderDataSourceSchema()
	ds := make([]tfprotov5.DataSourceMetadata, 0, len(sch))
	for k := range sch {
		ds = append(ds, tfprotov5.DataSourceMetadata{TypeName: k})
	}

	resp := &tfprotov5.GetMetadataResponse{
		Resources:   rs,
		DataSources: ds,
	}
	return resp, nil
}

func (p *Provider) GetProviderSchema(ctx context.Context, req *tfprotov5.GetProviderSchemaRequest) (*tfprotov5.GetProviderSchemaResponse, error) {
	p.logger.Trace("[Multi-GetProviderSchema][Request]\n%s\n", dump(*req))
	// This method is reimplemented because it's super simple and it's easier to inject ourselves in the public GetProviderConfigSchema and have our own than
	// intercepting the actual method.
	cfgSchema := GetProviderConfigSchema()
	resSchema := GetProviderResourceSchema()
	dsSchema := GetProviderDataSourceSchema()

	return &tfprotov5.GetProviderSchemaResponse{
		Provider:          cfgSchema,
		ResourceSchemas:   resSchema,
		DataSourceSchemas: dsSchema,
	}, nil
}

func (p *Provider) providerConfigHandler(config *tfprotov5.DynamicValue, f func(string, tftypes.Value) error) ([]*tfprotov5.Diagnostic, error){
	p.logger.Trace("[Multi-providerConfigHandler]", "dynvalue", config)
	var diagnostics []*tfprotov5.Diagnostic

	cfgType := manifestprovider.GetObjectTypeFromSchema(GetProviderConfigSchema())
	cfgVal, err := config.Unmarshal(cfgType)
	if err != nil {
		diagnostics = append(diagnostics, &tfprotov5.Diagnostic{
			Severity: tfprotov5.DiagnosticSeverityError,
			Summary:  "Failed to decode Provider Configuration parameter",
			Detail:   err.Error(),
		})
		return diagnostics, nil
	}

	var providerConfig map[string]tftypes.Value
	err = cfgVal.As(&providerConfig)
	if err != nil {
		// invalid configuration schema - this shouldn't happen, bail out now
		diagnostics = append(diagnostics, &tfprotov5.Diagnostic{
			Severity: tfprotov5.DiagnosticSeverityError,
			Summary:  "Provider configuration: failed to extract 'config_path' value",
			Detail:   err.Error(),
		})
		return diagnostics, nil
	}

	p.logger.Trace("[Multi-providerConfigHandler] Got Config", "config", providerConfig)

	var clusterBlock []tftypes.Value
	err = providerConfig["cluster"].As(&clusterBlock)
	if err != nil {
		// invalid attribute type - this shouldn't happen, bail out for now
		diagnostics = append(diagnostics, &tfprotov5.Diagnostic{
			Severity: tfprotov5.DiagnosticSeverityError,
			Summary:  "Provider configuration: failed to assert type of 'cluster' value",
			Detail:   err.Error(),
		})
		return diagnostics, nil
	}

	p.logger.Trace("[Multi-providerConfigHandler] Got ClusterBlock", "clusterBlock", clusterBlock)

	for _, clusterConfigAsValue := range clusterBlock {
		var clusterConfig map[string]tftypes.Value
		err := clusterConfigAsValue.As(&clusterConfig)
		if err != nil {
			// invalid attribute type - this shouldn't happen, bail out for now
			diagnostics = append(diagnostics, &tfprotov5.Diagnostic{
				Severity: tfprotov5.DiagnosticSeverityError,
				Summary:  "Provider configuration: failed to assert type of 'Cluster' value",
				Detail:   err.Error(),
			})
			return diagnostics, nil
		}
		var clusterName string
		err = clusterConfig["cluster_name"].As(&clusterName)
		if err != nil {
			// invalid attribute type - this shouldn't happen, bail out for now
			diagnostics = append(diagnostics, &tfprotov5.Diagnostic{
				Severity: tfprotov5.DiagnosticSeverityError,
				Summary:  "Provider configuration: failed to assert type of 'cluster_name' value",
				Detail:   err.Error(),
			})
			return diagnostics, nil
		}

		p.logger.Trace("[Multi-providerConfigHandler] Got cluster_name", "cluster_name", clusterName)
		if clusterName == "" {// remove me if we dont get here
			// invalid attribute type - this shouldn't happen, bail out for now
			diagnostics = append(diagnostics, &tfprotov5.Diagnostic{
				Severity: tfprotov5.DiagnosticSeverityError,
				Summary:  "Provider configuration: 'cluster_name' is empty",
				Detail:   err.Error(),
			})
			return diagnostics, nil
		}

		err = f(clusterName, clusterConfigAsValue)
		if err != nil {
			return nil, err
		}

	}

	return diagnostics, nil
}

func (p *Provider) PrepareProviderConfig(ctx context.Context, req *tfprotov5.PrepareProviderConfigRequest) (*tfprotov5.PrepareProviderConfigResponse, error) {
	p.logger.Trace("[Multi-PrepareProviderConfig][Request]\n", dump(*req))
	response := &tfprotov5.PrepareProviderConfigResponse{}

	diagnostics, err := p.providerConfigHandler(req.Config, func(clusterName string, clusterConfig tftypes.Value) error {
		p.logger.Trace("[Multi-PrepareProviderConfig]", "clusterConfig", clusterConfig, "clusterSchema", clusterSchemaBlock())
		singleClusterDynValue, err := tfprotov5.NewDynamicValue(manifestprovider.GetObjectTypeFromSchema(blockAsSchema(clusterSchemaBlock())), clusterConfig)
		if err != nil {
			return fmt.Errorf("Error converting to DynamicValue:%w",err)
		}

		// We create a RawProvider for each cluster right away.
		p.singleClusterProviders[clusterName] = newsingleClusterProvider(p.logLevel)

		// We call PrepareProviderConfig on each of them.
		singleClusterResponse, err := p.singleClusterProviders[clusterName].PrepareProviderConfig(ctx, &tfprotov5.PrepareProviderConfigRequest{
			Config: &singleClusterDynValue,
		})
		if err != nil {
			return fmt.Errorf("Error single cluster PrepareProviderConfig:%w",err)
		}
		if len(singleClusterResponse.Diagnostics) > 0 {
			response.Diagnostics = append(response.Diagnostics, singleClusterResponse.Diagnostics...)
		}
		// TODO we don't take the returned response.PreparedConfig element; luckily the current provider code doesn't leverage it.

		return nil
	})
	if err != nil {
		return nil, err
	}
	response.Diagnostics = append(response.Diagnostics, diagnostics...)
	p.logger.Trace("[Multi-PrepareProviderConfig] Created providers", "len", len(p.singleClusterProviders), "keys", maps.Keys(p.singleClusterProviders))

	return response, nil
}

func (p *Provider) ConfigureProvider(ctx context.Context, req *tfprotov5.ConfigureProviderRequest) (*tfprotov5.ConfigureProviderResponse, error) {
	p.logger.Trace("[ConfigureProvider][Request]\n", dump(*req))
	response := &tfprotov5.ConfigureProviderResponse{}

	diagnostics, err := p.providerConfigHandler(req.Config, func(clusterName string, clusterConfig tftypes.Value) error {
		singleClusterDynValue, err := tfprotov5.NewDynamicValue(manifestprovider.GetObjectTypeFromSchema(manifestprovider.GetProviderConfigSchema()), clusterConfig)
		if err != nil {
			return fmt.Errorf("Error converting to DynamicValue:%w",err)
		}

		singleClusterProvider, found :=  p.singleClusterProviders[clusterName]
		if !found {
			p.logger.Trace("[Multi-ConfigureProvider][Request]PrepareProviderConfig was not called, creating the single cluster provider now.")
			singleClusterProvider = newsingleClusterProvider(p.logLevel)
			p.singleClusterProviders[clusterName] = singleClusterProvider
		}

		singleClusterResponse, err := singleClusterProvider.ConfigureProvider(ctx, &tfprotov5.ConfigureProviderRequest{
			Config: &singleClusterDynValue,
			TerraformVersion: req.TerraformVersion,
		})
		if err != nil {
			return fmt.Errorf("Error single cluster ConfigureProvider:%w",err)
		}
		if len(singleClusterResponse.Diagnostics) > 0 {
			response.Diagnostics = append(response.Diagnostics, singleClusterResponse.Diagnostics...)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	response.Diagnostics = append(response.Diagnostics, diagnostics...)
	p.logger.Trace("[Multi-ConfigureProvider] Created providers", "len", len(p.singleClusterProviders), "keys", maps.Keys(p.singleClusterProviders))

	if len(response.Diagnostics) > 0 {
		p.logger.Trace("[Multi-ConfigureProvider] has diags", "diags", response.Diagnostics)
	}

	return response, nil
}

func (p *Provider) StopProvider(ctx context.Context, req *tfprotov5.StopProviderRequest) (*tfprotov5.StopProviderResponse, error) {
	// TODO(corentone): call StopProvider on each of the provider and exit.
	return nil, fmt.Errorf("Unimplemented 1")
}

// ResourceServer interface methods
// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-go@v0.21.0/tfprotov5#ResourceServer

func (p *Provider) ValidateResourceTypeConfig(ctx context.Context, req *tfprotov5.ValidateResourceTypeConfigRequest) (*tfprotov5.ValidateResourceTypeConfigResponse, error) {
	response := &tfprotov5.ValidateResourceTypeConfigResponse{}
	// We could check that the cluster string is set here. But normally the config validation enforces it already so it's not too interesting.

	// to call the config provider, we need to first convert the req content to the right schema (rename the resource)
	// return p.configOnlySingleClusterProvider.ValidateResourceTypeConfig(ctx, req)
	return response, nil
}
func (p *Provider) UpgradeResourceState(ctx context.Context, req *tfprotov5.UpgradeResourceStateRequest) (*tfprotov5.UpgradeResourceStateResponse, error) {
	response := &tfprotov5.UpgradeResourceStateResponse{}

	p.logger.Trace("[Multi-PlanResourceChange]", "RawStateJSON", string(req.RawState.JSON[:]))

	var rawState map[string]any
	err := json.Unmarshal(req.RawState.JSON, &rawState)
	if err != nil {
		return response, fmt.Errorf("could not decode RawStateJSON: %w", err)
	}

	clusterNameAsAny, found := rawState["cluster"]
	if !found {
		return response, fmt.Errorf("cluster not found in state")
	}
	clusterName := 	clusterNameAsAny.(string)

	p.logger.Trace("[Multi-PlanResourceChange]", "clusterName", clusterName)

	// For now only calling the PlanResourceChange on the new cluster since we will anyway just delete it from the old one...
	singleClusterProvider, found := p.singleClusterProviders[clusterName]
	if !found {
		return response, fmt.Errorf("Requested cluster was not found in the provider.")
	}

	delete(rawState, "cluster")

	singleClusterRawState , err := json.Marshal(rawState)
	if err != nil {
		return response, fmt.Errorf("could not encode RawState: %w", err)
	}

	singleClusterReq := &tfprotov5.UpgradeResourceStateRequest{
		TypeName: strings.TrimPrefix(req.TypeName, "multi"),
		Version: req.Version,
		RawState: &tfprotov5.RawState{JSON:singleClusterRawState},
	}

	singleClusterResp, err := singleClusterProvider.UpgradeResourceState(ctx, singleClusterReq)
  response.Diagnostics = append(response.Diagnostics, singleClusterResp.Diagnostics...)
	if err != nil {
		return response, err
	}

	// we may need to convert back the resp.
	response.UpgradedState, err = p.multiClusterResourceDynamicValueWithCluster(singleClusterResp.UpgradedState, strings.TrimPrefix(req.TypeName, "multi"), clusterName)
	if err != nil {
		// swallow the error for now.
		return response, nil
	}
	return response, nil
}
func (p *Provider) ReadResource(ctx context.Context, req *tfprotov5.ReadResourceRequest) (*tfprotov5.ReadResourceResponse, error) {
	response := &tfprotov5.ReadResourceResponse{}

	clusterName, diagnostics, err := p.extractClusterName(req.TypeName, req.CurrentState)
	response.Diagnostics = append(response.Diagnostics, diagnostics...)
	if err != nil {
		return response, err
	}

	p.logger.Trace("[Multi-ReadResource]", "clusterName", clusterName)

	singleClusterProvider, found := p.singleClusterProviders[clusterName]
	if !found {
		return response, fmt.Errorf("Requested cluster was not found in the provider.")
	}

	currentState, err := p.singleClusterResourceDynamicValue(req.CurrentState, req.TypeName)
	if err != nil {
		return response, err
	}

	singleClusterReq := &tfprotov5.ReadResourceRequest{
		TypeName: strings.TrimPrefix(req.TypeName, "multi"),
		CurrentState: currentState,
	}

	singleClusterResp, err := singleClusterProvider.ReadResource(ctx, singleClusterReq)
	response.Diagnostics = append(response.Diagnostics, singleClusterResp.Diagnostics...)
	if err != nil {
		return response, err
	}
	response.Private = singleClusterResp.Private

	p.logger.Trace("[Multi-ReadResource]", "singleClusterResp.NewState", singleClusterResp.NewState)
	newState, err := p.multiClusterResourceDynamicValueWithCluster(singleClusterResp.NewState, strings.TrimPrefix(req.TypeName, "multi"), clusterName)
	if err != nil {
		return response, nil
	}
	response.NewState = newState
	p.logger.Trace("[Multi-ReadResource]", "NewState", response.NewState)
	return response, nil
}
func (p *Provider) PlanResourceChange(ctx context.Context, req *tfprotov5.PlanResourceChangeRequest) (*tfprotov5.PlanResourceChangeResponse, error) {
	response := &tfprotov5.PlanResourceChangeResponse{}

	response.RequiresReplace = append(response.RequiresReplace,
		tftypes.NewAttributePath().WithAttributeName("cluster"),
	)

	oldClusterName, diagnostics, err := p.extractClusterName(req.TypeName, req.PriorState)
	response.Diagnostics = append(response.Diagnostics, diagnostics...)
	if err != nil {
		return response, err
	}
	newClusterName, diagnostics, err := p.extractClusterName(req.TypeName, req.ProposedNewState)
	response.Diagnostics = append(response.Diagnostics, diagnostics...)
	if err != nil {
		return response, err
	}

	p.logger.Trace("[Multi-PlanResourceChange]", "oldClusterName", oldClusterName, "newClusterName", newClusterName)

	// For now only calling the PlanResourceChange on the new cluster since we will anyway just delete it from the old one...
	singleClusterProvider, found := p.singleClusterProviders[newClusterName]
	if !found {
		return response, fmt.Errorf("Requested cluster was not found in the provider.")
	}

	// The DynamicValue contains our Schema, we need to rebuild it with the schema of the underlying provider
	remappedPriorState, err := p.singleClusterResourceDynamicValue(req.PriorState, req.TypeName)
	if err != nil {
		response.Diagnostics = append(response.Diagnostics, &tfprotov5.Diagnostic{
				Severity: tfprotov5.DiagnosticSeverityError,
				Summary:  "Failed to decode PriorState parameter",
				Detail:   err.Error(),
			})
	}
	req.PriorState = remappedPriorState

	remappedProposedNewState, err := p.singleClusterResourceDynamicValue(req.ProposedNewState, req.TypeName)
	if err != nil {
		response.Diagnostics = append(response.Diagnostics, &tfprotov5.Diagnostic{
				Severity: tfprotov5.DiagnosticSeverityError,
				Summary:  "Failed to decode ProposedNewState parameter",
				Detail:   err.Error(),
			})
	}
	req.ProposedNewState = remappedProposedNewState

	// The resource is named differently in our provider, we just need to change the TypeName to have it read properly (remove "multi")
	req.TypeName = strings.TrimPrefix(req.TypeName, "multi")

	resp, err := singleClusterProvider.PlanResourceChange(ctx, req)
	response.Diagnostics = append(response.Diagnostics, resp.Diagnostics...)
	if err != nil {
		return response, fmt.Errorf("PlanResourceChange for cluster %v failed:%w", newClusterName, err)
	}
	response.RequiresReplace = append(response.RequiresReplace, resp.RequiresReplace...)
	response.PlannedPrivate = resp.PlannedPrivate

	// The Planned state was written in single provider schema so we map it back
	remappedPlannedState, err := p.multiClusterResourceDynamicValueWithCluster(resp.PlannedState, req.TypeName, newClusterName)
	if err != nil {
		response.Diagnostics = append(response.Diagnostics, &tfprotov5.Diagnostic{
				Severity: tfprotov5.DiagnosticSeverityError,
				Summary:  "Failed to decode PlannedState parameter",
				Detail:   err.Error(),
			})
	}
	response.PlannedState = remappedPlannedState

	return response, nil
}

func (p *Provider) ApplyResourceChange(ctx context.Context, req *tfprotov5.ApplyResourceChangeRequest) (*tfprotov5.ApplyResourceChangeResponse, error) {
	response := &tfprotov5.ApplyResourceChangeResponse{}

	oldClusterName, diagnostics, err := p.extractClusterName(req.TypeName, req.PriorState)
	response.Diagnostics = append(response.Diagnostics, diagnostics...)
	if err != nil {
		return response, err
	}
	newClusterName, diagnostics, err := p.extractClusterName(req.TypeName, req.PlannedState)
	response.Diagnostics = append(response.Diagnostics, diagnostics...)
	if err != nil {
		return response, err
	}

	p.logger.Trace("[Multi-ApplyResourceChange]", "oldClusterName", oldClusterName, "newClusterName", newClusterName)

	var clusterToApply, clusterToDelete string
	if oldClusterName == newClusterName {
		// No change in cluster
		clusterToApply = newClusterName
	} else {
		// Move from oldCluster to NewCluster
		clusterToDelete = oldClusterName
		clusterToApply = newClusterName

		// if either one of the planned or prior cluster is empty, it's special cases
		if oldClusterName == "" {
			// Creation, no cluster to delete
			clusterToDelete = ""
		}
		if newClusterName == "" {
			// Deletion, no cluster to apply
			clusterToApply = ""
		}
	}

	singleClusterConfig, err := p.singleClusterResourceDynamicValue(req.Config, req.TypeName)
	if err != nil {
		return response, err
	}

	var clusterToDeleteResp, clusterToApplyResp *tfprotov5.ApplyResourceChangeResponse
	if clusterToDelete != "" {
		clusterToDeleteProvider, found := p.singleClusterProviders[clusterToDelete]
		if !found {
			return response, fmt.Errorf("Could not find cluster to delete provider %v", clusterToDelete)
		}

		var plannedState *tfprotov5.DynamicValue
		if clusterToApply == "" {
			// There is no new cluster so it's a pure deletion, the planned state is therefore good
			plannedState, err = p.singleClusterResourceDynamicValue(req.PlannedState, req.TypeName)
			if err != nil {
				return response, err
			}
		} else {
			// there is a new cluster, so we need to clean up the Planned State for this cluster to be deleted
			plannedState, err = singleClusterResourceEmptyDynamicValue(req.TypeName)
			if err != nil {
				return response, err
			}
		}

		priorState, err := p.singleClusterResourceDynamicValue(req.PriorState, req.TypeName)
		if err != nil {
			return response, err
		}

		// Request for this provider contains the type renamed, PlannedState
		applyReq := &tfprotov5.ApplyResourceChangeRequest{
			TypeName: strings.TrimPrefix(req.TypeName, "multi"),
			PriorState: priorState,
			PlannedState: plannedState,
			Config: singleClusterConfig,
		}
		clusterToDeleteResp, err = clusterToDeleteProvider.ApplyResourceChange(ctx, applyReq)
		if err != nil {
			return response, fmt.Errorf("Failed ApplyResourceChange on %v, req:%v, err:%v", clusterToApply, applyReq, err)
		}
		response.Diagnostics = append(clusterToDeleteResp.Diagnostics, diagnostics...)
		p.logger.Trace("[Multi-ApplyResourceChange] called delete", "clusterName", clusterToDeleteResp, "req", applyReq, "resp", clusterToDeleteResp)
	}

	if clusterToApply != "" {
		clusterToApplyProvider, found := p.singleClusterProviders[clusterToApply]
		if !found {
			return response, fmt.Errorf("Could not find cluster to apply provider %v", clusterToApply)
		}

		var applyPriorState *tfprotov5.DynamicValue
		if clusterToDelete == "" {
			// same cluster or full create, we can use the req.PriorState, which is valid for this cluster
			applyPriorState, err = p.singleClusterResourceDynamicValue(req.PriorState, req.TypeName)
			if err != nil {
				return response, err
			}

		} else {
			// there is an old cluster to clean, so we are transitioning clusters, so we need to pass an empty priorState to this cluster
			applyPriorState, err = singleClusterResourceEmptyDynamicValue(req.TypeName)
			if err != nil {
				return response, err
			}
		}

		plannedState, err := p.singleClusterResourceDynamicValue(req.PlannedState, req.TypeName)
		if err != nil {
			return response, err
		}

		// Request for this provider contains the type renamed, PriorState
		applyReq := &tfprotov5.ApplyResourceChangeRequest{
			TypeName: strings.TrimPrefix(req.TypeName, "multi"),
			PriorState: applyPriorState,
			PlannedState: plannedState,
			Config: singleClusterConfig,
			PlannedPrivate: req.PlannedPrivate,
		}
		clusterToApplyResp, err = clusterToApplyProvider.ApplyResourceChange(ctx, applyReq)
		if err != nil {
			return response, fmt.Errorf("Failed ApplyResourceChange on %v, req:%v, err:%v", clusterToApply, applyReq, err)
		}
		response.Diagnostics = append(clusterToApplyResp.Diagnostics, diagnostics...)
		p.logger.Trace("[Multi-ApplyResourceChange] called apply", "clusterName", clusterToApply, "req", applyReq, "resp", clusterToApplyResp)
	}

	var newState *tfprotov5.DynamicValue
	if clusterToApply != "" {
		newState, err = p.multiClusterResourceDynamicValueWithCluster(clusterToApplyResp.NewState, strings.TrimPrefix(req.TypeName, "multi"), clusterToApply)
		if err != nil {
			// slurps the error for now
			return response, nil
		}
	} else {
		// delete case
		newState = req.PlannedState
	}
	response.NewState = newState

	return response, nil
}
func (p *Provider) ImportResourceState(ctx context.Context, req *tfprotov5.ImportResourceStateRequest) (*tfprotov5.ImportResourceStateResponse, error) {
	// TODO(corentone): TBD
	return nil, fmt.Errorf("Unimplemented 7")
}

// DataSourceServer interface
// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-go@v0.21.0/tfprotov5#DataSourceServer

func (p *Provider) ValidateDataSourceConfig(ctx context.Context, req *tfprotov5.ValidateDataSourceConfigRequest) (*tfprotov5.ValidateDataSourceConfigResponse, error) {
	// TODO(corentone): TBD
	return nil, fmt.Errorf("Unimplemented 8")
}
func (p *Provider) ReadDataSource(ctx context.Context, req *tfprotov5.ReadDataSourceRequest) (*tfprotov5.ReadDataSourceResponse, error) {
	// TODO(corentone): TBD
	return nil, fmt.Errorf("Unimplemented 9")
}

// ---

func GetProviderResourceSchema() map[string]*tfprotov5.Schema {
	singleClusterSchema  := manifestprovider.GetProviderResourceSchema()

	singleClusterSchema["kubernetes_manifest"].Block.Attributes = append(singleClusterSchema["kubernetes_manifest"].Block.Attributes,
	&tfprotov5.SchemaAttribute{
		Name:        "cluster",
		Type:        tftypes.String,
		Required:    true,
		Description: "Cluster to which apply the resource.",
	})

	singleClusterSchema["multikubernetes_manifest"] = singleClusterSchema["kubernetes_manifest"]
	delete(singleClusterSchema, "kubernetes_manifest")

	return singleClusterSchema
}

func GetProviderDataSourceSchema() map[string]*tfprotov5.Schema {
	singleClusterSchema  := manifestprovider.GetProviderDataSourceSchema()

	singleClusterSchema["kubernetes_resource"].Block.Attributes = append(singleClusterSchema["kubernetes_resource"].Block.Attributes,
	&tfprotov5.SchemaAttribute{
		Name:        "cluster",
		Type:        tftypes.String,
		Required:    true,
		Description: "Cluster to which apply the resource.",
	})
	singleClusterSchema["multikubernetes_resource"] = singleClusterSchema["kubernetes_resource"]
	delete(singleClusterSchema, "kubernetes_resource")

	singleClusterSchema["kubernetes_resources"].Block.Attributes = append(singleClusterSchema["kubernetes_resources"].Block.Attributes,
	&tfprotov5.SchemaAttribute{
		Name:        "cluster",
		Type:        tftypes.String,
		Required:    true,
		Description: "Cluster to which apply the resources.",
	})
	singleClusterSchema["multikubernetes_resources"] = singleClusterSchema["kubernetes_resources"]
	delete(singleClusterSchema, "kubernetes_resources")

	return singleClusterSchema
}

func GetProviderConfigSchema() *tfprotov5.Schema {
	b := tfprotov5.SchemaBlock{
		BlockTypes: []*tfprotov5.SchemaNestedBlock{
			{
				TypeName: "cluster",
				Nesting:  tfprotov5.SchemaNestedBlockNestingModeList,
				MinItems: 0,
				Block: clusterSchemaBlock(),
			},
		},
	}

	return &tfprotov5.Schema{
		Version: 0,
		Block:   &b,
	}
}

func clusterSchemaBlock() *tfprotov5.SchemaBlock {
	singleClusterSchemaBlock := manifestprovider.GetProviderConfigSchema().Block
	return &tfprotov5.SchemaBlock{
		Attributes: append(singleClusterSchemaBlock.Attributes,
			&tfprotov5.SchemaAttribute{
				Name:            "cluster_name",
				Type:            tftypes.String,
				Required:        true,
				Optional:        false,
				Computed:        false,
				Sensitive:       false,
				Description:     "Cluster Name to which the config refers.",
				DescriptionKind: 0,
				Deprecated:      false,
			}),
		BlockTypes: singleClusterSchemaBlock.BlockTypes,
	}
}

func blockAsSchema(b *tfprotov5.SchemaBlock) *tfprotov5.Schema {
	return &tfprotov5.Schema{
		Version: 0,
		Block:   b,
	}
}

func dump(v interface{}) hclog.Format {
	return hclog.Fmt("%v", v)
}

func (p *Provider) extractClusterName(typeName string, val *tfprotov5.DynamicValue) (string, []*tfprotov5.Diagnostic, error){
	diagnostics := []*tfprotov5.Diagnostic{}
	cfgType := manifestprovider.GetObjectTypeFromSchema(GetProviderResourceSchema()[typeName])
	cfgVal, err := val.Unmarshal(cfgType)
	if err != nil {
		diagnostics = append(diagnostics, &tfprotov5.Diagnostic{
			Severity: tfprotov5.DiagnosticSeverityError,
			Summary:  "Failed to decode Provider Configuration parameter",
			Detail:   err.Error(),
		})
		return "", diagnostics, nil
	}

	var resourceConfig map[string]tftypes.Value
	err = cfgVal.As(&resourceConfig)
	if err != nil {
		// invalid configuration schema - this shouldn't happen, bail out now
		diagnostics = append(diagnostics, &tfprotov5.Diagnostic{
			Severity: tfprotov5.DiagnosticSeverityError,
			Summary:  "Provider configuration: failed to extract 'config_path' value",
			Detail:   err.Error(),
		})
		return "", diagnostics, nil
	}

	p.logger.Trace("[Multi-PlanResourceChange] Got old cluster config", "resourceConfig", resourceConfig)

	var clusterName string
	err = resourceConfig["cluster"].As(&clusterName)
	if err != nil {
		return "", nil, fmt.Errorf("couldn't unpack the clustername(val=%v), got:%w", resourceConfig["cluster"], err)
	}

	return clusterName, diagnostics, nil
}

func singleClusterResourceEmptyDynamicValue(typeName string) (*tfprotov5.DynamicValue, error){
	typeName = strings.TrimPrefix(typeName, "multi")
	val := tftypes.NewValue(manifestprovider.GetObjectTypeFromSchema(manifestprovider.GetProviderResourceSchema()[typeName]), nil)
	dn, err := tfprotov5.NewDynamicValue(manifestprovider.GetObjectTypeFromSchema(manifestprovider.GetProviderResourceSchema()[typeName]), val)
	return &dn, err
}

func (p *Provider) singleClusterResourceDynamicValue(config *tfprotov5.DynamicValue, typeName string) (*tfprotov5.DynamicValue, error) {
	// typename is implied to be the prefixed with multi one.
	cfgType := manifestprovider.GetObjectTypeFromSchema(GetProviderResourceSchema()[typeName])
	cfgVal, err := config.Unmarshal(cfgType)
	if err != nil {
		return nil, err
	}

	// since the typename is the prefixed one, we remove the prefix to find in the singleCluster schema.
	typeName = strings.TrimPrefix(typeName, "multi")
	dn, err := tfprotov5.NewDynamicValue(manifestprovider.GetObjectTypeFromSchema(manifestprovider.GetProviderResourceSchema()[typeName]), cfgVal)
	return &dn, err
}

func (p *Provider) multiClusterResourceDynamicValueWithCluster(config *tfprotov5.DynamicValue, typeName string, clusterName string) (*tfprotov5.DynamicValue, error) {
	p.logger.Trace("[Multi-multiClusterResourceDynamicValueWithCluster] Conversion", "config", config, "typeName", typeName, "clusterName", clusterName)
	if config == nil {
		return nil, fmt.Errorf("config should not be nil")
	}
	// typename is implied to be NOT prefixed
	cfgType := manifestprovider.GetObjectTypeFromSchema(manifestprovider.GetProviderResourceSchema()[typeName])
	cfgVal, err := config.Unmarshal(cfgType)
	if err != nil {
		return nil, err
	}

	// since the typename is not prefixed, we prefix it again.
	typeName = "multi"+typeName


	var cfgValAsMap map[string]tftypes.Value
	err = cfgVal.As(&cfgValAsMap)
	cfgValAsMap["cluster"] = tftypes.NewValue(tftypes.String, clusterName)

	val := tftypes.NewValue(manifestprovider.GetObjectTypeFromSchema(GetProviderResourceSchema()[typeName]), cfgValAsMap)
	dn, err := tfprotov5.NewDynamicValue(manifestprovider.GetObjectTypeFromSchema(GetProviderResourceSchema()[typeName]), val)
	return &dn, err
}

func NewRawState(jsonMap map[string]tftypes.Value ) *tfprotov5.RawState {
	rawStateJSON, _ := json.Marshal(jsonMap)
	return &tfprotov5.RawState{
		JSON: rawStateJSON,
	}
}
