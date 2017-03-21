FROM golang:alpine
COPY run.sh /runs/
CMD ["/runs/run.sh", "go"]
