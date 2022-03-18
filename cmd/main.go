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

type CmdArgs struct {
	Amount     float64
	SymbolFrom string
	SymbolTo   string
}

const requestTimeout = time.Second * 3

// todo: move secrets into config files
func main() {
	cmdArgs, err := parseArgs(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	coinMarketCapClient := coinmarketcap.NewClient(
		"ab730ace-7180-42a8-8d5d-6ac3a19de87a",
		"https://pro-api.coinmarketcap.com",
	)
	cmcProvider := provider.NewCoinMarketCapProvider(coinMarketCapClient)

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	symbolToPrice, err := cmcProvider.Convert(ctx, cmdArgs.Amount, cmdArgs.SymbolFrom, cmdArgs.SymbolTo)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(strconv.FormatFloat(symbolToPrice, 'f', 8, 64))
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
