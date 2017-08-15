FROM alpine:latest

RUN apk --no-cache add ca-certificates clamav clamav-libunrar \
    && mkdir /run/clamav \
    && chown clamav:clamav /run/clamav

RUN wget -O /var/lib/clamav/main.cvd http://database.clamav.net/main.cvd \
    && wget -O /var/lib/clamav/daily.cvd http://database.clamav.net/daily.cvd \
    && wget -O /var/lib/clamav/bytecode.cvd http://database.clamav.net/bytecode.cvd

ADD ./*.conf /etc/clamav/

COPY start.sh .

CMD ["/start.sh"]