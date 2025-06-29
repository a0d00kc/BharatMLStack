package index

type Index struct {
	idx map[string][]byte
}

func NewIndex() *Index {
	return &Index{
		idx: make(map[string][]byte),
	}
}

func (i *Index) Put(key string, value []byte) {
	i.idx[key] = value
}

func (i *Index) Get(key string) ([]byte, bool) {
	value, ok := i.idx[key]
	return value, ok
}

func (i *Index) Delete(key string) {
	delete(i.idx, key)
}
