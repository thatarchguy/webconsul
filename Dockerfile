FROM alpine:latest
ADD . /go/src/github.com/thatarchguy/webconsul
WORKDIR /go/src/github.com/thatarchguy/webconsul
RUN make linux

FROM alpine:latest
RUN apk add --update curl
COPY --from=0 /go/src/github.com/thatarchguy/webconsul/webconsul /usr/local/bin/webconsul
CMD ["webconsul"]
