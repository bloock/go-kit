package client

import "time"

type SQLConnOpts struct {
	MaxConnLifeTime time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
}
