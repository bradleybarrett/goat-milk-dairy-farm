# blue-green-deployment-poc

Build Spring App:
./gradlew build

Run Spring App (starts on default host and port - localhost:8080)
./gradlew bootrun

Run Consul container on local machine:
docker run -p 8400:8400  -p 8500:8500 -p 8600:8600/udp --name=consul consul:latest agent -server -bootstrap -ui -client=0.0.0.0

Values are read from consul:

config/application/hosts/betaService = localhost:8081

<prefix>/application/<optionalClassAnnotation>/<propertyField>
<prefix>/<spring.application.name>/<optionalClassAnnotation>/<propertyField>

default prefix = "config"

"application" properties are shared by all services.
<spring.application.name> properties take precedence over property values shared by all services.

Note, it can take some time for the KV update in consul to sync to the spring app!

Load Balancer:
build the lb docker image and run it:
1) cd into loadbalancer directory
2) docker build -t bbarrett:lb-alpha-1 .
3) docker-compose up &




