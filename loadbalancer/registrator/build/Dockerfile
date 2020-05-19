# Build registrator executable with GO in throwaway container
FROM golang:alpine as gobuilder

# Add git used by "go get" command
RUN apk add --update git

# Create a user and run as that user instead of root
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

RUN mkdir -p /go/src/build
WORKDIR /go/src/build
COPY ./register.go .
COPY ./service.go .
RUN go get ./...
RUN go build -o registrator register.go service.go

# Build registrator image to run the executable
FROM alpine
LABEL Description="Runs a registrator which registers an haproxy instance with consul" Version="0.1"

# Add compute-weight executable from throwaway container
COPY --from=gobuilder /go/src/build/registrator /go/bin/registrator

ENTRYPOINT ["./go/bin/registrator"]