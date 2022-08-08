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
- We can use prometheus to track metrics and jaeger to trace.

● Why throttling is useful (if it is)? How would you implement it here?
- ??

● What we have to change to make DNS be able to service several sectors at the same
time?
 - We can pass the sectorID param into the json request, or make a dynamic handler.
For example r.HandleFunc("/v1/locate/{sectorID}", locateDrone(sectorID))
   
● Our CEO wants to establish B2B integration with Mom's Friendly Robot Company by
allowing cargo ships of MomCorp to use DNS. The only issue is - MomCorp software
expects loc value in location field, but math stays the same. How would you
approach this? What’s would be your implementation strategy?
- ??

● Atlas Corp mathematicians made another breakthrough and now our navigation math is
even better and more accurate, so we started producing a new drone model, based on
new math. How would you enable scenario where DNS can serve both types of clients?
- We can keep the pld /v1/ handler and make a new /v2 one that accepts a new type of requests

● In general, how would you separate technical decision to deploy something from
business decision to release something?
 - ??