FROM ubuntu:20.04
ARG DEBIAN_FRONTEND=noninteractive
RUN apt update && apt install mysql-server -y && service mysql start && apt install git -y
RUN apk --update add --no-cache ca-certificates openssl git tzdata && update-ca-certificates
WORKDIR /root/go-app
COPY . .
RUN go env -w GOPROXY=direct GOFLAGS="-insecure" && go build -o go_auth
RUN echo "CREATE DATABASE go_auth" > file.sql && cat file.sql | mysql
CMD ["./go_auth"]
EXPOSE 8000
