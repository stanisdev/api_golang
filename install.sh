#!/bin/bash
echo "Donwloading Go compiler"
wget "https://dl.google.com/go/go1.11.1.linux-amd64.tar.gz"

echo "Extracting archive"
tar -C /usr/local -xzf go1.11.1.linux-amd64.tar.gz
rm go1.11.1.linux-amd64.tar.gz

export PATH=$PATH:/usr/local/go/bin
CPWD=`pwd`
export GOPATH=$CPWD

echo "Downloading dependencies, plase wait..."
go get github.com/gin-gonic/gin
go get github.com/jinzhu/gorm
go get github.com/spf13/viper
go get github.com/asaskevich/govalidator
go get github.com/dgrijalva/jwt-go
go get golang.org/x/crypto/bcrypt
go get github.com/jinzhu/gorm/dialects/mysql

echo "Moving project"
mkdir src/app
mv main.go src/app/main.go
mv controllers src/app/controllers
mv middlewares src/app/middlewares
mv models src/app/models
mv services src/app/services
mv structures src/app/structures
mv .git src/app/.git
mv .gitignore src/app/.gitignore
mv README.md src/app/README.md

go build app