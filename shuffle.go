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

type card struct {
	// Used for deciding if this card is identical to others
	id int
	// This card's position in the future, shuffled deck
	position int
}

func badArgument(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	fmt.Fprintln(os.Stderr,
		"shuffle must have at least 1 argument, all of them integers")
	os.Exit(1)
}

// true means pile A. false means pile B.
func whichPile(c card, median int) bool {
	return c.position < median
}

// Returns true if the positions in the pile are in consecutive ascending order,
// otherwise false.
func inOrder(pile []card) bool {
	for i, c := range pile[:len(pile)-1] {
		if c.position != pile[i+1].position-1 {
			return false
		}
	}
	return true
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

	deckSize := uniqueCards
	for _, n := range identicalGroups {
		deckSize += n
	}

	positions := rand.New(rand.NewSource(*seed)).Perm(deckSize)

	// The initial deck. The end of the slice represents the top of the pile.
	deck := make([]card, uniqueCards)
	id := 0
	for ; id < uniqueCards; id++ {
		deck[id] = card{id, positions[id]}
	}
	positions = positions[uniqueCards:]
	for _, n := range identicalGroups {
		id++
		sort.Ints(positions[:n])
		for i := 0; i < n; i++ {
			deck = append(deck, card{id, positions[0]})
			positions = positions[1:]
		}
	}

	piles := [][]card{deck}

	for len(piles) > 0 {
		lastI := len(piles) - 1
		hand := piles[lastI]
		piles = piles[:lastI]
		if len(hand) > 1 {
			fmt.Printf("Take pile of %d cards.\n", len(hand))
			if inOrder(hand) {
				fmt.Println("This pile is already shuffled.")
			} else {
				min := deckSize
				max := 0
				for _, c := range hand {
					if c.position < min {
						min = c.position
					}
					if c.position > max {
						max = c.position
					}
				}
				median := (min + max + 1) / 2
				pileA := []card{}
				pileB := []card{}
				for len(hand) > 0 {
					i := len(hand) - 1
					newPile := whichPile(hand[i], median)
					for ; i >= 0 && whichPile(hand[i], median) == newPile; i-- {
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
