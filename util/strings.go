package util

func ContainsOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
