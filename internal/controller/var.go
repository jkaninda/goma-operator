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
	FinalizerName        = "finalizer.gomaproj.jonaskaninda.com"
	ConfigName           = "goma.yml"
	TLSCertFile          = "/etc/goma/certs/tls.crt"
	TLSKeyFile           = "/etc/goma/certs/tls.key"
)

var (
	ReplicaCount int32 = 1
)
