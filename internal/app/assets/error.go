package assets

type ErrAssetNoID struct {
}

// Error function.
func (e ErrAssetNoID) Error() string {
	return "asset id cannot be empty"
}

func NewAssetNoIDError() *ErrAssetNoID {
	return &ErrAssetNoID{}
}
