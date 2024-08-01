package treeset

import "gurms/internal/infra/logging/core/appender/file/logfile"

type Iterator struct {
	tree     *Tree
	node     *Node
	position position
}

type position byte

const (
	begin, between, end position = 0, 1, 2
)

func (tree *Tree) Iterator() *Iterator {
	return &Iterator{tree: tree, node: nil, position: begin}
}

func (iterator *Iterator) Next() bool {
	if iterator.position == end {
		goto end
	}
	if iterator.position == begin {
		left := iterator.tree.Left()
		if left == nil {
			goto end
		}
		iterator.node = left
		goto between
	}
	if iterator.node.Right != nil {
		iterator.node = iterator.node.Right
		for iterator.node.Left != nil {
			iterator.node = iterator.node.Left
		}
		goto between
	}
	for iterator.node.Parent != nil {
		node := iterator.node
		iterator.node = iterator.node.Parent
		if node == iterator.node.Left {
			goto between
		}
	}

end:
	iterator.node = nil
	iterator.position = end
	return false

between:
	iterator.position = between
	return true
}

func (iterator *Iterator) Key() logfile.LogFile {
	return iterator.node.Key
}
