FROM golang:1.7.5

ENV PORT 80
EXPOSE 80

COPY . /go/src/gitlab.com/peragrin/api
WORKDIR /go/src/gitlab.com/peragrin/api
RUN go install

CMD ["api", "serve"]
