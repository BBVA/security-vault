FROM ubuntu

ADD ./cloudframe-security-vault /usr/bin/
RUN echo -n "ffe5c779-f23c-beac-7228-9a600a23b73f" > /etc/token

ENTRYPOINT /usr/bin/cloudframe-security-vault
