---
title: Install rgallery with Docker or Docker Compose
LinkTitle: Docker
weight: 100
logo: /logos/docker-logo.svg
---

# Install rgallery with Docker or Docker Compose

## Prerequisites

The following prerequisites are required to install rgallery with Docker or Docker Compose:

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Docker

In a terminal, configure the path to your media files and run the following command:

```shell
docker run \
  -v /path-to-your-media-files:/media:ro \
  -v ./data:/data \
  -v ./cache:/cache \
  -p 3000:3000 \
  robbymilo/rgallery:latest
```

The following volume is required for rgallery to access your media files:

- path to your media files: `/path-to-your-media-files:/media:ro`

The following volumes are required to persist the database and image thumbnail cache:

- path to the database directory: `./data:/data`
- path to the cache directory: `./cache:/cache`

If they are unchanged the data and cache directories will be created in the current directory where the command is run.

The application will be available at [http://localhost:3000](http://localhost:3000). The default username and password are both **admin**.

## Docker Compose

{{< get-resource "docker-compose.yml" >}}

and then run

```shell
docker compose up -d
```

The application will be available at [http://localhost:3000](http://localhost:3000).

With docker compose in scalable mode:

{{< get-resource "docker-compose-scalable.yml" >}}

and then run:

```shell
docker compose -d -f docker-compose-scalable.yml up --scale rgallery-resize=3
```

The application will be available at [http://localhost:3000](http://localhost:3000).
