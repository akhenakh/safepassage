FROM alpine:latest as builder

RUN apk add --update --no-cache musl-dev make git go
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN make safepassage-static

FROM gcr.io/distroless/static

WORKDIR /root/
COPY --from=builder  /build/safepassage .
ENTRYPOINT ["/root/safepassage"]
