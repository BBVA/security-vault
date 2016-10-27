FROM ubuntu

RUN apt update && \
    apt install -y fuse && \
    apt clean

ADD ./cloudframe-security-vault /usr/bin/
ADD ./certs/ca.cert /tmp/ca.cert
ADD ./certs/public.key /tmp/public.key
ADD ./certs/private.key /tmp/private.key

ENTRYPOINT /usr/bin/cloudframe-security-vault
