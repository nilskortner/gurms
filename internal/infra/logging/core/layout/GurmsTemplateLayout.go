package layout

import (
	"gurms/internal/infra/cluster/node"
	"gurms/internal/infra/lang"
	"gurms/internal/infra/logging/core/model/loglevel"
	"gurms/internal/supportpkgs/mathsupport"
	"strings"
)

var ESTIMATED_PATTERN_TEXT_LENGTH = 128
var LEVELS [][]byte
var NULL = []byte{'n', 'u', 'l', 'l'}
var COLON_SEPARATOR = []byte{' ', ':', ' '}
var TRACE_ID_LENGTH = 19
var STRUCT_NAME_LENGTH = 40

var NODE_TYPE_AI_SERVING = int('A')
var NODE_TYPE_GATEWAY = int('G')
var NODE_TYPE_SERVICE = int('S')
var NODE_TYPE_UNKNOWN = int('U')

func init() {
	var levels = []loglevel.LogLevel{0, 1, 2, 3, 4, 5}
	var levelCount = len(levels)
	LEVELS = make([][]byte, levelCount)
	maxLength := 0
	for _, level := range levels {
		maxLength = mathsupport.Max(len(level.String()), maxLength)
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

func FormatStructName(name string) []byte {
	rawBytes := []byte(name)
	if len(rawBytes) == STRUCT_NAME_LENGTH {
		return rawBytes
	}
	if len(rawBytes) > STRUCT_NAME_LENGTH {
		parts := lang.TokenizeToStringArray(name, ".")
		structName := parts[len(parts)-1]
		structNameLength := len(structName)
		if structNameLength >= STRUCT_NAME_LENGTH {
			return []byte(structName[:STRUCT_NAME_LENGTH])
		}
		result := make([]byte, STRUCT_NAME_LENGTH)
		writeIndex := STRUCT_NAME_LENGTH
		for i := len(parts) - 1; i >= 0; i-- {
			part := []byte(parts[i])
			if i == len(parts)-1 {
				writeIndex -= len(part)
				copy(result[writeIndex:], part)
			} else if writeIndex >= 2 {
				writeIndex -= 2
				result[writeIndex] = part[0]
				result[writeIndex+1] = '.'
			} else {
				break
			}
		}
		for i := 0; i < writeIndex; i++ {
			result[i] = ' '
		}
		return result
	}
	return []byte(lang.PadStart(name, STRUCT_NAME_LENGTH, ' '))
}
