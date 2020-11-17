package happa

type Headers []string

func (h *Headers) String() string {
	return "ok"
}

func (h *Headers) Set(s string) error {
	*h = append(*h, s)
	return nil
}