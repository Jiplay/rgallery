---
title: rgallery default account
LinkTitle: Default account
weight: 100
---

# Default account

The default user is **admin**, and the default password is **admin**.

Once logged in for the first time, create a new account. The default admin account will be deleted.

## Create a new account

1. Click the admin link in the nav.
1. Click **Add user**.
1. Enter a username and password.
1. Select **Admin** or **Viewer** role.

> Users with the **Admin** role can create users, API keys, and initiate scans.

## Remove users

Users must be removed via the command line.

To remove a user:

1. Obtain a shell connection to the rgallery deployment.
1. Run:

```bash
rgallery users delete <username>
```

## Reset all users

To remove all users and recreate the default user, run:

```bash
rgallery users reset
```
