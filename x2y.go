package gjxy

// Xstr2Y : XML string to YAML string
func Xstr2Y(xmlstr string) string {
	return Jstr2Y(Xstr2J(xmlstr))
}

// Xfile2Y : XML file to YAML string
func Xfile2Y(xmlfile string) string {
	//return string(Jb2Yb([]byte(Xfile2J(xmlfile))))
	return Jstr2Y(Xfile2J(xmlfile))
}
