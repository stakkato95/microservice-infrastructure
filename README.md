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

Istio Gateway == Envoy Edge Proxy

Istio dark release

Istio fault injection
1 first show in UI
2 was ist Chaos Engineering?






























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












==============================================================================================================================

### Kurzer Überblick über Technologien

### service-engineering-go-lib

Das ist eine von mir geschriebene nano-Bibliothek, die die Funktionen / Komponenten umfasst, die in der letzten Übung einfach von einem zum anderen Service dupliziert wurden. Diese Komponenten sind Logging (Wrapper für uber zap Bibliothek), Konfigurationsmanagement (Wrapper für Viper Bibliothek) und eine Funktion, die in HTTP Handlers einer HTTP Response erleichtert. Link: https://github.com/stakkato95/service-engineering-go-lib

![2_go_lib](/images/2_go_lib.jpg)

### gin

![13_gin](/images/13_gin.jpg)

In Go Community ist Gin das beliebteste Framework. In Microservices, die mit Go entwickelt werden, werden oft einfach Routers eingesetzt, aber Gin bietet zusätzlich Parametervalidierung, Serving statischer Dateien und mehrere Arten von Middleware. In meinem Projekt wurde Gin nicht in allen Services eingesetzt, nur im Tweets Service (im Users Service wird "chi" Router verwendet).

```go
func Start() {
	repo := domain.NewTweetsRepo()
	service := service.NewTweetsService(repo)

	h := TweetsHandler{service}

	router := gin.Default()
	router.POST("/tweets", h.addTweet)
	router.GET("/tweets/:userId", h.getTweets)
	router.Run(config.AppConfig.ServerPort)
}

type TweetsHandler struct {
	service service.TweetsService
}

func (h *TweetsHandler) addTweet(ctx *gin.Context) {
	var tweetDto dto.TweetDto
	if err := ctx.ShouldBindJSON(&tweetDto); err != nil {
		errorResponse(ctx, err)
		return
	}

	createdTweet := h.service.AddTweet(tweetDto)
	ctx.JSON(http.StatusOK, dto.ResponseDto{Data: *createdTweet})
}
```

### gqlgen (Graphql Generator)

gqlgen ermöglicht Generierung eines Graphql Services auf Basis des Graphql Schemas. Graphql Schema meines Backend-for-Frontend Services schaut wie folgt aus:

```graphql
type Tweet {
  id: Int!
  userId: Int!
  text: String!
}

type Query {
  tweets: [Tweet!]!
}

input NewUser {
  username: String!
  password: String!
}

input Login {
  username: String!
  password: String!
}

input NewTweet {
  text: String!
}

type Mutation {
  createUser(input: NewUser!): String!

  login(input: Login!): String!

  createTweet(input: NewTweet!): Tweet!
}
```

Mit dem oben dargestellten Schema kann man folgende Queries ausführen:

```graphql
mutation {
  createUser(input: {username: "user1", password: "pass"})
}

mutation {
  login(input: {username: "user1", password: "pass"})
}

mutation {
  createTweet(input: {userId: 1, text: "new tweet"}) {
    id
    userId
    text
  }
}

{
  tweets {
    id
    userId
    text
  }
}
```

Go Handler-Funktionen, die aus dem Schema erzeugt wurden, schauen so aus:
```go
func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	return r.UserService.Create(input)
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	return r.UserService.Authenticate(input)
}

func (r *mutationResolver) CreateTweet(ctx context.Context, input model.NewTweet) (*model.Tweet, error) {
	user := middleware.ForContext(ctx)
	if user == nil {
		return nil, errors.New("invalid authorization")
	}

	return r.TweetService.CreateTweet(input, int(user.Id))
}

func (r *queryResolver) Tweets(ctx context.Context) ([]*model.Tweet, error) {
	user := middleware.ForContext(ctx)
	if user == nil {
		return nil, errors.New("invalid authorization")
	}

	return r.TweetService.GetTweets(int(user.Id))
}
```

### gorm

Im letzten Projekt (zu Message Oriented Middleware) wurden keine ORMs ausprobiert, die das Mapping von Entities auf Tabellen einer Datenbank ermöglichen. Einer der Gründe dafür war geringe beliebtheit von ORMs in Go Community. Viele Projekte verwenden einfach SQL Driver für entsprechende Datenbanken.

Dieses Mal wurde von mir eine der bekanntesten (und warscheinlich auch wenigen) ORM Bibliotheken ausprobiert, und zwar gorm. Gorm bietet alle gewöhnliche Funktionen eines ORMs, inklusive Migrationen, Select mit Where-Bediengungen, Updates usw.

```go
func NewTweetsRepo() TweetsRepo {
	db, err := gorm.Open(postgres.Open(config.AppConfig.DbSource), &gorm.Config{})
	db.AutoMigrate(&Tweet{})
	return &postgresTweetsRepo{db}
}

func (r *postgresTweetsRepo) AddTweet(tweet Tweet) *Tweet {
	r.db.Create(&tweet)
	return &tweet
}

func (r *postgresTweetsRepo) GetAllTweets(userId int) []Tweet {
	tweets := []Tweet{}
	r.db.Where("user_id = ?", userId).Find(&tweets)
	return tweets
}
```

### jwt-go

Diese Bibliothek bietet Heilfsfunktionen zur Erstellung und Validierung der JWT Tokens. Die Struktur eines Tokens im Falle meiner Anwendung schut so aus:

![7_jwt](/images/7_jwt.jpg)

Zusätzliche Hilfsfunktionen auf Basis jwt-go in meinem Projekt:

```go
var (
	SecretKey = []byte(config.AppConfig.JwtSecret)
)

func GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		return username, nil
	} else {
		return "", err
	}
}
``` 

### Packaging für Deployment

Als Packaging für den Service wurde wie bei allen anderen Projekten Docker verwendet (Dockerfile im Rootverzeichnis jedes einzelnen Services). Build eines Images wird durch Makefile gestartet (auch wie im letzten Projekt). Jeder Service hat sein eigenes Repository auf Docker Hub.

![3_dockerhub](/images/3_dockerhub.jpg)

### Deployment der Infrastruktur

Die Infrastruktur dieses Projekts (MySQL und Prostgres) wird mittels Helm (Package Manager für Kubernetes) im k8s installiert (k8s Cluster im Docker Desktop). Zwecks Vereinfachung wiederkehrender Operationen wurden insgesamt vier PowerShell Skripts geschrieben (jeweils zwei pro Datenbamk). Skripts im Verzeichnis `k8s\mysql`:

`install.ps1`
```powershell
$name = "mysql"

# 1
kubectl delete $(kubectl get pvc -l app.kubernetes.io/name=$name -o=name)

# 2
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install $name bitnami/mysql --set auth.rootPassword=root --set auth.database=users

# 3
echo "`nwaiting for pod to be ready..."
kubectl wait --for=condition=Ready $(kubectl get pod -l app.kubernetes.io/name=$name -o=name)
kubectl port-forward svc/$name 3306:3306
```

Operationen im install Skript:
1) Persistence Volume Claim der zuletzt angelegten Datenbank löschen. Dadurch wird automatisch Persistence Volume und folglich alle Daten der Datenbank entfernt.
2) Helm Repo, falls noch nicht hinzugefügt, hinzufügen
3) Abwarten bis der Pod im Redy zustand ist und dann den Port auf localhost mappen (zwecks Visualisierung der Daten in einer Database Viewer App)

`uninstall.ps1`
```powershell
$name = "mysql"

#1
helm uninstall $name

#2
kubectl delete $(kubectl get pvc -l app.kubernetes.io/name=$name -o=name)
```

Operationen im uninstall Skript:
1) Das Deployment der Datenbank mit Helm löschen
2) Anschließend (einfach zur Sicherheit) ebenfalls Persistence Volume Claim löschen

### Deployment der Infrastruktur

Die Anwendung wird dieses Mal nicht nur mit PowerShell Skripts, sondern auch (zum Teil) mit Helm deployed. Im Verzeichnis jeden einzelnen Services wurde mithilfe Helm "helm" Verzeichnis angelegt. In dem Verzeichnis liegen so genannte "Templates". Zu Templates zählen Deployments, Services, Ingress und alle andere Ressourcen, die man mit `kubecetl apply -f XYZ.yaml` im k8s Cluster erstellen kann. 

![4_helm](/images/4_helm.jpg)

Deployment hat folgende Struktur:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Chart.Name }}-deployment
  labels:
    {{- include "helm.labels" . | nindent 4 }}
spec:
  replicas: 1
  selector:
    matchLabels:
      {{- include "helm.selectorLabels" . | nindent 6 }}
  template:
    metadata:
      labels:
        {{- include "helm.selectorLabels" . | nindent 8 }}
    spec:
      containers:
        - name: {{ .Chart.Name }}-container
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          ports:
            - name: {{ .Values.service.http.name }}-container
              containerPort: 8080
              protocol: TCP
```

Die Werte, die in eckigen Klammern `{{ }}` stehen, werden aus der Chart-Definition und aus der `values.yaml` Datei übernommen.

`Chart.yaml`
```yaml
apiVersion: v2
name: twitter-service-tweets
description: A Helm chart for twitter tweets service

home: https://github.com/stakkato95/service-engineering-simple-twitter

maintainers:
  - name: Artsiom Kaliaha
    email: stakkato95@gmail.com

type: application

# Chart version. Versions are expected to follow Semantic Versioning (https://semver.org/)
version: 0.1.0

# Application version
appVersion: "0.1.0"
```

`values.yaml`
```yaml
image:
  repository: stakkato95/twitter-service-tweets
  pullPolicy: Always
  tag: "latest"

nameOverride: ""
fullnameOverride: ""

service:
  type: ClusterIP
  http:
    name: http
    port: 80

ingress:
  enabled: true
```

Gerade bevor ein Helm Chart deployed wird, wird er kompiliert. Das Ergebniss für das Deployment:

`helm install tweets helm --dry-run`
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: twitter-service-tweets-deployment
  labels:
    helm.sh/chart: twitter-service-tweets-0.1.0
    app.kubernetes.io/name: twitter-service-tweets
    app.kubernetes.io/instance: tweets
    app.kubernetes.io/version: "0.1.0"
    app.kubernetes.io/managed-by: Helm
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: twitter-service-tweets
      app.kubernetes.io/instance: tweets
  template:
    metadata:
      labels:
        app.kubernetes.io/name: twitter-service-tweets
        app.kubernetes.io/instance: tweets
    spec:
      containers:
        - name: twitter-service-tweets-container
          image: "stakkato95/twitter-service-tweets:latest"
          imagePullPolicy: Always
          ports:
            - name: http-container
              containerPort: 8080
              protocol: TCP
```

Auf solche Art werden alle Services lokal deployed. Der Skript, der alle Services auf einmal installiert:

`app-install.ps1`
```powershell
cd ../twitter-service-graphql
helm install graphql helm
echo ""

cd ../twitter-service-tweets
helm install tweets helm
echo ""

cd ../twitter-service-users
helm install users helm

cd ../k8s
```

Anzahl der Pods, wenn Infrastruktur und alle Services deployed sind:

![5_pods](/images/5_pods.jpg)

### Testdurchlauf

Als Frontend für mein Projekt wird GraphiQL verwendet. 

Im ersten Schritt muss sich ein User im System anmelden. Das erfolgt mittels einer GraphQL Mutation Operation. Als Ergebniss bekommt man JWT Token. JWT Token muss man bei der Erstellung neuer Tweets oder beim Abfragen aller Tweets im Header mitgeben. In der Auth-Middleware wird der Token extrahiert und an Users Service geschickt, wo die Validierung erfolgt. GraphQL Service kriegt nur das Resultat der Validierung zurück (User oder Fehler Objekt). 

![8_create_user](/images/8_create_user.jpg)

Beim Verlust des Tokens kann ein neuer durch die Login-Operation ausgestellt werden.  

![9_login](/images/9_login.jpg)

Nach der Anmeldung kann ein User anfangen, Tweets zu posten. Dafür wird nur der Text des Tweets und der JWT Token im Auth-Header benötigt.

![10_create_tweet](/images/10_create_tweet.jpg)

Wenn keiner / falscher / ungültiger Token mitgegeben wird, bekommt User eine Fehlermeldung.

![11_auth_err](/images/11_auth_err.jpg)

Nach dem einige Tweets gepostet sind, können sie mittels einer Query abgefragt werden (dafür wird auch ein Token benötigt).

![12_all_tweets](/images/12_all_tweets.jpg)

### Con­clu­sio

Im vorliegenden Projekt wurde Folgendes ausprobiert:
- Graphql als "Frontend for Backend" Pattern
- Authentifizierung mittels JWT beim Graphql Server
- Go ORM gorm
- Go Web-Framework gin
- Packaging und Deployment der Microservices mittels Helm Charts
- Erstellung eigener Go Packages am Beispiel service-engineering-go-lib