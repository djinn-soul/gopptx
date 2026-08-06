package pptxxml

// A nested spec used to be poured into the template by walking its points in
// document order, which scattered children into the next entry's slots. Walking
// the spec tree and the template tree together instead keeps each child under
// the parent it belongs to — and keeps the template's presentation points, which
// is where node pictures have to hang.

type smartArtSlot struct {
	modelID string
	node    SmartArtNodeSpec
}

// assignSmartArtSlots pairs each node of the spec with the template data point
// that holds its level and position. Nodes beyond the template's slots are
// dropped; callers check capacity first.
func assignSmartArtSlots(spec SmartArtSpec, data string) []smartArtSlot {
	capacityByParent, rootID := smartArtTemplateCapacity(data)
	if rootID == "" {
		return nil
	}

	out := make([]smartArtSlot, 0, flattenSmartArtTextsInitCap)
	var walk func(nodes []SmartArtNodeSpec, slots []string)
	walk = func(nodes []SmartArtNodeSpec, slots []string) {
		for i, node := range nodes {
			if i >= len(slots) {
				return
			}
			modelID := slots[i]
			out = append(out, smartArtSlot{modelID: modelID, node: node})
			if len(node.Children) > 0 {
				walk(node.Children, capacityByParent[modelID])
			}
		}
	}
	walk(spec.Nodes, capacityByParent[rootID])
	return out
}

func smartArtSlotModelIDs(slots []smartArtSlot) []string {
	out := make([]string, 0, len(slots))
	for _, slot := range slots {
		out = append(out, slot.modelID)
	}
	return out
}

func smartArtSlotTexts(slots []smartArtSlot) []string {
	out := make([]string, 0, len(slots))
	for _, slot := range slots {
		out = append(out, slot.node.Text)
	}
	return out
}

func smartArtSlotNodes(slots []smartArtSlot) []SmartArtNodeSpec {
	out := make([]SmartArtNodeSpec, 0, len(slots))
	for _, slot := range slots {
		out = append(out, slot.node)
	}
	return out
}
