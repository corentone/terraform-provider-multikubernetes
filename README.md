# terraform-provider-multikubernetes

The multikubernetes Terraform provider is a shim over the kubernetes provider but supporting multiple clusters.
This provider is useful when a resource must be applied to a bunch of clusters with similar config. This is currently not possible
with the existing kubernetes provider due to the inflexibility of Terraform providers.
With this provider, one can create a few clusters and easily apply a resource to all of them and leverage for_each and other
iterative abilities of Terraform.

# Build & Install
Build with golang directly:
```sh
go build -o terraform-provider-multikubernetes
```

Installation must be done to your plugin folder, on linux, this is:
`~/.terraform.d/plugins/local/corentone/multikubernetes/1.0.0/linux_amd64/`

I recommend making a link to the binary directly so you can potentially modify the code and recompile.
```
mkdir -p ~/.terraform.d/plugins/local/corentone/multikubernetes/1.0.0/linux_amd64/
ln -s $PWD/terraform-provider-multikubernetes ~/.terraform.d/plugins/local/corentone/multikubernetes/1.0.0/linux_amd64/terraform-provider-multikubernetes
```

If you recompile the plugin, you need to remove `$PWD/.terraform.lock.hcl` as the signature of the binary will have changed.

# Usage

* Predefined resources from original Kubernetes Provider are supported, with an additional required parameter called `cluster`, their original name is prefixed with `multi`.
* Generic `kubernetes_manifest` resource is supported as `multikubernetes_manifest` and has an additional required parameter called `cluster`.
* Provider configuration is similar to the original Kubernetes Provider, but enclosed into a block called `cluster` and with a required field called `cluster_name`.

## Example

```
terraform {
  required_providers {
    multikubernetes = {
      #version = "~> 1.0.0"
      source  = "local/corentone/multikubernetes"
    }
  }
}

locals {
  connect_gateway_project = "1234" # project number
  clusters = {
   "quick-cluster" = "us-central1"
   "cluster-1" = "global"
  }
}

# This leverages Google's Connect Gateway with GKE clusters
provider multikubernetes {
  dynamic "cluster" {
    for_each = local.clusters
    content {
      cluster_name = cluster.key

      host  = "https://connectgateway.googleapis.com/v1/projects/${local.connect_gateway_project}/locations/${cluster.value}/gkeMemberships/${cluster.key}"
      exec {
        api_version = "client.authentication.k8s.io/v1beta1"
        command     = "gke-gcloud-auth-plugin"
      }
    }
  }
}

# Example of a predefined resource
resource "multikubernetes_namespace" "my-ns" {
  for_each = local.clusters
  cluster = each.key
  metadata {
    name = "my-ns"
  }
}

# Example of generic multikubernete_manifest
resource "multikubernetes_manifest" "test-configmap" {
  for_each = local.clusters
  cluster = each.key
  manifest = {
    "apiVersion" = "v1"
    "kind"       = "ConfigMap"
    "metadata" = {
      "name"      = "test-config"
      "namespace" = "default"
    }
    "data" = {
      "foo" = "bar"
    }
  }
}
```

# Disclaimer

This plugin is really a POC and written hastily to show multicluster examples.


There are current feature limitations:

* For resource `multikubernetes_manifest`, Import is not supported.
* Despite being defined, data sources `multikubernetes_resource` and `multikubernetes_resources` are not supported.
* Provider configuration leveraging environment variables doesn't work as it wouldn't be able to be mapped to a particular cluster. If environment variables are set, they would affect all clusters defined in the provider config.


Other general warnings:

* There is no automated testing besides limited manual testing.
* Plugin is not published to a registry so one has to import it locally.
* Plugin relies on a fork I made with very minor edits to github.com/hashicorp/terraform-provider-kubernetes as of 02/05/2024 (`408d1a9f35fb0b73530b367c8aa7f352b0b20070`)
* LICENSE is currently undefined as I'm not sure what it'd be. So for now assume Open source with restrictive copyleft. (The intent is to make it friendlier, but it's not clear yet which one.)
* On errors, you may have weird messages. The plugin doesn't properly use diagnostics/errors and will sometimes use one and sometimes the other.
