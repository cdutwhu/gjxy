package gjxy

// SchemaBuild : ignore identical field
func SchemaBuild(s Str, typename, field, fieldtype string) (schema string) {
	if ok, pos, _ := s.SearchStrsIgnore("type", typename, "{", " \t"); ok {
		content, _, r := s[pos:].BracketsPos(BCurly, 1, 1)
		// if sCtn(content.V(), field+":") {
		if content.Contains(field + ":") {
			return s.V()
		}
		schema = s.V()[:pos+r]
		tail := s.V()[pos+r+1:]
		add := fSf("\t%s: %s\n}", field, fieldtype)
		schema += add + tail
	} else {
		s += Str(fSf("\n\ntype %s {\n\t%s: %s\n}", typename, field, fieldtype))
		schema = s.V()
	}
	return schema
}
