from         ubuntu:12.04
maintainer   Matt Williamson "matt@aimatt.com"

# Upstart+Docker Hack
RUN dpkg-divert --local --rename --add /sbin/initctl
RUN ln -s /bin/true /sbin/initctl

ADD . /usr/local/webhookr
ADD etc_init_webhookr.conf /etc/init/webhookr.conf

RUN apt-get update -y
RUN apt-get upgrade -y
RUN apt-get install python-setuptools libevent-dev python-all-dev build-essential -y
RUN easy_install pip
RUN pip install -r /usr/local/webhookr/REQUIREMENTS.txt

RUN restart webhookr

EXPOSE 9000

