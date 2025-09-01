---
title: Developing rgallery
LinkTitle: Developing
weight: 800
---

# Developing rgallery

1. Ensure you have Go and Node.js installed on your system, as well as dependencies noted in [Deployment](#deployment).
1. Clone this repo.
1. Navigate to this repo.
1. Make sure you have an `./images` directory with a few test folders of images (and/or videos).
1. Run `npm i` to install frontend dependencies.
1. Run `make run`. This will start a scan of the images directory, and the application will be available at http://localhost:3000.

When editing CSS or JS, the browser will be reloaded when a file is saved. If changes are made to the application (any \*.go files), the application will need to be restarted.
