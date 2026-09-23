FROM debian:12-slim

RUN apt-get update && \
  apt-get install -y coreutils nano curl && \
  apt-get clean && apt-get autoclean && apt-get autoremove

WORKDIR /usr/src
COPY ./ssgo .
COPY ./config.yaml .
COPY ./public ./public
COPY entrypoint.sh .

RUN chmod 777 "/usr/src/ssgo"
RUN chmod 777 "/usr/src/entrypoint.sh"

ENTRYPOINT [ "/usr/src/entrypoint.sh" ]