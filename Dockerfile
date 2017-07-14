FROM golang:1.8-alpine

RUN apk add --no-cache git

RUN go get github.com/constabulary/gb/...

COPY . /sensebox-mailer

WORKDIR /sensebox-mailer

RUN gb build -ldflags "-s -w" all && \
  mv /sensebox-mailer/bin/sensebox-mailer /sensebox-mailer && \
  rm -rf /sensebox-mailer/bin /sensebox-mailer/src /sensebox-mailer/pkg /sensebox-mailer/vendor && \
  rm -rf /go

RUN apk del git

EXPOSE 3924

CMD ["/sensebox-mailer/sensebox-mailer"]
