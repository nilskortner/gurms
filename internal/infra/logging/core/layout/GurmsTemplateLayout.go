package layout

import (
	"gurms/internal/infra/cluster/node"
	"gurms/internal/infra/lang"
	"gurms/internal/infra/logging/core/model"
	"gurms/internal/supportpkgs/mathsupport"
	"strings"
)

var ESTIMATED_PATTERN_TEXT_LENGTH = 128
var LEVELS [][]byte
var NULL = []byte{'n', 'u', 'l', 'l'}
var COLON_SEPARATOR = []byte{' ', ':', ' '}
var TRACE_ID_LENGTH = 19
var CLASS_NAME_LENGTH = 40

var NODE_TYPE_AI_SERVING = int('A')
var NODE_TYPE_GATEWAY = int('G')
var NODE_TYPE_SERVICE = int('S')
var NODE_TYPE_UNKNOWN = int('U')

func init() {
	var levels = []model.LogLevel{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
	var levelCount = len(levels)
	LEVELS = make([][]byte, levelCount)
	maxLength := 0
	for _, level := range levels {
		maxLength = mathsupport.Max(len(level), maxLength)
	}
	for i := 0; i < levelCount; i++ {
		level := lang.PadStart(levels[i].String(), maxLength, ' ')
		level = strings.ToUpper(level)
		var err error
		LEVELS[i], err = lang.GetBytes(level)
		if err != nil {
			println(err)
		}
	}

}

type GurmsTemplateLayout struct {
	nodeType int
	nodeId   string
}

func NewGurmsTemplateLayout(nodeType node.NodeType, nodeId string) *GurmsTemplateLayout {
	var typ int
	noteId := nodeId
	switch nodeType {
	case node.AI_SERVING:
		typ = NODE_TYPE_AI_SERVING
	case node.GATEWAY:
		typ = NODE_TYPE_GATEWAY
	case node.SERVICE:
		typ = NODE_TYPE_SERVICE
	default:
		typ = NODE_TYPE_UNKNOWN
	}
	return &GurmsTemplateLayout{
		nodeType: typ,
		nodeId:   noteId,
	}
}
