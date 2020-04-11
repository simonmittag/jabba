package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

//Version is the server version
var Version string = "v0.2.16"

//ID is a unique server ID
var ID string = "unknown"

//Runtime struct defines runtime environment wrapper for a config.
type Runtime struct {
	Config
}

//Runner is the Live environment of the server
var Runner *Runtime

//BootStrap starts up the server from a ServerConfig
func BootStrap() {
	config := new(Config).
		parse("./jabba.json").
		reApplyResourceNames().
		addDefaultPolicy().
		setDefaultTimeouts()

	Runner = &Runtime{Config: *config}
	Runner.initUserAgent().
		assignHandlers().
		startListening()
}

func (runtime Runtime) startListening() {
	readTimeoutDuration := time.Second * time.Duration(runtime.Cnx.Dwn.ReadTimeoutSeconds)
	writeTimeoutDuration := time.Second * time.Duration(runtime.Cnx.Dwn.RoundTripTimeoutSeconds)
	idleTimeoutDuration := time.Second * time.Duration(runtime.Cnx.Dwn.IdleTimeoutSeconds)

	log.Debug().
		Float64("downstreamReadTimeoutSeconds", readTimeoutDuration.Seconds()).
		Float64("downstreamWriteTimeoutSeconds", writeTimeoutDuration.Seconds()).
		Float64("downstreamIdleConnTimeoutSeconds", idleTimeoutDuration.Seconds()).
		Msg("server derived downstream params")
	log.Info().Msgf("Jabba %s listening on port %d...", Version, runtime.Cnx.Dwn.Port)

	server := &http.Server{
		Addr:         ":" + strconv.Itoa(runtime.Cnx.Dwn.Port),
		Handler:      nil,
		ReadTimeout:  readTimeoutDuration,
		WriteTimeout: writeTimeoutDuration,
		IdleTimeout:  idleTimeoutDuration,
	}

	//this line blocks execution and the server stays up
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal().Err(err).Msgf("unable to start HTTP server on port %d, exiting...", runtime.Cnx.Dwn.Port)
		panic(err.Error())
	}
}

func (runtime Runtime) assignHandlers() Runtime {
	for _, route := range runtime.Routes {
		if route.Resource == AboutJabba {
			http.HandleFunc(route.Path, aboutHandler)
			log.Debug().Msgf("assigned about handler to path %s", route.Path)
		}
	}
	http.HandleFunc("/", proxyHandler)
	log.Debug().Msgf("assigned proxy handler to path %s", "/")
	return runtime
}

func (runtime Runtime) initUserAgent() Runtime {
	if httpClient == nil {
		httpClient = scaffoldHTTPClient(runtime)
	}
	return runtime
}

//TODO: this really needs to be configurable. it adds a lot of options to every response otherewise.
func writeStandardResponseHeaders(proxy *Proxy) {
	header := proxy.Dwn.Resp.Writer.Header()

	header.Set("Server", fmt.Sprintf("Jabba %s %s", Version, ID))
	header.Set("Content-Encoding", proxy.contentEncoding())
	header.Set("Cache-control:", "no-store, no-cache, must-revalidate, proxy-revalidate")
	//for TLS response, we set HSTS header see RFC6797
	if Runner.Cnx.Dwn.Mode == "TLS" {
		header.Set("Strict-Transport-Security", "max-age=31536000")
	}
	header.Set("X-xss-protection", "1;mode=block")
	header.Set("X-content-type-options", "nosniff")
	header.Set("X-frame-options", "sameorigin")
	//copy the X-REQUEST-ID from the request
	header.Set(XRequestID, proxy.XRequestID)
}

func sendStatusCodeAsJSON(proxy *Proxy) {

	if proxy.Dwn.Resp.StatusCode >= 299 {
		log.Warn().Int("downstreamResponseCode", proxy.Dwn.Resp.StatusCode).
			Str("downstreamResponseMessage", proxy.Dwn.Resp.Message).
			Str("path", proxy.Dwn.Req.URL.Path).
			Str(XRequestID, proxy.XRequestID).
			Str("method", proxy.Dwn.Method).
			Msgf("request not served")
	}
	writeStandardResponseHeaders(proxy)
	proxy.Dwn.Resp.Writer.Header().Set("Content-Type", "application/json")
	statusCodeResponse := StatusCodeResponse{
		Code:       proxy.Dwn.Resp.StatusCode,
		Message:    proxy.Dwn.Resp.Message,
		XRequestID: proxy.XRequestID,
	}
	proxy.Dwn.Resp.Writer.Write(statusCodeResponse.AsJSON())
}
