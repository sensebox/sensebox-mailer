FROM golang:1.9 as builder

WORKDIR /go/src/sensebox-mailer

COPY . ./



RUN export branch=$(git rev-parse --abbrev-ref HEAD) &&  \
  export ts=$(TZ=UTC git log --date=local --pretty=format:"%ct" -n 1) && \
  export hash=$(TZ=UTC git log --date=local --pretty=format:"%h" -n 1) && \
  CGO_ENABLED=0 go install -a -tags netgo -ldflags "-extldflags -static -X main.branch=$branch -X main.ts=$ts -X main.hash=$hash"

FROM scratch

COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs

COPY --from=builder /go/bin/sensebox-mailer /sensebox-mailer
COPY --from=builder /go/src/sensebox-mailer/templates /templates
COPY --from=builder /go/src/sensebox-mailer/translations.json /translations.json

CMD ["/sensebox-mailer"]
