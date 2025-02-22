package envs

type Config struct {
	IsDebug bool `env:"IS_DEBUG"`
	// URLS    []string `env:"URLS"`
	// Logger
	Logger string `env:"LOGGER,required"` // json, console, none
	// Tracer
	Tracer            string `env:"TRACER,required"` // jaeger, datadog, none
	TracerServiceName string `env:"TRACER_SERVICE_NAME,required"`
	TracerVersion     string `env:"TRACER_VERSION,required"`
	// Database
	DBMaxConnections  int32  `env:"DB_MAX_CONNECTIONS,required"`
	PostgresURL       string `env:"POSTGRES_URL,required"`
	InfluxdbURL       string `env:"INFLUXDB_URL,required"`
	InfluxdbToken     string `env:"INFLUXDB_TOKEN,required"`
	InfluxdbBucket    string `env:"INFLUXDB_BUCKET,required"`
	InfluxdbOrg       string `env:"INFLUXDB_ORG,required"`
	MongodbURL        string `env:"MONGODB_URL,required"`
	MongodbDB         string `env:"MONGODB_DB,required"`
	MongodbCollection string `env:"MONGODB_COLLECTION,required"`
	// Fetcher
	MaxWorkers int64 `env:"MAX_WORKERS,required"`
}
