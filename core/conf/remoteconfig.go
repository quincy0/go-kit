package conf

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/quincy0/go-kit/core/logx"
	"github.com/quincy0/go-kit/rest/httpc"
)

const (
	ConfigRedisChannelPrefix = "xconf/%s/%s/%s"

	ConfRedisClusterName = "common"
	ConfRedisServiceName = "config-redis"

	ENVDev        = "dev"
	ENVProduction = "production"
	ENVName       = "ENVIRONMENT"
	ENVProjectId  = "PROJECT_ID"
)

type ConfRedis struct {
	ConfSubscribeHost string
	ConfSubscribePass string
}

type RemoteConfig interface {
	Load(c interface{}) error
	WithConfig(clusterName, serviceName string)
}

type cItem struct {
	clusterName string
	serviceName string
}

type remoteConfig struct {
	sync.Mutex

	env       string
	projectId string
	redis     *redis.Client
	c         []cItem
}

func NewRemoteConfig(clusterName, serviceName string) (RemoteConfig, error) {
	env := os.Getenv(ENVName)
	if env != ENVDev && env != ENVProduction {
		logx.Errorw("ENVIRONMENT not set")
	}

	projectId := os.Getenv(ENVProjectId)

	var redisClient *redis.Client

	confRedis := ConfRedis{}
	if err := LoadFromConfigService(ConfRedisClusterName, ConfRedisServiceName, &confRedis); err == nil {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     confRedis.ConfSubscribeHost,
			Password: confRedis.ConfSubscribePass,
		})
	} else {
		logx.Errorw("LoadFromConfigService err",
			logx.Field("cluster", ConfRedisClusterName),
			logx.Field("service", ConfRedisServiceName))
	}

	return &remoteConfig{
		env:       env,
		projectId: projectId,
		redis:     redisClient,
		c: []cItem{
			{
				clusterName: clusterName,
				serviceName: serviceName,
			},
		},
	}, nil
}

func (r *remoteConfig) Load(c interface{}) error {
	if r.env != ENVDev && r.env != ENVProduction {
		logx.Errorw("ENVIRONMENT not set config is not load")
		return nil
	}

	tic := time.NewTicker(time.Minute * 5)
	pingT := time.NewTicker(time.Second * 10)

	chMap := make(map[string]cItem)
	chs := make([]string, 0)
	for _, v := range r.c {
		chKey := fmt.Sprintf(ConfigRedisChannelPrefix, v.serviceName, v.clusterName, r.env) // AppName, ClusterName, NamespaceName
		chs = append(chs, chKey)
		chMap[chKey] = v
	}

	if err := r.load(c); err != nil {
		logx.Errorw("load config failed.", logx.Field("err", err))
		return err
	}

	// 默认定时轮询逻辑
	go func() {
		if err := recover(); err != nil {
			logx.Errorw("remote config load failed", logx.Field("err", err))
			return
		}
		for {
			select {
			case <-tic.C:

				logx.Infow("time to refresh config...")

				for _, v := range r.c {
					if err := r.load(c); err != nil {
						logx.Errorw("load config failed.", logx.Field("err", err), logx.Field("v", v))
					}
				}

				logx.Infow("refresh config finished...")
			}
		}
	}()

	if r.redis != nil {
		rch := r.redis.Subscribe(context.Background(), chs...)

		logx.Errorw("redis subscribe start...", logx.Field("chan", rch.String()))

		go func() {
			if err := recover(); err != nil {
				logx.Errorw("remote config load failed", logx.Field("err", err))
				return
			}
			for {
				select {
				case msg := <-rch.Channel():
					logx.Infow("redis subscribe start",
						logx.Field("channel", msg.Channel),
						logx.Field("payload", msg.Payload))

					_, ok := chMap[msg.Channel]
					if !ok {
						logx.Errorw("config channel not found", logx.Field("msg", msg))
						continue
					}

					if err := r.load(c); err != nil {
						logx.Errorw("load config failed", logx.Field("err", err))
					}

					logx.Infow("redis subscribe finished",
						logx.Field("channel", msg.Channel),
						logx.Field("payload", msg.Payload))

				case <-pingT.C:
					if err := rch.Ping(context.Background()); err != nil {
						logx.Errorw("redis ping err", logx.Field("err", err))

						for _, v := range r.c {
							if err := r.load(c); err != nil {
								logx.Errorw("load config failed.", logx.Field("err", err), logx.Field("v", v))
							}
						}

						rch = r.redis.Subscribe(context.Background(), chs...)

						logx.Errorw("redis subscribe retry...", logx.Field("chan", rch.String()))
					}
				}
			}
		}()
	}

	return nil
}

func (r *remoteConfig) WithConfig(clusterName, serviceName string) {
	// conf-redis不能with
	if clusterName == ConfRedisClusterName && serviceName == ConfRedisServiceName {
		return
	}

	if r.c == nil {
		r.c = make([]cItem, 0)
	}
	r.c = append(r.c, cItem{
		clusterName: clusterName,
		serviceName: serviceName,
	})
}

func (r *remoteConfig) load(c interface{}) error {
	r.Lock()
	defer r.Unlock()
	return r.loadFromRemote(c)
}

func (r *remoteConfig) loadFromRemote(c interface{}) error {
	yamlContent := ""

	for _, v := range r.c {
		url := fmt.Sprintf("https://dev-xconf-m4kxkytjyq-nn.a.run.app/config/%s/dev/%s", v.clusterName, v.serviceName)
		//url := fmt.Sprintf("http://dev.xconf.newb.bio:8000/config/%s/dev/%s", clusterName, serviceName)
		if r.env == "production" {
			url = fmt.Sprintf("https://production-xconf-m4kxkytjyq-nn.a.run.app/config/%s/production/%s", v.clusterName, v.serviceName)
		}
		if r.projectId == "gexhub" {
			url = fmt.Sprintf("https://xconf-223026273590.northamerica-northeast1.run.app/config/%s/%s/%s", v.clusterName, r.env, v.serviceName)
		}

		ret, err := httpc.Do(context.Background(), http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		if ret != nil {
			defer ret.Body.Close()
		}
		confReply := ConfReply{}
		bodyStr, err := ioutil.ReadAll(ret.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(bodyStr, &confReply); err != nil {
			return err
		}

		if confReply.Code != 0 {
			return errors.New(fmt.Sprintf("config reply err: %s", confReply.Msg))
		}

		switch confReply.Data.Format {
		case "yaml":
			yamlContent += "\n"
			yamlContent += confReply.Data.Value
		case "json":
			err := LoadFromJsonBytes([]byte(confReply.Data.Value), c)
			if err != nil {
				return err
			}
		default:
			return errors.New(fmt.Sprintf("config data format invalid:%s", confReply.Data.Format))
		}
	}

	if len(yamlContent) > 0 {
		err := LoadFromYamlBytes([]byte(yamlContent), c)
		if err != nil {
			logx.Errorw("config load result error", logx.Field("err", err))
			return err
		}
	}

	if r.env == "dev" {
		logx.Infow("config load result", logx.Field("config", c))
	}

	return nil
}
