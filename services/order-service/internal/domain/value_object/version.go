package value_object

import "github.com/samber/lo"

type Version struct {
	Old int64
	new *int64
}

func (v *Version) Current() int64 {
	if v.new != nil {
		return *v.new
	}
	return v.Old
}

func (v *Version) Inc() {
	v.new = lo.ToPtr(v.Old + 1)
}

func (v *Version) Previous() int64 {
	return v.Old
}
