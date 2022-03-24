package types

type Order struct {
	Value         float64 `json:"value"`
	Quantity      int     `json:"quantity"`
	OperationType string  `json:"operationType"`
	UserId        int     `json:"userId"`
	RequestId     string  `json:"RequestId"`
}
