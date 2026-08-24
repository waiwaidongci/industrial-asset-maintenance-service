package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Address         string
	ShutdownTimeout time.Duration
	RequestTimeout  time.Duration
	RateLimit       int
}

func Load(path string) Config {
	c := Config{Address: ":8085", ShutdownTimeout: 10 * time.Second, RequestTimeout: 15 * time.Second, RateLimit: 200}
	if b, e := os.ReadFile(path); e == nil {
		parseYAML(&c, string(b))
	}
	applyEnv(&c)
	return c
}
func parseYAML(c *Config, data string) {
	for _, line := range strings.Split(data, "\n") {
		p := strings.SplitN(strings.TrimSpace(line), ":", 2)
		if len(p) != 2 {
			continue
		}
		key := strings.TrimSpace(p[0])
		value := strings.Trim(strings.TrimSpace(p[1]), "\"")
		switch key {
		case "address":
			if value != "" {
				c.Address = value
			}
		case "shutdown_timeout":
			if v, e := time.ParseDuration(value); e == nil {
				c.ShutdownTimeout = v
			}
		case "request_timeout":
			if v, e := time.ParseDuration(value); e == nil {
				c.RequestTimeout = v
			}
		case "rate_limit":
			if v, e := strconv.Atoi(value); e == nil && v > 0 {
				c.RateLimit = v
			}
		}
	}
}
func applyEnv(c *Config) {
	if v := os.Getenv("HTTP_ADDR"); v != "" {
		c.Address = v
	}
	if v := os.Getenv("SHUTDOWN_TIMEOUT"); v != "" {
		if d, e := time.ParseDuration(v); e == nil {
			c.ShutdownTimeout = d
		}
	}
	if v := os.Getenv("REQUEST_TIMEOUT"); v != "" {
		if d, e := time.ParseDuration(v); e == nil {
			c.RequestTimeout = d
		}
	}
	if v := os.Getenv("RATE_LIMIT"); v != "" {
		if n, e := strconv.Atoi(v); e == nil && n > 0 {
			c.RateLimit = n
		}
	}
}
