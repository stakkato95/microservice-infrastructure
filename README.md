# microservice-infrastructure
Playground to try out diferent microservice infrastructure components

https://istio.io/latest/about/faq/distributed-tracing/

https://github.com/stakkato95/gin-propagate-xheaders

https://itnext.io/safe-and-automation-friendly-canary-deployments-with-helm-669394d2c48a

https://istio.io/latest/docs/ops/deployment/requirements/

https://istio.io/latest/docs/reference/config/analysis/ist0118/


ISTIO UNGÜLTIGE VERSIONNAMEN
0.1.0
v0.1.0
0_1_0
v0_1_0

Virtual Services und Destination Rules
kubectl get vs
kubectl get dr

========================================================================
grpc context propagation
========================================================================
https://medium.com/@the.real.yushuf/propagate-trace-headers-with-istio-grpc-http-1-1-go-73e7f5382643

https://rakyll.medium.com/context-propagation-over-http-in-go-d4540996e9b0
========================================================================






========================================================================
canary
========================================================================
https://deliverybot.dev/2019/09/14/safe-and-automation-friendly-canary-deployments-with-helm/#:~:text=Safe%20and%20automation%20friendly%20canary%20deployments%20with%20Helm,only%20exposing%20it%20to%20a%20subset%20of%20traffic.

https://github.com/deliverybot/helm/blob/master/charts/app/templates/deployment.yaml

https://itnext.io/safe-and-automation-friendly-canary-deployments-with-helm-669394d2c48a
========================================================================





Monitoring
1 meine Bibliothek
2 jaeger
3 grafana

Virtual Service canary rollout
1 über UI
2 mein middle-service-rollout skript
3 YAML in Kiali UI bearbeiten

Load Balancing based on header
1 YAML only
https://istio.io/v1.1/docs/reference/config/networking/v1alpha3/destination-rule/#:~:text=LoadBalancerSettings.ConsistentHashLB%20Consistent%20Hash-based%20load%20balancing%20can%20be%20used,balancing%20policy%20is%20applicable%20only%20for%20HTTP%20connections.

Istio Gateway == Envoy Edge Proxy

Istio dark release

Istio fault injection
1 first show in UI
2 was ist Chaos Engineering?
