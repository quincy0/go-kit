package trace

// TraceName represents the tracing name.
const TraceName = "go-kit"

// A Config is an opentelemetry config.
type Config struct {
	Name             string  `json:",optional"`
	Endpoint         string  `json:",optional"`
	Sampler          float64 `json:",default=1.0"`
	Batcher          string  `json:",default=jaeger,options=jaeger|zipkin|otlpgrpc|otlphttp"`
	ProjectId        string  `json:",default=hubbuy-2025"`
	GoogleCloudTrace bool    `json:",default=false"`
}
