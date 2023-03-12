package woss

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

type options struct {
	endpoint  string
	keyID     string
	keySecret string
	bucket    string
}

type injects struct {
	client *oss.Client
	bucket *oss.Bucket
}

// Option option for installation
type Option = func(opts *options)

func WithEndpoint(s string) Option {
	return func(opts *options) {
		opts.endpoint = s
	}
}

func WithKeyID(s string) Option {
	return func(opts *options) {
		opts.keyID = s
	}
}

func WithKeySecret(s string) Option {
	return func(opts *options) {
		opts.keySecret = s
	}
}

func WithBucket(s string) Option {
	return func(opts *options) {
		opts.bucket = s
	}
}
