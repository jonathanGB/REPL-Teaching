FROM python:alpine
COPY run.sh /runs/
CMD ["/runs/run.sh", "py"]
