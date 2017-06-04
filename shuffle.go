package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func badArgument(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	fmt.Fprintln(os.Stderr, "shuffle must have 1 or 2 integer arguments")
	os.Exit(1)
}

func main() {
	seed := flag.Int64("s", time.Now().Unix(),
		"Random number seed. Default is the current time.")
	flag.Parse()

	rand.Seed(*seed)

	if flag.NArg() < 1 {
		badArgument(nil)
	}
	// n is the total number of cards
	n, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		badArgument(err)
	}

	// nonRandomCards is the number of cards that are not yet randomized
	nonRandomCards := n
	if flag.NArg() > 1 {
		nonRandomCards, err = strconv.Atoi(flag.Arg(1))
		if err != nil {
			badArgument(err)
		}
	}
	// The number of cards to move from the random pile
	cardsToMove := 0

	for ; n > 0; n-- {
		card := rand.Intn(n)
		if card < nonRandomCards {
			if cardsToMove > 0 {
				fmt.Printf("Take %d cards.", cardsToMove)
				cardsToMove = 0
				fmt.Scanln()
			}
			nonRandomCards--
			fmt.Printf("%2d %2d", card, nonRandomCards-card)
			fmt.Scanln()
		} else {
			cardsToMove++
		}
	}

	// Make sure we have the correct number of cards left.
	fmt.Printf("%2d cards left.\n", cardsToMove)
}
