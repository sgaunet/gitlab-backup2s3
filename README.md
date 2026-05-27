[![Go Report Card](https://goreportcard.com/badge/github.com/sgaunet/gitlab-backup2s3)](https://goreportcard.com/report/github.com/sgaunet/gitlab-backup2s3)
[![GitHub release](https://img.shields.io/github/release/sgaunet/gitlab-backup2s3.svg)](https://github.com/sgaunet/gitlab-backup2s3/releases/latest)
![GitHub Downloads](https://img.shields.io/github/downloads/sgaunet/gitlab-backup2s3/total)
![Test Coverage](https://raw.githubusercontent.com/wiki/sgaunet/gitlab-backup2s3/coverage-badge.svg)
[![linter](https://github.com/sgaunet/gitlab-backup2s3/actions/workflows/linter.yml/badge.svg)](https://github.com/sgaunet/gitlab-backup2s3/actions/workflows/linter.yml)
[![coverage](https://github.com/sgaunet/gitlab-backup2s3/actions/workflows/coverage.yml/badge.svg)](https://github.com/sgaunet/gitlab-backup2s3/actions/workflows/coverage.yml)
[![Snapshot Build](https://github.com/sgaunet/gitlab-backup2s3/actions/workflows/snapshot.yml/badge.svg)](https://github.com/sgaunet/gitlab-backup2s3/actions/workflows/snapshot.yml)
[![Release Build](https://github.com/sgaunet/gitlab-backup2s3/actions/workflows/release.yml/badge.svg)](https://github.com/sgaunet/gitlab-backup2s3/actions/workflows/release.yml)
[![GoDoc](https://godoc.org/github.com/sgaunet/gitlab-backup2s3?status.svg)](https://godoc.org/github.com/sgaunet/gitlab-backup2s3)
[![License](https://img.shields.io/github/license/sgaunet/gitlab-backup2s3.svg)](LICENSE)

# gitlab-backup2s3

gitlab-backup2s3 is an enhanced docker image to export gitlab projects, encrypt the archive (optional) and save them in a S3.

You can use the binary but it will need some prerequisites :

* [gitlab-backup](https://github.com/sgaunet/gitlab-backup) >= v1.19.0 (native [age](https://age-encryption.org) encryption built in)
* [gocrypt](https://github.com/sgaunet/gocrypt) >= v2.0.0 (legacy AES encryption — kept for backward compatibility)

## Archive encryption

Two encryption options are bundled in the image:

* **[age](https://age-encryption.org) (recommended)** — public-key encryption built natively into
  `gitlab-backup` since v1.19.0. Recipient public keys are configured via env vars; the matching
  private identity stays offline and is only used for restore. No shell required, no hook needed.
  See the [Configuration](#configuration) section below.
* **[gocrypt](https://github.com/sgaunet/gocrypt)** — symmetric AES, invoked as a `POSTBACKUP` hook.
  Still supported for existing setups.

The `age` CLI is also bundled in the image so you can compose custom `POSTBACKUP` pipelines if
needed, but for most users the native env-var configuration is enough.

## Restoring a backup

The [gitlab-backup](https://github.com/sgaunet/gitlab-backup) project ships a companion binary
named **`gitlab-restore`**, which restores a GitLab project backup from the archive produced by
`gitlab-backup` back into a GitLab instance.

If the archive was encrypted, decrypt it locally first:

```bash
# age (recommended) — uses your offline private identity:
age -d -i backup-key.txt -o myproject-42.tar.gz s3-downloaded.tar.gz

# gocrypt (legacy):
gocrypt dec --i archive.tar.gz
```

Then pass the plaintext archive to `gitlab-restore --archive ...`. `gitlab-restore` does not
decrypt automatically.

## Version Compatibility

⚠️ **Important Breaking Change** ⚠️

Version 2 of **gocrypt** (v2) introduced AES GCM (Galois/Counter Mode) encryption, which breaks compatibility with files encrypted using version 1 (v1).

- Files encrypted with v1 **cannot** be decrypted with v2
- Files encrypted with v2 **cannot** be decrypted with v1

This incompatibility is due to the fundamental change in the encryption mode from v1 to v2. AES GCM provides better security with authenticated encryption but requires a different format that is not backwards compatible.

Version 2 of **gocrypt** is not compatible with version 1. If you have files encrypted with v1, you will need to decrypt them using the v1 version of **gocrypt** before you can use them with v2. Version 2 of gitlab-backup2s3 uses v2 of gocrypt.
Version 1 of **gitlab-backup2s3** is compatible with version 1 of **gocrypt**. 

## Configuration

It needs some environement variables to run:

**GitLab / target selection**
* GITLAB_TOKEN
* GITALB_URI (if the endpoint differs from https://gitlab.com)
* GITLABPROJECTID: id of the project to export
* GITLABGROUPID: id of the group to export (will export all sub projects)

**Storage**
* LOCALPATH: if you want to export archives locally (let empty if you prefer to copy archives to s3)
* S3ENDPOINT: Example https://s3.eu-west-3.amazonaws.com   or http://localhost:9090 for a local minio instance
* S3REGION: region of s3
* S3BUCKETNAME
* S3BUCKETPATH
* AWS_SECRET_ACCESS_KEY: not mandatory if you associate an IAM role to the pod or ec2
* AWS_ACCESS_KEY_ID: not mandatory too

**Encryption — age (recommended)**
* AGE_RECIPIENTS: comma-separated public keys (e.g. `age1ql3z7hjy54pw3...,ssh-ed25519 AAAA... user@host`)
* AGE_RECIPIENTS_FILE: alternative — path to a recipients file (one per line, `#` comments allowed)
* AGE_ARMOR: `true` for ASCII-armored output (`false` by default, produces binary `.age` payload)

Encryption is enabled as soon as `AGE_RECIPIENTS` *or* `AGE_RECIPIENTS_FILE` is set. The archive
filename is preserved (still ends in `.tar.gz`) — only the bytes change.

**Encryption — gocrypt (legacy)**
* GOCRYPT_KEY (if you want to encrypt archives with AES)
* POSTBACKUP (if you want to encrypt archives, set it to: `gocrypt enc --i %INPUTFILE%`)

**Misc**
* DEBUGLEVEL: info by default
* TMPDIR: /tmp by default
* EXPORT_TIMEOUT_MIN: default timeout export in minutes (default "10")

### Generating an age key pair

Use the `age-keygen` CLI from any [age release](https://github.com/FiloSottile/age/releases),
**run outside the cluster** so the private key stays offline:

```bash
age-keygen -o backup-key.txt
# backup-key.txt contains:
#   # public key: age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p
#   AGE-SECRET-KEY-1...
```

Store `backup-key.txt` in a vault / password manager. Copy the `# public key:` line into
`AGE_RECIPIENTS`. For multi-recipient setups (primary + offline recovery), list both keys.

# Development

This project is using :

* Golang
* [task for development](https://taskfile.dev/)
* docker
* [docker buildx](https://github.com/docker/buildx)
* docker manifest
* [goreleaser](https://goreleaser.com/)
* [pre-commit](https://pre-commit.com/)

There are hooks executed in the precommit stage. Once the project cloned on your disk, please install pre-commit:

```
brew install pre-commit
```

Install tools:

```
task dev:install-prereq
```

And install the hooks:

```
task dev:install-pre-commit
```

If you like to launch manually the pre-commmit hook:

```
task dev:pre-commit
```

## Example of deployment

### raw kubernetes manifests

In the [deploy/k8s folder](deploy/k8s/), you will find manifests to deploy a cronjob in kubernetes.

### helm

Another github project contains the helm chart. This is [https://github.com/sgaunet/helm-gitlab-backup2s3](https://github.com/sgaunet/helm-gitlab-backup2s3), check the README.

[Configuration of the helm chart is available here.](https://github.com/sgaunet/helm-gitlab-backup2s3/blob/main/charts/gitlab-backup2s3/README.md)
