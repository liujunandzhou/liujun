package main

import "flag"
import "fmt"
import "os"
import "time"
import "net"
import "net/http"

import "polling_server/encoding"
import "polling_server/channel"
import "polling_server/uniq"

const (
	defaultListenHost  string = ":6872"
	defaultMessageHost string = ":6873"
)

var g_cm *manager.ChManager = nil

func init() {

	g_cm = manager.NewManager()
}

func checkError(hint string, err error, exit bool) bool {

	if err != nil {

		fmt.Printf("%s:%s\n", hint, err)

		if exit {
			os.Exit(1)
		}

		return false
	}

	return true
}

func handleConn(conn net.Conn) {

	defer conn.Close()

	pkg, err := proto.Decode(conn)

	if !checkError("proto.Decode", err, false) {
		return
	}

	pkg.Debug()

	sessionId := uniq.Uniq()

	pkg.Id = sessionId

	messageChannel := make(chan string, 2)

	g_cm.Add(sessionId, messageChannel)

	//发送探活包
	timeoutTimer := time.NewTimer(time.Second * 10)

	for {

		var errEnc error = nil

		select {
		case msg := <-messageChannel:
			//发送到客户端
			pkg.Type = proto.TYPE_MSG
			pkg.Msg = msg
			errEnc = proto.Encode(conn, pkg)

		case <-timeoutTimer.C:
			//探活
			pkg.Type = proto.TYPE_ECHO
			pkg.Msg = "echo"
			errEnc = proto.Encode(conn, pkg)

			timeoutTimer.Reset(time.Second * 10)
		}

		if !checkError("proto.Encode", errEnc, false) {
			break
		}
	}

	g_cm.Close(sessionId)
}

func handlePub(w http.ResponseWriter, req *http.Request) {

	idStr := req.FormValue("id")

	msgStr := req.FormValue("msg")

	fmt.Printf("id:%s msg:%s\n", idStr, msgStr)

	retErr := g_cm.Send(idStr, msgStr)

	callBack := "ok"

	if retErr != nil {
		callBack = fmt.Sprintf("%s", retErr)
	}

	w.Write([]byte(callBack))
}

//消息监听crontine
func messageServer(listen string) {

	http.HandleFunc("/pub", handlePub)

	http.ListenAndServe(listen, nil)
}

func main() {

	var listenHost string
	var messageHost string

	flag.StringVar(&listenHost, "phost", defaultListenHost, "polling host")
	flag.StringVar(&messageHost, "mhost", defaultMessageHost, "message host")

	flag.Parse()

	fmt.Printf("message server listen at :%s\n", messageHost)
	go messageServer(messageHost)

	fmt.Printf("polling server listen at: %s\n", listenHost)
	listenConn, err := net.Listen("tcp", listenHost)

	defer listenConn.Close()

	checkError("net.Listen", err, true)

	for {
		conn, err := listenConn.Accept()

		if err != nil {
			fmt.Println("Accpet failed")
			continue
		}

		go handleConn(conn)
	}
}
