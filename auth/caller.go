package auth

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/amkulikov/extrpc"
)

const (
	HttpHeaderAuthType = "Iguan-AuthType"
	HttpHeaderLogin    = "Iguan-Login"
	HttpHeaderPassword = "Iguan-Password"
)

// Implements extrpc.Caller
type Caller struct {
	Foo     bool
	a       *AuthCredentials
	ip      net.IP
	headers map[string]string
	dataMu  sync.RWMutex
}

// New instance of defaultCaller
func NewCaller() extrpc.Caller {
	return &Caller{
		headers: make(map[string]string),
	}
}

// get client IP
func (c *Caller) AuthData() *AuthCredentials {
	c.dataMu.RLock()
	defer c.dataMu.RUnlock()

	return c.a
}

// get client IP
func (c *Caller) IP() net.IP {
	c.dataMu.RLock()
	defer c.dataMu.RUnlock()

	return c.ip
}

// header value by key if ok
func (c *Caller) Header(key string) (value string, ok bool) {
	c.dataMu.RLock()
	defer c.dataMu.RUnlock()

	value, ok = c.headers[key]
	return
}

// Parse Caller data from net.Conn
func (c *Caller) ParseConn(conn net.Conn) error {
	c.dataMu.Lock()
	defer c.dataMu.Unlock()

	// get ip as remote addr
	c.ip = net.ParseIP(conn.RemoteAddr().String())

	r := bufio.NewReader(conn)
	authType, err := r.ReadByte()
	if err != nil {
		return ErrParseAuthData
	}

	c.a = &AuthCredentials{
		authType: authType,
	}
	if (authType & AuthTypeToken) != 0 {
		tokenSize, err := r.ReadByte()
		if err != nil || tokenSize == 0 {
			return ErrParseAuthData
		}
		token := make([]byte, tokenSize)
		_, err = io.ReadFull(conn, token)
		if err != nil {
			return ErrParseAuthData
		}
		c.a.login = token
	}

	if (authType & AuthTypeName) != 0 {
		nameSize, err := r.ReadByte()
		if err != nil || nameSize == 0 {
			return ErrParseAuthData
		}
		name := make([]byte, nameSize)
		_, err = io.ReadFull(conn, name)
		if err != nil {
			return ErrParseAuthData
		}
		c.a.password = name
	}

	if !c.a.Valid() {
		return ErrAccessDenied
	}

	return nil
}

// Parse Caller data from http.Request
func (c *Caller) ParseHTTP(req *http.Request) error {
	c.dataMu.Lock()
	defer c.dataMu.Unlock()

	authType, err := strconv.ParseUint(req.Header.Get(HttpHeaderAuthType), 10, 8)
	if err != nil {
		return ErrParseAuthData
	}

	c.a = &AuthCredentials{
		authType: uint8(authType),
		login:    []byte(req.Header.Get(HttpHeaderLogin)),
		password: []byte(req.Header.Get(HttpHeaderPassword)),
	}

	// get ip as remote addr
	c.ip = net.ParseIP(req.RemoteAddr)

	if !c.a.Valid() {
		return ErrAccessDenied
	}
	return nil
}
