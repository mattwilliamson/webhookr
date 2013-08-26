from         ubuntu:12.04
maintainer   Matt Williamson "matt@aimatt.com"


ADD . /usr/local/webhookr

RUN apt-get update -y
RUN apt-get upgrade -y
RUN apt-get install python-setuptools libevent-dev python-all-dev build-essential -y
RUN easy_install pip
RUN pip install -r /usr/local/webhookr/REQUIREMENTS.txt

EXPOSE 9000

CMD ["/usr/bin/python", "/usr/local/webhookr/manage.py", "runserver_socketio", "0.0.0.0:9000"]
