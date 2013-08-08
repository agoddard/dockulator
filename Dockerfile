# Build an Ubuntu machine with calculators 
FROM ubuntu 
MAINTAINER Chuck Ha 

RUN apt-get install -y ruby 
ADD calculators /opt/dockulator 

FROM centos 
RUN yum install -y ruby 
ADD calculators /opt/dockulator 
