package constants

type TaskStatus string

const (
	StatusPending         TaskStatus = "pending"   //статус "Получен"
	StatusRunning         TaskStatus = "running"   //статус "Выполняется"
	StatusCompleted       TaskStatus = "completed" //статус "Выполнен"
	StatusFailed          TaskStatus = "failed"    //статус "Не выполнен"
	StatusNotFullCompleat TaskStatus = "partial"   //статус "Выполнен частично"
)
