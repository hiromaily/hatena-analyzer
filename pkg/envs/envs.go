package envs

type Config struct {
	IsDebug bool `env:"IS_DEBUG"`
	// URLS    []string `env:"URLS"`
	// Database
	PostgresURL       string `env:"POSTGRES_URL,required"`
	InfluxdbURL       string `env:"INFLUXDB_URL,required"`
	InfluxdbToken     string `env:"INFLUXDB_TOKEN,required"`
	InfluxdbBucket    string `env:"INFLUXDB_BUCKET,required"`
	InfluxdbOrg       string `env:"INFLUXDB_ORG,required"`
	MongodbURL        string `env:"MONGODB_URL,required"`
	MongodbDB         string `env:"MONGODB_DB,required"`
	MongodbCollection string `env:"MONGODB_COLLECTION,required"`
}
