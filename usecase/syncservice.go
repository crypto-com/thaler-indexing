package usecase

type SyncService interface {
	Sync() error
	GetStatus() SyncStatus
}

type SyncStatus struct {
	TendermintBlockHeight uint64
	SyncBlockHeight       uint64
}
