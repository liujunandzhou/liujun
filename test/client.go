package main

import "fmt"
import "net"
import "flag"
import "polling_server/encoding"

func checkError(hint string, err error) {
	if err != nil {
		fmt.Printf("%s:%s\n", hint, err)
	}
}

func main() {
	Dest := flag.String("host", "127.0.0.1:6872", "set host to connect")

	flag.Parse()

	conn, errDial := net.Dial("tcp", *Dest)

	checkError("net.Dial", errDial)

	defer conn.Close()

	pkg := proto.Package{"12345", "liujun", "hello", proto.TYPE_MSG}

	errEn := proto.Encode(conn, pkg)

	checkError("proto.Encode", errEn)

	for {

		pkg, errDec := proto.Decode(conn)

		checkError("proto.Decode", errDec)

		pkg.Debug()

		if errDec != nil {
			break
		}
	}
}
