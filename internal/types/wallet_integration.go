package types

const (
	credit = "CREDIT"
	debit  = "DEBIT"
)

type Order struct {
	Value         float64 `json:"value"`
	Quantity      int     `json:"quantity"`
	OperationType string  `json:"operationType"`
	UserId        int     `json:"userId"`
	TraceId       string  `json:"traceId"`
}
