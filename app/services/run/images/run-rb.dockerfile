FROM ruby:alpine
COPY run.sh /runs/
CMD ["/runs/run.sh", "rb"]
