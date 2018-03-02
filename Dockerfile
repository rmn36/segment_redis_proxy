FROM golang:1.9

RUN apt-get update -qq && apt-get install -y git

WORKDIR $GOPATH/src/redis-proxy
COPY ./src/redis-proxy .

RUN go-wrapper download   
RUN go-wrapper install  

CMD ["go-wrapper", "run"] # ["redis-proxy"]