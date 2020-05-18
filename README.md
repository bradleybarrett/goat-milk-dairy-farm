## "Goat Milk?" Dairy Farm

## Project Overview

The goal of this project was to learn more about loadbalancing tools that support rolling deployments (blue-green and canary).

The project implements a dairy farm that produces bottles of milk. Each incoming request for milk is a received by a farmer, who milks a goat, then returns a milk bottle stamped with the farmer and goat who serviced the request. 

These labels are important because they change as milk requests are loadbalanced across the farmers and goats. Also, "Goat Milk?" Dairy Farm, is all about that farm-to-table.

The diary farm has multiple instances of farmers and goats, each with a service name and version number. Ex: farmer-v1, goat-v1, goat-v2. Internal loadbalancing is implemented to allow service-level canary deployments by service name and version. A routing weight can be assigned to each service-version pair. Routing rules are automically updated when services start and stop.

For example, introducing a new goat is easy: start the goat with the new version number and update the service-version routing weight to a non-zero value. To stop routing traffic to all services of a specific version, set the service-version routing weight to zero.

The dairy farm uses loadbalancing tools that can all be run in a non-orchestrated and non-cloud environment. 

This makes it easy to: test the tools, get a feel for the concepts, and see what additional functionality is desired in a full deployment pipeline.

See the section on deployment pipeline features to see what's included in diary farm implementation and what's missing.


## Implementation Overview

* A loadbalancer for service-level, canary deployments by service name and version.
* Routing is implemented using server-side loadbalancing with one loadbalancer per service cluster.
* Routing rules are updated in real-time using information stored in consul: service registrations and traffic weights.
* HA-Proxy loadbalancers access information in consul and populate routing rules using consul-template.
* Traffic weights are stored in a git repository and synced to the consul kv-store in real-time using gonsul.
* HA-Proxy hot reload and consul-template polling allow routing rules to be updated automatically without dropping existing traffic.
* Registration side-car containers register each loadbalancer with consul. As a result, the loadbalancer address is flexible and client-side loadbalancing can be implemented to support redundant loadbalncers for each service.
* This loadbalancing implementation attempts to externalize all routing logic from the application services (hence the use of server-side loadbalancing).

Note: There a pros and cons to each loadbalancing implementation pattern. See the section loadbalancing implementation patterns for the patterns types, examples from well-known tools, and the pattern used for the dairy farm.

TODO: add a block diagram here!!
```
farmer LB       service registry        gonsul      git repo with weights
farmers         kv store: weights
goat LB
goats
```

## Deployment Pipeline Features

Ideally, everything in this list should be automated - even the commits to git which update the deployment config!
Automation and resource management is where cloud native tools/pipelines and orchestration tools really come in handy.

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


## Loadbalancing Implementation Patterns

#### Three key elements of loadbalancing

Loadbalancing implementation patterns can be characterized by three key elements:
1. Populating the service registry
2. Populating the routing rules
3. Location of the loadbalancing decision

Definition: Load balancing - balancing traffic across a set of resources.

Two pieces of information are need to make a loadbalancing decision:

1. What resources are available? - service registry
2. Which resource should be served next? - routing rules

Once we have all the information, the last question is:

3. Where is the loadbalancing decision made? - location of the decision (server-side or client side)

Implementation options for each loadbalancing element

1. **Populating the service registry**
    * Smart orchestrator, simple clients
        - orchestrator keeps track of where it deploys client apps and checks up on their health
    * Simple orchestrator, smart clients
        - orchestrator deploys and forgets, client apps routinely register themselves

2. **Populating routing rules**
    * Locally sourced
        - rules generated by loadbalancing host
    * Externally sourced
        - rules served from an api to the loadbalancing host

3. **Location of loadbalancing decision** (where to send the request)
    * Server-side loadbalancing
        - decision made by a load balancing server
    * Client-side loadbalancing
        - decision can be made in app or in side-car

Existing loadbalancing implementations exhibit some combination of these elements.

#### Implementation patterns of well-known tools (and this project)

#### Kubernetes
1. Populate service registry
    - **Smart orchestrator, simple clients**
        (kubernetes populates IP tables for the services it deploys)
2. Populate routing rules
    - **Locally sourced rules**
        (kubernetes uses round-robin for each service based on IP tables)
3. Location of loadbalancing decision
    - **Server-side**
        (resolves service name using IP tables - clients are unaware of loadbalancing)

#### Istio
1. Populate service registry
    - **Simple orchestrator, smart clients**
        (envoy side-cars register and report health to the management API, forming the data plane)
2. Populate routing rules
    - **Externally sourced rules**
        (envoy management API serves rules to Envoy side-cars through the control plane)
3. Location of loadbalancing decision
    - **Client-side (in side-car)**
        (envoy side-car proxy routes traffic from the service)

#### Netflix Eureka + Ribbon
1. Populate service registry
    - **Simple orchestrator, smart clients**
        (client apps register and report health to the eureka server)
2. Populate routing rules
    - **Either: Locally out-of-the-box, but external is possible if you provide the api**
        (app uses the Ribbon client to query eureka and select a service instance, selection and query logic can be customized with app code)
3. Location of loadbalancing decision
    - **Client-side (in app)**
        (app choses a service based on selection logic in the Ribbon client)

#### HA-Proxy + Consul + consul-template (this project)
1. Populate service registry
    - **Simple orchestrator, smart clients**
        (client apps register and provide a health check endpoint for consul to monitor their status)
2. Populate routing rules
    - **Locally sourced rules**
        (the template used to create routing rules is part of the loadbalancer image)
3. Location of loadbalancing decision
    - **Server-side**
        (routing decision is made by the HA-Proxy instance sitting in front of the service cluster)



