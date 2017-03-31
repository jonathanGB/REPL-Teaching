FROM php:alpine
COPY run.sh /runs/
CMD ["/runs/run.sh", "php"]
