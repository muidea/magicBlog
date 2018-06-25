FROM nginx

LABEL Author="rangh"
LABEL Application="magicBlogUI"


RUN mkdir -p /etc/nginx/logs

COPY magicBlogUI /tmp/
RUN cd /tmp \
 && tar -xvf magicBlogUI -C /usr/share/nginx/html/ \
 && rm -f magicBlogUI

ADD nginx.conf /etc/nginx/nginx.conf
