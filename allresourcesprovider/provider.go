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
		Schema: p.schemaFromOriginal(),
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
	if configClusters, ok := d.GetOk("cluster"); ok {
		for _, configCluster := range configClusters.([]interface{}) {
			perClusterConfigData, ok := configCluster.(map[string]interface{});
			if  !ok {
				return nil, diag.FromErr(fmt.Errorf("Failed to parse cluster %v", configCluster))
			}

			clusterName := perClusterConfigData["cluster_name"].(string)

			gettableResourceData := localResourceData{perClusterConfigResourceData:perClusterConfigData}
			clientset, diags := duplicatedProviderConfigure(ctx, gettableResourceData, terraformVersion)
			if diags.HasError() {
				return nil, diags
			}
			clientSets[clusterName] = clientset.(kubernetes.KubeClientsets)
		}
	} else {
		return nil, diag.FromErr(fmt.Errorf("cant find cluster"))
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


/*
The schema as being defined by the original provider and added to the list doesn't seem to work.
I'm guessing I need to simplify it sadly.

Error: Incorrect attribute value type

  on main.tf line 12, in provider "multikubernetes":
  12:   per_cluster = {
  13:     corentin = [{
  14:       host  = "host.com"
  15:       exec = {
  16:         api_version = "client.authentication.k8s.io/v1beta1"
  17:         command     = "gke-gcloud-auth-plugin"
  18:       }
  19:     }]
  20:   }

Inappropriate value for attribute "per_cluster": element "corentin": element
0: attributes "client_certificate", "client_key", "cluster_ca_certificate",
"config_context", "config_context_auth_info", "config_context_cluster",
"config_path", "config_paths", "experiments", "ignore_annotations",
"ignore_labels", "insecure", "password", "proxy_url", "tls_server_name",
"token", and "username" are required.


*/
func (p *Provider) schemaFromOriginal() map[string]*schema.Schema {

	// Our schema is literally the existing schema, with an additional key called cluster_name,
	// similarly to every resource we have.
	perClusterSchema := p.singleClusterProvider.Schema
	perClusterSchema["cluster_name"] = &schema.Schema{
		Type: schema.TypeString,
		Required: true,
		Description: "Cluster Name to which the config refers",
	}

	// Hotfix so that ConflictsWith works.
	// sadly that removes it but it's not supported in a TypeList
	// https://github.com/hashicorp/terraform-plugin-sdk/issues/71
	perClusterSchema["config_path"].ConflictsWith = nil

	return map[string]*schema.Schema{
		"cluster": {
			Type: schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: perClusterSchema,
			},
		},
	}
}
