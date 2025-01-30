package controller

const (
	AppImageName         = "jkaninda/goma-gateway"
	ConfigPath           = "/etc/goma"
	CertsPath            = "/etc/goma/certs"
	BelongsTo            = "goma-gateway"
	GatewayConfigVersion = "2"
	FinalizerName        = "gomaproj.github.io/resources.finalizer"
	ConfigName           = "goma.yml"
	TLSCertFile          = "/etc/goma/certs/tls.crt"
	TLSKeyFile           = "/etc/goma/certs/tls.key"
)

// Middlewares type
const (
	AccessMiddleware = "access" // access middlewares
	BasicAuth        = "basic"  // basic authentication middlewares
	JWTAuth          = "jwt"    // JWT authentication middlewares
	OAuth            = "oauth"  // OAuth authentication middlewares
	accessPolicy     = "accessPolicy"
	addPrefix        = "addPrefix"
	RateLimit        = "rateLimit"
	redirectRegex    = "redirectRegex"
	rewriteRegex     = "rewriteRegex"
	forwardAuth      = "forwardAuth"
	httpCache        = "httpCache"
	redirectScheme   = "redirectScheme"
	bodyLimit        = "bodyLimit"
)

var (
	ReplicaCount int32 = 1
)
