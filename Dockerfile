FROM ubuntu

ADD ./cloudframe-security-vault /usr/bin/

ENTRYPOINT /usr/bin/cloudframe-security-vault
