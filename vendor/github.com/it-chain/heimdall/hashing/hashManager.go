package hashing

type HashManager interface {

	Hash(data []byte, tail []byte, opts HashOpts) ([]byte, error)

}