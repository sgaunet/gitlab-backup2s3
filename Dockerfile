FROM alpine:3.24.1 AS alpine
RUN apk add --no-cache age

FROM ghcr.io/sgaunet/gocrypt:2.0.2 AS gocrypt

FROM ghcr.io/sgaunet/gitlab-backup:1.19.1 AS gitlab-backup-image

FROM scratch
LABEL org.opencontainers.image.authors="sgaunet"
LABEL org.opencontainers.image.description="Backup gitlab projects to a S3"
LABEL org.opencontainers.image.documentation="https://github.com/sgaunet/gitlab-backup2s3"
LABEL org.opencontainers.image.source="https://github.com/sgaunet/gitlab-backup2s3"
LABEL org.opencontainers.image.licenses="MIT"

COPY --from=alpine --chown=1000:1000 /tmp /tmp
COPY resources /
COPY --from=alpine /etc/ssl /etc/ssl
COPY --from=gitlab-backup-image --chown=1000:1000 /usr/local/bin/gitlab-backup /usr/bin/gitlab-backup
COPY --from=gocrypt /gocrypt /usr/bin/gocrypt
COPY --from=alpine /usr/bin/age /usr/bin/age
COPY gitlab-backup2s3 /usr/bin/gitlab-backup2s3
# WORKDIR /usr/bin
USER gitlab-backup

# since giltab-backup 1.8.0 (to avoid duplicate timestamp in logs)
ENV NOLOGTIME=true

VOLUME [ "/tmp" ]
CMD ["/usr/bin/gitlab-backup2s3"]