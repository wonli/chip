package datasource

type Handler struct {
	Csv *csvHandler
}

func Use() *Handler {
	return &Handler{
		Csv: &csvHandler{},
	}
}
