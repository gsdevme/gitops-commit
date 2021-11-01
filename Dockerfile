FROM alpine:3 as builder

RUN ls /tmp/

FROM scratch as base

ENV SLACK_SIGNING_SECRET=""
ENV PORT=8080

COPY --from=builder /tmp /tmp

ENTRYPOINT ["/gitops-commit"]
COPY gitops-commit /