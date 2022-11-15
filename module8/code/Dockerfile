FROM ubuntu
MAINTAINER Darren darren94me@163.com
RUN echo '开始制作Dockerfile'
RUN mkdir -p /data/
COPY main /data/main
RUN cd /data/
RUN mkdir -p /data/logs/

ENTRYPOINT [ "/data/main" ]