FROM golang:1.16

RUN apt-get update \
    && apt-get -y install locales --no-install-recommends \
    && rm -rf /var/lib/apt/lists/*

#日本語設定
RUN dpkg-reconfigure locales && \
    locale-gen C.UTF-8 && \
    /usr/sbin/update-locale LANG=C.UTF-8

ENV LANG ja_JP.UTF-8
ENV TZ Asia/Tokyo
