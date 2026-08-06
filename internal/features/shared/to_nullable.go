package shared

func ToNullable[K any, T any](ptr *K, constructor func(K) (T, error)) (*T, error) {
	if ptr == nil {
		return nil, nil
	}

	v, err := constructor(*ptr)
	if err != nil {
		return nil, err
	}

	return &v, nil
}
