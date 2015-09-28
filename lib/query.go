package lib

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"strings"
)

type ServiceOption struct {
	HostIp string
	Port   string
}
type ServiceIdle struct {
	Port  int
	Hosts []string
}

func GetServicesIdle(machines []string, name string) (*ServiceIdle, error) {
	registrys, err := getRegistrys(machines)
	if err != nil {
		return nil, err
	}
	var (
		services     map[string]string = make(map[string]string)
		pods         map[string]string = make(map[string]string)
		idle         *ServiceIdle      = new(ServiceIdle)
		serviceValue string
	)
	for k, v := range registrys {
		switch {
		case strings.Contains(k, "/services/specs"):
			services[k] = v
		case strings.Contains(k, "pods") && strings.Contains(k, name):
			pods[k] = v
		}
	}
	for k, v := range services {
		if strings.Contains(k, name) {
			serviceValue = v
		}
	}
	idle.Port, _ = resolvePort(serviceValue)
	for _, v := range pods {
		ip, _ := resolveIp(v)
		if len(strings.TrimSpace(ip)) > 0 {
			idle.Hosts = append(idle.Hosts, ip)
		}

	}
	return idle, nil
}
func resolveIp(meta string) (string, error) {
	type (
		Labels struct {
			Run string `json:"run"`
		}
		MetaData struct {
			Name              string `json:"name"`
			GenerateName      string `json:"generateName"`
			NameSpace         string `json:"namespace"`
			SelfLink          string `json:"selfLink"`
			Uid               string `json:"uid"`
			ResourceVersion   string `json:"resourceVersion"`
			CreationTimestamp string `json:"creationTimestamp"`
			Labels            Labels `json:"labels"`
		}

		Status struct {
			Phase  string `json:"phase"`
			HostIP string `json:"hostIP"`
		}
		ServicePod struct {
			Kind       string   `json:"kind"`
			ApiVersion string   `json:"apiVersion"`
			Metadata   MetaData `json:"metadata"`
			Status     Status   `json:"status"`
		}
	)
	var model ServicePod
	err := json.Unmarshal([]byte(meta), &model)
	if err != nil {
		return "", err
	}
	ip := model.Status.HostIP
	if model.Status.Phase == "Running" {
		return ip, nil
	}
	return "", nil
}
func resolvePort(meta string) (int, error) {
	type (
		Port struct {
			Protocol   string `json:"protocol"`
			XPort      int    `json:"port"`
			TargetPort int    `json:"targetPort"`
			NodePort   int    `json:"nodePort"`
		}
		Spec struct {
			Ports           []Port `json:"ports,omitempty"`
			ClusterIP       string `json:"clusterIP"`
			Type            string `json:"type"`
			SessionAffinity string `json:"sessionAffinity"`
		}
		MetaData struct {
			Name              string `json:"name"`
			NameSpace         string `json:"namespace"`
			Uid               string `json:"uid"`
			CreationTimestamp string `json:"creationTimestamp"`
		}
		ServiceSpec struct {
			Kind       string   `json:"kind"`
			ApiVersion string   `json:"apiVersion"`
			Metadata   MetaData `json:"metadata,omitempty"`
			Spec       Spec     `json:"spec,omitempty"`
		}
	)

	var model ServiceSpec
	err := json.Unmarshal([]byte(meta), &model)
	if err != nil {
		return 0, err
	}
	return model.Spec.Ports[0].NodePort, nil
}
func GetServicesOption(machines []string) ([]ServiceOption, error) {
	registrys, err := getRegistrys(machines)
	if err != nil {
		return nil, err
	}
	var (
		services map[string]string = make(map[string]string)
		pods     map[string]string = make(map[string]string)
		options  []ServiceOption
	)
	for k, v := range registrys {
		switch {
		case strings.Contains(k, "/services/specs"):
			services[k] = v
		case strings.Contains(k, "pods"):
			pods[k] = v
		}
	}
	for _, v := range services {
		fmt.Println(v)
		port, node, _ := getServicePort(v)
		fmt.Println(node)
		fmt.Println(port)
		pod := searchMap(node, pods)
		ip, isRun, _ := getServiceIp(pod)
		if isRun {
			options = append(options, ServiceOption{HostIp: ip, Port: port})
		}
	}
	return options, nil
}
func searchMap(filter string, data map[string]string) string {
	for k, v := range data {
		if strings.Contains(k, filter) {
			return v
		}
	}
	return ""
}
func getServiceIp(meta string) (string, bool, error) {
	type (
		Labels struct {
			Run string `json:"run"`
		}
		MetaData struct {
			Name              string `json:"name"`
			GenerateName      string `json:"generateName"`
			NameSpace         string `json:"namespace"`
			SelfLink          string `json:"selfLink"`
			Uid               string `json:"uid"`
			ResourceVersion   string `json:"resourceVersion"`
			CreationTimestamp string `json:"creationTimestamp"`
			Labels            Labels `json:"labels"`
		}

		Status struct {
			Phase  string `json:"phase"`
			HostIP string `json:"hostIP"`
		}
		ServicePod struct {
			Kind       string   `json:"kind"`
			ApiVersion string   `json:"apiVersion"`
			Metadata   MetaData `json:"metadata"`
			Status     Status   `json:"status"`
		}
	)
	var (
		model ServicePod
		flags bool = false
	)
	err := json.Unmarshal([]byte(meta), &model)
	if err != nil {
		return "", false, err
	}
	ip := model.Status.HostIP

	if model.Status.Phase == "Running" {
		flags = true
	}
	return ip, flags, nil
}
func getServicePort(meta string) (string, string, error) {
	type (
		Port struct {
			Protocol   string `json:"protocol"`
			Port       string `json:"port"`
			TargetPort string `json:"targetPort"`
			NodePort   string `json:"nodePort"`
		}
		Labels struct {
			Component string `json:"component,omitempty"`
			Provider  string `json:"provider,omitempty"`
			Run       string `json:"run,omitempty"`
		}
		Selector struct {
			Run string `json:"run"`
		}
		LoadBalancer struct {
		}
		Status struct {
			LoadBalancer LoadBalancer `json:"loadBalancer,omitempty"`
		}
		Spec struct {
			Ports           []Port   `json:"ports,omitempty"`
			Selector        Selector `json:"selector,omitempty"`
			ClusterIP       string   `json:"clusterIP"`
			Type            string   `json:"type"`
			SessionAffinity string   `json:"sessionAffinity"`
		}
		MetaData struct {
			Name              string `json:"name"`
			NameSpace         string `json:"namespace"`
			Uid               string `json:"uid"`
			CreationTimestamp string `json:"creationTimestamp"`
			Labels            Labels `json:"labels,,omitempty"`
		}
		ServiceSpec struct {
			Kind       string   `json:"kind"`
			ApiVersion string   `json:"apiVersion"`
			Metadata   MetaData `json:"metadata,omitempty"`
			Spec       Spec     `json:"spec,omitempty"`
			Status     Status   `json:"status,omitempty"`
		}
	)

	var model ServiceSpec
	err := json.Unmarshal([]byte(meta), &model)
	if err != nil {
		return "", "", err
	}
	fmt.Println(model.Spec.Ports)
	fmt.Println("我叫MT")
	port := model.Spec.Ports[0].NodePort
	label := model.Metadata.Labels.Run
	return port, label, nil
}
func getValueMapByKind(key string, data map[string]string) (map[string]interface{}, error) {
	var (
		list map[string]interface{}
		err  error
	)
	if v, ok := data[key]; ok {
		err = json.Unmarshal([]byte(v), &list)
	}
	return list, err
}
func getRegistrys(machines []string) (map[string]string, error) {
	c := etcd.NewClient(machines)
	resp, err := c.Get("/registry", true, true)
	if err != nil {
		return nil, err
	}
	list := make(map[string]string)
	recursive(resp.Node, list)
	return list, nil
}

func recursive(node *etcd.Node, values map[string]string) {
	if node == nil {
		return
	}
	key := node.Key
	if !node.Dir {
		values[key] = node.Value
	} else {
		for _, subNode := range node.Nodes {
			recursive(subNode, values)
		}
	}
}
