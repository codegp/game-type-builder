FROM python:2-onbuild

COPY sourcemanager /app/app
COPY run.sh /app/run.sh
COPY teamrunner.py /app/teamrunner.py
COPY gen-py /app/gen-py

ENV PYTHONPATH $PYTHONPATH:/usr/src/app/gen-py
ENV PYTHONPATH $PYTHONPATH:/usr/src/app/myBot

WORKDIR /usr/src/app
CMD ["./app"]
