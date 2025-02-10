package envs

// TODO: AppEnvの環境によって、必須パラメータが異なる状況をどう制御すべきか
// 今後は、productionやstagingのenvファイルに関しても値が追加されていくので、適宜`required`を追加していく
// 最終的には、production,stagingの全項目が`required`属性となるはず

type Config struct {
	IsDebug bool `env:"IS_DEBUG"`
	// Database
	InfluxdbURL       string `env:"INFLUXDB_URL,required"`
	InfluxdbToken     string `env:"INFLUXDB_TOKEN,required"`
	InfluxdbBucket    string `env:"INFLUXDB_BUCKET,required"`
	InfluxdbOrg       string `env:"INFLUXDB_ORG,required"`
	MongodbURL        string `env:"MONGODB_URL,required"`
	MongodbDB         string `env:"MONGODB_DB,required"`
	MongodbCollection string `env:"MONGODB_COLLECTION,required"`
}
