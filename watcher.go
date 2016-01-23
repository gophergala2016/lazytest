package lazytest

type FileMatch struct {
	PathPrefix string
	Extensions []string
}

type Mod struct {
	Package  string
	FilePath string
	Function string
	Line     int
}

func Watch(include []FileMatch, exclude []FileMatch) chan Mod {
	mods := make(chan Mod, 50)

	return mods
}
