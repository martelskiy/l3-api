package stake

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	"github.com/martelskiy/l3-api/config"
	"github.com/martelskiy/l3-api/internal/shared/api/response"
	"github.com/martelskiy/l3-api/internal/shared/logger"
)

var log = logger.Get()

// @Summary      Get stakes
// @Tags  		 stake
// @Accept       json
// @Produce      json
// @Param        wallet path string true "wallet address"
// @Failure      400 {object} response.BadRequestProblemDetails
// @Failure      500 {object} response.InternalServerErrorProblemDetails
// @Success      200 
// @Router       /api/stakes/{wallet} [get]
func GetStakesHandler(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	wallet := vars["wallet"]
	if wallet == "" {
		responseWriter.WriteHeader(http.StatusBadRequest)
		errorResponse, _ := json.Marshal(response.NewBadRequestProblemDetails([]response.ProblemDetailsError{
			{
				Title:  "invalid wallet",
				Detail: "wallet address must be specified",
			},
		}))
		log.Info("wallet address was empty")
		_, _ = responseWriter.Write(errorResponse)
		return
	}
	if !common.IsHexAddress(wallet) {
		responseWriter.WriteHeader(http.StatusBadRequest)
		errorResponse, _ := json.Marshal(response.NewBadRequestProblemDetails([]response.ProblemDetailsError{
			{
				Title:  "invalid wallet",
				Detail: "specific wallet address is not valid Ethereum address",
			},
		}))
		log.Info("wallet address was not Ethereum address")
		_, _ = responseWriter.Write(errorResponse)
		return
	}

	log.Info("connecting to the node")
	conn, err := ethclient.Dial(config.Configuration.RPCUrl)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		log.Error("failed to connect to the Ethereum client: %v", err)
		errorResponse, _ := json.Marshal(response.NewInternalServerErrorProblemDetails())
		_, _ = responseWriter.Write(errorResponse)
		return
	}

	log.Info("connection succeeded. Instantiating Stake contract")
	stakeContract, err := NewStake(common.HexToAddress(config.Configuration.StakeContractAddr), conn)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		log.Error("failed to instantiate Stake contract: %v", err)
		_, _ = responseWriter.Write([]byte("internal server error"))
		return

	}

	log.Infof("reading Stakes for wallet address: '%s'", wallet)
	stakes, err := stakeContract.Stakes(nil, common.HexToAddress(wallet))
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		log.Error("failed to fetch stakes: %v", err)
		_, _ = responseWriter.Write([]byte("internal server error"))
		return
	}

	log.Infof("stakes were fetched. Amount: '%v', StartTime: '%v'", stakes.Amount, stakes.StartTime)

	responseWriter.WriteHeader(http.StatusOK)
	statusResponse, _ := json.Marshal(stakes.Amount)
	_, _ = responseWriter.Write([]byte(statusResponse))
}
