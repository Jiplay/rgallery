---
title: Admin
---

# Admin

The admin UI allows a user to initiate scans, create users, and create API keys.

## API keys

> API keys are granted admin access.

To create an API key:

1. Click the admin link in the nav.
1. Click **Add API key**.
1. Enter a name for the API key.
1. Click **Add key**.

Save the API key to your secrets manager as it will never be shown again.

## How to schedule scans

To initiate a scan daily at, for example, 00:30:

1. Create an API key
1. Create a cron job on a server, such as:

```shell
30 0 * * * curl -H 'api-key: $(API_KEY)' 'https://<replace-with-rgallery-url>/scan'
```

> You may need to restart the cron service - refer to your distribution's documentation.
