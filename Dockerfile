FROM golang:1.8
COPY       micro-user-rest /bin/micro-user-rest
ENTRYPOINT ["/bin/micro-user-rest"]
