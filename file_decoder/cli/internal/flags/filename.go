package flags

type FilenameFlag struct {
	set   bool
	value string
}

func (cf *FilenameFlag) Set(x string) error {
	cf.value = x
	cf.set = true
	return nil
}

func (cf *FilenameFlag) String() string {
	return cf.value
}

func (cf *FilenameFlag) Filename() string {
	return cf.value
}

func (cf *FilenameFlag) IsSet() bool {
	return cf.set
}
