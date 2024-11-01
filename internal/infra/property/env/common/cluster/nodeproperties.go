package cluster

const NODE_ID_MAX_LENGTH int = 32
const NODE_NAME_MAX_LENGTH int = 128

type NodeProperties struct {
	// "The node ID must start with a letter or underscore, "
	//   - "and matches zero or more of characters [a-zA-Z0-9_] after the beginning. "
	//   - "e.g. \"turms001\", \"turms_002\". "
	//   - "A node must have a unique ID. "
	//   - "If not specified, Turms server will generate a random unique ID"
	Id string
	// "The node name must start with a letter or underscore, "
	//   - "and matches zero or more of characters [a-zA-Z0-9_] after the beginning. "
	//   - "e.g. \"turms001\", \"turms_002\". "
	//   - "The node name can be duplicate in the cluster. "
	//   - "If not specified, Turms server will use the node ID as the node name"
	Name string
	// "e.g. \"us-east-1\" and \"ap-east-1\""
	Zone            string
	ActiveByDefault bool
	// "Only works when it is a turms-service node"
	LeaderEligible bool
	// "The priority to be a leader"
	Priority int
}

func InitNodeProperties() *NodeProperties {
	return &NodeProperties{
		ActiveByDefault: true,
		LeaderEligible:  true,
	}
}
