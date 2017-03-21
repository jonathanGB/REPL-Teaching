FROM node:alpine
COPY run.sh /runs/
CMD ["/runs/run.sh", "js"]
