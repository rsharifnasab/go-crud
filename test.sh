#!/usr/bin/sh

curl --request POST \
    --url http://localhost:1323/customers \
    --header 'content-type: application/json' \
    --data '{
    "cName"    : "ali",
    "cTel"     : 9123012345,
    "cAddress" : "Valiasr St. ,Tehran"
}'


curl --request POST \
    --url http://localhost:1323/customers \
    --header 'content-type: application/json' \
    --data '{
    "cName"    : "taghi",
    "cTel"     : 9123012345,
    "cAddress" : "Valiasr St. ,Tehran"
}'


curl --request PUT \
    --url http://localhost:1323/customers/2 \
    --header 'content-type: application/json' \
    --data '{
    "cName"    : "hasan2",
    "cTel"     : 9123012345,
    "cAddress" : "Valiasr St. ,Tehran"
}'


curl --request GET \
    --url http://localhost:1323/customers

curl --request GET \
    --url http://localhost:1323/report/1

curl --request GET \
    --url 'http://localhost:1323/customers?cName=ha'
