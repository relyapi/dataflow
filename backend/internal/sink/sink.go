package sink

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tomeai/dataflow/api/v1/flow"
	model2 "github.com/tomeai/dataflow/internal/model"
	"net/url"

	"encoding/base64"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	client "github.com/zinclabs/sdk-go-zincsearch"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"net/http"
	"time"
)

type DataHubSink interface {
	Sink(resources []*model2.Resource) (err error)
}

type SourceStrategy struct {
	sinkProvider DataHubSink
}

func (source *SourceStrategy) setStrategy(sinkProvider DataHubSink) {
	source.sinkProvider = sinkProvider
}

func NewMongoSink(sinkInfo *flow.Sink, logger log.Logger) (DataHubSink, error) {
	loggerHelper := log.NewHelper(log.With(logger, "module", "sink/mongodb"))

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", sinkInfo.Source.Username, sinkInfo.Source.Password, sinkInfo.Source.Host, sinkInfo.Source.Port)
	clientOpts := options.Client().ApplyURI(uri)

	var mongoClient *mongo.Client

	// 设置上下文和超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 使用 backoff 重试机制连接 MongoDB
	err := backoff.Retry(func() error {
		var err error
		mongoClient, err = mongo.Connect(ctx, clientOpts)
		if err != nil {
			return err
		}

		// ping 检查是否连接成功
		err = mongoClient.Ping(ctx, nil)
		return err
	}, backoff.NewExponentialBackOff())

	if err != nil {
		return nil, fmt.Errorf("连接 MongoDB 失败: %w", err)
	}

	return &MongoSink{
		database:  mongoClient.Database(sinkInfo.Source.DbName),
		tableName: sinkInfo.TableName,
		log:       loggerHelper,
	}, nil
}

type authTransport struct {
	wrapped http.RoundTripper
	header  string
}

func (a *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", a.header)
	return a.wrapped.RoundTrip(req)
}

func NewZincSearchSink(sinkInfo *flow.Sink, logger log.Logger) (DataHubSink, error) {
	loggerHelper := log.NewHelper(log.With(logger, "module", "sink/zincsearch"))

	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", sinkInfo.Source.Username, sinkInfo.Source.Password)))

	httpClient := &http.Client{
		Transport: &authTransport{
			wrapped: http.DefaultTransport,
			header:  basicAuth,
		},
	}

	configuration := client.NewConfiguration()
	configuration.Servers = client.ServerConfigurations{
		{
			URL: fmt.Sprintf("http://%s:%s", sinkInfo.Source.Host, sinkInfo.Source.Port),
		},
	}
	configuration.HTTPClient = httpClient

	apiClient := client.NewAPIClient(configuration)

	return &ZincSearchSink{
		zincDocument: apiClient.Document,
		tableName:    sinkInfo.TableName,
		log:          loggerHelper,
	}, nil
}

func NewMysqlSink(sinkInfo *flow.Sink, logger log.Logger) (DataHubSink, error) {
	loggerHelper := log.NewHelper(log.With(logger, "module", "sink/mysql"))
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True&loc=Local", sinkInfo.Source.Username, sinkInfo.Source.Password, sinkInfo.Source.Host, sinkInfo.Source.Port, sinkInfo.Source.DbName)

	var db *gorm.DB
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	retryOp := func() error {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		return err
	}

	bo := backoff.NewExponentialBackOff()
	boCtx := backoff.WithContext(bo, ctx)

	err = backoff.RetryNotify(
		retryOp,
		boCtx,
		func(err error, d time.Duration) {
			loggerHelper.Infof("MySQL 连接失败，%s 后重试：%v", d, err)
		},
	)

	if err != nil {
		return nil, fmt.Errorf("连接 MySQL 失败: %w", err)
	}

	if err := db.Exec(fmt.Sprintf(model2.MysqlDSL, sinkInfo.TableName)).Error; err != nil {
		loggerHelper.Infof("MySQL %s %v", sinkInfo.TableName, err)
	}

	return &MysqlSink{
		db:        db,
		tableName: sinkInfo.TableName,
		log:       loggerHelper,
	}, nil
}

func NewCosSink(sinkInfo *flow.Sink, logger log.Logger) (DataHubSink, error) {
	loggerHelper := log.NewHelper(log.With(logger, "module", "sink/cos"))

	u, err := url.Parse(sinkInfo.Source.Host)
	if err != nil {
		return nil, err
	}
	b := &cos.BaseURL{BucketURL: u}
	return &CosSink{
		threadPool: 5,
		log:        loggerHelper,
		cosClient: cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  sinkInfo.Source.Username,
				SecretKey: sinkInfo.Source.Password,
			},
		}),
	}, nil
}
