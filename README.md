## "Goat Milk?" Dairy Farm
* [Project Overview](#1)
* [Implementation Overview](#2)
* [Deployment Pipeline Features](#3)
* [Load Balancing Implementation Patterns](#4)
   * [Populating the service registry](#4-1)
   * [Populating routing rules](#4-2)
   * [Location of load balancing decision](#4-3)
* [Implementation Patterns of Well-Known Tools](#5)
   * [Kubernetes](#5-1)
   * [Istio](#5-2)
   * [Netflix Eureka + Ribbon](#5-3)
   * [HA-Proxy + Consul + consul-template (this project)](#5-4)

## Project Overview <a name="1"></a>

The goal of this project was to learn more about load balancing tools that support rolling deployments (blue-green and canary).

The project implements a dairy farm that produces bottles of milk. Each incoming request for milk is a received by a farmer, who milks a goat, then returns a milk bottle stamped with the farmer and goat who serviced the request. 

These labels are important because they change as milk requests are loadbalanced across the farmers and goats. Also, "Goat Milk?" Dairy Farm, is all about that farm-to-table.

The diary farm has multiple instances of farmers and goats, each with a service name and version number. Ex: farmer-v1, goat-v1, goat-v2. 

Internal load balancing is implemented to allow service-level canary deployments by service name and version. A routing weight can be assigned to each service-version pair. Routing rules are automically updated as services scale up and down.

For example, introducing a new goat is easy: start the goat with the new version number and update the service-version routing weight to a non-zero value. To stop routing traffic to all services of a specific version, set the service-version routing weight to zero.

The dairy farm uses load balancing tools that can all be run in a non-orchestrated and non-cloud environment. 

This makes it easy to: test the tools, get a feel for the concepts, and see what additional functionality is desired in a full deployment pipeline.

See the section on deployment pipeline features to see what's included in diary farm implementation and what's missing.


## Implementation Overview <a name="2"></a>

* A loadbalancer for service-level, canary deployments by service name and version.
* Routing is implemented using server-side load balancing with one loadbalancer per service cluster.
* Routing rules are updated in real-time using information stored in consul: service registrations and traffic weights.
* HA-Proxy loadbalancers access information in consul and populate routing rules using consul-template.
* Traffic weights are stored in a git repository and synced to the consul kv-store in real-time using gonsul.
* HA-Proxy hot reload and consul-template polling allow routing rules to be updated automatically without dropping existing traffic.
* Registration side-car containers register each loadbalancer with consul. As a result, the loadbalancer address is flexible and client-side load balancing can be implemented to support redundant loadbalncers for each service.
* This load balancing implementation attempts to externalize all routing logic from the application services (hence the use of server-side load balancing).

Note: There a pros and cons to each load balancing implementation pattern. See the section load balancing implementation patterns for the patterns types, examples from well-known tools, and the pattern used by this dairy farm.

TODO: add a block diagram here!!
```
farmer LB       service registry        gonsul      git repo with weights
farmers         kv store: weights
goat LB
goats
```

## Deployment Pipeline Features <a name="3"></a>

Ideally, everything in this list would be automated - even the commits to git which update the deployment config!
Automation and resource management is where cloud native pipelines and orchestration tools really come in handy.

| Deployment Feature | Implementation Status    |
| :---               | :---                     |
| Provision new instances and start/stop running apps       | Manual    |
| Test the newly provisioned services                       | Manual    |
| Update deployment config in version control               | Manual    |
| Update routing based on deployment config and app status  | Automated |
| Single environment view of apps and infrstructure         | Automated<sup>1</sup> |
| Database migration on upgrade or rollback                 | Missing<sup>2</sup>   |

1. All parts of the application are not registered with consul: goats, farmers and LBs are registered, but gonsul is not.
2. Apps in the dairy farm are simple and don't have any data storage.

Note on Data Migration: 
* Data migration may need downtime and definitely needs separate tooling.
* If API changes are not backwards compatible, then you might not be able to have two versions of the app running at once.
* In this case, you need blue-green and can't go with canary.


## Load Balancing Implementation Patterns <a name="4"></a>

Definition: Load balancing - balancing traffic across a set of resources.

Load balancing implementations can be characterized by three key elements:
1. Populating the service registry
   * What resources are available?
2. Populating the routing rules
   * Which resource should be served next?
3. Location of the load balancing decision
   * Where is the load balancing decision made? (server-side or client side)

Implementation options for each load balancing element:

#### 1. Populating the service registry <a name="4-1"></a>
 * Smart orchestrator, simple clients
     - orchestrator keeps track of where it deploys client apps and checks up on their health
 * Simple orchestrator, smart clients
     - orchestrator deploys and forgets, client apps routinely register themselves

#### 2. Populating routing rules <a name="4-2"></a>
 * Locally sourced
     - rules generated by load balancing host
 * Externally sourced
     - rules served from an api to the load balancing host

#### 3. Location of load balancing decision (where to send the request) <a name="4-3"></a>
 * Server-side load balancing
     - decision made by a load balancing server
 * Client-side load balancing
     - decision can be made in app or in side-car

Existing load balancing implementations exhibit some combination of these elements.

## Implementation Patterns of Well-Known Tools (and this project) <a name="5"></a>

#### Kubernetes <a name="5-1"></a>
1. Populate service registry
    - **Smart orchestrator, simple clients**
        (kubernetes populates IP tables for the services it deploys)
2. Populate routing rules
    - **Locally sourced rules**
        (kubernetes uses round-robin for each service based on IP tables)
3. Location of load balancing decision
    - **Server-side**
        (resolves service name using IP tables - clients are unaware of load balancing)

#### Istio <a name="5-2"></a>
1. Populate service registry
    - **Simple orchestrator, smart clients**
        (envoy side-cars register and report health to the management API, forming the data plane)
2. Populate routing rules
    - **Externally sourced rules**
        (envoy management API serves rules to Envoy side-cars through the control plane)
3. Location of load balancing decision
    - **Client-side (in side-car)**
        (envoy side-car proxy routes traffic from the service)

#### Netflix Eureka + Ribbon <a name="5-3"></a>
1. Populate service registry
    - **Simple orchestrator, smart clients**
        (client apps register and report health to the eureka server)
2. Populate routing rules
    - **Either: Locally out-of-the-box, but external is possible if you provide the api**
        (app uses the Ribbon client to query eureka and select a service instance, selection and query logic can be customized with app code)
3. Location of load balancing decision
    - **Client-side (in app)**
        (app choses a service based on selection logic in the Ribbon client)

#### HA-Proxy + Consul + consul-template (this project) <a name="5-4"></a>
1. Populate service registry
    - **Simple orchestrator, smart clients**
        (client apps register and provide a health check endpoint for consul to monitor their status)
2. Populate routing rules
    - **Locally sourced rules**
        (the template used to create routing rules is part of the loadbalancer image)
3. Location of load balancing decision
    - **Server-side**
        (routing decision is made by the HA-Proxy instance sitting in front of the service cluster)
