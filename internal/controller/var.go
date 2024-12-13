package controller

const (
	AppImageName         = "jkaninda/goma-gateway"
	ConfigPath           = "/etc/goma"
	CertsPath            = "/etc/goma/certs"
	BasicAuth            = "basic" // basic authentication middlewares
	JWTAuth              = "jwt"   // JWT authentication middlewares
	OAuth                = "oauth"
	ratelimit            = "ratelimit"
	RateLimit            = "rateLimit"
	BelongsTo            = "goma-gateway"
	GatewayConfigVersion = "1.0"
	FinalizerName        = "gomaproj.github.io/resources.finalizer"
	ConfigName           = "goma.yml"
	TLSCertFile          = "/etc/goma/certs/tls.crt"
	TLSKeyFile           = "/etc/goma/certs/tls.key"
	accessPolicy         = "accessPolicy"
)

var (
	ReplicaCount int32 = 1
)
