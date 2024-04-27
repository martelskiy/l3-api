package config

var Configuration *AppConfig

type AppConfig struct {
	StakeContractAddr string `json:"stake_contract_addr"`
	RPCUrl            string `json:"rpc_url"`
	Api               Api    `json:"api"`
}

type Api struct {
	Port string `json:"port"`
}
