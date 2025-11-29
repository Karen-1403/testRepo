package conn

import (
	"io"
	"net"
	"sync"

	"github.com/zGate-Team/zGate-Platform/internal/utils"
)

// Forward forwards data bidirectionally between two connections
func Forward(clientConn, serverConn net.Conn) error {
	var wg sync.WaitGroup
	wg.Add(2)

	// Client → Server
	go func() {
		defer wg.Done()
		if _, err := io.Copy(serverConn, clientConn); err != nil {
			utils.Logger.Error("client to server copy failed", "error", err)
		}
	}()

	// Server → Client
	go func() {
		defer wg.Done()
		if _, err := io.Copy(clientConn, serverConn); err != nil {
			utils.Logger.Error("server to client copy failed", "error", err)
		}
	}()

	wg.Wait()
	return nil
}