# kubepf
interactive client to proxy to a pod in kubernetes cluster
It will take first available port after 8080 on localhost and then open the browser. 

### Installation

```
go get -u github.com/jonaz/kubepf
```

### Usage
```
kubepf -h
Usage of kubepf:
  -pod string
    	Specify pod filter
```


### Example
```
kubepf -pod traefik
0 - traefik-ingress-controller-9lxzc kube-system
1 - traefik-ingress-controller-bjfap kube-system
2 - traefik-ingress-controller-dsxtw kube-system
3 - traefik-ingress-controller-lkyja kube-system
Select pod number and press enter: 0
Chosen pod: traefik-ingress-controller-9lxzc
0 - 80
1 - 8081
Select port in the pod: 1
Chosen port:  8081
TCP Port "8080" is availableStarting kubectl port-forward...
kubectl port-forward traefik-ingress-controller-9lxzc 8080:8081 --namespace kube-system
Forwarding from 127.0.0.1:8080 -> 8081

```
