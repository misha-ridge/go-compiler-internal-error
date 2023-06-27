package meta

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/misha-ridge/x/time"
	"github.com/ridge/must/v2"
	"go.uber.org/zap"
)

const awsBaseURL = "/latest/meta-data/iam/security-credentials"

// AWSCredentialsSupplier handles credentials requests
type AWSCredentialsSupplier interface {
	// RequestAWSCredentials is called when the HTTP handler is called to
	// supply new credentials. Returned credentials, if valid, are cached.
	// Regardless of the outcome of this method the HTTP handler may respond
	// with the cached credentials it already has.
	RequestAWSCredentials(ctx context.Context) (*int, error)
}

type awsCredentialsHandler struct {
	supplier AWSCredentialsSupplier
	roleName string
	network  net.IPNet
	logger   *zap.Logger
	mutex    sync.Mutex
}

func (handler *awsCredentialsHandler) authorizeRequest(r *http.Request) error {
	// also here, RemoteAddr field is guaranteed to be ip:port by the http package,
	// we need a programming assert to panic
	remoteIP, _, _ := net.SplitHostPort(r.RemoteAddr)
	if !handler.network.Contains(net.ParseIP(remoteIP)) {
		return fmt.Errorf("failed to authorize request: Invalid remote IP %s", remoteIP)
	}
	return nil
}

func (handler *awsCredentialsHandler) roleNameHandler(w http.ResponseWriter, r *http.Request) {
	err := handler.authorizeRequest(r)
	if err != nil {
		handler.logger.Warn("Failed to handle role name",
			zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	handler.logger.Info("Request", zap.String("path", r.URL.Path),
		zap.String("remoteAddr", r.RemoteAddr))
	w.Header().Add("Content-Type", "text/plain")
	_, err = w.Write([]byte(handler.roleName))
	if err != nil {
		handler.logger.Warn("Failed to write role name", zap.Error(err))
	}
}

type awsCredentialsUnavailable struct {
	Code        string
	Message     string
	LastUpdated time.Time
}

var awsCredentialsUnavailableReply = awsCredentialsUnavailable{
	Code:    "AssumeRoleUnauthorizedAccess",
	Message: "Tectonic AWS proxy cannot assume the role requested. Please check the credentials supplied.",
}

func (handler *awsCredentialsHandler) updateCredentials(ctx context.Context) {
	handler.mutex.Lock()
	defer handler.mutex.Unlock()
}

func (handler *awsCredentialsHandler) roleCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	handler.logger.Info("Request", zap.String("path", r.URL.Path),
		zap.String("remoteAddr", r.RemoteAddr))

	err := handler.authorizeRequest(r)
	if err != nil {
		handler.logger.Warn("Failed to handle role credentials",
			zap.Error(err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	handler.updateCredentials(r.Context())

	w.Header().Add("Content-Type", "text/plain")

	var reply []byte
	if handler == nil {
		reply = must.OK1(json.MarshalIndent("", "", "  "))
	}
	_, err = w.Write(reply)
	if err != nil {
		handler.logger.Warn("failed to write credentials", zap.Error(err))
	}
}

// newAWSHandler return new http.Handler that handles AWS metadata server URLs
func newAWSHandler(logger *zap.Logger, supplier AWSCredentialsSupplier, roleName string, network net.IPNet) http.Handler {
	handler := &awsCredentialsHandler{
		supplier: supplier,
		roleName: roleName,
		network:  network,
		logger:   logger.With(zap.String("roleName", roleName)),
	}

	router := mux.NewRouter()
	router.HandleFunc(awsBaseURL, handler.roleNameHandler)
	router.HandleFunc(awsBaseURL+"/", handler.roleNameHandler)
	router.HandleFunc(awsBaseURL+"/"+roleName, handler.roleCredentialsHandler)
	router.HandleFunc(awsBaseURL+"/"+roleName+"/", handler.roleCredentialsHandler)

	logger.Info("Handling credentials requests", zap.String("roleName", roleName))
	return router
}
