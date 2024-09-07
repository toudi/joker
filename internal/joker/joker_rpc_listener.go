package joker

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/phuslu/log"
)

const rpcListenerFile = "joker.sock"

// we're only supporting the above rpc's and 1024 bytes is more than enough
// to handle the command and one argument.
const rpcMaxLength = 1024

// https://dev.to/douglasmakey/understanding-unix-domain-sockets-in-golang-32n8
func (j *Joker) startRPCListener() error {
	socket, err := net.Listen("unix", rpcListenerFile)
	if err != nil {
		log.Error().Err(err).Msg("could not bind to unix socket file")
		return err
	}

	// register the cleanup function
	j.Defer(func() {
		socket.Close()
		os.Remove(rpcListenerFile)
	})

	go func() {
		for {
			// Accept an incoming connection.
			conn, err := socket.Accept()
			if err != nil {
				log.Fatal().Err(err)
			}

			if conn == nil || err != nil {
				// this essentially means that the deferred cleanup function
				// was called and unless we'd exit here we'd have a race
				// condition.
				return
			}

			// Handle the connection in a separate goroutine.
			go func(conn net.Conn) {
				defer conn.Close()

				// Create a buffer for incoming data.
				buf := make([]byte, rpcMaxLength)

				// Read data from the connection.
				n, err := conn.Read(buf)
				if err != nil {
					log.Fatal().Err(err)
				}

				if err = j.handleRpcCall(buf[:n]); err != nil {
					_, _ = conn.Write([]byte(fmt.Sprintf("ERROR: %v\n", err)))
				} else {
					_, _ = conn.Write([]byte("OK\n"))
				}

			}(conn)
		}
	}()

	return nil
}

var availableRPCs map[string]func(*Joker, string) error

var errUnknownRPC = errors.New("unrecognized RPC")

func (j *Joker) handleRpcCall(buf []byte) error {
	rpcInput := strings.Replace(string(buf), "\n", "", 1)

	for rpcName, rpcHandler := range availableRPCs {
		if strings.HasPrefix(rpcInput, rpcName) {
			return rpcHandler(j, strings.TrimLeft(strings.TrimPrefix(rpcInput, rpcName), " "))
		}
	}
	return errUnknownRPC
}

func init() {
	availableRPCs = map[string]func(*Joker, string) error{
		rpcCmdShutdown:       rpcCmdShutdownHandler,
		rpcCmdCall:           rpcCmdCallHandler,
		rpcCmdStopService:    rpcCmdStopServiceHandler,
		rpcCmdStartService:   rpcCmdStartServiceHandler,
		rpcCmdRestartService: rpcCmdRestartServiceHandler,
	}
}
