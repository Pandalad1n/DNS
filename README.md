# DNS
drone navigation service

## Prerequisites
Docker - https://docs.docker.com/engine/install

## Architecture

- [cmd](./cmd) - executables
- [internal](./internal) - internal project files
- [prom](./prom) - prometheus config folder

## Testing and Development

To run all tests run `make test`.
To run linters run `make lint`.

To run the project locally execute `make build`, then `make start`.
To run the web part only execute `make start-web`.

● What instrumentation this service would need to ensure its observability and operational
transparency?

- There are many ways to ensure observability and operational
transparency. The basic way is to write logs.
We can use prometheus to record metrics.
We can use jaeger to track the trace.
We can also send alerts to our engineers then monitored parameters reach critical values, 
  as a good practice we should always aim to agree with SLO.

● Why throttling is useful (if it is)? How would you implement it here?

- We can make ddos protection at the gateway level. We could also implement leaky bucket.

● What we have to change to make DNS be able to service several sectors at the same
time?
 - We can pass the sectorID param into the json request, or make a dynamic handler.
For example r.HandleFunc("/v1/locate/{sectorID}", locateDrone(sectorID))
   
● Our CEO wants to establish B2B integration with Mom's Friendly Robot Company by
allowing cargo ships of MomCorp to use DNS. The only issue is - MomCorp software
expects loc value in location field, but math stays the same. How would you
approach this? What’s would be your implementation strategy?
- It depends a lot on the details of the project. 
  The main thing is that the business logic has to stay the same.
  We can implement a new handler for this company or even make a separate service to bring all different clients to the same format
  
● Atlas Corp mathematicians made another breakthrough and now our navigation math is
even better and more accurate, so we started producing a new drone model, based on
new math. How would you enable scenario where DNS can serve both types of clients?
- We can keep the old /v1/ handler and make a new /v2 one that accepts a new type of requests

● In general, how would you separate technical decision to deploy something from
business decision to release something?
- I will use a "feature flag" to deploy a feature as soon as it is ready and then activate it when its needed.
  We can also use split.io to make feature available to some group of users or other logic.