package main

import (
	"crypto/rand"

	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/pkg/errors"

	"golang.org/x/crypto/nacl/box"
)

type secureReader struct {
	r         io.Reader
	priv, pub *[32]byte
}

func (s *secureReader) Read(p []byte) (n int, err error) {
	nonce := [24]byte{}

	_, err = io.ReadFull(s.r, nonce[:])
	if err != nil {
		return 0, err
	}
	boxed := make([]byte, len(p)+box.Overhead)
	n, err = s.r.Read(boxed)
	if err != nil {
		return n, err
	}
	dec, ok := box.Open(nil, boxed[:n], &nonce, s.pub, s.priv)
	if !ok {
		return n, errors.New("error while decrypting")
	}
	copy(p, dec)
	return len(dec), nil
}

// NewSecureReader instantiates a new SecureReader
func NewSecureReader(r io.Reader, priv, pub *[32]byte) io.Reader {
	return &secureReader{r, priv, pub}
}

type secureWriter struct {
	w         io.Writer
	priv, pub *[32]byte
}

func (s *secureWriter) Write(p []byte) (n int, err error) {
	nonce := [24]byte{}

	_, err = io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return 0, err
	}
	boxed := box.Seal(nonce[:], p, &nonce, s.pub, s.priv)
	n, err = s.w.Write(boxed)
	if err != nil {
		return n, err
	}
	return len(p), nil
}

// NewSecureWriter instantiates a new SecureWriter
func NewSecureWriter(w io.Writer, priv, pub *[32]byte) io.Writer {
	return &secureWriter{w, priv, pub}
}

type secureCon struct {
	io.Reader
	io.Writer
	io.Closer
}

// Dial generates a private/public key pair,
// connects to the server, perform the handshake
// and return a reader/writer.
func Dial(addr string) (io.ReadWriteCloser, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "nacl.Dial")
	}
	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, errors.Wrap(err, "Dial.GenerateKey")
	}
	sReader := &secureReader{conn, pub, priv}
	sWriter := &secureWriter{conn, pub, priv}
	return &secureCon{sReader, sWriter, conn}, nil
}

// Serve starts a secure echo server on the given listener.
func Serve(l net.Listener) error {
	for {
		client, err := l.Accept()
		if err != nil {
			return errors.Wrap(err, "nacl.Serve")
		}
		defer client.Close()
		io.Copy(client, client)
	}
}

func main() {
	port := flag.Int("l", 0, "Listen mode. Specify port")
	flag.Parse()

	// Server mode
	if *port != 0 {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
		if err != nil {
			log.Fatal(err)
		}
		defer l.Close()
		log.Fatal(Serve(l))
	}

	// Client mode
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <port> <message>", os.Args[0])
	}
	conn, err := Dial("localhost:" + os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	if _, err := conn.Write([]byte(os.Args[2])); err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, len(os.Args[2]))
	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", buf[:n])
}
