package filereader

type Reader interface {
	Read(interface{}) error
}
