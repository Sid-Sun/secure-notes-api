FROM alpine

# Install the required packages
RUN apk add --update git go musl-dev
# Install the required dependencies
RUN go get github.com/gorilla/mux
RUN go get golang.org/x/crypto/sha3
RUN go get github.com/lib/pq
RUN go get firebase.google.com/go
# Setup the proper workdir
WORKDIR /root/go/src/secure-notes-api
# Copy indivisual files at the end to leverage caching
COPY ./LICENSE ./
COPY ./README.md ./
COPY ./*.go ./
COPY db db
RUN go build
COPY start.sh start.sh
RUN chmod +x start.sh

#Executable command needs to be static
CMD ["sh","/root/go/src/secure-notes-api/start.sh"]
