FROM golang:1.9 as builder

WORKDIR /go/src/sensebox-mailer

COPY . ./

RUN CGO_ENABLED=0 go install -a -tags netgo -ldflags '-extldflags "-static"'


FROM scratch

COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs

COPY --from=builder /go/bin/sensebox-mailer /sensebox-mailer
COPY --from=builder /go/src/sensebox-mailer/templates /templates
COPY --from=builder /go/src/sensebox-mailer/translations.json /translations.json

CMD ["/sensebox-mailer"]
