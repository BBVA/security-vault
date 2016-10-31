FROM ubuntu

ADD ./cloudframe-security-vault /usr/bin/
ADD ./certs/ca.cert /tmp/ca.cert
ADD ./certs/public.key /tmp/public.key
ADD ./certs/private.key /tmp/private.key

ENTRYPOINT /usr/bin/cloudframe-security-vault
