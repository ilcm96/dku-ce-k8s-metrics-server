package metadata

import "os"

var NodeName string
var Namespace string

func init() {
	NodeName = os.Getenv("NODE_NAME")
	if NodeName == "" {
		NodeName = "unknown-node"
	}
	Namespace = os.Getenv("POD_NAMESPACE")
	if Namespace == "" {
		Namespace = "unknown-namespace"
	}
}
