// Problem: Wanted to expose nodeport on host.


// Solution 1: Tried simple port forwarding from host machine to VM machine.
// Challenge 1: simple port forwarding is setup to forward request from host machine to VM machine
// 		directs request to user-space of VM not kernal space 
//      where simple application can listen.
//  	lets say if i have golang server, it will listen in user space and
//      you can find if any app is listening on that particular port using netstat.
// 		Normal app server will work fine with this port forwarding and will receive request from host machine.

// 		but WHEN IT COMES TO K8s, we have to use nodeport to expose the service to the outside cluster.
// 		let say port 30080, but Kubernetes NodePort (30080) is handled by KERNEL-LEVEL IPTABLES, not a user-space process.
// 		and thats why you won't find any app listening on that port (30080) using netstat. it will be han
// 		so when request made from host machine to port 30080, it will try to forward to VM machine user space, 
// 		but in VM there is no app listening to 30080 port in user space... so request fails.

//LIMA (vm of type vz) LIMITATION OF BRIDGE NETWORK: other solution is to use port forwarding with IP not only locahost... 
// 		here comes the lima/MAC vz type VM LIMITATION, it do not support bridge or host network access.
//   	otherwise with other vm provider OR lima vm with QEMU(via socket_vmnet), we can make the request forward from host to vm by enabling bridge n/w



// so next solution is to use reverse proxy (kind of): user-space reverse proxy


// what it will do:
// Mac (curl localhost:29999) 
//  → Lima forwards to VM:29999 - (within user space)
//    → Go app listens on 29999, running within cluser, forwards to localhost:30080 (NodePort) - (to kernal iptable)
//      → iptables → K8s Service → Pod:8080





package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	// Target is localhost:30080 (NodePort service)
	targetURL, err := url.Parse("http://localhost:30080")
	if err != nil {
		log.Fatalf("Failed to parse target URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Log incoming request
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		proxy.ServeHTTP(w, r)
	})

	log.Println("Proxy listening on :29999 → forwarding to :30080")
	err = http.ListenAndServe(":29999", nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}


// FROM YOUR HOST: YOU WILL ALWAYS DO LOCALHOST (because, with LIMA, only localhost port forward is possible):
// to make it accessible via ip http://<ip>:29999, you can use qemu with bridge network 

// curl http://localhost:29999
// http://localhost:29999/api/users/1