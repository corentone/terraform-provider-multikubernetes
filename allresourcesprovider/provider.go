package allresourcesprovider

 import(
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-kubernetes/kubernetes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
 )


type Provider struct {
	singleClusterProvider *schema.Provider

	clientSets map[string]kubernetes.KubeClientsets
}

func NewProvider() *Provider {
	return &Provider{
		singleClusterProvider: kubernetes.Provider(),
		clientSets: map[string]kubernetes.KubeClientsets{},
	}
}

func (p *Provider) TFProvider() *schema.Provider {
	tfprovider := &schema.Provider{
		Schema: nil, // We don't have config for our provider just yet.
		ResourcesMap: p.buildResourcesMap(p.singleClusterProvider),
		DataSourcesMap: p.buildDataSourcesMap(p.singleClusterProvider),
	}
	return tfprovider
}

func (p *Provider) buildResourcesMap( singleClusterProvider *schema.Provider ) map[string]*schema.Resource{
	resourceMap := make(map[string]*schema.Resource, len(singleClusterProvider.ResourcesMap))
	for name, definition := range(singleClusterProvider.ResourcesMap) {
		resourceMap["multi_"+name] = p.buildNewResourceDefinition(definition)
	}
	return resourceMap
}

func (p *Provider) buildDataSourcesMap( singleClusterProvider *schema.Provider ) map[string]*schema.Resource{
	resourceMap := make(map[string]*schema.Resource, len(singleClusterProvider.DataSourcesMap))
	for name, definition := range(singleClusterProvider.DataSourcesMap) {
		resourceMap["multi_"+name] = p.buildNewResourceDefinition(definition)
	}
	return resourceMap
}

func (p *Provider) buildNewResourceDefinition(originalDefinition *schema.Resource ) *schema.Resource {
	return &schema.Resource{
		CreateContext: p.wrapOriginalActionContext(originalDefinition.CreateContext),
		ReadContext:   p.wrapOriginalActionContext(originalDefinition.ReadContext),
		UpdateContext: p.wrapOriginalActionContext(originalDefinition.UpdateContext),
		DeleteContext: p.wrapOriginalActionContext(originalDefinition.DeleteContext),
		Importer: originalDefinition.Importer,
		Timeouts: originalDefinition.Timeouts,
		SchemaVersion: originalDefinition.SchemaVersion,
		Schema: AddClusterToSchema(originalDefinition.Schema),
	}
}

func (p *Provider) wrapOriginalActionContext(originalFunc func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics) func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if originalFunc == nil {
		return nil
	}
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		clientSet, err := p.clientSetFromSchema(d)
		if err != nil {
			return diag.FromErr(err)
		}
		return originalFunc(ctx, d, clientSet)
	}
}

func (p *Provider) clientSetFromSchema(d *schema.ResourceData) (kubernetes.KubeClientsets, error) {
	clusterName := d.Get("cluster").(string)
	clientSet, found := p.clientSets[clusterName]
	if !found {
		return nil, fmt.Errorf("Cluster not Found")
	}
	return clientSet, nil
}

func AddClusterToSchema (originalSchema map[string]*schema.Schema) map[string]*schema.Schema {
	originalSchema["cluster"] = &schema.Schema{
		Type:         schema.TypeString,
		Description:  "Cluster to which apply the resource",
		ForceNew:     true,
	}
	return originalSchema
}
