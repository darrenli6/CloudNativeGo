
FROM ubuntu
MAINTAINER Darren darren94me@163.com
RUN echo '开始制作Dockerfile'
RUN mkdir -p /data/
COPY main /data/main
RUN cd /data/

EXPOSE 8081
ENTRYPOINT [ "/data/main" ]


