package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"go-coinconv/internal/provider"
	"go-coinconv/pkg/coinmarketcap"
)

func main() {
	cmdArgs, err := parseArgs(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	// todo: move secrets into config files
	coinMarketCapClient := coinmarketcap.NewClient(
		"ab730ace-7180-42a8-8d5d-6ac3a19de87a",
		"https://pro-api.coinmarketcap.com",
	)
	cmcProvider := provider.NewCoinMarketCapProvider(coinMarketCapClient)
	conversionVal, err := convertCMD(cmdArgs, cmcProvider)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(strconv.FormatFloat(conversionVal, 'f', 8, 64))
}

func convertCMD(cmdArgs *CmdArgs, convertor provider.Convertor) (float64, error) {
	const requestTimeout = time.Second * 3

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	return convertor.Convert(ctx, cmdArgs.Amount, cmdArgs.SymbolFrom, cmdArgs.SymbolTo)
}

type CmdArgs struct {
	Amount     float64
	SymbolFrom string
	SymbolTo   string
}

// todo: consider symbol args validation, avoiding extra api call
// todo: add --usage description on error
func parseArgs(args []string) (*CmdArgs, error) {
	if len(args) < 3 {
		return nil, errors.New("not enough args")
	}

	amount, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		return nil, err
	}

	return &CmdArgs{
		Amount:     amount,
		SymbolFrom: args[1],
		SymbolTo:   args[2],
	}, nil
}
