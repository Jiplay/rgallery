---
title: Install rgallery within a Kubernetes cluster
LinkTitle: Kubernetes
weight: 200
logo: /logos/k8s-logo.svg
---

# Install rgallery within a Kubernetes cluster

> Helm chart coming soon.

## Prerequisites

The following prerequisites are required to install rgallery within a Kubernetes cluster:

- [Kubernetes](https://kubernetes.io/)
- [kubectl](https://kubernetes.io/docs/reference/kubectl/)

## Install rgallery manifests

You can install rgallery in your Kubernetes cluster using the provided manifest files.

### Download individual manifest files

If you prefer not to clone the entire repository, you can download and apply each manifest individually:

Create a directory for the manifests:

```shell
mkdir -p rgallery-manifests
cd rgallery-manifests
```

Download all manifest files:

```shell
curl -O https://raw.githubusercontent.com/robbymilo/rgallery/main/production/manifests/v1.Namespace-rgallery.yaml
curl -O https://raw.githubusercontent.com/robbymilo/rgallery/main/production/manifests/apps-v1.StatefulSet-rgallery.yaml
curl -O https://raw.githubusercontent.com/robbymilo/rgallery/main/production/manifests/v1.ConfigMap-rgallery-config.yaml
curl -O https://raw.githubusercontent.com/robbymilo/rgallery/main/production/manifests/v1.PersistentVolume-rgallery-cache.yaml
curl -O https://raw.githubusercontent.com/robbymilo/rgallery/main/production/manifests/v1.PersistentVolume-rgallery-media.yaml
curl -O https://raw.githubusercontent.com/robbymilo/rgallery/main/production/manifests/v1.PersistentVolumeClaim-rgallery-cache.yaml
curl -O https://raw.githubusercontent.com/robbymilo/rgallery/main/production/manifests/v1.PersistentVolumeClaim-rgallery-media.yaml
curl -O https://raw.githubusercontent.com/robbymilo/rgallery/main/production/manifests/v1.Service-rgallery.yaml
```

Apply all manifests:

> **Note:** The manifests run rgallery as a user with UID `1000` and GID `1000`. The data, media, and cache directories must have appropriate read/write permissions for this user.

```shell
kubectl apply -f .
```

### Manifest Details

Below are the contents of each manifest file:

#### Namespace

{{< get-resource "production/manifests/v1.Namespace-rgallery.yaml" >}}

#### StatefulSet

{{< get-resource "production/manifests/apps-v1.StatefulSet-rgallery.yaml" >}}

#### ConfigMap

For config information, see [configuration file](/docs/configure/#configuration-file).

{{< get-resource "production/manifests/v1.ConfigMap-rgallery-config.yaml" >}}

#### PersistentVolume for data directory

> Note: This volume should be on a high speed SSD for best performance.

{{< get-resource "production/manifests/v1.PersistentVolume-rgallery-data.yaml" >}}

#### PersistentVolumeClaim for data directory

{{< get-resource "production/manifests/v1.PersistentVolumeClaim-rgallery-data.yaml" >}}

#### PersistentVolume for cached image thumbnails and video transcodes

{{< get-resource "production/manifests/v1.PersistentVolume-rgallery-cache.yaml" >}}

#### PersistentVolumeClaim for cached image thumbnails and video transcodes

> Note: This volume should be on a high speed SSD for best performance.

{{< get-resource "production/manifests/v1.PersistentVolumeClaim-rgallery-cache.yaml" >}}

#### PersistentVolume for media directory

> Note: This volume can be on slower storage, such as a hard drive.

{{< get-resource "production/manifests/v1.PersistentVolume-rgallery-media.yaml" >}}

#### PersistentVolumeClaim for media directory

{{< get-resource "production/manifests/v1.PersistentVolumeClaim-rgallery-media.yaml" >}}

#### Service

{{< get-resource "production/manifests/v1.Service-rgallery.yaml" >}}
