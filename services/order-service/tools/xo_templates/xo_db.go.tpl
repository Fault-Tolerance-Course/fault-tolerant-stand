// Jsonb ...
type Jsonb string

// JSON ...
type JSON string

// Oid ...
type Oid oid.Oid

func argFn(fields []string, tm stom.ToMappable) (res []interface{}, err error) {
	m, err := tm.ToMap()
	if err != nil {
		return nil, err
	}
	for _, v := range fields {
		if value, ok := m[v]; ok {
			if iv, ok := value.(interface{ Valid() bool }); ok {
				if !iv.Valid() {
					res = append(res, nil)
					continue
				}
			}
			if value == nil {
				res = append(res, nil)
				continue
			}
			res = append(res, value)
		} else {
			res = append(res, nil)
		}
	}
	return
}

var (
	reQuestion = regexp.MustCompile(`\$\d`)

	ErrRowAlreadyExists = errors.New("db insert failed: already exists")
	ErrRowMarkedForDeletion = errors.New("db update failed: marked for deletion")
	ErrRowDoesNotExists = errors.New("update failed: does not exist")
)
