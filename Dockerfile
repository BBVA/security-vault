FROM ubuntu

RUN apt update && \
    apt install -y fuse && \
    apt clean && \
    echo "private" > /tmp/private.key && \
    echo "public" > /tmp/public.key && \
    echo "cacert" > /tmp/ca.cert

ADD ./cloudframe-security-vault /usr/bin/

ENTRYPOINT /usr/bin/cloudframe-security-vault
