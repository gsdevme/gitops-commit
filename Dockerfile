FROM alpine:3.15

ENV SLACK_SIGNING_SECRET=""
ENV PORT=8080
ENV SSH_KNOWN_HOSTS="/etc/ssh/ssh_known_hosts"

RUN apk add --no-cache openssh-client && ssh-keyscan -H github.com >> /etc/ssh/ssh_known_hosts

COPY entrypoint.sh .

ENTRYPOINT ["/entrypoint.sh"]
COPY gitops-commit /
