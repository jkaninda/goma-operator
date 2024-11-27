package controller

import gomaprojv1beta1 "github.com/jkaninda/goma-operator/api/v1beta1"

func mapToGateway(g gomaprojv1beta1.GatewaySpec) Gateway {
	return Gateway{
		SSLKeyFile:                   "",
		SSLCertFile:                  "",
		Redis:                        Redis{},
		WriteTimeout:                 g.Server.WriteTimeout,
		ReadTimeout:                  g.Server.ReadTimeout,
		IdleTimeout:                  g.Server.IdleTimeout,
		LogLevel:                     g.Server.LogLevel,
		Cors:                         g.Server.Cors,
		DisableHealthCheckStatus:     g.Server.DisableHealthCheckStatus,
		DisableRouteHealthCheckError: g.Server.DisableHealthCheckStatus,
		DisableKeepAlive:             g.Server.DisableKeepAlive,
		InterceptErrors:              g.Server.InterceptErrors,
		EnableMetrics:                g.Server.EnableMetrics,
	}

}
