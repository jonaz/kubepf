package kubectl

import (
	"os/exec"
	"strings"
)

type Pod struct {
	Name      string
	Namespace string
	Ports     []string
}

func GetPods(filter string) ([]*Pod, error) {
	out, err := Exec("kubectl", "get", "po", "--all-namespaces", `--template={{range .items}}{{ .metadata.namespace}} {{.metadata.name}} {{range .spec.containers}}{{range .ports}}{{.containerPort }},{{end}}{{end}} {{"\n"}}{{end}}`)
	pods := parsePods(out, filter)
	return pods, err
}

func parsePods(data []byte, filter string) []*Pod {

	rows := strings.Split(string(data), "\n")
	var pods []*Pod
	for _, row := range rows {
		tmp := strings.Split(row, " ")
		tmp = removeEmpty(tmp)

		if len(tmp) < 2 {
			continue
		}
		ports := []string{}
		if len(tmp) == 3 {
			for _, port := range removeEmpty(strings.Split(tmp[2], ",")) {
				ports = append(ports, port)
			}
		}
		pod := &Pod{
			Namespace: tmp[0],
			Name:      tmp[1],
			Ports:     ports,
		}

		if filter != "" && !strings.Contains(pod.Name, filter) {
			continue
		}
		pods = append(pods, pod)
	}
	return pods
}

func removeEmpty(data []string) []string {
	var newData []string
	for _, row := range data {
		if row == "" {
			continue
		}
		newData = append(newData, row)
	}
	return newData
}

func Exec(cmd string, args ...string) ([]byte, error) {
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return nil, err
	}

	return out, nil
}
