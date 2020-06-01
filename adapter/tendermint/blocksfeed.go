package tendermint

// Tendermint latest block height news feed
type BlocksFeed interface {
	// Subscribe to block height updates via channel. It does not guarantee
	// every block height wll be sent to the channel. It is the downstream
	// responsibility to fill all those gaps
	Subscribe(chan<- uint64) error
}
