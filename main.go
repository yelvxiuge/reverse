// etcd-utils project main.go
package main

import (
	"etcd-utils/lib"
	"flag"

	//"fmt"
	"log"
	"net/http"
	//	"strings"
)

func main() {
	port := flag.String("p", "2000", "local port")
	flag.Parse()
	
	
//	machines := []string{"http://192.168.50.30:2379"}
//	idle, err := lib.GetServicesIdle(machines, "ms-nginx")
//	if err != nil {
//		log.Fatalln(err)
//	}
//	var array []string
//	for _, v := range idle.Hosts {
//		link := fmt.Sprintf("%s:%v", v, idle.Port)
//		array = append(array, link)
//	}
	
	
	
	
	//s := []string{"127.0.0.1:10001", "127.0.0.1:10002"}
	//h := &handle{targets:s, Counter: 1}
//	h := &lib.Handle{Targets: array, Counter: 1}

	h := &lib.Handle{}
	h.Counter =1
	print("Listen on " + *port)
	err := http.ListenAndServe(":"+*port, h)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
	

}

//func main() {
//	machines := []string{"http://192.168.50.30:2379"}
//	idle, err := lib.GetServicesIdle(machines, "ms-nginx")
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(idle)
//}
