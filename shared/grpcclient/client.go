package grpcclient

import (
	"crypto/tls"
	"errors"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type SecurityMode string

const (
	SecurityTLS      SecurityMode = "tls"
	SecurityInsecure SecurityMode = "insecure"
)

type Options struct {
	Address            string
	DialTimeout        time.Duration
	Block              bool
	Security           SecurityMode
	TLSConfig          *tls.Config
	UnaryInterceptors  []grpc.UnaryClientInterceptor
	StreamInterceptors []grpc.StreamClientInterceptor
	DialOptions        []grpc.DialOption
}

type Client struct {
	conn      *grpc.ClientConn
	closeOnce sync.Once
	closeErr  error
}

func New(opts Options) (*Client, error) {
	if opts.Address == "" {
		return nil, errors.New("grpcclient: address is required")
	}

	dialTimeout := opts.DialTimeout
	if dialTimeout <= 0 {
		dialTimeout = 10 * time.Second
	}

	dialOpts := make([]grpc.DialOption, 0, 8+len(opts.DialOptions))
	transportCreds, err := transportCredentials(opts.Security, opts.TLSConfig)
	if err != nil {
		return nil, err
	}
	dialOpts = append(dialOpts, grpc.WithTransportCredentials(transportCreds))

	if opts.Block {
		dialOpts = append(dialOpts, grpc.WithBlock())
	}
	if len(opts.UnaryInterceptors) > 0 {
		dialOpts = append(dialOpts, grpc.WithChainUnaryInterceptor(opts.UnaryInterceptors...))
	}
	if len(opts.StreamInterceptors) > 0 {
		dialOpts = append(dialOpts, grpc.WithChainStreamInterceptor(opts.StreamInterceptors...))
	}
	dialOpts = append(dialOpts, opts.DialOptions...)

	// ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	// defer cancel()

	conn, err := grpc.NewClient(opts.Address, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("grpcclient: dial %q failed: %w", opts.Address, err)
	}

	return &Client{conn: conn}, nil
}

func (c *Client) Conn() *grpc.ClientConn {
	return c.conn
}

func (c *Client) Close() error {
	if c == nil {
		return nil
	}

	c.closeOnce.Do(func() {
		if c.conn == nil {
			return
		}
		if err := c.conn.Close(); err != nil {
			c.closeErr = fmt.Errorf("grpcclient: close failed: %w", err)
		}
	})

	return c.closeErr
}

func transportCredentials(mode SecurityMode, tlsConfig *tls.Config) (credentials.TransportCredentials, error) {
	switch mode {
	case "", SecurityTLS:
		return credentials.NewTLS(tlsConfig), nil
	case SecurityInsecure:
		return insecure.NewCredentials(), nil
	default:
		return nil, fmt.Errorf("grpcclient: unsupported security mode %q", mode)
	}
}
