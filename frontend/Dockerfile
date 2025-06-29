ARG NODE_IMAGE=node:12.18.0
ARG BASE_IMAGE="debian:bookworm-slim"
ARG DISTROLESS_IMAGE="gcr.io/distroless/static-debian12"

FROM ${BASE_IMAGE} AS angie-build

ENV ANGIE_VERSION=1.9.1
ENV SOURCE=https://git.angie.software/web-server/angie/archive/Angie-${ANGIE_VERSION}.tar.gz
ENV PREFIX=/opt/angie
ENV CONFIGURE_ARGS="--prefix=$PREFIX \
    --sbin-path=$PREFIX/sbin/angie \
    --modules-path=$PREFIX/lib/angie/modules \
    --conf-path=$PREFIX/etc/angie/angie.conf \
    --error-log-path=/dev/stderr \
    --http-log-path=/dev/stdout \
    --pid-path=$PREFIX/var/run/angie.pid \
    --lock-path=$PREFIX/var/run/angie.lock \
    --user=angie \
    --group=angie \
    --http-client-body-temp-path=$PREFIX/cache/client_body_temp \
    --http-proxy-temp-path=$PREFIX/cache/proxy_temp \
    --http-fastcgi-temp-path=$PREFIX/cache/fastcgi_temp \
    --http-uwsgi-temp-path=$PREFIX/cache/uwsgi_temp \
    --http-scgi-temp-path=$PREFIX/cache/scgi_temp \
    --with-file-aio \
    --with-threads \
    --with-http_addition_module \
    --with-http_auth_request_module \
    --with-http_gunzip_module \
    --with-http_gzip_static_module \
    --with-http_realip_module \
    --with-http_ssl_module \
    --with-http_stub_status_module \
    --with-http_sub_module \
    --with-http_v2_module\
    "

ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /src

RUN set -x \
    && useradd angie -u 1000 -s /usr/bin/nologin -d /var/cache/nginx \
    && apt-get update \
    && apt-get install -yq curl gcc make libpcre2-dev libssl-dev zlib1g-dev \
    && :

RUN set -x \
    && curl -fL "$SOURCE" -o angie.tar.gz \
    && tar xzvf angie.tar.gz \
    && rm angie.tar.gz \
    && :

RUN set -x \
    && cd angie \
    && ./configure $CONFIGURE_ARGS --with-cc-opt="-O2" --with-ld-opt="-s -static" \
    && make -j2 \
    && make install \
    && :

RUN set -x \
    && mkdir -p $PREFIX/cache \
    && chown -R angie $PREFIX/cache $PREFIX/var/run/ \
    && :

FROM ${DISTROLESS_IMAGE} AS webserver

USER 1000

COPY --from=angie-build /opt/ /opt/

EXPOSE 8080
STOPSIGNAL SIGTERM
CMD ["/opt/angie/sbin/angie", "-g", "daemon off;"]

FROM ${NODE_IMAGE} AS frontend

WORKDIR /code
ADD package.json .
ADD yarn.lock .
RUN yarn install
ADD . .
RUN yarn build

FROM webserver

COPY --from=frontend /code/build/ /opt/angie/html
COPY angie.conf /opt/angie/etc/angie/angie.conf
