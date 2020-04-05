package server

import (
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/rs/zerolog/log"
)

func scaffoldHTTPClient() *http.Client {
	if httpClient == nil {
		idleConnTimeoutDuration := time.Duration(Runner.
			Connection.
			Client.
			IdleTimeoutSeconds) * time.Second

		tLSHandshakeTimeoutDuration := time.Duration(Runner.
			Connection.
			Client.
			SocketTimeoutSeconds) * time.Second

		socketTimeoutDuration := time.Duration(Runner.
			Connection.
			Client.
			SocketTimeoutSeconds) * time.Second

		httpClient = &http.Client{
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   socketTimeoutDuration,
					KeepAlive: getKeepAliveIntervalDuration(),
				}).Dial,
				//TLS handshake timeout is the same as connection timeout
				TLSHandshakeTimeout: tLSHandshakeTimeoutDuration,
				MaxIdleConns:        Runner.Connection.Client.PoolSize,
				MaxIdleConnsPerHost: Runner.Connection.Client.PoolSize,
				IdleConnTimeout:     idleConnTimeoutDuration,
			},
		}

		log.Debug().
			Int("upstreamMaxIdleConns", Runner.Connection.Client.PoolSize).
			Int("upstreamMaxIdleConnsPerHost", Runner.Connection.Client.PoolSize).
			Float64("upstreamTransportDialTimeoutSeconds", socketTimeoutDuration.Seconds()).
			Float64("upstreamTlsHandshakeTimeoutSeconds", tLSHandshakeTimeoutDuration.Seconds()).
			Float64("upstreamIdleConnTimeoutSeconds", idleConnTimeoutDuration.Seconds()).
			Float64("upstreamTransportDialKeepAliveIntervalSeconds", getKeepAliveIntervalDuration().Seconds()).
			Msg("server derived upstream params")
	}
	return httpClient
}

// getKeepAliveIntervalSecondsDuration. KeepAlive is effectively: initial delay + interval * TCP_KEEPCNT (9 on linux, 8 ox OSX).
// The KeepAliveIntervalSecondsDuration here defines interval, i.e. default 15s * 9 = 135s on linux
// See: https://github.com/golang/go/issues/23459#issuecomment-374777402
// The OS uses zero payload TCP segments to attempt to keep the connection alive.
// after the total number of unacknowledged TCP_KEEPCNT is reached, the dialer kills the
// connection.
func getKeepAliveIntervalDuration() time.Duration {
	return time.Duration(float64(Runner.
		Connection.
		Client.
		IdleTimeoutSeconds) / float64(getTCPKeepCnt()) * float64(time.Second))
}

func getTCPKeepCnt() int {
	switch runtime.GOOS {
	case "windows":
		return 5
	case "darwin", "freebsd", "openbsd":
		return 8
	case "linux":
		return 9
	//if we don't know, assume some kind of linux
	default:
		return 9
	}
}