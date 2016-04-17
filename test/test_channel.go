package main

import "fmt"
import "time"
import "polling_server/channel"

func consume(id string, ch chan string) {
	for {

		str := <-ch

		fmt.Printf("Received: id:%s msg:%s\n", id, str)
	}
}

func checkRet(str string, ret bool) {

	if ret {
		fmt.Printf("%s is ok\n", str)
	} else {
		fmt.Printf("%s is fail\n", str)
	}
}

func main() {

	cm := manager.NewManager()

	ch1 := make(chan string, 2)
	checkRet("cm.Add", cm.Add("liujun", ch1))
	go consume("liujun", ch1)

	ch2 := make(chan string, 2)
	checkRet("cm.Add", cm.Add("liujun", ch2))
	go consume("linlianhuan", ch2)

	for i := 0; i < 1000; i++ {
		for _, id := range [2]string{"liujun", "linlianhuan"} {
			cm.Send(id, id)
		}

		time.Sleep(time.Second * 1)
	}

	cm.Close("liujun")
	cm.Close("linlianhuan")
}
