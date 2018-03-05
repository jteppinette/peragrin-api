FROM golang:1.8.3

ENV PORT 80
EXPOSE 80

COPY . /go/src/github.com/jteppinette/peragrin-api
WORKDIR /go/src/github.com/jteppinette/peragrin-api
RUN go install

ENTRYPOINT ["peragrin-api"]
CMD ["serve"]
