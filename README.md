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

* [gocrypt](https://github.com/sgaunet/gocrypt) >= v2.0.0 (if you like to encrypt archives with AES)
* [gitlab-backup](https://github.com/sgaunet/gitlab-backup) >= v1.0.0

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

* GOCRYPT_KEY (if you want to encrypt archives)
* POSTBACKUP (if you want to encrypt archives, set it to: gocrypt enc --i %INPUTFILE% )
* GITLAB_TOKEN
* GITALB_URI (if the endpoint differs from https://gitlab.com)
* GITLABPROJECTID: id of the project to export
* GITLABGROUPID: id of the group to export (will export all sub projects)
* DEBUGLEVEL: info by default
* TMPDIR: /tmp by default
* LOCALPATH: if you want to export archives locally (let empty if you prefer to copy archives to s3)
* S3ENDPOINT: Example https://s3.eu-west-3.amazonaws.com   or http://localhost:9090 for a local minio instance
* S3REGION: region of s3
* S3BUCKETNAME
* S3BUCKETPATH
* AWS_SECRET_ACCESS_KEY: not mandatory if you associate an IAM role to the pod or ec2
* AWS_ACCESS_KEY_ID: not mandatory too

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
