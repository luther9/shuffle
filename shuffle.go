/*
2017 Luther Thompson. This program is public domain. See COPYING for details.
*/

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

func badArgument(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	fmt.Fprintln(os.Stderr,
		"shuffle must have at least 1 argument, all of them integers")
	os.Exit(1)
}

// Return a slice of integers representing the original state of the physical
// deck. The integers themselves represent the positions of the cards in the
// final, shuffled deck. seed is the random number seed. unique is the number of
// unique, unshuffled cards. groups is a slice of numbers of identical or pre-
// shuffled cards.
//
// This function is the only place in the program where we use an RNG, so we
// initialize it here.
func newDeck(seed int64, unique int, groups []int) sort.IntSlice {
	size := unique
	for _, n := range groups {
		size += n
	}

	deck := sort.IntSlice(rand.New(rand.NewSource(seed)).Perm(size))

	i := unique
	for _, n := range groups {
		deck[i : i+n].Sort()
		i += n
	}

	return deck
}

func main() {
	seed := flag.Int64("s", time.Now().Unix(),
		"Random number seed. Default is the current time.")
	flag.Parse()

	if flag.NArg() < 1 {
		badArgument(nil)
	}
	// The number of unique, unshuffled cards
	uniqueCards, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		badArgument(err)
	}

	// A slice of numbers of identical or pre-shuffled cards
	identicalGroups := make([]int, flag.NArg()-1)
	for i, nStr := range flag.Args()[1:] {
		identicalGroups[i], err = strconv.Atoi(nStr)
		if err != nil {
			badArgument(err)
		}
	}

	piles := []sort.IntSlice{newDeck(*seed, uniqueCards, identicalGroups)}

	for len(piles) > 0 {
		lastI := len(piles) - 1
		hand := piles[lastI]
		piles = piles[:lastI]
		if len(hand) > 1 {
			fmt.Printf("Take pile of %d cards.\n", len(hand))
			if sort.IsSorted(hand) {
				fmt.Println("This pile is already shuffled.")
			} else {
				min := hand[0]
				max := min
				for _, c := range hand {
					if c < min {
						min = c
					}
					if c > max {
						max = c
					}
				}
				median := (min + max + 1) / 2
				pileA := sort.IntSlice{}
				pileB := sort.IntSlice{}
				for len(hand) > 0 {
					i := len(hand) - 1
					newPile := hand[i] < median
					for ; i >= 0 && hand[i] < median == newPile; i-- {
					}
					i++
					transfer := hand[i:]
					hand = hand[:i]
					pile := &pileB
					pileName := 'B'
					if newPile {
						pile = &pileA
						pileName = 'A'
					}
					fmt.Printf("%d to %c", len(transfer), pileName)
					*pile = append(*pile, transfer...)
					fmt.Scanln()
				}
				piles = append(piles, pileB, pileA)
			}
		}
	}
}
