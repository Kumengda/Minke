package minke

func (m *Minke) matchContainerName(containerNames []string) bool {
	for _, cname := range containerNames {
		for _, nr := range m.containerNameRule {
			if nr.MatchString(cname) {
				return true
			}
		}
	}
	return false
}

func (m *Minke) matchImageName(imageName string) bool {
	for _, ir := range m.imageNameRule {
		if ir.MatchString(imageName) {
			return true
		}
	}
	return false
}
