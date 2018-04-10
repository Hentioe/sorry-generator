FROM ubuntu:16.04


ARG TAG=latest
ARG FILE_NAME=sorry-gen.tar.gz
ARG DL_ADDRESS="https://github.com/Hentioe/sorry-generator/releases/download/$TAG/$FILE_NAME"


ARG DIST_DIR=/data/dist


RUN apt-get update && apt-get install -y wget ffmpeg \
    && mkdir -p $DIST_DIR \
    && wget $DL_ADDRESS -O "/data/$FILE_NAME" \
    && (cd /data && tar -zxvf $FILE_NAME) \
    && ln -s /data/bin /usr/bin/sorry-gen \
    && rm "/data/$FILE_NAME" \
    && apt-get purge -y wget \
    && rm -rf /var/lib/apt/lists/* \
    && rm -rf /var/lib/apt/lists/partial/*


WORKDIR /data


EXPOSE 8080


ENTRYPOINT ["bash"]