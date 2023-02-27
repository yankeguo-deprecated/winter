package wwx

type HandlerFunc func(c Context)

type options struct {
	appID       string
	appSecret   string
	token       string
	aesKey      string
	redisKey    []string
	msgHandlers map[string]HandlerFunc
	evtHandlers map[string]HandlerFunc
}

type Option = func(opts *options)

// WithRedisKey set redis to use
func WithRedisKey(altKeys ...string) Option {
	return func(opts *options) {
		opts.redisKey = altKeys
	}
}

// WithAppID set AppID
func WithAppID(appID string) Option {
	return func(opts *options) {
		opts.appID = appID
	}
}

// WithAppSecret set AppSecret
func WithAppSecret(appSecret string) Option {
	return func(opts *options) {
		opts.appSecret = appSecret
	}
}

// WithToken set Token
func WithToken(token string) Option {
	return func(opts *options) {
		opts.token = token
	}
}

// WithEncodingAESKey set EncodingAESKey
func WithEncodingAESKey(aesKey string) Option {
	return func(opts *options) {
		opts.aesKey = aesKey
	}
}

// WithMessageHandler set message handler
func WithMessageHandler(typ string, h HandlerFunc) Option {
	return func(opts *options) {
		if opts.msgHandlers == nil {
			opts.msgHandlers = make(map[string]HandlerFunc)
		}
		opts.msgHandlers[typ] = h
	}
}

// WithEventHandler set message handler
func WithEventHandler(evt string, h HandlerFunc) Option {
	return func(opts *options) {
		if opts.evtHandlers == nil {
			opts.evtHandlers = make(map[string]HandlerFunc)
		}
		opts.evtHandlers[evt] = h
	}
}
