package cluster

const NODE_ID_MAX_LENGTH int = 32
const NODE_NAME_MAX_LENGTH int = 128

type NodeProperties struct {
	id              string
	name            string
	zone            string
	activeByDefault bool
	leaderEligible  bool
	priority        int
}

// "The node ID must start with a letter or underscore, "
//   - "and matches zero or more of characters [a-zA-Z0-9_] after the beginning. "
//   - "e.g. \"turms001\", \"turms_002\". "
//   - "A node must have a unique ID. "
//   - "If not specified, Turms server will generate a random unique ID"
var id string = ""

// "The node name must start with a letter or underscore, "
//   - "and matches zero or more of characters [a-zA-Z0-9_] after the beginning. "
//   - "e.g. \"turms001\", \"turms_002\". "
//   - "The node name can be duplicate in the cluster. "
//   - "If not specified, Turms server will use the node ID as the node name"
var name string = ""

// "e.g. \"us-east-1\" and \"ap-east-1\""
var zone string = ""

var activeByDefault bool = true

// "Only works when it is a turms-service node"
var leaderEligible bool = true

// "The priority to be a leader"
var priority int
