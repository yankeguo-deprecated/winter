package woss

import (
	"context"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/guoyk93/winter/wext"
)

var (
	ext = wext.New[options, *injects]("oss").
		Startup(func(ctx context.Context, opt *options) (inj *injects, err error) {
			inj = &injects{}
			if inj.client, err = oss.New(opt.endpoint, opt.keyID, opt.keySecret); err != nil {
				return
			}
			if inj.bucket, err = inj.client.Bucket(opt.bucket); err != nil {
				return
			}
			return
		})
)

func Client(ctx context.Context, altKeys ...string) *oss.Client {
	return ext.Instance(altKeys...).Get(ctx).client
}

func Bucket(ctx context.Context, altKeys ...string) *oss.Bucket {
	return ext.Instance(altKeys...).Get(ctx).bucket
}

func Installer(opts ...Option) wext.Installer {
	return ext.Installer(opts...)
}
