package models

type MetricsUpdate struct {
	ID    string   `json:"id" binding:"required"`
	MType string   `json:"type" binding:"required,oneof=counter gauge"`
	Delta *int64   `json:"delta,omitempty" binding:"required_if=MType counter"`
	Value *float64 `json:"value,omitempty" binding:"required_if=MType gauge"`
}

type MetricsValue struct {
	ID    string   `json:"id" binding:"required"`
	MType string   `json:"type" binding:"required,oneof=counter gauge"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
