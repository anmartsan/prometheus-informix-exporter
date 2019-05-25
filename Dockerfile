FROM ubuntu:bionic
#FROM golang:1.12.2-alpine3.9
#FROM golang:1.12.2-stretch
MAINTAINER Antonio Martinez Sanchez-Escalonilla  <a.martinez@scmsi.es>

RUN apt-get update && apt-get install -y vim \
	unixodbc golang-go supervisor
#RUN apk update && apk add unixodbc && rm -rf /var/cache/apk/*
ADD CSDK_4.10.FC12W1_LIN-x86_64_IFix.tar /tmp
ADD supervisord.conf /etc/supervisor/conf.d/supervisord.conf
RUN mkdir /opt/exporter
RUN /tmp/installclientsdk -i silent  -DLICENSE_ACCEPTED=TRUE && rm -rf /tmp
RUN  /bin/echo  -e  "/opt/IBM/Informix_Client-SDK/lib\n/opt/IBM/Informix_Client-SDK/lib/esql\n/opt/IBM/Informix_Client-SDK/lib/cli\n">> /etc/ld.so.conf 
RUN ldconfig
ENV INFORMIXDIR=/opt/IBM/Informix_Client-SDK/
ENV ODBCINI=/opt/exporter/odbc.ini
ENV INFORMIXSQLHOSTS=/opt/exporter/sqlhosts
RUN ln -s /opt/IBM/Informix_Client-SDK/ /opt/IDS12
VOLUME     [ "/opt/exporter" ]
EXPOSE 8080
#ENTRYPOINT [ "/opt/exporter/informixcollector", "-configfile", "/opt/exporter/config.yaml" ]
ENTRYPOINT [ "/usr/bin/supervisord" ]
