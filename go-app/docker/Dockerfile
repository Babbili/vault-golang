FROM golang:1.22.5-bookworm as build
WORKDIR /work
COPY . .
RUN go build -o getvault

FROM scratch as runtime
COPY --from=build /work/getvault .

CMD ./getvault
