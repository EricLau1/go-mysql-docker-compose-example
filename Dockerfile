FROM golang

WORKDIR /app

ENV GOPATH=/app

RUN go get github.com/gorilla/mux
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/jinzhu/gorm
RUN go get github.com/badoux/checkmail
RUN go get github.com/joho/godotenv
RUN go get github.com/dgrijalva/jwt-go
RUN go get golang.org/x/crypto/bcrypt

COPY . .

RUN go install main

EXPOSE 9000

ENTRYPOINT ["./bin/main"]