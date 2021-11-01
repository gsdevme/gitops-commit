FROM alpine:3

ENV SLACK_SIGNING_SECRET=""
ENV PORT=8080

ENTRYPOINT ["/gitops-commit"]
COPY gitops-commit /