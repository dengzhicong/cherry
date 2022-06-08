package cherryConnector

import (
	"crypto/tls"
	cherryFacade "github.com/cherry-game/cherry/facade"
	"github.com/cherry-game/cherry/logger"
	"net"
	"net/http"
)

type connector struct {
	address           string
	listener          net.Listener
	certFile          string
	keyFile           string
	onConnectListener []cherryFacade.OnConnectListener
	connChan          chan cherryFacade.INetConn
}

func (w *connector) OnConnectListener(listener ...cherryFacade.OnConnectListener) {
	w.onConnectListener = append(w.onConnectListener, listener...)
}

func (w *connector) IsSetListener() bool {
	return len(w.onConnectListener) > 0
}

func (w *connector) GetConnChan() chan cherryFacade.INetConn {
	return w.connChan
}

func (w *connector) inChan(conn cherryFacade.INetConn) {
	w.connChan <- conn
}

func (w *connector) executeListener() {
	go func() {
		for conn := range w.GetConnChan() {
			for _, listener := range w.onConnectListener {
				listener(conn)
			}
		}
	}()
}

func newConnector(address, certFile, keyFile string) connector {
	return connector{
		address:           address,
		listener:          nil,
		certFile:          certFile,
		keyFile:           keyFile,
		onConnectListener: make([]cherryFacade.OnConnectListener, 0),
		connChan:          make(chan cherryFacade.INetConn),
	}
}

// GetNetListener 证书构造 net.Listener
func GetNetListener(address, certFile, keyFile string) (net.Listener, error) {
	if certFile == "" || keyFile == "" {
		return net.Listen("tcp", address)
	}

	crt, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		cherryLogger.Fatalf("failed to listen: %s", err.Error())
	}
	tlsCfg := &tls.Config{Certificates: []tls.Certificate{crt}}

	return tls.Listen("tcp", address, tlsCfg)
}

//CheckOrigin 请求检查函数 防止跨站请求伪造 true则不检查
func CheckOrigin(_ *http.Request) bool {
	return true
}
