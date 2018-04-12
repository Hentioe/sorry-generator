FROM jrottenberg/ffmpeg:3.3


ARG TAG=0.1.4
ARG FILE_NAME=sorry-gen.tar.gz
ARG DL_ADDRESS="https://github.com/Hentioe/sorry-generator/releases/download/$TAG/$FILE_NAME"


ARG DIST_DIR=/data/dist


RUN apt-get install -y wget ttf-wqy-microhei \
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


ENTRYPOINT ["sorry-gen"]