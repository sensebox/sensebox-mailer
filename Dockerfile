FROM golang:1.8-alpine

RUN apk add --no-cache git

RUN go get github.com/constabulary/gb/...

COPY . /sensebox-mailer

WORKDIR /sensebox-mailer

RUN gb build -f -F -ldflags "-s -w" all

# Second stage
FROM alpine

EXPOSE 3924

COPY --from=0 /sensebox-mailer/bin/sensebox-mailer /sensebox-mailer
COPY --from=0 /sensebox-mailer/templates /templates

CMD ["/sensebox-mailer"]
