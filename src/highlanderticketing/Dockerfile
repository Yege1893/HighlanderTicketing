FROM golang:1.20-buster
# non-go modules dependencies
RUN apt update

# copy code and protobuf
WORKDIR /go/src/app
COPY ./highlanderticketing .

RUN go mod download
RUN go install

RUN wget https://raw.githubusercontent.com/vishnubob/wait-for-it/81b1373f17855a4dc21156cfe1694c31d7d1792e/wait-for-it.sh
RUN chmod +x ./wait-for-it.sh ./docker-entrypoint.sh

ENTRYPOINT ["./docker-entrypoint.sh"]
CMD ["highlanderticketing"]

EXPOSE 8000