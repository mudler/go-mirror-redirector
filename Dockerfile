
ARG BASE_IMAGE=scratch

FROM golang:1.14 as build
ADD . /go-mirror-redirect
WORKDIR /go-mirror-redirect
RUN CGO_ENABLED=0 go build -o go-mirror-redirect
RUN chmod +x go-mirror-redirect

FROM $BASE_IMAGE
COPY --from=build /go-mirror-redirect/go-mirror-redirect /bin/
ENTRYPOINT ["/bin/go-mirror-redirect"]