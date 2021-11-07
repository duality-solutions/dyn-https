package models

// BlockchainOverview stores the blockchain node sync status
type BlockchainOverview struct {
	Network           string  `json:"network"`
	ClientVersion     int     `json:"client_version"`
	Peers             int     `json:"peers"`
	Blocks            int     `json:"blocks"`
	TotalBlocks       int     `json:"total_blocks"`
	SyncProgress      float64 `json:"sync_progress"`
	StatusDescription string  `json:"status_description"`
	FullySynced       bool    `json:"fully_synced"`
}

// WalletOverview stores the wallet balance overview values
type WalletOverview struct {
	Transactions     int     `json:"transactions"`
	Encrypted        bool    `json:"encrypted"`
	UnlockedEpoch    int64   `json:"unlockedepoch"`
	AvailableBalance float64 `json:"available_balance"`
	PendingBalance   float64 `json:"pending_balance"`
	TotalBalance     float64 `json:"total_balance"`
	Credits          float64 `json:"credits"`
	Deposits         float64 `json:"deposits"`
}

// AccountOverview stores the current account overview values
type AccountOverview struct {
	Users         int `json:"users"`
	CompleteLinks int `json:"complete_links"`
	PendingLinks  int `json:"pending_links"`
	Certificates  int `json:"certificates"`
	Audits        int `json:"audits"`
}

// OverviewResponse is used to store the chain, wallet, balance, account and bridge overviews
type OverviewResponse struct {
	Chain    BlockchainOverview `json:"chain"`
	Wallet   WalletOverview     `json:"wallet"`
	Accounts AccountOverview    `json:"accounts"`
}
