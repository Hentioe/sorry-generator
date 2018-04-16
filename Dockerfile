FROM jrottenberg/ffmpeg:3.3


ARG FILE_NAME=sorry-gen.tar.gz
ARG DIST_DIR=/data/dist
ARG TMP_DIR=/data/tmp


COPY "build/$FILE_NAME" /data/


WORKDIR /data


RUN apt-get install -y ttf-wqy-microhei \
    && mkdir -p $DIST_DIR \
    && tar -zxvf $FILE_NAME \
    && ln -s /data/sorry-gen /usr/bin/sorry-gen \
    && rm "/data/$FILE_NAME" \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/* \
    && rm -rf /var/lib/apt/lists/partial/*


EXPOSE 8080


VOLUME $TMP_DIR


ENTRYPOINT ["sorry-gen"]
