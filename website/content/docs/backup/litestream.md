---
title: Backup rgallery with Litestream
weight: 200
Draft: true
---

# Continuously backup rgallery with Litestream

It is recommended to replicate the SQLite database with [Litestream](https://litestream.io/) when deploying rgallery in a Kubernetes cluster.

Litestream enables a statefulset to be deployed to any node in a cluster as the SQLite database is restored before rgallery is started.

As rgallery scans media, or users or API keys are added, the changes to the database are streamed to SFTP or object storage.

## Prerequisites

- [Kubernetes](https://kubernetes.io/)
- [kubectl](https://kubernetes.io/docs/reference/kubectl/)
- [S3 compatible object storage such as Minio](https://github.com/minio/minio)
