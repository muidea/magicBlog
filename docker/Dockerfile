FROM centos

LABEL Author="rangh"

RUN mkdir -p /var/app

WORKDIR /var/app

ADD magicBlog /var/app/magicBlog
ADD bootstrap.sh /var/app/bootstrap.sh

RUN chmod +x /var/app/magicBlog
RUN chmod +x /var/app/bootstrap.sh

ENTRYPOINT [ "/var/app/bootstrap.sh" ]

EXPOSE 8080
