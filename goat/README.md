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

Build docker image with tag "latest", can specify anything as tag:
./gradlew jibDockerBuild --image=bbarrett/goat:latest

Run docker image as container on server port 8877 with version 0.0.0.3:
    version is provided for ease of demo, this would not normally be overriden.
docker run -p 8102:8101 bbarrett/goat:latest --host.ip=10.100.41.156 --host.port=8102 --consul.host=10.100.41.156 --version=0.0.0.3

Run the app as a docker container on port with version (used for demo)
./docker/run-goat.sh -p 8101 -v 0.0.0.1
