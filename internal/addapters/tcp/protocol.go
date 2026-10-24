package addapters

type Method string

const (
	CreateTicket Method = "create-ticket"
)

type Payload struct {
	Method Method `json:"method"`
	Data   any    `data:"data"`
}
