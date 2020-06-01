package httpapi

type RoutesRegistry struct {
	router Router

	statusHandler       *StatusHandler
	chainStatusHandler  *ChainStatusHandler
	activitiesHandler   *ActivitiesHandler
	blocksHandler       *BlocksHandler
	councilNodesHandler *CouncilNodesHandler
	searchHandler       *SearchHandler
}

func NewRoutesRegistry(
	router Router,

	statusHandler *StatusHandler,
	chainStatusHandler *ChainStatusHandler,
	activitiesHandler *ActivitiesHandler,
	blocksHandler *BlocksHandler,
	councilNodesHandler *CouncilNodesHandler,
	searchHandler *SearchHandler,
) *RoutesRegistry {
	return &RoutesRegistry{
		router,

		statusHandler,
		chainStatusHandler,
		activitiesHandler,
		blocksHandler,
		councilNodesHandler,
		searchHandler,
	}
}

func (api *RoutesRegistry) RegisterHandlers() {
	api.router.Get("/health", api.statusHandler.Health)
	api.router.Get("/status", api.statusHandler.Status)

	api.router.Get("/chain/status", api.chainStatusHandler.GetChainStatus)

	api.router.Get("/chain/blocks", api.blocksHandler.ListBlocks)
	api.router.Get("/chain/blocks/{hash_or_height}", api.blocksHandler.FindBlock)
	api.router.Get("/chain/blocks/{hash_or_height}/transactions", api.blocksHandler.ListBlockTransactions)
	api.router.Get("/chain/blocks/{hash_or_height}/events", api.blocksHandler.ListBlockEvents)

	api.router.Get("/chain/transactions", api.activitiesHandler.ListTransactions)
	api.router.Get("/chain/transactions/{txid}", api.activitiesHandler.FindTransactionByTxId)
	api.router.Get("/chain/events", api.activitiesHandler.ListEvents)
	api.router.Get("/chain/events/{height}-{position}", api.activitiesHandler.FindEventByBlockHeightEventPosition)

	api.router.Get("/chain/council-nodes", api.councilNodesHandler.ListActiveCouncilNodes)
	api.router.Get("/chain/council-nodes/{id}", api.councilNodesHandler.FindCouncilNodeById)
	api.router.Get("/chain/council-nodes/{id}/activities", api.councilNodesHandler.ListCouncilNodeActivitiesById)

	api.router.Get("/chain/search/all", api.searchHandler.All)
}
