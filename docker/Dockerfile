FROM centos

LABEL Author="rangh"

RUN mkdir -p /var/app

WORKDIR /var/app

ADD magicBlog /var/app/magicBlog
ADD bootstrap.sh /var/app/bootstrap.sh
ADD wait-for-it.sh /var/app/wait-for-it.sh
ADD config.xml /var/app/config.xml
ADD setupTool /var/app/setupTool

RUN chmod +x /var/app/magicBlog
RUN chmod +x /var/app/bootstrap.sh
RUN chmod +x /var/app/wait-for-it.sh
RUN chmod +x /var/app/setupTool

ENTRYPOINT [ "/var/app/bootstrap.sh" ]

EXPOSE 8080
