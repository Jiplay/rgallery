---
title: Scanning and importing media with rgallery
LinkTitle: Scanning
weight: 200
---

# Scanning and importing media

When rgallery is run for the first time, a default scan is started. During a default scan, the media directory is walked recursively, and media files are read and persisted, and thumbnails and video transcode files are generated.

## Supported extensions

rgallery will attempt to import files with the following extensions:

Images:

- .jpg
- .jpeg
- .heic
- .gif
- .png

Videos:

- .mp4
- .mov

When a subsequent scan is started, rgallery removes references to any files in the database that are no longer on disk. If any files have been updated they are reimported and thumbnails are regenerated.

There are three types of scans:

1. Default scan - regenerates thumbnails of modified media only.
1. Metadata scan - re-imports all existing media items without recreating thumbnails.
1. Deep scan - re-imports all existing media items and recreates thumbnails.

> Only users with the role of admin can initiate scans.

## Scan Errors

If an error occurs with an image or video during scan, it is skipped on subsequent scans. If the file is updated, it will be scanned again.
