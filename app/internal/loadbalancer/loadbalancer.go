package loadbalancer

// Forked by https://github.com/kasvith/simplelb

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/ortisan/router-go/internal/config"
	"github.com/ortisan/router-go/internal/constant"
	errApp "github.com/ortisan/router-go/internal/error"
	"github.com/ortisan/router-go/internal/repository"
)

const (
	Attempts int = iota
	Retry
	HealthCheckTicksTime = 10 * time.Second
	PrefixConfig         = "/services/prefix/"
	MaxRetries           = 3
	BackoffTimeout       = 10 * time.Millisecond
	StatusUp             = "up"
	statusDown           = "down"

	HealthCheckDialUp = 1
	HealthCheckHttp   = 2
)

type HealthCheck struct {
	Type     int
	Endpoint string
}

// Backend holds the data about a server
type Backend struct {
	ServicePrefix string
	URL           *url.URL
	ZoneAws       string
	Alive         bool
	HealthCheck   HealthCheck
	mux           sync.RWMutex
}

func newHealthcheck(typeStr string, endpoint string) HealthCheck {
	if typeStr == "tcp" {
		return HealthCheck{Type: HealthCheckDialUp, Endpoint: endpoint}
	} else {
		return HealthCheck{Type: HealthCheckHttp, Endpoint: endpoint}
	}
}

// SetAlive for this backend
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()
}

// IsAlive returns true when backend is alive
func (b *Backend) IsAlive() (alive bool) {
	b.mux.RLock()
	alive = b.Alive
	b.mux.RUnlock()
	return
}

// ServerPool holds information about reachable backends
type ServerPool struct {
	ServicePrefix string
	backends      []*Backend
	current       uint64
}

// AddBackend to the server pool
func (s *ServerPool) AddBackend(backend *Backend) {
	s.backends = append(s.backends, backend)
}

// Pools by prefix
type ServerPools struct {
	ServerPoolByPrefix  map[string]*ServerPool
	ServerPoolByAwsZone map[string]*ServerPool
}

func NewServerPools() *ServerPools {
	return &ServerPools{ServerPoolByPrefix: make(map[string]*ServerPool), ServerPoolByAwsZone: make(map[string]*ServerPool)}
}

func (s *ServerPools) AddServerPoolByPrefix(servicePrefix string, serverPool *ServerPool) {
	s.ServerPoolByPrefix[servicePrefix] = serverPool
}

func (s *ServerPools) AddServerPoolByAwsZone(prefix string, serverPool *ServerPool) {
	s.ServerPoolByAwsZone[prefix] = serverPool
}

func (s *ServerPools) GetServerPoolByPrefix(prefix string) *ServerPool {
	return s.ServerPoolByPrefix[prefix]
}

// NextIndex atomically increase the counter and return an index
func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.backends)))
}

// MarkBackendStatus changes a status of a backend
func (s *ServerPool) MarkBackendStatus(backendUrl *url.URL, alive bool) {
	for _, b := range s.backends {
		if b.URL.String() == backendUrl.String() {
			b.SetAlive(alive)
			break
		}
	}
}

// GetNextPeer returns next active peer to take a connection
func (s *ServerPool) GetNextPeer() *Backend {
	// loop entire backends to find out an Alive backend
	next := s.NextIndex()
	l := len(s.backends) + next // start from next and move a full cycle
	for i := next; i < l; i++ {
		idx := i % len(s.backends)     // take an index by modding
		if s.backends[idx].IsAlive() { // if we have an alive backend, use it and store if its not the original one
			if i != next {
				atomic.StoreUint64(&s.current, uint64(idx))
			}
			return s.backends[idx]
		}
	}
	return nil
}

func (s *ServerPool) HandleRequest(c *gin.Context, pathUri string, method string, headers map[string][]string) {

	peer := s.GetNextPeer()
	if peer == nil {
		panic(errApp.NewGenericError("No backend servers was found", nil))
	}

	requestUri := fmt.Sprintf("%s%s", peer.URL.String(), pathUri)

	// Tracing this request
	ctx, span := tracer.Start(c.Request.Context(), "HandleRequest", trace.WithAttributes(
		attribute.String("ServicePrefix", s.ServicePrefix)),
	)

	defer span.End()

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, method, requestUri, c.Request.Body)
	if err != nil {
		panic(errApp.NewGenericError("Error to create request", err))
	}

	// Set trace id
	req.Header.Set(constant.TraceIdHeaderName, c.GetString(constant.TraceIdHeaderName))
	// Copying headers
	for name, values := range headers {
		for _, value := range values {
			log.Debug().Str(name, value).Msg("Iterating headers...")
			if found := HeadersDisabledInRedirection()(name); !found {
				req.Header.Set(name, value)
			}
		}
	}

	resp, err := client.Do(req) // Call API
	if err != nil {
		panic(errApp.NewIntegrationError("Error to call API", err))
	}

	defer resp.Body.Close() // Defer will close after this function ends
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(errApp.NewIntegrationError("Error read response body", err))
	}

	c.Data(resp.StatusCode, resp.Header.Get(constant.ContentTypeHeaderName), body) // Data is returned
}

// doHealthCheck pings the backends and update the status
func (s *ServerPool) doHealthCheck() {
	for _, b := range s.backends {
		b.doHealthCheck()
	}
}

// GetAttemptsFromContext returns the attempts for request
func GetAttemptsFromContext(r *http.Request) int {
	if attempts, ok := r.Context().Value(Attempts).(int); ok {
		return attempts
	}
	return 1
}

// GetAttemptsFromContext returns the attempts for request
func GetRetryFromContext(r *http.Request) int {
	if retry, ok := r.Context().Value(Retry).(int); ok {
		return retry
	}
	return 0
}

var tracer = otel.Tracer(config.ConfigObj.App.Name)

// isAlive checks whether a backend is Alive by establishing a TCP connection
func (b *Backend) doHealthCheck() {
	var alive = false

	if b.HealthCheck.Type == HealthCheckDialUp {
		timeout := 2 * time.Second
		var address = b.URL.Host
		if !strings.Contains(address, ":") {
			address = fmt.Sprintf("%s:%s", address, "80")
		}
		conn, err := net.DialTimeout("tcp", address, timeout)
		if err != nil {
			log.Warn().Err(err).Msg("Server healthcheck error.")
		} else {
			alive = true
		}
		defer conn.Close()
	} else {
		resp, err := http.Get(b.HealthCheck.Endpoint)
		if err != nil {
			log.Warn().Err(err).Msg("Server healthcheck error.")
		} else {
			alive = resp.StatusCode == http.StatusOK
		}
	}
	b.SetAlive(alive)
	var status string
	if alive {
		status = StatusUp
	} else {
		status = statusDown
	}
	serverId := fmt.Sprintf("servers-%s-%s", b.ServicePrefix, b.URL.String())
	repository.PutCacheValue(serverId, status)
}

// healthCheck runs a routine for check status of the backends every 30 seconds
func healthCheck() {
	t := time.NewTicker(HealthCheckTicksTime)
	for {
		select {
		case <-t.C:
			log.Debug().Msg("Starting health check...")
			for servicePrefix, serverPool := range ServerPoolsObj.ServerPoolByPrefix {
				log.Debug().Str("prefix", servicePrefix).Msg("Health checking services of this prefix")
				serverPool.doHealthCheck()
			}
			log.Debug().Msg("Health check completed")
		}
	}
}

func HeadersDisabledInRedirection() func(string) bool {
	innerMap := map[string]int{
		"Accept-Encoding": 1, // This header transform encodings
	}
	return func(key string) bool {

		_, found := innerMap[key]
		return found
	}
}

func Config() {

	serversConfig := config.ConfigObj.Servers

	for _, server := range serversConfig {

		serverUrl, err := url.Parse(server.EndpointUrl)
		if err != nil {
			panic(err)
		}

		var serverPool *ServerPool
		serverPool = ServerPoolsObj.GetServerPoolByPrefix(server.ServicePrefix)
		if serverPool == nil {
			// Init Server Pool
			serverPool = &ServerPool{ServicePrefix: server.ServicePrefix}
			ServerPoolsObj.AddServerPoolByPrefix(server.ServicePrefix, serverPool)
		}

		// Update status status from cache db
		serverId := fmt.Sprintf("servers-%s-%s", server.ServicePrefix, server.EndpointUrl)
		serverStatus, err := repository.GetCacheValue(serverId)
		var alive = false
		if err != nil {
			switch err.(type) {
			case errApp.NotFoundError:
				alive = true
			default:
				panic(err)
			}
		}
		alive = alive || serverStatus == StatusUp
		// Add server to serverpool
		serverPool.AddBackend(&Backend{ServicePrefix: server.ServicePrefix, URL: serverUrl, Alive: alive, ZoneAws: server.ZoneAws, HealthCheck: newHealthcheck(server.HealthCheck.Type, server.HealthCheck.Endpoint)})
	}

	// start health checking
	go healthCheck()
}

var ServerPoolsObj *ServerPools = NewServerPools()
