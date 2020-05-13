FROM alpine

# Install the required packages
RUN apk add --update git go musl-dev
# Install the required dependencies
RUN go get github.com/gorilla/mux
RUN go get golang.org/x/crypto/sha3
RUN go get github.com/lib/pq
# Setup the proper workdir
WORKDIR /root/api
# Copy indivisual files at the end to leverage caching
COPY ./LICENSE ./
COPY ./README.md ./
COPY ./*.go ./
RUN go build

#Executable command needs to be static
CMD ["/root/api/api"]
