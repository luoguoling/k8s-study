# k8s增删改查

## 一.获取clientset对象

### 1.1创建admin账户

```yaml
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: k8s-authorize
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: k8s-authorize
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: k8s-authorize
  namespace: kube-system
```

### 1.2 获取admin的token

```bash
 #获取token,对api进行操作
 kubectl describe secrets $(kubectl get secrets -n kube-system |grep admin |cut -f1 -d ' ') -n kube-system |grep -E '^token' |cut -f2 -d':'|tr -d '\t'|tr -d ' '
```

### 1.3创建代理并在k8s-master运行

```go
package main
import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)
func main() {
	var tlsConfig = &tls.Config{
		InsecureSkipVerify: true, //忽略证书验证
	}
	var transport http.RoundTripper = &http.Transport{
		Proxy:                  nil,
		DialContext:            nil,
		Dial:                   nil,
		DialTLSContext:         nil,
		DialTLS:                nil,
		TLSClientConfig:        tlsConfig,
		TLSHandshakeTimeout:    0,
		DisableKeepAlives:      false,
		DisableCompression:     true,
		MaxIdleConns:           0,
		MaxIdleConnsPerHost:    0,
		MaxConnsPerHost:        0,
		IdleConnTimeout:        0,
		ResponseHeaderTimeout:  0,
		ExpectContinueTimeout:  0,
		TLSNextProto:           nil,
		ProxyConnectHeader:     nil,
		MaxResponseHeaderBytes: 0,
		WriteBufferSize:        0,
		ReadBufferSize:         0,
		ForceAttemptHTTP2:      false,
	}
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		//server,_ := url.Parse("https://10.206.16.18:16443")
		server, _ := url.Parse("https://10.0.12.9:8443")
		log.Println(request.URL.Path)
		p := httputil.NewSingleHostReverseProxy(server)
		p.Transport = transport
		p.ServeHTTP(writer, request)

	})
	log.Println("开始反向代理k8sapi")
	err := http.ListenAndServe("0.0.0.0:9090", nil)
	if err != nil {
		log.Fatal(err)
	}
}
```





### 1.3获取clientset对象

```go
package main
import (
	"context"
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)
// 实例化clientset对象
	config := rest.Config{
		Host: "http://121.5.106.67:9090",  //代理地址
		BearerToken: "eyJhbGciOiJSUzI1NiIsImtpZCI6IkNTaDRNUk1aSEs4YnBEVm5fZGw4RFZoN3VZQ3pkdV9mRHVmOGctWEVhVGsifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrOHMtYXV0aG9yaXplLXRva2VuLWdqNWJtIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6Ims4cy1hdXRob3JpemUiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiIzYTcwNzI4Ni1mZDMyLTQ4ZDEtYmU2Yi03YjAxOGFmNWUxMmIiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6a3ViZS1zeXN0ZW06azhzLWF1dGhvcml6ZSJ9.VAnQsm2oLxNIab0SpmAkKO3FgaSGjSWs24LZ_gh08nXcsps40_DDTJzUG2jFjOCAluOOUz2EzbuVbud7EN9wOSbkA7-DaBDe6v009HrFWZ0mWt3MUG2uEzFJCRP7v5ySYMtNGb8ORX-68UvVvOCGHrN0dHH2IAwtke6U9npg_sWU_wHX835C-NF05qWGk2n3dlVBFsCq6U6ntVFhEJnq48vAZA3RfMPHkEha8xKroSERSVQkbi28EVKaepimF9-LV5RBY4bzbjz8fCcC9ikvW2goggcQx4getIC9DR0NmB3qybfPdZ7ltWCOiE3lFWwELk0Rd4geb9CpWdbLojn_ug",
	}
	clientset, err := kubernetes.NewForConfig(&config)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("operation is %v\n", *operate)

```

## 二.k8s查询操作

### 2.1获取namespace中的pod信息

```go
pods, err := clientset.CoreV1().Pods("default").List(context.TODO(),metav1.ListOptions{})
if err != nil {
    panic(err)
}
// 循环打印pod的信息
for _,pod := range pods.Items {
    fmt.Println(pod.ObjectMeta.Name,pod.Status.Phase)
}
```

