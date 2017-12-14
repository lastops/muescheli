FROM alpine:3.7

RUN apk --no-cache add ca-certificates wget clamav clamav-libunrar \
    && mkdir /run/clamav \
    && chown clamav:clamav /run/clamav

# Download virus definitions
RUN wget -O /var/lib/clamav/main.cvd http://database.clamav.net/main.cvd \
    && wget -O /var/lib/clamav/daily.cvd http://database.clamav.net/daily.cvd \
    && wget -O /var/lib/clamav/bytecode.cvd http://database.clamav.net/bytecode.cvd

# Download dumb-init 1.2.1
RUN wget -nv -O /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.1/dumb-init_1.2.1_amd64 && chmod 755 /usr/local/bin/dumb-init

ADD ./*.conf /etc/clamav/

# copy the needed files
COPY readiness.sh .
COPY eicar.com .
COPY start.sh .

ENTRYPOINT ["/usr/local/bin/dumb-init", "--"]

CMD ["start.sh"]