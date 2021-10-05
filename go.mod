module nationalparks-rest

go 1.16

require (
	github.com/XSAM/otelsql v0.7.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gorilla/mux v1.8.0
	github.com/sirupsen/logrus v1.8.1
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.24.0 // indirect
	go.opentelemetry.io/otel v1.0.0
	go.opentelemetry.io/otel/exporters/jaeger v1.0.0 // indirect
	go.opentelemetry.io/otel/sdk v1.0.0
	golang.org/x/sys v0.0.0-20211001092434-39dca1131b70 // indirect
)
