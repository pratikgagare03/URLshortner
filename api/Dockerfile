FROM golang:1.22.2-alpine3.19
WORKDIR /urlshortner
COPY . /urlshortner
RUN go build /urlshortner
EXPOSE 8080
ENTRYPOINT ["./urlshortner"]

