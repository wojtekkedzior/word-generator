
# FROM golang:latest
#FROM scratch

#WORKDIR /
#COPY src/word-generator/word-generator /

#ENTRYPOINT ["./word-generator youngster"]
#ENTRYPOINT ["./word-generator"]

############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

# Git is required for fetching the dependencies.
#RUN apk update && apk add --no-cache gitWORKDIR $GOPATH/src/mypackage/myapp/
#COPY . .# Fetch dependencies.# Using go get.

#RUN go get -d -v# Build the binary.
#RUN go build -o /go/bin/hello############################
# STEP 2 build a small image
############################

COPY src/word-generator/word-generator /
RUN chmod 775 /word-generator

# RUN ls /

# FROM scratch
FROM amazonlinux:latest

COPY --from=builder /word-generator /word-generator

EXPOSE 8081

CMD word-generator

# ENTRYPOINT ["/word-generator"]
# CMD [ "echo", "$PATH" ]