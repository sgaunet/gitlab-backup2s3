[![Go Report Card](https://goreportcard.com/badge/github.com/sgaunet/gitlab-backup2s3)](https://goreportcard.com/report/github.com/sgaunet/gitlab-backup2s3)


# gitlab-backup2s3

gitlab-backup2s3 is an enhanced docker image to export gitlab projects, encrypt the archive and save them in a S3.

You can use the binary but it will need some prerequisites :

* [gocrypt](https://github.com/sgaunet/gocrypt) >= v1.2.0
* [gitlab-backup](https://github.com/sgaunet/gitlab-backup) >= v1.0.0

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

## Example of deployment

In the deploy folder, you will find manifests to deploy a cronjob in kubernetes.