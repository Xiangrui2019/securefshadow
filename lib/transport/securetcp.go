package transport

import (
	"io"
	"log"
	"net"

	"github.com/xiangrui2019/securefshadow/lib/encryption"
)

const (
	BUF_SIZE = 1024
)

type SecureTCPConnection struct {
	io.ReadWriteCloser
	Cipher *encryption.Cipher
}

func (secureSocket *SecureTCPConnection) DecodeRead(bs []byte) (n int, err error) {
	n, err = secureSocket.Read(bs)

	if err != nil {
		return
	}

	secureSocket.Cipher.Decode(bs[:n])
	return
}

func (secureSocket *SecureTCPConnection) EncodeWrite(bs []byte) (int, error) {
	secureSocket.Cipher.Encode(bs)
	return secureSocket.Write(bs)
}

func (secureSocket *SecureTCPConnection) EncodeCopy(dst io.ReadWriteCloser) error {
	buf := make([]byte, BUF_SIZE)

	for {
		readCount, err := secureSocket.Read(buf)

		if err != nil {
			if err != io.EOF {
				return err
			} else {
				return nil
			}
		}

		if readCount > 0 {
			writeCount, err := (&SecureTCPConnection{
				ReadWriteCloser: dst,
				Cipher:          secureSocket.Cipher,
			}).EncodeWrite(buf[0:readCount])

			if err != nil {
				return err
			}

			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
	}
}

func (secureSocket *SecureTCPConnection) DecodeCopy(dst io.Writer) error {
	buf := make([]byte, BUF_SIZE)

	for {
		readCount, err := secureSocket.DecodeRead(buf)

		if err != nil {
			if err != io.EOF {
				return err
			} else {
				return nil
			}
		}

		if readCount > 0 {
			writeCount, err := dst.Write(buf[0:readCount])

			if err != nil {
				return err
			}

			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
	}
}

func DialSecureTCP(readdr *net.TCPAddr, cipher *encryption.Cipher) (*SecureTCPConnection, error) {
	remoteConn, err := net.DialTCP("tcp", nil, readdr)
	if err != nil {
		return nil, err
	}

	remoteConn.SetLinger(0)

	return &SecureTCPConnection{
		ReadWriteCloser: remoteConn,
		Cipher:          cipher,
	}, nil
}

func ListenSecureTCP(laddr *net.TCPAddr, cipher *encryption.Cipher, handleConn func(localConn *SecureTCPConnection), didListen func(listenAddr *net.TCPAddr)) error {
	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	if didListen != nil {
		go didListen(listener.Addr().(*net.TCPAddr))
	}

	for {
		localConn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}

		localConn.SetLinger(0)
		go handleConn(&SecureTCPConnection{
			ReadWriteCloser: localConn,
			Cipher:          cipher,
		})
	}
}
