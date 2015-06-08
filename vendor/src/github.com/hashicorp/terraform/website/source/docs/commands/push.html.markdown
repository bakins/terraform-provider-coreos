---
layout: "docs"
page_title: "Command: push"
sidebar_current: "docs-commands-push"
description: |-
  The `terraform push` command is used to upload the Terraform configuration to HashiCorp's Atlas service for automatically managing your infrastructure in the cloud.
---

# Command: push

The `terraform push` command uploads your Terraform configuration to
be managed by HashiCorp's [Atlas](https://atlas.hashicorp.com).
By uploading your configuration to Atlas, Atlas can automatically run
Terraform for you, will save all state transitions, will save plans,
and will keep a history of all Terraform runs.

This makes it significantly easier to use Terraform as a team: team
members modify the Terraform configurations locally and continue to
use normal version control. When the Terraform configurations are ready
to be run, they are pushed to Atlas, and any member of your team can
run Terraform with the push of a button.

Atlas can also be used to set ACLs on who can run Terraform, and a
future update of Atlas will allow parallel Terraform runs and automatically
perform infrastructure locking so only one run is modifying the same
infrastructure at a time.

## Usage

Usage: `terraform push [options] [path]`

The `path` argument is the same as for the
[apply](/docs/commands/apply.html) command.

The command-line flags are all optional. The list of available flags are:

* `-atlas-address=<url>` - An alternate address to an Atlas instance.
  Defaults to `https://atlas.hashicorp.com`.

* `-upload-modules=true` - If true (default), then the
  [modules](/docs/modules/index.html)
  being used are all locked at their current checkout and uploaded
  completely to Atlas. This prevents Atlas from running `terraform get`
  for you.

* `-name=<name>` - Name of the infrastructure configuration in Atlas.
  The format of this is: "username/name" so that you can upload
  configurations not just to your account but to other accounts and
  organizations. This setting can also be set in the configuration
  in the
  [Atlas section](/docs/configuration/atlas.html).

* `-no-color` - Disables output with coloring

* `-token=<token>` - Atlas API token to use to authorize the upload.
  If blank or unspecified, the `ATLAS_TOKEN` environmental variable
  will be used.

* `-vcs=true` - If true (default), then Terraform will detect if a VCS
  is in use, such as Git, and will only upload files that are comitted to
  version control. If no version control system is detected, Terraform will
  upload all files in `path` (parameter to the command).

## Packaged Files

The files that are uploaded and packaged with a `push` are all the
files in the `path` given as the parameter to the command, recursively.
By default (unless `-vcs=false` is specified), Terraform will automatically
detect when a VCS such as Git is being used, and in that case will only
upload the files that are comitted. Because of this built-in intelligence,
you don't have to worry about excluding folders such as ".git" or ".hg" usually.

If Terraform doesn't detect a VCS, it will upload all files.

The reason Terraform uploads all of these files is because Terraform
cannot know what is and isn't being used for provisioning, so it uploads
all the files to be safe. To exclude certain files, specify the `-exclude`
flag when pushing, or specify the `exclude` parameter in the
[Atlas configuration section](/docs/configuration/atlas.html).

## Remote State Requirement

`terraform push` requires that
[remote state](/docs/commands/remote-config.html)
is enabled. The reasoning for this is simple: `terraform push` sends your
configuration to be managed remotely. For it to keep the state in sync
and for you to be able to easily access that state, remote state must
be enabled instead of juggling local files.

While `terraform push` sends your configuration to be managed by Atlas,
the remote state backend _does not_ have to be Atlas. It can be anything
as long as it is accessible by the public internet, since Atlas will need
to be able to communicate to it.

**Warning:** The credentials for accessing the remote state will be
sent up to Atlas as well. Therefore, we recommend you use access keys
that are restricted if possible.
