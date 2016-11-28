package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/jonaz/kubepf/kubectl"
)

var podFilter string

func init() {
	flag.StringVar(&podFilter, "pod", "", "Specify pod filter")
	flag.Parse()
}

func main() {
	pods, err := kubectl.GetPods(podFilter)

	if err != nil {
		fmt.Println(err)
		return
	}

	for k, pod := range pods {
		fmt.Printf("%d - %s %s\n", k, pod.Name, pod.Namespace)
	}

	podKey, err := readInput("Select pod number and press enter: ")
	if err != nil {
		fmt.Println(err)
		return
	}

	pod := pods[podKey]
	fmt.Println("Chosen pod:", pod.Name)

	for k, port := range pod.Ports {
		fmt.Printf("%d - %s\n", k, port)
	}
	portKey, err := readInput("Select port in the pod: ")
	if err != nil {
		fmt.Println(err)
		return
	}
	containerPort := pod.Ports[portKey]
	fmt.Println("Chosen port: ", containerPort)

	startPort := 8080
	for {
		if isPortFree(startPort) {
			break
		}
		startPort = startPort + 1
	}
	portString := strconv.Itoa(startPort)
	go openBrowser(portString)

	PortForward(pod, portString, containerPort)
}

func openBrowser(port string) {
	time.Sleep(time.Millisecond * 300)
	switch runtime.GOOS {
	case "linux":
		kubectl.Exec("xdg-open", "http://localhost:"+port)
	case "darwin":
		kubectl.Exec("open", "http://localhost:"+port)
	case "windows":
		kubectl.Exec("start", "http://localhost:8081"+port)
	}
}

func isPortFree(p int) bool {
	port := strconv.Itoa(p)
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't listen on port %q: %s", port, err)
		return false
	}
	ln.Close()
	fmt.Printf("TCP Port %q is available", port)
	return true
}

func PortForward(pod *kubectl.Pod, port, containerPort string) {
	fmt.Println("Starting kubectl port-forward...")
	fmt.Sprintf("%#v\n", pod)
	fmt.Println("kubectl", "port-forward", pod.Name, port+":"+containerPort, "--namespace", pod.Namespace)
	cmd := exec.Command("kubectl", "port-forward", pod.Name, port+":"+containerPort, "--namespace", pod.Namespace)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatalf("kubectl error: %v", err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatalf("kubectl error: %v", err)
	}
}

func readInput(txt string) (int, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(txt)
	tmp, _ := reader.ReadString('\n')
	if tmp == "\n" {
		return 0, nil
	}
	podKey, err := strconv.Atoi(strings.Trim(tmp, " \n"))
	if err != nil {
		return 0, err
	}
	return podKey, nil
}
