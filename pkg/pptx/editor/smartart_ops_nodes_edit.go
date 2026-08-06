package editor

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
)

// Every SmartArt edit used to be a whole-tree replace: a caller who wanted to
// rename one node had to read the tree, mutate it, and write all of it back.
// These do that read-mutate-write themselves, addressing a node by its path —
// the 0-based index of each step down from the top level, so [1 0] is the first
// child of the second entry.

// errNoSmartArtNodePath is returned when an op names no node to act on.
var errNoSmartArtNodePath = errors.New("no SmartArt node path given")

// SmartArtNodeChange describes what to set on a node. An empty field is left
// as it is.
type SmartArtNodeChange struct {
	Text      string
	Color     string
	ImagePath string
}

// AddSmartArtNode inserts a node under the node at parentPath, or at the top
// level when the path is empty. A negative index appends.
func (e *PresentationEditor) AddSmartArtNode(
	slideIndex, shapeID int,
	parentPath []int,
	index int,
	node smartart.Node,
) error {
	return e.mutateSmartArtNodes(slideIndex, shapeID, func(nodes []smartart.Node) ([]smartart.Node, error) {
		if len(parentPath) == 0 {
			return insertSmartArtNode(nodes, index, node), nil
		}
		parent, err := smartArtNodeAt(nodes, parentPath)
		if err != nil {
			return nil, err
		}
		parent.Children = insertSmartArtNode(parent.Children, index, node)
		return nodes, nil
	})
}

// RemoveSmartArtNode deletes the node at path, along with its children.
func (e *PresentationEditor) RemoveSmartArtNode(slideIndex, shapeID int, path []int) error {
	return e.mutateSmartArtNodes(slideIndex, shapeID, func(nodes []smartart.Node) ([]smartart.Node, error) {
		if len(path) == 0 {
			return nil, errNoSmartArtNodePath
		}
		siblings := &nodes
		if len(path) > 1 {
			parent, err := smartArtNodeAt(nodes, path[:len(path)-1])
			if err != nil {
				return nil, err
			}
			siblings = &parent.Children
		}
		last := path[len(path)-1]
		if last < 0 || last >= len(*siblings) {
			return nil, fmt.Errorf("SmartArt node path %v is out of range", path)
		}
		*siblings = append((*siblings)[:last], (*siblings)[last+1:]...)
		return nodes, nil
	})
}

// UpdateSmartArtNode changes the text, colour or picture of the node at path.
func (e *PresentationEditor) UpdateSmartArtNode(
	slideIndex, shapeID int,
	path []int,
	change SmartArtNodeChange,
) error {
	return e.mutateSmartArtNodes(slideIndex, shapeID, func(nodes []smartart.Node) ([]smartart.Node, error) {
		node, err := smartArtNodeAt(nodes, path)
		if err != nil {
			return nil, err
		}
		if change.Text != "" {
			node.Text = change.Text
		}
		if change.Color != "" {
			node.Color = change.Color
		}
		if change.ImagePath != "" {
			node.ImagePath = change.ImagePath
		}
		return nodes, nil
	})
}

// mutateSmartArtNodes reads the diagram's tree, applies the change, and writes
// the whole tree back — which is what keeps the data model consistent.
func (e *PresentationEditor) mutateSmartArtNodes(
	slideIndex, shapeID int,
	change func([]smartart.Node) ([]smartart.Node, error),
) error {
	info, err := e.GetSmartArt(slideIndex, shapeID)
	if err != nil {
		return err
	}
	nodes, err := change(info.Nodes)
	if err != nil {
		return err
	}
	return e.SetSmartArtNodes(slideIndex, shapeID, nodes)
}

// smartArtNodeAt walks a path down the tree and returns the node it names, so
// callers can change it in place.
func smartArtNodeAt(nodes []smartart.Node, path []int) (*smartart.Node, error) {
	if len(path) == 0 {
		return nil, errNoSmartArtNodePath
	}
	level := nodes
	var node *smartart.Node
	for depth, index := range path {
		if index < 0 || index >= len(level) {
			return nil, fmt.Errorf("SmartArt node path %v is out of range at step %d", path, depth)
		}
		node = &level[index]
		level = node.Children
	}
	return node, nil
}

func insertSmartArtNode(nodes []smartart.Node, index int, node smartart.Node) []smartart.Node {
	if index < 0 || index > len(nodes) {
		return append(nodes, node)
	}
	out := make([]smartart.Node, 0, len(nodes)+1)
	out = append(out, nodes[:index]...)
	out = append(out, node)
	return append(out, nodes[index:]...)
}

// registerSmartArtNodeImages copies each node picture into the package and
// relates it from the diagram's data part, which is where SmartArt looks for
// one. Without this a node could name a picture the deck does not carry.
func (e *PresentationEditor) registerSmartArtNodeImages(nodes []smartart.Node, dataPath string) error {
	relIDByPath := map[string]string{}
	var walk func(items []smartart.Node) error
	walk = func(items []smartart.Node) error {
		for i := range items {
			node := &items[i]
			if node.ImagePath != "" {
				relID, ok := relIDByPath[node.ImagePath]
				if !ok {
					var err error
					relID, err = e.relateSmartArtImage(node.ImagePath, dataPath)
					if err != nil {
						return err
					}
					relIDByPath[node.ImagePath] = relID
				}
				node.ImageRelID = relID
			}
			if err := walk(node.Children); err != nil {
				return err
			}
		}
		return nil
	}
	return walk(nodes)
}

func (e *PresentationEditor) relateSmartArtImage(imagePath, dataPath string) (string, error) {
	partPath, err := e.RegisterImageFromFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("register SmartArt node picture: %w", err)
	}
	target, err := filepath.Rel(filepath.Dir(dataPath), partPath)
	if err != nil {
		return "", fmt.Errorf("relate SmartArt node picture: %w", err)
	}
	target = filepath.ToSlash(target)

	relsPath := common.RelsPathFor(dataPath)
	if data, ok := e.parts.Get(relsPath); ok {
		rels, parseErr := parseRelationshipsXML(data)
		if parseErr != nil {
			return "", parseErr
		}
		for _, rel := range rels {
			if rel.Target == target {
				return rel.ID, nil
			}
		}
	}
	relID := e.nextSmartArtImageRelID(relsPath)
	if err := e.addRelationship(dataPath, relID, common.RelTypeImage, target); err != nil {
		return "", err
	}
	return relID, nil
}

func (e *PresentationEditor) nextSmartArtImageRelID(relsPath string) string {
	next := 1
	if data, ok := e.parts.Get(relsPath); ok {
		if rels, err := parseRelationshipsXML(data); err == nil {
			next += len(rels)
		}
	}
	return "rId" + strconv.Itoa(next)
}
