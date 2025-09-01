---
title: Install rgallery with Tanka
LinkTitle: Kubernetes (Tanka)
logo: /logos/tanka.svg
weight: 201
---

# Install rgallery with Tanka

## Prerequisites

The following prerequisites are required to install rgallery within a Kubernetes cluster:

- [Kubernetes](https://kubernetes.io/)
- [kubectl](https://kubernetes.io/docs/reference/kubectl/)
- [Tanka](https://tanka.dev/)

### Setup folder structure

```shell
tk init
tk env add environments/rgallery --namespace=rgallery --server=<Kubernetes API server>
```

### Setup rgallery libary

> jb module coming soon.

```shell
mkdir -p lib/rgallery
touch lib/rgallery/rgallery.libsonnet
```

In `lib/rgallery/rgallery.libsonnet`, copy the following contents:

{{< get-resource "production/tanka/rgallery.libsonnet" "jsonnet" >}}

In `environments/rgallery/main.jsonnet` copy the following contents:

{{< get-resource "production/tanka/main.jsonnet" "jsonnet" >}}

If you want to add [lens aliases](/docs/configure/lens-aliases), create the file `environments/rgallery/files/config.yml` and copy the following contents (replace with the necessary lens values):

{{< get-resource "production/tanka/files/config.yml" "yaml" >}}

### Apply changes

> **Note:** The manifests run rgallery as a user with UID `1000` and GID `1000`. The data, media, and cache directories must have appropriate read/write permissions for this user.

```shell
tk apply environments/rgallery
```
