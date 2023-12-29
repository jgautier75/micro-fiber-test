package config

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"go.uber.org/zap"
)

type Configuration struct {
	ServerPort            string
	TenantId              int64
	LogsMetrics           string
	LogsStd               string
	OAuthClientId         string
	OAuthClientSecret     string
	OAuthCallbackUrl      string
	OAuthRedirectUri      string
	OAuthGithub           string
	OAuthDebug            bool
	RdbmsUrl              string
	RdbmsPoolMin          int
	RdbmsPoolMax          int
	RedisHost             string
	RedisPort             int
	RedisUser             string
	RedisPass             string
	PrometheusMetricsPath string
	PrometheusEnabled     bool
	BasicAuthUser         string
	BasicAuthPass         string
}

func LoadConfigFile(configPath string) *Configuration {
	var kConfig = koanf.New(".")
	errLoadCfg := kConfig.Load(file.Provider(configPath), yaml.Parser())
	if errLoadCfg != nil {
		panic(errLoadCfg)
	}
	config := Configuration{
		ServerPort:            kConfig.String("http.server.port"),
		TenantId:              kConfig.Int64("app.tenant"),
		LogsMetrics:           kConfig.String("app.accessLogFile"),
		LogsStd:               kConfig.String("app.stdLogFile"),
		OAuthCallbackUrl:      kConfig.String("app.oauthCallback"),
		OAuthClientId:         kConfig.String("app.oauthClientId"),
		OAuthClientSecret:     kConfig.String("app.oauthClientSecret"),
		OAuthDebug:            kConfig.Bool("app.oauthDebug"),
		OAuthGithub:           kConfig.String("app.oauthGithub"),
		OAuthRedirectUri:      kConfig.String("app.oauthRedirectUri"),
		PrometheusEnabled:     kConfig.Bool("app.prometheusEnabled"),
		PrometheusMetricsPath: kConfig.String("app.metricsPath"),
		RdbmsUrl:              kConfig.String("app.pgUrl"),
		RdbmsPoolMin:          kConfig.Int("app.pgPoolMin"),
		RdbmsPoolMax:          kConfig.Int("app.pgPoolMax"),
		RedisHost:             kConfig.String("app.redisHost"),
		RedisPort:             kConfig.Int("app.redisPort"),
		RedisUser:             kConfig.String("app.redisUser"),
		RedisPass:             kConfig.String("app.redisPass"),
		BasicAuthUser:         kConfig.String("app.basicAuthUser"),
		BasicAuthPass:         kConfig.String("app.basicAuthPass"),
	}
	return &config
}

func (c Configuration) SetupCnxPool(zapLogger *zap.Logger) (*pgxpool.Pool, error) {
	zapLogger.Info("CnxPool -> Parse configuration")
	dbConfig, errDbCfg := pgxpool.ParseConfig(c.RdbmsUrl)
	if errDbCfg != nil {
		panic(errDbCfg)
	}
	dbConfig.MinConns = int32(c.RdbmsPoolMin)
	dbConfig.MaxConns = int32(c.RdbmsPoolMax)
	zapLogger.Info("CnxPool -> Initialize pool")
	return pgxpool.NewWithConfig(context.Background(), dbConfig)
}
