package controller

const (
	AppImageName         = "jkaninda/goma-gateway"
	ExtraConfigPath      = "/etc/goma/extra/"
	BasicAuth            = "basic" // basic authentication middlewares
	JWTAuth              = "jwt"   // JWT authentication middlewares
	OAuth                = "oauth"
	ratelimit            = "ratelimit"
	RateLimit            = "rateLimit"
	BelongsTo            = "goma-gateway"
	GatewayConfigVersion = "1.0"
	FinalizerName        = "finalizer.gomaproj.jonaskaninda.com"
	ConfigName           = "goma.yml"
)

var (
	ReplicaCount int32 = 1
)
