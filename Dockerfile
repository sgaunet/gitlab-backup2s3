FROM alpine:3.17.2 AS alpine
RUN wget -q https://github.com/sgaunet/gocrypt/releases/download/v1.2.0/gocrypt_1.2.0_linux_amd64 -O gocrypt && \
    chmod +x gocrypt

FROM sgaunet/gitlab-backup:1.2.0 AS gitlab-backup-image

FROM scratch
LABEL description="Backup gitlab projects to a S3"

COPY --from=alpine --chown=1000:1000 /tmp /tmp
COPY --from=gitlab-backup-image /etc /etc
COPY --from=gitlab-backup-image --chown=1000:1000 /usr/local/bin/gitlab-backup /usr/bin/gitlab-backup
COPY --from=alpine /gocrypt /usr/bin/gocrypt
COPY gitlab-backup2s3 /usr/bin/gitlab-backup2s3
# WORKDIR /usr/bin
USER gitlab-backup

VOLUME [ "/tmp" ]
CMD ["/usr/bin/gitlab-backup2s3"]