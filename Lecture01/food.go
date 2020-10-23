package main

import (
	"fmt"
	"sort"
)

// Food Module
type Food struct {
	name     string
	value    float64
	calories float64
}

func (f Food) getValue() float64 {
	return f.value
}

func (f Food) getCost() float64 {
	return f.calories
}

func (f Food) density() float64 {
	return f.value / f.calories
}

func (f Food) String() string {
	return fmt.Sprintf("%s: <%.2f, %.2f>", f.name, f.value, f.calories)
}

// names, vallues, calories list the same length.
// name a list of strings
// values and calories lists of numbers
// return list of Foods
func buildMenu(names []string, values []float64, calories []float64) []Food {
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

func greedy(items []Food, maxCost float64, compFunc func(x, y Food) bool) ([]Food, float64) {
	//var itemsCopy []Food
	//copy(itemsCopy, items)
	itemsCopy := items[:]

	sort.Sort(customSort{itemsCopy, compFunc})
	fmt.Println(itemsCopy)

	var result []Food
	totalValue := 0.0
	totalCost := 0.0

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

func testGreedy(items []Food, constraint float64, compFunc func(x, y Food) bool) {
	taken, val := greedy(items, constraint, compFunc)
	fmt.Println("Total value of items taken =", val)
	for _, item := range taken {
		fmt.Println("   ", item)
	}
}

func testGreedys(foods []Food, maxUnits float64) {
	fmt.Println("Use greedy by value to allocate", maxUnits, "calories")
	testGreedy(foods, maxUnits, func(x, y Food) bool {
		return x.getValue() > y.getValue()
	})

	fmt.Println()
	fmt.Println("Use greedy by cost to allocate", maxUnits, "calories")
	testGreedy(foods, maxUnits, func(x, y Food) bool {
		return x.getCost() < y.getCost()
	})

	fmt.Println()
	fmt.Println("Use greedy by density to allocate", maxUnits, "calories")
	testGreedy(foods, maxUnits, func(x, y Food) bool {
		return x.density() > y.density()
	})
}

func main() {
	names := [...]string{"wine", "beer", "pizza", "burger", "fries", "cola", "apple", "donut", "cake"}
	values := [...]float64{89, 90, 95, 100, 90, 79, 50, 10}
	calories := [...]float64{123, 154, 258, 354, 365, 150, 95, 195}
	foods := buildMenu(names[:], values[:], calories[:])
	for _, f := range foods {
		fmt.Println(f)
	}

	testGreedys(foods, 750.0)
	testGreedys(foods, 1000.0)
}
