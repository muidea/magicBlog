FROM alpine:3.13.4

LABEL Author="rangh"
LABEL Application="magicBlog"

RUN mkdir -p /var/app

WORKDIR /var/app

ADD magicBlog /var/app/magicBlog
ADD bootstrap.sh /var/app/bootstrap.sh
ADD wait-for-it.sh /var/app/wait-for-it.sh

COPY static /tmp/
RUN cd /tmp \
 && tar -xvf static -C /var/app/ \
 && rm -f static


RUN chmod +x /var/app/magicBlog
RUN chmod +x /var/app/bootstrap.sh
RUN chmod +x /var/app/wait-for-it.sh

ENTRYPOINT [ "/var/app/bootstrap.sh" ]

EXPOSE 8080
