FROM golang:1.8.3

ENV PORT 80
EXPOSE 80

COPY . /go/src/gitlab.com/peragrin/api
WORKDIR /go/src/gitlab.com/peragrin/api
RUN go install

ENTRYPOINT ["api"]
CMD ["serve"]
