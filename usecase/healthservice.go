package usecase

type HealthService interface {
	IsTendermintHealthy() bool
	IsServerHealthy() bool
}
