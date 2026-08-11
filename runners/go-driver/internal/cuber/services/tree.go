package services

type Node struct {
	value int
	left  *Node
	right *Node
}

type BinaryTree struct {
	root *Node
}

func (bt *BinaryTree) Insert(value int) {
	if bt.root == nil {
		bt.root = &Node{value, nil, nil}
		return
	}

	queue := []*Node{bt.root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if node.left == nil {
			node.left = &Node{value, nil, nil}
			return
		} else {
			queue = append(queue, node.left)
		}

		if node.right == nil {
			node.right = &Node{value, nil, nil}
			return
		} else {
			queue = append(queue, node.right)
		}
	}
}

func GetBranches(root *Node, path []int, branches *[][]int) {
	if root == nil {
		return
	}

	path = append(path, root.value)

	if root.left == nil && root.right == nil {
		// This is a leaf node
		*branches = append(*branches, path)
		return
	}

	GetBranches(root.left, path, branches)
	GetBranches(root.right, path, branches)
}
