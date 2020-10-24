package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

// Food Module
type Food struct {
	name     string
	value    int
	calories int
}

func (f Food) getValue() int    { return f.value }
func (f Food) getCost() int     { return f.calories }
func (f Food) density() float64 { return float64(f.value) / float64(f.calories) }
func (f Food) String() string {
	return fmt.Sprintf("%s: <%d, %d>", f.name, f.value, f.calories)
}

// names, vallues, calories list the same length.
// name a list of strings
// values and calories lists of numbers
// return list of Foods
func buildMenu(names []string, values []int, calories []int) []Food {
	var menu []Food
	for i, v := range values {
		var f Food
		f.name = names[i]
		f.value = v
		f.calories = calories[i]
		menu = append(menu, f)
	}

	return menu
}

type customSort struct {
	t    []Food
	less func(x, y Food) bool
}

func (x customSort) Len() int           { return len(x.t) }
func (x customSort) Less(i, j int) bool { return x.less(x.t[i], x.t[j]) }
func (x customSort) Swap(i, j int)      { x.t[i], x.t[j] = x.t[j], x.t[i] }

type CompFunction func(x, y Food) bool

func greedy(items []Food, maxCost int, compFunc CompFunction) ([]Food, int) {
	itemsCopy := items[:]

	sort.Sort(customSort{itemsCopy, compFunc})
	//fmt.Println(itemsCopy)

	var result []Food
	totalValue := 0
	totalCost := 0

	for _, item := range itemsCopy {
		if (totalCost + item.getCost()) <= maxCost {
			result = append(result, item)
			totalCost += item.getCost()
			totalValue += item.getValue()
		}
	}
	//fmt.Println("total cost=", totalCost)

	return result, totalValue
}

func testGreedy(items []Food, constraint int, compFunc CompFunction) {
	taken, val := greedy(items, constraint, compFunc)
	fmt.Println("Total value of items taken =", val)
	for _, item := range taken {
		fmt.Println("   ", item)
	}
}

func testGreedys(foods []Food, maxUnits int) {
	fmt.Println("Use greedy by value to allocate", maxUnits, "calories")
	testGreedy(foods, maxUnits, func(x, y Food) bool {
		return x.value > y.value
	})
	fmt.Println()

	fmt.Println("Use greedy by cost to allocate", maxUnits, "calories")
	testGreedy(foods, maxUnits, func(x, y Food) bool {
		return x.calories < y.calories
	})
	fmt.Println()

	fmt.Println("Use greedy by density to allocate", maxUnits, "calories")
	testGreedy(foods, maxUnits, func(x, y Food) bool {
		return x.density() > y.density()
	})
	fmt.Println()
}

func maxVal(toConsider []Food, avail int) (int, []Food) {
	var taken []Food
	var val int
	if len(toConsider) == 0 || avail == 0.0 {
		val, taken = 0, nil
	} else if toConsider[0].getCost() > avail {
		// Explore right branch only
		val, taken = maxVal(toConsider[1:], avail)
	} else {
		nextItem := toConsider[0]
		// Explore left branch
		withVal, withToTake := maxVal(toConsider[1:], avail-nextItem.getCost())
		withVal += nextItem.getValue()
		// Explore right branch
		withoutVal, withoutToTake := maxVal(toConsider[1:], avail)
		// Choose better branch
		if withVal > withoutVal {
			val, taken = withVal, append(withToTake, nextItem)
		} else {
			val, taken = withoutVal, withoutToTake
		}
	}
	return val, taken
}

func testMaxVal(foods []Food, maxUnits int, printItems bool) {
	fmt.Println("Use search tree to allocate", maxUnits, "calories")
	val, taken := maxVal(foods, maxUnits)
	fmt.Println("Total value of items taken =", val)
	if printItems {
		for _, item := range taken {
			fmt.Println("   ", item)
		}
	}
}

func buildLargeMenu(numItems int, maxVal int, maxCost int) []Food {
	var items []Food

	for i := 0; i < numItems; i++ {
		name := strconv.Itoa(i)
		val := rand.Intn(maxVal) + 1
		calories := rand.Intn(maxCost) + 1
		items = append(items, Food{name, val, calories})
	}

	return items
}

func main() {
	names := [...]string{"wine", "beer", "pizza", "burger", "fries", "cola", "apple", "donut", "cake"}
	values := [...]int{89, 90, 95, 100, 90, 79, 50, 10}
	calories := [...]int{123, 154, 258, 354, 365, 150, 95, 195}
	foods := buildMenu(names[:], values[:], calories[:])
	fmt.Println("menu =", foods)
	fmt.Println()

	testGreedys(foods, 750)
	//testGreedys(foods, 800)
	//testGreedys(foods, 1000)

	testMaxVal(foods, 750, true)

	rand.Seed(time.Now().UTC().UnixNano())
	for numItems := 5; numItems < 50; numItems += 5 {
		fmt.Println("Try a menu with", numItems, "items")
		items := buildLargeMenu(numItems, 90, 250)
		start := time.Now()
		testMaxVal(items, 750, false)
		fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
	}
}
