package gocialite

//Interface for structs to provide state saving of Dispatcher
type GocialStorage interface {
	Get(key string) (*Gocial, error)
	Set(key string, value *Gocial) error
	Delete(key string) error
}
