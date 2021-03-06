package core

import (
	"github.com/stretchr/testify/assert"
	"github.com/zhenghaoz/gorse/base"
	"sort"
	"testing"
)

func TestTop(t *testing.T) {
	model := NewEvaluatorTesterModel([]int{0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]int{0, 1, 2, 3, 4, 5, 6, 7, 8},
		[]float64{1, 9, 2, 8, 3, 7, 4, 6, 5})
	testSet := NewDataSet(NewDataTable([]int{0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]int{0, 1, 2, 3, 4, 5, 6, 7, 8},
		[]float64{1, 9, 2, 8, 3, 7, 4, 6, 5}))
	exclude := map[int]float64{7: 0, 8: 0}
	items := Items(testSet)
	top, _ := Top(items, 0, 5, exclude, model)
	assert.Equal(t, []int{1, 3, 5, 6, 4}, top)
}

func TestNeighbors(t *testing.T) {
	// Generate test data
	//     1 2 3 4 5 6 7 8 9 (item)
	//   +--------------------
	// 1 | 1 1 1 1 1 1 1 1 1
	// 2 | 1 1 1 1 1 1 1 1 0
	// 3 | 1 1 1 1 1 1 1 0 0
	// 4 | 1 1 1 1 1 1 0 0 0
	// 5 | 1 1 1 1 1 0 0 0 0
	// 6 | 1 1 1 1 0 0 0 0 0
	// 7 | 1 1 1 0 0 0 0 0 0
	// 8 | 1 1 0 0 0 0 0 0 0
	// 9 | 1 0 0 0 0 0 0 0 0
	// (user)
	users := make([]int, 0, 81)
	items := make([]int, 0, 81)
	ratings := make([]float64, 0, 81)
	for i := 1; i < 10; i++ {
		for j := 1; j < 10; j++ {
			users = append(users, i)
			items = append(items, j)
			if i+j < 9 {
				ratings = append(ratings, 1)
			} else {
				ratings = append(ratings, 0)
			}
		}
	}
	dataSet := NewDataSet(NewDataTable(users, items, ratings))
	// Find N nearest neighbors
	neighbors, _ := Neighbors(dataSet, 1, 5, base.MSDSimilarity)
	sort.Sort(sort.IntSlice(neighbors))
	assert.Equal(t, []int{2, 3, 4, 5, 6}, neighbors)
}
