package naming

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	_registerURL = "http://%s/api/register"
	_cancelURL   = "http://%s/api/cancel"
	_renewURL    = "http://%s/api/renew"
	_fetchURL    = "http://%s/api/fetch"
	_nodesURL    = "http://%s/api/nodes"
)

const (
	NodeInterval  = 90 * time.Second
	RenewInterval = 60 * time.Second
)

type Config struct {
	Nodes []string
	Env   string
}

type Discovery struct {
	once       *sync.Once
	conf       *Config
	ctx        context.Context
	cancelFunc context.CancelFunc

	//local cache
	mutex    sync.RWMutex
	apps     map[string]*FetchData
	registry map[string]struct{}

	//registry center node
	idx  uint64       //node index
	node atomic.Value //node list
}

func New(conf *Config) *Discovery {
	if len(conf.Nodes) == 0 {
		panic("conf nodes empty!")
	}
	ctx, cancel := context.WithCancel(context.Background())
	dis := &Discovery{
		ctx:        ctx,
		cancelFunc: cancel,
		conf:       conf,
		apps:       map[string]*FetchData{},
		registry:   map[string]struct{}{},
	}
	//from conf get node list
	dis.node.Store(conf.Nodes)
	go dis.updateNode()
	return dis
}

func (dis *Discovery) Register(ctx context.Context, instance *Instance) (context.CancelFunc, error) {
	var err error
	//check local cache
	dis.mutex.Lock()
	if _, ok := dis.registry[instance.AppId]; ok {
		err = errors.New("instance duplicate register")
	} else {
		dis.registry[instance.AppId] = struct{}{} //register local cache
	}
	dis.mutex.Unlock()

	if err != nil {
		return nil, err
	}

	//http register
	ctx, cancel := context.WithCancel(dis.ctx)
	if err = dis.register(instance); err != nil {
		//fail
		dis.mutex.Lock()
		delete(dis.registry, instance.AppId)
		dis.mutex.Unlock()
		return cancel, err
	}

	ch := make(chan struct{}, 1)
	cancelFunc := context.CancelFunc(func() {
		cancel()
		<-ch
	})

	//renew&cancel
	go func() {
		ticker := time.NewTicker(RenewInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := dis.renew(instance); err != nil {
					dis.register(instance)
				}
			case <-ctx.Done():
				dis.cancel(instance)
				ch <- struct{}{}
			}
		}

	}()
	return cancelFunc, nil
}

func (dis *Discovery) register(instance *Instance) error {
	uri := fmt.Sprintf(_registerURL, dis.pickNode())
	log.Println("discovery - request register url:" + uri)
	params := make(map[string]interface{})
	params["env"] = dis.conf.Env
	params["appid"] = instance.AppId
	params["hostname"] = instance.Hostname
	params["addrs"] = instance.Addrs
	params["version"] = instance.Version
	params["status"] = 1
	resp, err := HttpPost(uri, params) //ctx
	if err != nil {
		log.Println(err)
		return err
	}
	res := Response{}
	err = json.Unmarshal([]byte(resp), &res)
	if err != nil {
		log.Println(err)
		return err
	}
	if res.Code != 200 {
		log.Printf("uri is (%v), response code (%v)\n", uri, res.Code)
		return errors.New("conflict")
	}
	return nil
}

func (dis *Discovery) renew(instance *Instance) error {
	uri := fmt.Sprintf(_renewURL, dis.pickNode())
	log.Println("discovery - request renew url:" + uri)
	params := make(map[string]interface{})
	params["env"] = dis.conf.Env
	params["appid"] = instance.AppId
	params["hostname"] = instance.Hostname

	resp, err := HttpPost(uri, params)
	if err != nil {
		log.Println(err)
		return err
	}
	res := Response{}
	err = json.Unmarshal([]byte(resp), &res)
	if err != nil {
		log.Println(err)
		return err
	}
	if res.Code != 200 {
		log.Printf("uri is (%v), response code (%v)\n", uri, res.Code)
		return err
	}
	return nil
}

func (dis *Discovery) cancel(instance *Instance) error {
	//local cache
	uri := fmt.Sprintf(_cancelURL, dis.pickNode())
	log.Println("discovery - request cancel url:" + uri)
	params := make(map[string]interface{})
	params["env"] = dis.conf.Env
	params["appid"] = instance.AppId
	params["hostname"] = instance.Hostname

	resp, err := HttpPost(uri, params)
	if err != nil {
		log.Println(err)
		return err
	}
	res := Response{}
	err = json.Unmarshal([]byte(resp), &res)
	if err != nil {
		log.Println(err)
		return err
	}
	if res.Code != 200 {
		log.Printf("uri is (%v), response code (%v)\n", uri, res.Code)
		return err
	}
	return nil
}

func (dis *Discovery) updateNode() {
	ticker := time.NewTicker(NodeInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			uri := fmt.Sprintf(_nodesURL, dis.pickNode())
			log.Println("discovery - request and update node, url:" + uri)
			params := make(map[string]interface{})
			params["env"] = dis.conf.Env
			resp, err := HttpPost(uri, params)
			if err != nil {
				log.Println(err)
				continue
			}
			res := ResponseFetch{}
			err = json.Unmarshal([]byte(resp), &res)
			if err != nil {
				log.Println(err)
				continue
			}
			newNodes := []string{}
			for _, ins := range res.Data.Instances {
				for _, addr := range ins.Addrs {
					newNodes = append(newNodes, strings.TrimPrefix(addr, "http://"))
				}
			}
			if len(newNodes) == 0 {
				continue

			}
			curNodes := dis.node.Load().([]string)
			if !compareNodes(curNodes, newNodes) {
				dis.node.Store(newNodes)
				log.Println("nodes list changed!")
				log.Println(newNodes)
			} else {
				log.Println("nodes list not change")
				log.Println(curNodes)
			}
		}
	}
}

func compareNodes(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	mapB := make(map[string]struct{}, len(b))
	for _, node := range b {
		mapB[node] = struct{}{}
	}
	for _, node := range a {
		if _, ok := mapB[node]; !ok {
			return false
		}
	}
	return true
}

func (dis *Discovery) pickNode() string {
	nodes, ok := dis.node.Load().([]string)
	if !ok || len(nodes) == 0 {
		return dis.conf.Nodes[dis.idx%uint64(len(dis.conf.Nodes))]
	}
	return nodes[dis.idx%uint64(len(nodes))]
}

func (dis *Discovery) Fetch(ctx context.Context, appId string) ([]*Instance, bool) {
	//from local
	dis.mutex.RLock()
	fetchData, ok := dis.apps[appId]
	dis.mutex.RUnlock()
	if ok {
		log.Println("get data from local memory, appid:" + appId)
		return fetchData.Instances, ok
	}
	//from remote
	uri := fmt.Sprintf(_fetchURL, dis.pickNode())
	params := make(map[string]interface{})
	params["env"] = dis.conf.Env
	params["appid"] = appId
	params["status"] = 1 //up
	log.Println(params)
	resp, err := HttpPost(uri, params)
	res := ResponseFetch{}
	err = json.Unmarshal([]byte(resp), &res)
	if err != nil {
		log.Println(err)
		return nil, false
	}
	var result []*Instance
	for _, ins := range res.Data.Instances {
		result = append(result, ins)

	}
	if len(result) > 0 {
		ok = true
		dis.mutex.Lock()
		dis.apps[appId] = &res.Data
		dis.mutex.Unlock()
	}
	return result, ok
}

func (dis *Discovery) Close() error {
	//todo close myself
	return nil
}

type Response struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type FetchData struct {
	Instances       []*Instance `json:"instances"`
	LatestTimestamp int64       `json:"latest_timestamp"`
}

type ResponseFetch struct {
	Response
	Data FetchData `json:"data"`
}
