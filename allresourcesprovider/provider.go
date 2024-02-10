package allresourcesprovider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-kubernetes/kubernetes"
)

type Provider struct {
	singleClusterProvider *schema.Provider
	TerraformVersion string
}

func NewProvider() *Provider {
	return &Provider{
		singleClusterProvider: kubernetes.Provider(),
		TerraformVersion: "corentin",
	}
}

func (p *Provider) TFProvider() *schema.Provider {
	tfprovider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"per_cluster": {
				Type: schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema:p.singleClusterProvider.Schema,
					},
				},
			},
		},
		ResourcesMap:   p.buildResourcesMap(p.singleClusterProvider),
		DataSourcesMap: p.buildDataSourcesMap(p.singleClusterProvider),
	}
	tfprovider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return providerConfigure(ctx, d, p.TerraformVersion, p.singleClusterProvider)
	}

	return tfprovider
}

func providerConfigure(ctx context.Context, d *schema.ResourceData, terraformVersion string, singleClusterProvider *schema.Provider) (interface{}, diag.Diagnostics) {

	clientSets := map[string]kubernetes.KubeClientsets{}
	if v, ok := d.GetOk("per_cluster"); ok {
		for clusterName, clusterConfig := range v.(map[string]interface{}) {
			if perClusterConfigResourceData, ok := clusterConfig.([]interface{})[0].(map[string]interface{}); ok {
				perClusterConfigResourceDataAsResourceData := perClusterConfigResourceData.(schema.ResourceData)
				clientset, diags := singleClusterProvider.ConfigureContextFunc(ctx, perClusterConfigResourceDataAsResourceData)
				if diags.HasError() {
					return nil, diags
				}
				clientSets[clusterName] = clientset.(kubernetes.KubeClientsets)
			} else {
				return nil, fmt.Errorf("Failed to parse per_cluster")
			}
		}
	} else {
		return nil, fmt.Errorf("cant find per_cluster")
	}


	return ProviderContextMeta{
		clientSets: clientSets,
	}, diag.Diagnostics{}
}

func (p *Provider) buildResourcesMap(singleClusterProvider *schema.Provider) map[string]*schema.Resource {
	resourceMap := make(map[string]*schema.Resource, len(singleClusterProvider.ResourcesMap))
	for name, definition := range singleClusterProvider.ResourcesMap {
		resourceMap["multi"+name] = p.buildNewResourceDefinition(definition)
	}
	return resourceMap
}

func (p *Provider) buildDataSourcesMap(singleClusterProvider *schema.Provider) map[string]*schema.Resource {
	resourceMap := make(map[string]*schema.Resource, len(singleClusterProvider.DataSourcesMap))
	for name, definition := range singleClusterProvider.DataSourcesMap {
		resourceMap["multi"+name] = p.buildNewResourceDefinition(definition)
	}
	return resourceMap
}

func (p *Provider) buildNewResourceDefinition(originalDefinition *schema.Resource) *schema.Resource {
	return &schema.Resource{
		CreateContext: wrapOriginalActionContext(originalDefinition.CreateContext),
		ReadContext:   wrapOriginalActionContext(originalDefinition.ReadContext),
		UpdateContext: wrapOriginalActionContext(originalDefinition.UpdateContext),
		DeleteContext: wrapOriginalActionContext(originalDefinition.DeleteContext),
		Importer:      originalDefinition.Importer,
		Timeouts:      originalDefinition.Timeouts,
		SchemaVersion: originalDefinition.SchemaVersion,
		Schema:        AddClusterToSchema(originalDefinition.Schema),
	}
}

func wrapOriginalActionContext(originalFunc func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics) func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if originalFunc == nil {
		return nil
	}
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		pdm := meta.(ProviderContextMeta)
		clientSet, err := pdm.clientSetFromSchema(d)
		if err != nil {
			return diag.FromErr(err)
		}
		return originalFunc(ctx, d, clientSet)
	}
}

func AddClusterToSchema(originalSchema map[string]*schema.Schema) map[string]*schema.Schema {
	originalSchema["cluster"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Cluster to which apply the resource",
		ForceNew:    true,
		Required:     true,
	}
	return originalSchema
}

type ProviderContextMeta struct {
	clientSets map[string]kubernetes.KubeClientsets
}

func (p *ProviderContextMeta) clientSetFromSchema(d *schema.ResourceData) (kubernetes.KubeClientsets, error) {
	clusterName := d.Get("cluster").(string)
	clientSet, found := p.clientSets[clusterName]
	if !found {
		return nil, fmt.Errorf("Cluster not Found")
	}
	return clientSet, nil
}
