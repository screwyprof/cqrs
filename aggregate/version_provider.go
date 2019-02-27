package aggregate

type versionProvider struct {
	version uint64
}

func newVersionProvider() *versionProvider {
	return &versionProvider{}
}

func (v *versionProvider) Version() uint64 {
	return v.version
}

func (v *versionProvider) UpdateVersion(version uint64) {
	v.version = version
}
