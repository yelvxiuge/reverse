package lib

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Handle struct {
	Counter int
	Targets []string
}

func getSer(url string) string {
	ele := strings.Split(url, "/")
	if ele[2] != "server" {

		fmt.Println("is not server")
		return "err"
	}
	if ele[3] != "ms" {
		fmt.Println("is not ms")
		return "err"
	}
	return ele[len(ele)-1]

}
func (this *Handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	machines := []string{"http://192.168.50.30:2379"}
	//	idle, err := GetServicesIdle(machines, "ms-nginx")
	if r.RequestURI == "/" {
		http.Redirect(w, r, " ", 404)
		r.Body.Close()
		return

	}

	ser := getSer(r.RequestURI)
	if ser == "err" {
		http.Redirect(w, r, " ", 404)
		r.Body.Close()
		return
	}
	r.URL.Path = "/"
	idle, err := GetServicesIdle(machines, ser)
	if err != nil {
		log.Fatalln(err)
	}
	var array []string
	for _, v := range idle.Hosts {
		link := fmt.Sprintf("%s:%v", v, idle.Port)
		array = append(array, link)
	}

	this.Targets = array
	remote, err := url.Parse("http://" + this.getTarget())
	if err != nil {
		panic(err)

	}
	//	print(r.Proto)
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)

	//	resp ,err:=http.Get(remote.String()+":443")
	//	if err != nil{
	//		panic(err)
	//	}
	//	body,err := ioutil.ReadAll(resp.Body)
	//	if err != nil{
	//		panic(err)
	//	}
	//	w.Write(body)
	//	defer resp.Body.Close()

}
func (this *Handle) getTarget() string {
	if this.Counter >= len(this.Targets)-1 {
		this.Counter = 0
	} else {
		this.Counter += 1
	}
	return this.Targets[this.Counter]

}
