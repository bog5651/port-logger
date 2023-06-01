package main

import (
	"bytes"
	"context"
	"crypto/aes"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"os"
	"port-logger/pkg/logger"
	"time"
)

//	func main() {
//		ctx := context.Background()
//		server := app.Server{}
//
//		conf := config.NewConfig()
//		conf.Init()
//
//		r := server.Setup()
//
//		addr := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
//		logger.Infof(ctx, "Listening on: %s", addr)
//
//		err := http.ListenAndServe(addr, r)
//		if err != nil {
//			logger.Errorf(ctx, "Listening error: %s", err)
//		}
//	}

var InData = []byte{0x62, 0xE8, 0xD8, 0xCD, 0xB1, 0x6A, 0x3B, 0xBB, 0x02, 0xB4, 0xDE, 0x94, 0x2E, 0x73, 0xB1, 0x2D}
var NeedData = []byte{0x00, 0x00, 0x00, 0x09, 0x2E, 0x2E, 0x2F, 0x2E, 0x2E, 0x2F, 0x2E, 0x2E, 0x2F, 0x00, 0x00, 0x00} //
var NeedSHA1 = []byte{0x3F, 0x87, 0x5C, 0x20, 0x43, 0xFD, 0x1A, 0xB1, 0xE2, 0xFA, 0xE3, 0x72, 0xF4, 0xE3, 0x60, 0x92, 0xA0, 0x21, 0x5C, 0x79}

var InDataStr = "62E8D8CDB16A3BBB02B4DE942E73B12D"
var NeedDataStr = "000000092E2E2F2E2E2F2E2E2F000000"
var NeedSHA1Str = "3F875C2043FD1AB1E2FAE372F4E36092A0215C79"

var key = []byte{0, 0, 0, 0x7C, 0x0C, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

//stopped at  0000007C0C000000000000000000000000000000000000000000000000000000

func inc(index int) error {
	if key[index] == 255 {
		if index+1 >= 32 {
			logger.Infof(context.Background(), "next inc may out of array %d, %X", index, key)
			return fmt.Errorf("next inc may out of array %d, %X", index, key)
		}
		key[index] = 0
		err := inc(index + 1)
		//if index == 2 {
		//	logger.Infof(context.Background(), "reach key %X", key)
		//}
		return err
	} else {
		key[index] = key[index] + 1
		return nil
	}
}

func KeyTimer() {
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-ticker.C:
			logger.Infof(context.Background(), "[%d] reach key %X", time.Now().Unix(), key)
		}
	}
}

func main() {
	ctx := context.Background()

	if len(os.Args) < 2 || len(os.Args) > 3 {
		logger.Criticalf(ctx, "Use exe [from hex]")
		return
	}
	from := os.Args[1]

	fromKey, err := hex.DecodeString(from)
	if err != nil {
		logger.Criticalf(ctx, "invalid from")
		return
	}

	if len(fromKey) != 32 {
		logger.Criticalf(ctx, "invalid from len")
		return
	}
	key = fromKey
	logger.Infof(ctx, "Starting from hex %X", key)

	l, e := big.NewInt(32), big.NewInt(255)

	l.Exp(l, e, nil)

	go KeyTimer()

	//logger.Infof(ctx, "len %s", l.String())

	one := big.NewInt(1)

	for i := big.NewInt(0); i.Cmp(l) < 1; i.Add(i, one) {
		//logger.Infof(ctx, "testing key %X", key)
		//e := Encrypt(key, InData)

		//logger.Infof(ctx, "enc data %X", e)
		//if bytes.Equal(e, NeedData) {
		//	logger.Warningf(ctx, "found key %s", string(key))
		//	return
		//}
		d := Decrypt(key, InData)
		//logger.Infof(ctx, "decr data %X", d)
		//return
		if bytes.Equal(d, NeedData) {
			logger.Infof(ctx, "decr data %X", d)
			logger.Warningf(ctx, "found key %X", key)
			return
		}
		err := inc(0)
		if err != nil {
			logger.Warningf(ctx, "error inc %s", err)
			return
		}
	}

	//plainText := "Hello 8gwifi.org"
	//ct := Encrypt([]byte(key), plainText)
	//fmt.Printf("Original Text:  %s\n", plainText)
	//fmt.Printf("AES Encrypted Text:  %s\n", ct)
}

func Decrypt(key []byte, ct []byte) []byte {
	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Errorf("NewCipher(%d bytes) = %s", len(key), err)
		panic(err)
	}
	plain := make([]byte, len(ct))
	c.Decrypt(plain, ct)
	return plain
}

func Encrypt(key []byte, d []byte) []byte {
	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Errorf("NewCipher(%d bytes) = %s", len(key), err)
		panic(err)
	}
	out := make([]byte, len(d))
	c.Encrypt(out, d)

	return out
}

//func main() {
//	conf := config.NewConfig()
//	conf.Init()
//	listener, _ := net.Listen("tcp", "127.0.0.1:27312") // открываем слушающий сокет
//	for {
//		conn, err := listener.Accept() // принимаем TCP-соединение от клиента и создаем новый сокет
//		if err != nil {
//			continue
//		}
//		go handleClient(conn) // обрабатываем запросы клиента в отдельной го-рутине
//	}
//}

func handleClient(conn net.Conn) {
	defer conn.Close() // закрываем сокет при выходе из функции

	buf := make([]byte, 32) // буфер для чтения клиентских данных
	for {
		//conn.Write([]byte{1}) // пишем в сокета

		rLen, err := conn.Read(buf) // читаем из сокета
		logger.Infof(context.Background(), "reded len %d", rLen)
		if err != nil {
			fmt.Println(err)
			break
		}
		logger.Infof(context.Background(), "reded len %s", string(buf))

		wlen, err := conn.Write([]byte("pal-is-up")) // пишем в сокет
		if err != nil {
			fmt.Println(err)
			break
		}
		logger.Infof(context.Background(), "send len %d", wlen)

		conn.Write([]byte("pal-is-up")) // пишем в сокет
		//for {
		//	wlen, err = conn.Write([]byte{1}) // пишем в сокет
		//	if err != nil {
		//		fmt.Println(err)
		//		break
		//	}
		//	logger.Infof(context.Background(), "send len %d", wlen)
		//}
	}
}
