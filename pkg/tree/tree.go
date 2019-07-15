package tree

// 树的结构
type Tree struct {
	//节点ID
	Id string `json:"id"`
	//显示节点文本
	Text string `json:"text"`
	//节点状态，open closed
	State map[string]interface{} `json:"state"`
	//节点是否被选中 true false
	Checked bool `json:"checked"`
	//节点属性
	Attributes map[string]interface{} `json:"attributes"`
	//节点的子节点
	Children []*Tree `json:"children"`
	//父ID
	ParentId string `json:"parentId"`
	//是否有父节点
	HasParent bool `json:"hasParent"`
	//是否有子节点
	HasChildren bool `json:"hasChildren"`
}

func NewTree(id string, parentId string, text string, state map[string]interface{}) *Tree {
	return &Tree{
		Id:       id,
		ParentId: parentId,
		Text:     text,
		State:    state,
	}
}

func Build(nodes []*Tree) *Tree {
	if nodes == nil {
		return nil
	}

	topNodes := make([]*Tree, 0, len(nodes))
	for _, children := range nodes {
		pid := children.ParentId
		if len(pid) == 0 || pid == "0" {
			topNodes = append(topNodes, children)
			continue
		}

		for _, parent := range nodes {
			id := parent.Id
			if len(id) != 0 && id == pid {
				parent.Children = append(parent.Children, children)
				children.HasParent = true
				parent.HasChildren = true
				continue
			}
		}
	}

	root := &Tree{}
	if len(topNodes) == 1 {
		root = topNodes[0]
	} else {
		root.Id = "-1"
		root.ParentId = ""
		root.HasParent = false
		root.HasChildren = true
		root.Checked = true
		root.Children = topNodes
		root.Text = "顶级节点"
		root.State = map[string]interface{}{"opened": true}
	}
	return root
}

func BuildList(nodes []*Tree, idParam string) []*Tree {
	if nodes == nil {
		return nil
	}
	topNodes := make([]*Tree, 0, len(nodes))
	for _, children := range nodes {
		pid := children.ParentId
		if len(pid) == 0 || idParam == pid {
			topNodes = append(topNodes, children)
			continue
		}

		for _, parent := range nodes {
			id := parent.Id
			if len(id) != 0 && id == pid {
				parent.Children = append(parent.Children, children)
				children.HasParent = true
				parent.HasChildren = true
				continue
			}
		}
	}
	return topNodes
}
