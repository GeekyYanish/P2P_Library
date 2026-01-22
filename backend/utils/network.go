/*
================================================================================
NETWORK UTILITIES - P2P Academic Library "The Knowledge Exchange"
================================================================================
This file provides network-related helper functions for P2P communication.

Go Concepts Used:
- net package: Network operations
- Error handling: Go's error patterns
- Pointers: For connection handling
================================================================================
*/

package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// ============================================================================
// CONSTANTS
// ============================================================================

const (
	// DefaultPort is the default listening port for peers
	DefaultPort = 8080

	// ConnectionTimeout is the timeout for establishing connections
	ConnectionTimeout = 10 * time.Second

	// ReadTimeout is the timeout for read operations
	ReadTimeout = 30 * time.Second

	// WriteTimeout is the timeout for write operations
	WriteTimeout = 30 * time.Second

	// MaxMessageSize is the maximum size of a network message
	MaxMessageSize = 10 * 1024 * 1024 // 10 MB
)

// ============================================================================
// MESSAGE TYPES
// ============================================================================

// Message represents a P2P network message
type Message struct {
	Type    string          `json:"type"`
	Sender  string          `json:"sender"`
	Payload json.RawMessage `json:"payload"`
}

// MessageType constants
const (
	MsgTypePing      = "PING"
	MsgTypePong      = "PONG"
	MsgTypeSearch    = "SEARCH"
	MsgTypeResponse  = "RESPONSE"
	MsgTypeTransfer  = "TRANSFER"
	MsgTypeRequest   = "REQUEST"
	MsgTypeHandshake = "HANDSHAKE"
)

// ============================================================================
// CONNECTION HELPERS
// ============================================================================

// CreateListener creates a TCP listener on the specified port
// Parameters:
//   - port: The port number to listen on
//
// Returns:
//   - net.Listener: The created listener
//   - error: Error if listener cannot be created
func CreateListener(port int) (net.Listener, error) {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to create listener on port %d: %w", port, err)
	}
	return listener, nil
}

// Connect establishes a TCP connection to a remote peer
// Parameters:
//   - address: The address to connect to (ip:port)
//
// Returns:
//   - net.Conn: The established connection
//   - error: Error if connection fails
func Connect(address string) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", address, ConnectionTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", address, err)
	}
	return conn, nil
}

// SendMessage sends a message over a connection
// Parameters:
//   - conn: The connection to send on
//   - msg: The message to send
//
// Returns:
//   - error: Error if send fails
func SendMessage(conn net.Conn, msg *Message) error {
	// Set write deadline
	conn.SetWriteDeadline(time.Now().Add(WriteTimeout))

	// Marshal message to JSON
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Add newline as message delimiter
	data = append(data, '\n')

	// Send the data
	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// ReceiveMessage receives a message from a connection
func ReceiveMessage(conn net.Conn) (*Message, error) {
	// Set read deadline
	conn.SetReadDeadline(time.Now().Add(ReadTimeout))

	// Read data
	buffer := make([]byte, MaxMessageSize)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	// Unmarshal message
	var msg Message
	if err := json.Unmarshal(buffer[:n], &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return &msg, nil
}

// ============================================================================
// ADDRESS HELPERS
// ============================================================================

// GetLocalIP returns the local IP address
func GetLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

// ParseAddress parses an address string into IP and port
func ParseAddress(address string) (string, int, error) {
	host, portStr, err := net.SplitHostPort(address)
	if err != nil {
		return "", 0, err
	}

	var port int
	fmt.Sscanf(portStr, "%d", &port)
	return host, port, nil
}

// FormatAddress formats an IP and port into an address string
func FormatAddress(ip string, port int) string {
	return fmt.Sprintf("%s:%d", ip, port)
}

// IsValidPort checks if a port number is valid
func IsValidPort(port int) bool {
	return port > 0 && port <= 65535
}

// IsPortAvailable checks if a port is available for listening
func IsPortAvailable(port int) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	listener.Close()
	return true
}

// ============================================================================
// PING/PONG HELPERS
// ============================================================================

// Ping sends a ping message and waits for pong
func Ping(address string) (time.Duration, error) {
	start := time.Now()

	conn, err := Connect(address)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	// Send ping
	pingMsg := &Message{
		Type:   MsgTypePing,
		Sender: "local",
	}
	if err := SendMessage(conn, pingMsg); err != nil {
		return 0, err
	}

	// Receive pong
	pongMsg, err := ReceiveMessage(conn)
	if err != nil {
		return 0, err
	}

	if pongMsg.Type != MsgTypePong {
		return 0, fmt.Errorf("expected PONG, got %s", pongMsg.Type)
	}

	return time.Since(start), nil
}
