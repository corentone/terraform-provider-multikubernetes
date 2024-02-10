---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "multikubernetes_persistent_volume Resource - terraform-provider-multikubernetes"
subcategory: ""
description: |-
  
---

# multikubernetes_persistent_volume (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cluster` (String) Cluster to which apply the resource
- `metadata` (Block List, Min: 1, Max: 1) Standard persistent volume's metadata. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#metadata (see [below for nested schema](#nestedblock--metadata))
- `spec` (Block List, Min: 1) Spec of the persistent volume owned by the cluster (see [below for nested schema](#nestedblock--spec))

### Optional

- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--metadata"></a>
### Nested Schema for `metadata`

Optional:

- `annotations` (Map of String) An unstructured key value map stored with the persistent volume that may be used to store arbitrary metadata. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
- `labels` (Map of String) Map of string keys and values that can be used to organize and categorize (scope and select) the persistent volume. May match selectors of replication controllers and services. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/
- `name` (String) Name of the persistent volume, must be unique. Cannot be updated. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names

Read-Only:

- `generation` (Number) A sequence number representing a specific generation of the desired state.
- `resource_version` (String) An opaque value that represents the internal version of this persistent volume that can be used by clients to determine when persistent volume has changed. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency
- `uid` (String) The unique in time and space value for this persistent volume. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#uids


<a id="nestedblock--spec"></a>
### Nested Schema for `spec`

Required:

- `access_modes` (Set of String) Contains all ways the volume can be mounted. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#access-modes
- `capacity` (Map of String) A description of the persistent volume's resources and capacity. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#capacity
- `persistent_volume_source` (Block List, Min: 1, Max: 1) The specification of a persistent volume. (see [below for nested schema](#nestedblock--spec--persistent_volume_source))

Optional:

- `claim_ref` (Block List, Max: 1) A reference to the persistent volume claim details for statically managed PVs. More Info: https://kubernetes.io/docs/concepts/storage/persistent-volumes/#binding (see [below for nested schema](#nestedblock--spec--claim_ref))
- `mount_options` (Set of String) A list of mount options, e.g. ["ro", "soft"]. Not validated - mount will simply fail if one is invalid.
- `node_affinity` (Block List, Max: 1) A description of the persistent volume's node affinity. More info: https://kubernetes.io/docs/concepts/storage/volumes/#local (see [below for nested schema](#nestedblock--spec--node_affinity))
- `persistent_volume_reclaim_policy` (String) What happens to a persistent volume when released from its claim. Valid options are Retain (default) and Recycle. Recycling must be supported by the volume plugin underlying this persistent volume. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#reclaiming
- `storage_class_name` (String) A description of the persistent volume's class. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes/#class
- `volume_mode` (String) Defines if a volume is intended to be used with a formatted filesystem. or to remain in raw block state.

<a id="nestedblock--spec--persistent_volume_source"></a>
### Nested Schema for `spec.persistent_volume_source`

Optional:

- `aws_elastic_block_store` (Block List, Max: 1) Represents an AWS Disk resource that is attached to a kubelet's host machine and then exposed to the pod. More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore (see [below for nested schema](#nestedblock--spec--persistent_volume_source--aws_elastic_block_store))
- `azure_disk` (Block List, Max: 1) Represents an Azure Data Disk mount on the host and bind mount to the pod. (see [below for nested schema](#nestedblock--spec--persistent_volume_source--azure_disk))
- `azure_file` (Block List, Max: 1) Represents an Azure File Service mount on the host and bind mount to the pod. (see [below for nested schema](#nestedblock--spec--persistent_volume_source--azure_file))
- `ceph_fs` (Block List, Max: 1) Represents a Ceph FS mount on the host that shares a pod's lifetime (see [below for nested schema](#nestedblock--spec--persistent_volume_source--ceph_fs))
- `cinder` (Block List, Max: 1) Represents a cinder volume attached and mounted on kubelets host machine. More info: https://examples.k8s.io/mysql-cinder-pd/README.md (see [below for nested schema](#nestedblock--spec--persistent_volume_source--cinder))
- `csi` (Block List, Max: 1) Represents a CSI Volume. More info: https://kubernetes.io/docs/concepts/storage/volumes#csi (see [below for nested schema](#nestedblock--spec--persistent_volume_source--csi))
- `fc` (Block List, Max: 1) Represents a Fibre Channel resource that is attached to a kubelet's host machine and then exposed to the pod. (see [below for nested schema](#nestedblock--spec--persistent_volume_source--fc))
- `flex_volume` (Block List, Max: 1) Represents a generic volume resource that is provisioned/attached using an exec based plugin. This is an alpha feature and may change in future. (see [below for nested schema](#nestedblock--spec--persistent_volume_source--flex_volume))
- `flocker` (Block List, Max: 1) Represents a Flocker volume attached to a kubelet's host machine and exposed to the pod for its usage. This depends on the Flocker control service being running (see [below for nested schema](#nestedblock--spec--persistent_volume_source--flocker))
- `gce_persistent_disk` (Block List, Max: 1) Represents a GCE Disk resource that is attached to a kubelet's host machine and then exposed to the pod. Provisioned by an admin. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk (see [below for nested schema](#nestedblock--spec--persistent_volume_source--gce_persistent_disk))
- `glusterfs` (Block List, Max: 1) Represents a Glusterfs volume that is attached to a host and exposed to the pod. Provisioned by an admin. More info: https://examples.k8s.io/volumes/glusterfs/README.md (see [below for nested schema](#nestedblock--spec--persistent_volume_source--glusterfs))
- `host_path` (Block List, Max: 1) Represents a directory on the host. Provisioned by a developer or tester. This is useful for single-node development and testing only! On-host storage is not supported in any way and WILL NOT WORK in a multi-node cluster. More info: https://kubernetes.io/docs/concepts/storage/volumes#hostpath (see [below for nested schema](#nestedblock--spec--persistent_volume_source--host_path))
- `iscsi` (Block List, Max: 1) Represents an ISCSI Disk resource that is attached to a kubelet's host machine and then exposed to the pod. Provisioned by an admin. (see [below for nested schema](#nestedblock--spec--persistent_volume_source--iscsi))
- `local` (Block List, Max: 1) Represents a mounted local storage device such as a disk, partition or directory. Local volumes can only be used as a statically created PersistentVolume. Dynamic provisioning is not supported yet. More info: https://kubernetes.io/docs/concepts/storage/volumes#local (see [below for nested schema](#nestedblock--spec--persistent_volume_source--local))
- `nfs` (Block List, Max: 1) Represents an NFS mount on the host. Provisioned by an admin. More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs (see [below for nested schema](#nestedblock--spec--persistent_volume_source--nfs))
- `photon_persistent_disk` (Block List, Max: 1) Represents a PhotonController persistent disk attached and mounted on kubelets host machine (see [below for nested schema](#nestedblock--spec--persistent_volume_source--photon_persistent_disk))
- `quobyte` (Block List, Max: 1) Quobyte represents a Quobyte mount on the host that shares a pod's lifetime (see [below for nested schema](#nestedblock--spec--persistent_volume_source--quobyte))
- `rbd` (Block List, Max: 1) Represents a Rados Block Device mount on the host that shares a pod's lifetime. More info: https://examples.k8s.io/volumes/rbd/README.md (see [below for nested schema](#nestedblock--spec--persistent_volume_source--rbd))
- `vsphere_volume` (Block List, Max: 1) Represents a vSphere volume attached and mounted on kubelets host machine (see [below for nested schema](#nestedblock--spec--persistent_volume_source--vsphere_volume))

<a id="nestedblock--spec--persistent_volume_source--aws_elastic_block_store"></a>
### Nested Schema for `spec.persistent_volume_source.aws_elastic_block_store`

Required:

- `volume_id` (String) Unique ID of the persistent disk resource in AWS (Amazon EBS volume). More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore

Optional:

- `fs_type` (String) Filesystem type of the volume that you want to mount. Tip: Ensure that the filesystem type is supported by the host operating system. Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified. More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore
- `partition` (Number) The partition in the volume that you want to mount. If omitted, the default is to mount by volume name. Examples: For volume /dev/sda1, you specify the partition as "1". Similarly, the volume partition for /dev/sda is "0" (or you can leave the property empty).
- `read_only` (Boolean) Whether to set the read-only property in VolumeMounts to "true". If omitted, the default is "false". More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore


<a id="nestedblock--spec--persistent_volume_source--azure_disk"></a>
### Nested Schema for `spec.persistent_volume_source.azure_disk`

Required:

- `caching_mode` (String) Host Caching mode: None, Read Only, Read Write.
- `data_disk_uri` (String) The URI the data disk in the blob storage
- `disk_name` (String) The Name of the data disk in the blob storage

Optional:

- `fs_type` (String) Filesystem type to mount. Must be a filesystem type supported by the host operating system. Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
- `kind` (String) The type for the data disk. Expected values: Shared, Dedicated, Managed. Defaults to Shared
- `read_only` (Boolean) Whether to force the read-only setting in VolumeMounts. Defaults to false (read/write).


<a id="nestedblock--spec--persistent_volume_source--azure_file"></a>
### Nested Schema for `spec.persistent_volume_source.azure_file`

Required:

- `secret_name` (String) The name of secret that contains Azure Storage Account Name and Key
- `share_name` (String) Share Name

Optional:

- `read_only` (Boolean) Whether to force the read-only setting in VolumeMounts. Defaults to false (read/write).
- `secret_namespace` (String) The namespace of the secret that contains Azure Storage Account Name and Key. For Kubernetes up to 1.18.x the default is the same as the Pod. For Kubernetes 1.19.x and later the default is "default" namespace.


<a id="nestedblock--spec--persistent_volume_source--ceph_fs"></a>
### Nested Schema for `spec.persistent_volume_source.ceph_fs`

Required:

- `monitors` (Set of String) Monitors is a collection of Ceph monitors. More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it

Optional:

- `path` (String) Used as the mounted root, rather than the full Ceph tree, default is /
- `read_only` (Boolean) Whether to force the read-only setting in VolumeMounts. Defaults to `false` (read/write). More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it
- `secret_file` (String) The path to key ring for User, default is `/etc/ceph/user.secret`. More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it
- `secret_ref` (Block List, Max: 1) Reference to the authentication secret for User, default is empty. More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it (see [below for nested schema](#nestedblock--spec--persistent_volume_source--ceph_fs--secret_ref))
- `user` (String) User is the rados user name, default is admin. More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it

<a id="nestedblock--spec--persistent_volume_source--ceph_fs--secret_ref"></a>
### Nested Schema for `spec.persistent_volume_source.ceph_fs.secret_ref`

Optional:

- `name` (String) Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
- `namespace` (String) Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names



<a id="nestedblock--spec--persistent_volume_source--cinder"></a>
### Nested Schema for `spec.persistent_volume_source.cinder`

Required:

- `volume_id` (String) Volume ID used to identify the volume in Cinder. More info: https://examples.k8s.io/mysql-cinder-pd/README.md

Optional:

- `fs_type` (String) Filesystem type to mount. Must be a filesystem type supported by the host operating system. Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified. More info: https://examples.k8s.io/mysql-cinder-pd/README.md
- `read_only` (Boolean) Whether to force the read-only setting in VolumeMounts. Defaults to false (read/write). More info: https://examples.k8s.io/mysql-cinder-pd/README.md


<a id="nestedblock--spec--persistent_volume_source--csi"></a>
### Nested Schema for `spec.persistent_volume_source.csi`

Required:

- `driver` (String) the name of the volume driver to use. More info: https://kubernetes.io/docs/concepts/storage/volumes/#csi
- `volume_handle` (String) A string value that uniquely identifies the volume. More info: https://kubernetes.io/docs/concepts/storage/volumes/#csi

Optional:

- `controller_expand_secret_ref` (Block List, Max: 1) A reference to the secret object containing sensitive information to pass to the CSI driver to complete the CSI ControllerExpandVolume call. (see [below for nested schema](#nestedblock--spec--persistent_volume_source--csi--controller_expand_secret_ref))
- `controller_publish_secret_ref` (Block List, Max: 1) A reference to the secret object containing sensitive information to pass to the CSI driver to complete the CSI ControllerPublishVolume and ControllerUnpublishVolume calls. (see [below for nested schema](#nestedblock--spec--persistent_volume_source--csi--controller_publish_secret_ref))
- `fs_type` (String) Filesystem type to mount. Must be a filesystem type supported by the host operating system. Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
- `node_publish_secret_ref` (Block List, Max: 1) A reference to the secret object containing sensitive information to pass to the CSI driver to complete the CSI NodePublishVolume and NodeUnpublishVolume calls. (see [below for nested schema](#nestedblock--spec--persistent_volume_source--csi--node_publish_secret_ref))
- `node_stage_secret_ref` (Block List, Max: 1) A reference to the secret object containing sensitive information to pass to the CSI driver to complete the CSI NodeStageVolume and NodeStageVolume and NodeUnstageVolume calls. (see [below for nested schema](#nestedblock--spec--persistent_volume_source--csi--node_stage_secret_ref))
- `read_only` (Boolean) Whether to set the read-only property in VolumeMounts to "true". If omitted, the default is "false". More info: https://kubernetes.io/docs/concepts/storage/volumes#csi
- `volume_attributes` (Map of String) Attributes of the volume to publish.

<a id="nestedblock--spec--persistent_volume_source--csi--controller_expand_secret_ref"></a>
### Nested Schema for `spec.persistent_volume_source.csi.controller_expand_secret_ref`

Optional:

- `name` (String) Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
- `namespace` (String) Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names


<a id="nestedblock--spec--persistent_volume_source--csi--controller_publish_secret_ref"></a>
### Nested Schema for `spec.persistent_volume_source.csi.controller_publish_secret_ref`

Optional:

- `name` (String) Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
- `namespace` (String) Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names


<a id="nestedblock--spec--persistent_volume_source--csi--node_publish_secret_ref"></a>
### Nested Schema for `spec.persistent_volume_source.csi.node_publish_secret_ref`

Optional:

- `name` (String) Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
- `namespace` (String) Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names


<a id="nestedblock--spec--persistent_volume_source--csi--node_stage_secret_ref"></a>
### Nested Schema for `spec.persistent_volume_source.csi.node_stage_secret_ref`

Optional:

- `name` (String) Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
- `namespace` (String) Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names



<a id="nestedblock--spec--persistent_volume_source--fc"></a>
### Nested Schema for `spec.persistent_volume_source.fc`

Required:

- `lun` (Number) FC target lun number
- `target_ww_ns` (Set of String) FC target worldwide names (WWNs)

Optional:

- `fs_type` (String) Filesystem type to mount. Must be a filesystem type supported by the host operating system. Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
- `read_only` (Boolean) Whether to force the read-only setting in VolumeMounts. Defaults to false (read/write).


<a id="nestedblock--spec--persistent_volume_source--flex_volume"></a>
### Nested Schema for `spec.persistent_volume_source.flex_volume`

Required:

- `driver` (String) Driver is the name of the driver to use for this volume.

Optional:

- `fs_type` (String) Filesystem type to mount. Must be a filesystem type supported by the host operating system. Ex. "ext4", "xfs", "ntfs". The default filesystem depends on FlexVolume script.
- `options` (Map of String) Extra command options if any.
- `read_only` (Boolean) Whether to force the ReadOnly setting in VolumeMounts. Defaults to false (read/write).
- `secret_ref` (Block List, Max: 1) Reference to the secret object containing sensitive information to pass to the plugin scripts. This may be empty if no secret object is specified. If the secret object contains more than one secret, all secrets are passed to the plugin scripts. (see [below for nested schema](#nestedblock--spec--persistent_volume_source--flex_volume--secret_ref))

<a id="nestedblock--spec--persistent_volume_source--flex_volume--secret_ref"></a>
### Nested Schema for `spec.persistent_volume_source.flex_volume.secret_ref`

Optional:

- `name` (String) Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
- `namespace` (String) Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names



<a id="nestedblock--spec--persistent_volume_source--flocker"></a>
### Nested Schema for `spec.persistent_volume_source.flocker`

Optional:

- `dataset_name` (String) Name of the dataset stored as metadata -> name on the dataset for Flocker should be considered as deprecated
- `dataset_uuid` (String) UUID of the dataset. This is unique identifier of a Flocker dataset


<a id="nestedblock--spec--persistent_volume_source--gce_persistent_disk"></a>
### Nested Schema for `spec.persistent_volume_source.gce_persistent_disk`

Required:

- `pd_name` (String) Unique name of the PD resource in GCE. Used to identify the disk in GCE. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk

Optional:

- `fs_type` (String) Filesystem type of the volume that you want to mount. Tip: Ensure that the filesystem type is supported by the host operating system. Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
- `partition` (Number) The partition in the volume that you want to mount. If omitted, the default is to mount by volume name. Examples: For volume /dev/sda1, you specify the partition as "1". Similarly, the volume partition for /dev/sda is "0" (or you can leave the property empty). More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
- `read_only` (Boolean) Whether to force the ReadOnly setting in VolumeMounts. Defaults to false. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk


<a id="nestedblock--spec--persistent_volume_source--glusterfs"></a>
### Nested Schema for `spec.persistent_volume_source.glusterfs`

Required:

- `endpoints_name` (String) The endpoint name that details Glusterfs topology. More info: https://examples.k8s.io/volumes/glusterfs/README.md#create-a-pod
- `path` (String) The Glusterfs volume path. More info: https://examples.k8s.io/volumes/glusterfs/README.md#create-a-pod

Optional:

- `read_only` (Boolean) Whether to force the Glusterfs volume to be mounted with read-only permissions. Defaults to false. More info: https://examples.k8s.io/volumes/glusterfs/README.md#create-a-pod


<a id="nestedblock--spec--persistent_volume_source--host_path"></a>
### Nested Schema for `spec.persistent_volume_source.host_path`

Optional:

- `path` (String) Path of the directory on the host. More info: https://kubernetes.io/docs/concepts/storage/volumes#hostpath
- `type` (String) Type for HostPath volume. Allowed values are "" (default), DirectoryOrCreate, Directory, FileOrCreate, File, Socket, CharDevice and BlockDevice


<a id="nestedblock--spec--persistent_volume_source--iscsi"></a>
### Nested Schema for `spec.persistent_volume_source.iscsi`

Required:

- `iqn` (String) Target iSCSI Qualified Name.
- `target_portal` (String) iSCSI target portal. The portal is either an IP or ip_addr:port if the port is other than default (typically TCP ports 860 and 3260).

Optional:

- `fs_type` (String) Filesystem type of the volume that you want to mount. Tip: Ensure that the filesystem type is supported by the host operating system. Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified. More info: https://kubernetes.io/docs/concepts/storage/volumes#iscsi
- `iscsi_interface` (String) iSCSI interface name that uses an iSCSI transport. Defaults to 'default' (tcp).
- `lun` (Number) iSCSI target lun number.
- `read_only` (Boolean) Whether to force the read-only setting in VolumeMounts. Defaults to false.


<a id="nestedblock--spec--persistent_volume_source--local"></a>
### Nested Schema for `spec.persistent_volume_source.local`

Optional:

- `path` (String) Path of the directory on the host. More info: https://kubernetes.io/docs/concepts/storage/volumes#local


<a id="nestedblock--spec--persistent_volume_source--nfs"></a>
### Nested Schema for `spec.persistent_volume_source.nfs`

Required:

- `path` (String) Path that is exported by the NFS server. More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs
- `server` (String) Server is the hostname or IP address of the NFS server. More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs

Optional:

- `read_only` (Boolean) Whether to force the NFS export to be mounted with read-only permissions. Defaults to false. More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs


<a id="nestedblock--spec--persistent_volume_source--photon_persistent_disk"></a>
### Nested Schema for `spec.persistent_volume_source.photon_persistent_disk`

Required:

- `pd_id` (String) ID that identifies Photon Controller persistent disk

Optional:

- `fs_type` (String) Filesystem type to mount. Must be a filesystem type supported by the host operating system. Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.


<a id="nestedblock--spec--persistent_volume_source--quobyte"></a>
### Nested Schema for `spec.persistent_volume_source.quobyte`

Required:

- `registry` (String) Registry represents a single or multiple Quobyte Registry services specified as a string as host:port pair (multiple entries are separated with commas) which acts as the central registry for volumes
- `volume` (String) Volume is a string that references an already created Quobyte volume by name.

Optional:

- `group` (String) Group to map volume access to Default is no group
- `read_only` (Boolean) Whether to force the Quobyte volume to be mounted with read-only permissions. Defaults to false.
- `user` (String) User to map volume access to Defaults to serivceaccount user


<a id="nestedblock--spec--persistent_volume_source--rbd"></a>
### Nested Schema for `spec.persistent_volume_source.rbd`

Required:

- `ceph_monitors` (Set of String) A collection of Ceph monitors. More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
- `rbd_image` (String) The rados image name. More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it

Optional:

- `fs_type` (String) Filesystem type of the volume that you want to mount. Tip: Ensure that the filesystem type is supported by the host operating system. Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified. More info: https://kubernetes.io/docs/concepts/storage/volumes#rbd
- `keyring` (String) Keyring is the path to key ring for RBDUser. Default is /etc/ceph/keyring. More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
- `rados_user` (String) The rados user name. Default is admin. More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
- `rbd_pool` (String) The rados pool name. Default is rbd. More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it.
- `read_only` (Boolean) Whether to force the read-only setting in VolumeMounts. Defaults to false. More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
- `secret_ref` (Block List, Max: 1) Name of the authentication secret for RBDUser. If provided overrides keyring. Default is nil. More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it (see [below for nested schema](#nestedblock--spec--persistent_volume_source--rbd--secret_ref))

<a id="nestedblock--spec--persistent_volume_source--rbd--secret_ref"></a>
### Nested Schema for `spec.persistent_volume_source.rbd.secret_ref`

Optional:

- `name` (String) Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
- `namespace` (String) Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names



<a id="nestedblock--spec--persistent_volume_source--vsphere_volume"></a>
### Nested Schema for `spec.persistent_volume_source.vsphere_volume`

Required:

- `volume_path` (String) Path that identifies vSphere volume vmdk

Optional:

- `fs_type` (String) Filesystem type to mount. Must be a filesystem type supported by the host operating system. Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.



<a id="nestedblock--spec--claim_ref"></a>
### Nested Schema for `spec.claim_ref`

Required:

- `name` (String) The name of the PersistentVolumeClaim

Optional:

- `namespace` (String) The namespace of the PersistentVolumeClaim. Uses 'default' namespace if none is specified.


<a id="nestedblock--spec--node_affinity"></a>
### Nested Schema for `spec.node_affinity`

Optional:

- `required` (Block List, Max: 1) (see [below for nested schema](#nestedblock--spec--node_affinity--required))

<a id="nestedblock--spec--node_affinity--required"></a>
### Nested Schema for `spec.node_affinity.required`

Required:

- `node_selector_term` (Block List, Min: 1) (see [below for nested schema](#nestedblock--spec--node_affinity--required--node_selector_term))

<a id="nestedblock--spec--node_affinity--required--node_selector_term"></a>
### Nested Schema for `spec.node_affinity.required.node_selector_term`

Optional:

- `match_expressions` (Block List) A list of node selector requirements by node's labels. The requirements are ANDed. (see [below for nested schema](#nestedblock--spec--node_affinity--required--node_selector_term--match_expressions))
- `match_fields` (Block List) A list of node selector requirements by node's fields. The requirements are ANDed. (see [below for nested schema](#nestedblock--spec--node_affinity--required--node_selector_term--match_fields))

<a id="nestedblock--spec--node_affinity--required--node_selector_term--match_expressions"></a>
### Nested Schema for `spec.node_affinity.required.node_selector_term.match_expressions`

Required:

- `key` (String) The label key that the selector applies to.
- `operator` (String) A key's relationship to a set of values. Valid operators ard `In`, `NotIn`, `Exists`, `DoesNotExist`, `Gt`, and `Lt`.

Optional:

- `values` (Set of String) An array of string values. If the operator is `In` or `NotIn`, the values array must be non-empty. If the operator is `Exists` or `DoesNotExist`, the values array must be empty. This array is replaced during a strategic merge patch.


<a id="nestedblock--spec--node_affinity--required--node_selector_term--match_fields"></a>
### Nested Schema for `spec.node_affinity.required.node_selector_term.match_fields`

Required:

- `key` (String) The label key that the selector applies to.
- `operator` (String) A key's relationship to a set of values. Valid operators ard `In`, `NotIn`, `Exists`, `DoesNotExist`, `Gt`, and `Lt`.

Optional:

- `values` (Set of String) An array of string values. If the operator is `In` or `NotIn`, the values array must be non-empty. If the operator is `Exists` or `DoesNotExist`, the values array must be empty. This array is replaced during a strategic merge patch.






<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)