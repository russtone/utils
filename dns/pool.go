package dns

import (
	"net"
	"time"
)

type server struct {
	ip net.IP

	queriesCount int
	rateLimit    float64
	createdAt    time.Time

	lastUsedAt time.Time
}

func newServer(ip net.IP, rateLimit float64) server {
	return server{
		ip:        ip,
		rateLimit: rateLimit,
		createdAt: time.Now(),
	}
}

func (s *server) query(name string, qtype string) ([]string, error) {

	defer func() {
		s.lastUsedAt = time.Now()
	}()

	res, err := Query(name, s.ip, qtype)

	if err != nil {
		return nil, err
	}

	s.queriesCount++

	return res, err
}

func (s *server) rate() float64 {
	return float64(s.queriesCount) / float64(time.Since(s.createdAt).Seconds())
}

func (s *server) delay() time.Duration {
	return time.Duration(1/s.rateLimit)*time.Second - time.Since(s.lastUsedAt)
}

type pool struct {
	servers   chan server
	rateLimit float64
}

func newPool(rateLimit float64, capacity int) *pool {
	return &pool{
		servers:   make(chan server, capacity),
		rateLimit: rateLimit,
	}
}

func (p *pool) add(ip net.IP) {
	p.servers <- newServer(ip, p.rateLimit)
}

func (p *pool) take() server {
	return <-p.servers
}

func (p *pool) release(s server) {
	go func() {
		if delay := s.delay(); delay > 0 {
			time.Sleep(delay)
		}

		p.servers <- s
	}()
}
