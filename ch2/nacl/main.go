// -*- mode:go;mode:go-playground -*-
// snippet of code @ 2017-04-02 10:11:58

// === Go Playground ===
// Execute the snippet with Ctl-Return
// Remove the snippet completely with its dir and all files M-x `go-playground-rm`

// Writing secure reader and secure wirter
package main

import (
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/nacl/box"
)

type secureReader struct {
	r         io.Reader
	pub, priv *[32]byte
}

func NewSecureReader(r io.Reader, priv, pub *[32]byte) io.Reader {
	return &secureReader{
		r, pub, priv,
	}
}

func (sr *secureReader) Read(p []byte) (int, error) {
	var nonce [24]byte

	if n, err := io.ReadFull(sr.r, nonce[:]); err != nil {
		return n, err
	}

	msg := make([]byte, len(p)+box.Overhead)

	n, err := sr.r.Read(msg)
	if err != nil {
		return n, err
	}
	dec, ok := box.Open(p[:0], msg[:n], &nonce, sr.pub, sr.priv)
	if !ok {
		return 0, errors.New("secure-read: error decrypting message")
	}
	return len(dec), nil
}

type secureWriter struct {
	w         io.Writer
	pub, priv *[32]byte
}

func NewSecureWriter(w io.Writer, priv, pub *[32]byte) io.Writer {
	return &secureWriter{
		w, pub, priv,
	}
}

func (sw *secureWriter) Write(p []byte) (int, error) {
	var nonce [24]byte

	if n, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return n, err
	}

	enc := box.Seal(nonce[:], p, &nonce, sw.pub, sw.priv)

	if n, err := sw.w.Write(enc); err != nil {
		return n, err
	}
	return len(p), nil

}

func Dial(addrs string) (io.ReadWriteCloser, error) {
	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial("tcp", addrs)
	if err != nil {
		return nil, err
	}

	// exchange keys

	// read server public key

	serverPub := new([32]byte)
	if _, err := io.ReadFull(conn, serverPub[:]); err != nil {
		return nil, err
	}

	// write client public key

	if _, err := conn.Write(pub[:]); err != nil {
		return nil, err
	}

	return struct {
		io.Reader
		io.Writer
		io.Closer
	}{
		NewSecureReader(conn, priv, serverPub),
		NewSecureWriter(conn, priv, serverPub),
		conn,
	}, nil
}

func Serve(l net.Listener) error {
	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go handleConn(conn, pub, priv)
	}
}

func handleConn(conn net.Conn, pub, priv *[32]byte) {
	defer conn.Close()

	// write server public key
	if _, err := conn.Write(pub[:]); err != nil {
		log.Println("public key write error", err)
		return
	}

	// receive client's public key
	clientPub := new([32]byte)

	if _, err := io.ReadFull(conn, clientPub[:]); err != nil {
		log.Println("client public key read error", err)
	}

	sW := NewSecureWriter(conn, priv, clientPub)
	sR := NewSecureReader(conn, priv, clientPub)

	if _, err := io.Copy(sW, sR); err != nil {
		log.Println("conn failed", err)
		return
	}
}

func main() {
	port := flag.Int("l", 0, "Listen mode. Specify port")
	flag.Parse()

	// server mode
	if *port != 0 {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
		if err != nil {
			log.Fatal(err)
		}
		defer l.Close()
		log.Fatal(Serve(l))
	}

	// client mode
	if (len(os.Args)) != 3 {
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
	fmt.Println(string(buf[:n]))
}
