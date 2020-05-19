# blue-green-deployment-poc

Build Spring App:
./gradlew build

Run Spring App (starts on default host and port - localhost:8080)
./gradlew bootrun

Run Consul container on local machine:
docker run -d -p 8400:8400  -p 8500:8500 -p 8600:8600/udp --name=consul consul:latest agent -server -bootstrap -ui -client=0.0.0.0

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


TODO:
2) Try to add a user to your containers so they don't run as root
    Turns out this is really hard:
    - Gonsul: need to create a user and bind mount the ssh key into the home directory of that user 
        (currently, the process runs as root and the ssh key is mounted in the /root directory)
    - HA-Proxy: may need root access during reload when it starts and stops processes.
        Also needs write access to all haproxy config/management files.
        Refer: https://github.com/docker-library/haproxy/issues/6#issuecomment-457205041
        Refer: https://github.com/kubevirt/kubevirt/pull/270/commits/4c04de9d3392ac935e73399a78542cf41f23fbb3#diff-f0a09713727f1628ee4191f412158eeaR27
    - Registrator: this one is acrually easy! No weird permission things for this one!
3) Add a depended_on clause and wait-for-consul script to docker-compose for lb-related services
4) Read about dangers of alpine linux in article from alex




