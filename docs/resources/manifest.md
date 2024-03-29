---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "multikubernetes_manifest Resource - terraform-provider-multikubernetes"
subcategory: ""
description: |-
  
---

# multikubernetes_manifest (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cluster` (String) Cluster to which apply the resource.
- `manifest` (Dynamic) A Kubernetes manifest describing the desired state of the resource in HCL format.

### Optional

- `computed_fields` (List of String) List of manifest fields whose values can be altered by the API server during 'apply'. Defaults to: ["metadata.annotations", "metadata.labels"]
- `field_manager` (Block List, Max: 1) Configure field manager options. (see [below for nested schema](#nestedblock--field_manager))
- `object` (Dynamic) The resulting resource state, as returned by the API server after applying the desired state from `manifest`.
- `timeouts` (Block List, Max: 1) (see [below for nested schema](#nestedblock--timeouts))
- `wait` (Block List, Max: 1) Configure waiter options. (see [below for nested schema](#nestedblock--wait))
- `wait_for` (Object, Deprecated) A map of attribute paths and desired patterns to be matched. After each apply the provider will wait for all attributes listed here to reach a value that matches the desired pattern. (see [below for nested schema](#nestedatt--wait_for))

<a id="nestedblock--field_manager"></a>
### Nested Schema for `field_manager`

Optional:

- `force_conflicts` (Boolean) Force changes against conflicts.
- `name` (String) The name to use for the field manager when creating and updating the resource.


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) Timeout for the create operation.
- `delete` (String) Timeout for the delete operation.
- `update` (String) Timeout for the update operation.


<a id="nestedblock--wait"></a>
### Nested Schema for `wait`

Optional:

- `condition` (Block List) (see [below for nested schema](#nestedblock--wait--condition))
- `fields` (Map of String) A map of paths to fields to wait for a specific field value.
- `rollout` (Boolean) Wait for rollout to complete on resources that support `kubectl rollout status`.

<a id="nestedblock--wait--condition"></a>
### Nested Schema for `wait.condition`

Optional:

- `status` (String) The condition status.
- `type` (String) The type of condition.



<a id="nestedatt--wait_for"></a>
### Nested Schema for `wait_for`

Optional:

- `fields` (Map of String)
