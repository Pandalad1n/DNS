# DNS

drone navigation service

## Prerequisites

Docker - https://docs.docker.com/engine/install

## Architecture

- [cmd](./cmd) - executables
- [internal](./internal) - internal project files
- [prom](./prom) - prometheus config folder
- [k6](./k6) - k6 benchmark config folder

## Testing and Development

To run all tests run `make test`. To run linters run `make lint`. To run benchmark run `make start` and `make bench`. 

To run the project locally execute `make start`. 
To run the project with telemetry execute `make start-telemetry` then go to `localhost:9000`.

## Additional questions

#### What instrumentation this service would need to ensure its observability and operational transparency?

There four main pillars of observability:
  - Logs. Our application will write logs in json format to stdout.
  - Metrics. Our application exposes Prometheus metrics on `/metrics` endpoint.
  - Traces. Our application does not include traces, but production application will push metrics to services like Jaeger, Opentelemetry, etc.
  - Alerts. Our application does not include alerts, but production application will have alert manager based on metrics (ideally SLO based).

#### Why throttling is useful (if it is)? How would you implement it here?

Throttling is useful to protect application from overload or fraud (bruteforce), cost saving (sms costs), etc. Overload protection is usually done at the gateway level. In a simple application like this we could just implement "leaky bucket" on the application level.

#### What we have to change to make DNS be able to service several sectors at the same time?

We can pass the sectorID param into the json request, or make a dynamic handler. For example `r.HandleFunc("/v1/locate/{sectorID}", locateDrone)`.

#### Our CEO wants to establish B2B integration with Mom's Friendly Robot Company by allowing cargo ships of MomCorp to useDNS. The only issue is - MomCorp software expects loc value in location field, but math stays the same. How would you approach this? What’s would be your implementation strategy?

It depends a lot on the details of the project. The main thing is that the business logic for the existing endpoint has to stay the same to avoid increasing complexity. We can implement a new handler for this company or even make a separate service to bring all different clients to the same format.

#### Atlas Corp mathematicians made another breakthrough and now our navigation math is even better and more accurate, so we started producing a new drone model, based on new math. How would you enable scenario where DNS can serve both types of clients?

We can keep the old `/v1` handler and make a new `/v2` one that accepts a new type of requests.

#### In general, how would you separate technical decision to deploy something from business decision to release something?

I will use a "feature flag" to deploy a feature as soon as it is ready and then activate it when it's needed. We can also use split.io to make feature available to some group of users or do A/B testing.