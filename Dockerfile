FROM ubuntu

RUN apt update && \
    apt install -y fuse && \
    apt clean

ADD ./cloudframe-security-vault /usr/bin/

ENTRYPOINT /usr/bin/cloudframe-security-vault
