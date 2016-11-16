FROM ubuntu

ADD ./cloudframe-security-vault /usr/bin/
RUN echo -n "86666040-1b49-35e8-5bb7-e4c323f48df3" > /etc/token

ENTRYPOINT /usr/bin/cloudframe-security-vault
