From alpine:latest as builder

ARG VERSION=dev

ARG SKIP_INIT_WEB

ENV SKIP_INIT_WEB=${SKIP_INIT_WEB}

WORKDIR /synctv

COPY ./ ./

RUN apk add --no-cache bash curl git go musl-dev g++ && \
    bash script/build.sh -Mv ${VERSION}

From alpine:latest

ENV PUID=0 PGID=0 UMASK=022

COPY --from=builder /synctv/build/synctv /usr/local/bin/synctv

COPY script/entrypoint.sh /entrypoint.sh

RUN apk add --no-cache bash ca-certificates su-exec tzdata && \
    rm -rf /var/cache/apk/* && \
    chmod +x /entrypoint.sh && \
    mkdir -p /root/.synctv

WORKDIR /root/.synctv

EXPOSE 8080/tcp

VOLUME [ "/root/.synctv" ]

ENTRYPOINT [ "/entrypoint.sh" ]

CMD [ "server" ]
