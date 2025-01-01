package controller

const (
	AppImageName         = "jkaninda/goma-gateway"
	ConfigPath           = "/etc/goma"
	CertsPath            = "/etc/goma/certs"
	RateLimit            = "rateLimit"
	BelongsTo            = "goma-gateway"
	GatewayConfigVersion = "1.0"
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
	rateLimit        = "rateLimit"
	redirectRegex    = "redirectRegex"
	forwardAuth      = "forwardAuth"
)

var (
	ReplicaCount int32 = 1
)
