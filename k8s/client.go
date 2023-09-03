package k8s

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"math/rand"
	"path/filepath"
	"reflect"
	"time"
)

func NewKubeClient() (*kubernetes.Clientset, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	return kubernetes.NewForConfig(config)
}

func PatchNodeStatus() {
	clientset, err := NewKubeClient()
	if err != nil {
		panic(err)
	}
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	if len(nodes.Items) == 0 {
		panic(fmt.Errorf("no node in cluster"))
	}
	node := nodes.Items[0]
	newNode := node.DeepCopy()

	// shuffle condition
	var a []int
	for i := 0; i < len(node.Status.Conditions); i++ {
		a = append(a, i)
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })

	var newConditions []v1.NodeCondition
	for i := 0; i < len(a); i++ {
		newConditions = append(newConditions, node.Status.Conditions[a[i]])
	}

	for t, update := range newConditions {
		if !reflect.DeepEqual(newNode.Status.Conditions[t], update) {
			newNode.Status.Conditions[t] = update
		}
	}
	if err := SetConditions(node.Name, newNode.Status.Conditions); err != nil {
		panic(err)
	}

}

func SetConditions(nodeName string, newConditions []v1.NodeCondition) error {
	for i := range newConditions {
		// Each time we update the conditions, we update the heart beat time
		newConditions[i].LastHeartbeatTime = metav1.NewTime(time.Now())
	}
	patch, err := generatePatch(newConditions)
	if err != nil {
		return err
	}

	clientset, err := NewKubeClient()
	if err != nil {
		panic(err)
	}

	return clientset.RESTClient().Patch(types.StrategicMergePatchType).Resource("nodes").
		Name(nodeName).SubResource("status").Body(patch).Do(context.TODO()).Error()
}

// generatePatch generates condition patch
func generatePatch(conditions []v1.NodeCondition) ([]byte, error) {
	raw, err := json.Marshal(&conditions)
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf(`{"status":{"conditions":%s}}`, raw)), nil
}
