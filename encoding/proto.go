package proto

import "errors"
import "bytes"
import "io"
import "net"
import "fmt"
import "encoding/json"
import "encoding/binary"

const MAX_PACKAGE = 1 * 1024 * 1024

const (
	TYPE_HELLO int = iota + 1
	TYPE_ACK
	TYPE_ECHO
	TYPE_MSG
	TYPE_INFO
)

type Package struct {
	Id    string
	Topic string
	Msg   string
	Type  int
}

var BadPackage = errors.New("bad package")

func (p *Package) Valid() bool {

	if p.Type == TYPE_ECHO || p.Type == TYPE_MSG {
		if len(p.Id) <= 0 {
			return false
		}
	}

	return true
}

func (p *Package) Debug() {
	fmt.Printf("Id:%s Topic:%s Msg:%s Type:%d\n", p.Id, p.Topic, p.Msg, p.Type)
}

func Encode(conn net.Conn, pkg Package) error {

	msg, err := json.Marshal(pkg)

	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(nil)

	var bodyLen uint32 = uint32(len(msg))

	binary.Write(buffer, binary.BigEndian, &bodyLen)

	buffer.Write(msg)

	ret, err := conn.Write(buffer.Bytes())

	if err != nil && ret != buffer.Len() {
		return err
	}

	return nil
}

func Decode(conn net.Conn) (Package, error) {

	bbodyLen := make([]byte, 4, 4)

	pkg := Package{}

	n, err := io.ReadFull(conn, bbodyLen)

	if err != nil || n != 4 {
		return pkg, err
	}

	buffer := bytes.NewBuffer(bbodyLen)

	var ibodyLen uint32 = 0

	binary.Read(buffer, binary.BigEndian, &ibodyLen)

	if ibodyLen <= 0 || ibodyLen > MAX_PACKAGE {
		return pkg, BadPackage
	}

	bjsonBody := make([]byte, ibodyLen)

	_, errRead := io.ReadFull(conn, bjsonBody)

	if errRead != nil {
		return pkg, err
	}

	errUnmarshal := json.Unmarshal(bjsonBody, &pkg)

	return pkg, errUnmarshal
}
