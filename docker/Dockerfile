FROM alpine:3.7

RUN apk add --update --no-cache ca-certificates wget

# TEMPORARY: use edge repositories to get newest clam (to fix security issue)
RUN apk add --update --no-cache --repository http://dl-4.alpinelinux.org/alpine/edge/main clamav clamav-libunrar \
    && mkdir /run/clamav \
    && chown clamav:clamav /run/clamav

# download virus definitions
RUN wget -O /var/lib/clamav/main.cvd http://database.clamav.net/main.cvd \
    && wget -O /var/lib/clamav/daily.cvd http://database.clamav.net/daily.cvd \
    && wget -O /var/lib/clamav/bytecode.cvd http://database.clamav.net/bytecode.cvd

# download dumb-init 1.2.1
RUN wget -nv -O /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.1/dumb-init_1.2.1_amd64 && chmod 755 /usr/local/bin/dumb-init

ADD ./*.conf /etc/clamav/

# setup environment
ENV PATH "/app:${PATH}"
WORKDIR /app/
USER clamav

# copy the needed files
COPY eicar.com .
COPY *.sh /app/

ENTRYPOINT ["/usr/local/bin/dumb-init", "--"]

CMD ["start.sh"]