package proxy

import (
	"context"

	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/proxy"
)

// NewRequestBuilderMiddleware creates a proxy middleware that parses the request params received
// from the outter layer and generates the path to the backend endpoints
func NewRequestBuilderMiddleware(remote *config.Backend) proxy.Middleware {
	return func(next ...proxy.Proxy) proxy.Proxy {
		if len(next) > 1 {
			panic(proxy.ErrTooManyProxies)
		}
		return func(ctx context.Context, request *proxy.Request) (*proxy.Response, error) {
			r := request.Clone()
			// for wildcard routes handling
			if v, ok := remote.ExtraConfig["wildcard"].(map[string]interface{}); v != nil && ok {
				if keepPath, ok := v["keep_original_path"].(bool); keepPath && ok {
					r.Method = remote.Method
					return next[0](ctx, &r)
				}
			}
			r.GeneratePath(remote.URLPattern)
			r.Method = remote.Method
			return next[0](ctx, &r)
		}
	}
}
