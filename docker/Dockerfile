FROM centos

LABEL Author="rangh"
LABEL Application="magicBlog"

RUN mkdir -p /var/app

WORKDIR /var/app

ADD magicBlog /var/app/magicBlog
ADD bootstrap.sh /var/app/bootstrap.sh
ADD wait-for-it.sh /var/app/wait-for-it.sh

RUN chmod +x /var/app/magicBlog
RUN chmod +x /var/app/bootstrap.sh
RUN chmod +x /var/app/wait-for-it.sh

ENTRYPOINT [ "/var/app/bootstrap.sh" ]

EXPOSE 8080
