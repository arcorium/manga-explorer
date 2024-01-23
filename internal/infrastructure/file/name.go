package file

type Name string

func (f Name) String() string {
	return string(f)
}
