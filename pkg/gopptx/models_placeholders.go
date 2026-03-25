package gopptx

func placeholderTextType(index int) string {
	if index == 0 {
		return "title"
	}
	return "body"
}

func placeholderImageType(index int) string {
	if index == 0 {
		return "title"
	}
	return "pic"
}
