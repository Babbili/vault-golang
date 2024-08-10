FROM golang:1.22.5-bookworm as build
WORKDIR /work
COPY . .
RUN go build -o getvault

FROM debian:bookworm-slim as runtime
COPY --from=build /work/getvault .

EXPOSE 5000
CMD ./getvault
