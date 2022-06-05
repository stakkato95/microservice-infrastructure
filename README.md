# microservice-infrastructure
Playground to try out diferent microservice infrastructure components

https://itnext.io/safe-and-automation-friendly-canary-deployments-with-helm-669394d2c48a

https://istio.io/latest/docs/ops/deployment/requirements/

https://istio.io/latest/docs/reference/config/analysis/ist0118/

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


Load Balancing based on header
1 YAML only
https://istio.io/v1.1/docs/reference/config/networking/v1alpha3/destination-rule/#:~:text=LoadBalancerSettings.ConsistentHashLB%20Consistent%20Hash-based%20load%20balancing%20can%20be%20used,balancing%20policy%20is%20applicable%20only%20for%20HTTP%20connections.































# Microservice-Infrastruktur (Kaliaha Artsiom, s2110455009)

Das Ziel dieses Projekts bestand darin, unterschiedliche Technologien für die Microservice-Infrastruktur auszuprobieren und zu evaluieren.

Übersicht über die verwendeten Technologien:
- Programmiersprache: Golang
- Web Framework: Gin
- Bibliotheken:
    + chi (Routing)
    + gin (Web-Framework)
    + uber zap (strukturiertes Logging)
    + viper (Konfigurationsmanagement)
    + service-engineering-go-lib (eine eigene Sammlung von Hilfsfunktionen)
    + gin-propagate-xheaders (eigene Bibliothek zum Propagieren von X-* Headers)
- Deployment: k8s, Helm
- Service Mesh: Istio

### Architektur der Lösung

Bei dieser Übung wurden drei einfache Services ohnme Geschäftslogik entwickelt:

- Frontend. Dieser Service dient als API-Gateway für die anderen Services. Frontend ruft Middle Service auf.
- Middle. Middle Service könnte einen Microservice mit der Geschäftslogik darstellen. Das einzoge was er macht ist, Request an Backend weiterzuleiten.
- Backend. Der Service macht keine weitere Aufrufe.

Kommumikation zwischen Service erfolgt über REST API. Jeder dieser Services stellt nur einen Endpoint zu Verfügung, und zwar `GET /request`. Außerdem ist der Inhalt der Request-Handlers bei allen Service (fast) gleich. Der Fokus bei dieser Übung lag viel mehr auf Istio.


![1_architecture](/images/1_architecture.jpg)


### Istio. Überblick

Istio ist eine auf Open Source basierende Service Mesh-Plattform, die steuert, wie Microservices Daten miteinander teilen. Das Produkt enthält APIs, über die Istio in beliebige Logging-Plattformen, Telemetrie- oder Richtliniensysteme integriert werden kann. Es ist für die Ausführung in den verschiedensten Umgebungen ausgelegt: On-Premise, Cloud, Kubernetes Container, Services auf virtuellen Maschinen und mehr.

Die Architektur von Istio teilt sich in eine Data Plane und eine Control Plane auf. Auf der Data Plane fügen Sie Istio Support zu einem Service hinzu, indem Sie einen Sidecar Proxy innerhalb Ihrer Umgebung bereitstellen. Dieser Sidecar-Proxy befindet sich neben einem Microservice und leitet Anfragen an andere Proxies weiter. Zusammen bilden diese Proxies ein Mesh-Netzwerk, das Netzwerkkommunikationen zwischen Microservices abfängt. Auf der Control Plane werden Proxies verwaltet und konfiguriert, um den Datenverkehr zu regeln. Darüber hinaus werden darüber Komponenten konfiguriert, um Richtlinien durchzusetzen und Telemetriedaten zu sammeln.

Mit einem Service Mesh wie Istio haben DevOps-Teams die notwendigen Tools an der Hand, um den Umstieg von monolithischen Anwendungen auf cloudnative Apps – Ansammlungen aus kleinen, unabhängigen und lose gekoppelten Microservice-Applikationen – zu bewerkstelligen. Istio bietet Einblicke in Verhaltensmuster sowie betriebliche Kontrolle über das Service Mesh und die mit ihm gesteuerten Microservices. Durch den Einsatz eines Service Meshs lassen sich die Komplexität der Bereitstellungen sowie die Arbeitslast Ihrer Entwicklungs-Teams reduzieren.

Istio bietet folgende Fuktionalität (keine ausführliche Liste):
- Telemetriedaten
- Traffic Management (Canary Deployments)
- Load Balancing
- API Gateway (Canary Deployments)
- Dark Releases
- Fault Injection (in diesem Projekt nicht betrachtet)
- Circuit Breaker (in diesem Projekt nicht betrachtet)
- Mutual TLS


### Telemetriedaten

Nachdem Istio deployet wurde, findet man folgende Pods im `istio-system` Namespace.

![2_istio_system](/images/2_istio_system.jpg)

Es ist ersichtlich, dass Istio ein komplettes Monitoring Stack mit sich installiert: Grafana für Dashboards, Jaeger für Traces, Prometheus für die Sammlung von Metriken und Kiali für die Visualisierung und die Verwaltung des Service Meshes.

`Kiali`
![3_kiali_1](/images/3_kiali_1.jpg)

![4_kiali_2](/images/4_kiali_2.jpg)

![5_kiali_3](/images/5_kiali_3.jpg)

`Jaeger`
![6_jaeger_1](/images/6_jaeger_1.jpg)

![7_jaeger_2](/images/7_jaeger_2.jpg)

`Grafana`
![8_grafana_1](/images/8_grafana_1.jpg)

![9_grafana_2](/images/9_grafana_2.jpg)

Out-of-the-Box kann Prometheus / Jaeger Spanes zu Traces nicht zusamenstellen. Dafür muss ein Service entsprechend instrumentiert werden, d.h. die dafür nötigten Headers bei jeder weiteren Anfrage mitgeben. Laut der  [Istio Dokumentation](https://istio.io/latest/about/faq/distributed-tracing/#how-to-support-tracing) müssen folgende Headers weitergeleitet werden:

- x-request-id
- x-b3-traceid
- x-b3-spanid
- x-b3-parentspanid
- x-b3-sampled
- x-b3-flags
- b3

Ausgelesen werden diese Headers mittels Middleware, die speziell dafür geschrieben wurde. Der Code dieser Middleware liegt in diesem [Repository](https://github.com/stakkato95/gin-propagate-xheaders). Die ganze Middleware wird im Grunde in diesen zwei Funktionen implementiert:

```go
func XHeadersPropagation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		headers := http.Header{}

		for key, val := range ctx.Request.Header {
			if strings.ToLower(key)[0] == 'x' {
				headers.Set(key, val[0])
			}
		}

		ctx.Set(XHeadersKey, headers)
		ctx.Next()
	}
}

func GetXHeaders(ctx *gin.Context) http.Header {
	headersRaw, ok := ctx.Get(XHeadersKey)
	if !ok {
		logger.Fatal("XHeadersPropagation is not used")
	}
	headers, ok := headersRaw.(http.Header)
	if !ok {
		logger.Fatal("XHeadersPropagation cast error")
	}
	return headers
}
```


### Traffic Management (Canary Deployment)

Istio bietet zwei Wege, um Canary Deployments zu implementieren:

- über UI
- mittels YAML Konfigurationsdateien

Über UI kann ein Canary Deployment in wenigen Schritten umgesetzt werden:

- Zuerst müssen zwei Deployments mit den gleichen `app` und unterschiedlichen `version` Labels deployet werden. Außerdem muss noch ein Service deployet werden, bei dem `app` Label mit den Labels des Deployments übereinstimmt (ein normaler Kubernetes Service).

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: frontend-canary-deployment
  labels:
    helm.sh/chart: frontend-0.1.0
    app: frontend
    app.kubernetes.io/managed-by: Helm
spec:
  replicas: 1
  selector:
    matchLabels:
      app: frontend
  template:
    metadata:
      labels:
        app: frontend
        version: "canary"
    spec:
      containers:
        - name: frontend-container
          image: "stakkato95/microservice-infrastructure-frontend:canary"
          imagePullPolicy: Always
          ports:
            - name: http-container
              containerPort: 8080
              protocol: TCP
```

- In Falls des Middle Services, findet man diese zwei Deployments auf der Kiali Dashboard. Versionierte Deployments werden mittels eines Services gruppiert.

![10_canary_1](/images/10_canary_1.jpg)

- Dann kommt man auf die Seite des Services und in `Actions` Dropdown wählt `Create Weighted Routing`

![11_canary_2](/images/11_canary_2.jpg)

- Danach kann man die Gewichte anpassen.

![12_canary_3](/images/12_canary_3.jpg)

Dasselbe kann auch mittels YAML Konfigurationen erreicht werden. Dafür sind VirtualService und eine DestinationRule nötig (unten sind nicht kompilierte Helm Templates zu sehen):

```yaml
kind: VirtualService
apiVersion: networking.istio.io/v1alpha3
metadata:
  name: {{ .Chart.Name }}
  namespace: default
spec:
  hosts:
    - {{ .Chart.Name }}.default.svc.cluster.local
  http:
  - route:
    - destination:
        host: {{ .Chart.Name }}.default.svc.cluster.local
        subset: {{ .Values.app.version }}
      weight: {{ .Values.release.weight.new }}
    - destination:
        host: {{ .Chart.Name }}.default.svc.cluster.local
        subset: {{ .Values.app.oldVersion }}
      weight: {{ .Values.release.weight.old }}
```

```yaml
kind: DestinationRule
apiVersion: networking.istio.io/v1alpha3
metadata:
  name: {{ .Chart.Name }}
  namespace: default
spec:
  host: {{ .Chart.Name }}.default.svc.cluster.local
  subsets:
    - labels:
        version: {{ .Values.app.version }}
      name: {{ .Values.app.version }}
    - labels:
        version: {{ .Values.app.oldVersion }}
      name: {{ .Values.app.oldVersion }}
```

Wenn im Cluster eine instanz von Version 1 und eine Instanz von Version 2 nebeneinander laufen, dann arbeitet das erste Deployment (genauso wie das zweite) 50% aller Anfragen ab. Wenn es z.B. 4 Instanzen vom ersten und 1 Instanz vom zweiten Deployment gäbe, wäre das ein 80-20 Release. Da man in diesem Fall mehrere Deployment Instanzen (Pods) starten muss, um ein gewichtetes Release umzusetzen, ist es sehr umständlich (reine verschwendung der Rechenleistung furch Instanzierung mehrerer Pods).

Die Canary Version des Middle Services wurde in diesem Projekt allerdings nach dem effizienterem Istio Prinzip ausgerollt. Auf der Kiali UI werden zwei nebeneinander laufende Versionen des Middle Services folgendermaßen dargestellt:

![13_canary_4](/images/13_canary_4.jpg)

Wenn mehrere Anfragen an Middle Service nacheinander geschickt werden, sieht man, dass die Mehrheit aller Anfragen von der `stable` Version abgearbeitet wird. Das liegt an den Werten, die in `values.yaml` des Helm Cahrts definiert sind. Diese Werte werden beim Ausrollen des Canary Deployments überschrieben:

```yaml
app:
  version: "stable"
  oldVersion: null

release:
  weight:
    old: 90
    new: 10
```

Canary Release des Middle Services wird mittels Helm installiert. 90% aller Anfragen (auch wenn es auf dem Screenshot unten nicht genau 90% sind) werden immer noch von Stable Version abgearbeitet:

`middle-canary-install.ps1`
```powershell
helm install middle-canary helm --set app.version="canary" --set app.oldVersion="stable" --set image.tag=canary
```

![14_canary_5](/images/14_canary_5.jpg)

Wenn man rollout Skript ausführt, werden 90% aller Anfragen von der Canary Version abgearbeitet. Das sieht man auch an den Ergebnissen des Load-Testing Skripts `middle-canary-test.ps1`: 

`middle-canary-rollout.ps1`
```powershell
helm upgrade middle-canary helm --set release.weight.old=10 --set release.weight.new=90 --set app.version="canary" --set app.oldVersion="stable" --set image.tag=canary
```

![15_canary_6](/images/15_canary_6.jpg)

![16_canary_7](/images/16_canary_7.jpg)

Im Laufe der Arbeit an diesem Projekt wurde folgende Besonderheit von Istio entdeckt. Istio Erlaubt keine numerische Werte für `version` Label. Mit diesen Werten konnte Canary Release nicht konfiguriert werden:

- 0.1.0
- v0.1.0
- 0_1_0
- v0_1_0

Deswegen wurden `stable` und `canary` als Releasenamen gewählt.


### Load Balancing

Load Balancing wurde bei Backend Service mittels Hashing des Wertes eines Headers Implementiert.  Der Header, auf dessen Basis Load Balancing geschieht, ist `x-api-user-id`. Die Konfiguration wurde im DestinationRule YAML beschrieben. DestinationRule benötigt natürlich auch einen VirtualService.

```yaml
kind: DestinationRule
apiVersion: networking.istio.io/v1alpha3
metadata:
  name: {{ .Chart.Name }}
  namespace: default
spec:
  host: {{ .Chart.Name }}.default.svc.cluster.local
  trafficPolicy:
    loadBalancer:
      consistentHash:
        httpHeaderName: "x-api-user-id"
  subsets:
    - labels:
        {{- include "helm.matchLabels" . | nindent 8 }}
      name: {{ .Chart.Name }}-subset
```

```yaml
kind: VirtualService
apiVersion: networking.istio.io/v1alpha3
metadata:
  name: {{ .Chart.Name }}
  namespace: default
spec:
  hosts:
    - {{ .Chart.Name }}.default.svc.cluster.local
  http:
    - route:
        - destination:
            host: {{ .Chart.Name }}.default.svc.cluster.local
            subset: {{ .Chart.Name }}-subset

```

Diese Konfiguration wird in Kiali so abgebildet:

![17_load_balancing_kiali](/images/17_load_balancing_kiali.jpg)

Zum Testen wurden zwei Skripts geschrieben, die unterschiedliche Werte des `x-api-user-id` Headers mitschicken und deren Anfragen dadurch nur von der neuen oder von der alten Version des Backend Services abgearbeitet werden. 

![18_load_balancing_install](/images/18_load_balancing_install.jpg)

`backend-1-canary-test-1.ps1`
```powershell
while(1) {
    curl http://localhost/request -sS -H "X-Api-User-Id: 4" | jq ".nested.nested.service"    
    sleep(0.5);
}
```

![19_load_balancing_test_canary](/images/19_load_balancing_test_canary.jpg)

`backend-1-canary-test-2.ps1`
```powershell
while(1) {
    curl http://localhost/request -sS -H "X-Api-User-Id: 3" | jq ".nested.nested.service"    
    sleep(0.5);
}
```

![20_load_balancing_test_stable](/images/20_load_balancing_test_stable.jpg)


### API Gateway (Canary Deployments)

Da Istio größtenteils auf Basis Envoy aufgebaut wurde, entspricht Istio Gateway dem Envoy Edge Proxy. Istio bietet im Vergleich mit anderen API Gateways / IngressControllers keine extra Funktionen, außer Fehler-Injezierung (zwecks Chaos Engineering) und Canary Deployments (intelligentes Load Balancing und Traffic Routing). Canary Release, das in diesem Abschnitt umgesetzt wurde, kann auch ohne Istio implementiert werden, allerdings mit einem höheren Aufwand (hauptsächlich Aufgrung der Notwendigkeit mehr Pods laufen zu lassen, um Canary Release richtig zu implementieren).

Canary Release des Frontend Services funktioniert 1-zu-1 wie der Release des Middle Services mit einem einzigen Unterschied, dass Frontend Service am Istio API Gateway hängt. Dafür muss zuierst ein Istio API Gateway installiert werden und dann der Service. Alle diese Schritte erfolgen automatisch mit dem Ausrollen des Canary Releases des Frontend Services.

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: ingress-gateway
spec:
  selector:
    istio: ingressgateway
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "*"
```

```yaml
kind: VirtualService
apiVersion: networking.istio.io/v1alpha3
metadata:
  name: {{ .Chart.Name }}
  namespace: default
spec:
  hosts:
    - "*"
  gateways:
  - ingress-gateway
  http:
  - route:
    - destination:
        host: {{ .Chart.Name }}.default.svc.cluster.local
        subset: {{ .Values.app.version }}
      weight: {{ .Values.release.weight.new }}
    - destination:
        host: {{ .Chart.Name }}.default.svc.cluster.local
        subset: {{ .Values.app.oldVersion }}
      weight: {{ .Values.release.weight.old }}

```

Canary Release des Frontend Services kann auch in zwei Schritten ausgerollt werden, und zwar

- zuerst 90% `stable` und 10% `canary` (`frontend-canary-install.ps1`)
- und später 10% `stable` und 90% `canary` (`frontend-canary-rollout.ps1`)

`frontend-canary-install.ps1`
```powershell
helm install frontend-canary helm --set app.version="canary" --set app.oldVersion="stable" --set image.tag=canary
```

`frontend-canary-install.ps1`
```powershell
helm upgrade frontend-canary helm --set release.weight.old=10 --set release.weight.new=90 --set app.version="canary" --set app.oldVersion="stable" --set image.tag=canary
```

Der Mesh aller Services (inklusive Canary Releases der Middle und Backend Services) schaut dann so aus:

![21_gateway_kiali](/images/21_gateway_kiali.jpg)

Wenn alle Helm Charts installiert werden, dann sieht die Ausgabe von `helm ls` so aus:

![22_gateway_helm](/images/22_gateway_helm.jpg)


### Dark Releases

Ein Dark Release auf Basis Istio kann sehr einfach mittels Routing umgesetzt werden, der vom Inhalt eines Headers abhängt. Dark Release wurde in diesem Projekt im Frontend Service implementiert. So wenn ein Client eine Anfrage an Frontend mit dem Header `X-New-Frontend: 100500` schickt, wird die Ausführung der Anfrage um 3 Sekunden verzögert. 

```yaml
kind: VirtualService
apiVersion: networking.istio.io/v1alpha3
metadata:
  name: {{ .Chart.Name }}
  namespace: default
spec:
  hosts:
    - "*"
  gateways:
  - ingress-gateway
  http:
  - name: canary-dark-release
    match:
    - headers:
        X-New-Frontend:
          exact: "100500"
    fault:
      delay:
        percentage:
          value: 100.0
        fixedDelay: 3s
    route:
    - destination:
        host: {{ .Chart.Name }}.default.svc.cluster.local
        subset: {{ .Values.app.version }}
      weight: {{ .Values.release.weight.new }}
  - route:
    - destination:
        host: {{ .Chart.Name }}.default.svc.cluster.local
        subset: {{ .Values.app.version }}
      weight: {{ .Values.release.weight.new }}
    - destination:
        host: {{ .Chart.Name }}.default.svc.cluster.local
        subset: {{ .Values.app.oldVersion }}
      weight: {{ .Values.release.weight.old }}

```


### Mutual TLS

TLS ist bei Istio standardmäßig eingeschaltet, dafür muss man nichts tun.

![23_tls](/images/23_tls.jpg)


### Con­clu­sio

Im vorliegenden Projekt wurden folgende Funktionen von Istio ausprobiert:

- Sammlung der Telemetriedaten
- Canary Release mittels Traffic Management
- Load Balancing mittels Header-Hashing 
- Canary Release für Frontend Services (eng verbinden mit API Gateway)
- API Gateway
- Dark Releases mittels Header Matching
- Mutual TLS