---
title: Backup rgallery with sqlite3
weight: 100
---

# Backup rgallery with sqlite3

Use sqlite3 to dump the database to a file:

```bash
sqlite3 data/sqlite.db ".dump" > rgallery.sql
```

You can then backup that file to a a location separate from the rgallery data directory.

## Restore from sqlite3 backup

```bash
mkdir data
touch data/sqlite.db
sqlite3 data/sqlite.db < rgallery.sql
```
