#!/bin/bash 

cd ../test/ ; docker-compose up --build -d redis-proxy
curl -X GET 'http://192.168.99.100:8080/set?key=HarryPotter&value=JKRowling' & curl -X GET 'http://192.168.99.100:8080/set?key=LordOfTheRings&value=JRTolkein'
curl -X GET 'http://192.168.99.100:8080/get?key=LordOfTheRings' & curl -X GET 'http://192.168.99.100:8080/get?key=HarryPotter'