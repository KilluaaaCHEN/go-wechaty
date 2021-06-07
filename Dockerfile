FROM alpine:latest
COPY ./runtime/main /
WORKDIR /
CMD ["/main"]
