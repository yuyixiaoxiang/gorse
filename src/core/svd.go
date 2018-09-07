// The famous SVD algorithm, as popularized by Simon Funk during the
// Netflix Prize. When baselines are not used, this is equivalent to
// Probabilistic Matrix Factorization [SM08] (see note below). The
// prediction r^ui is set as:
//               \hat{r}_{ui} = μ + b_u + b_i + q_i^Tp_u
// If user u is unknown, then the bias bu and the factors pu are
// assumed to be zero. The same applies for item i with bi and qi.

package core

import (
	"github.com/gonum/floats"
)

type SVD struct {
	userFactor [][]float64 // p_u
	itemFactor [][]float64 // q_i
	userBias   []float64   // b_u
	itemBias   []float64   // b_i
	globalBias float64     // mu
	trainSet   TrainSet
}

func NewSVD() *SVD {
	return new(SVD)
}

func (svd *SVD) Predict(userId int, itemId int) float64 {
	innerUserId := svd.trainSet.ConvertUserId(userId)
	innerItemId := svd.trainSet.ConvertItemId(itemId)
	ret := svd.globalBias
	// + b_u
	if innerUserId != newId {
		ret += svd.userBias[innerUserId]
	}
	// + b_i
	if innerItemId != newId {
		ret += svd.itemBias[innerItemId]
	}
	// + q_i^Tp_u
	if innerItemId != newId && innerUserId != newId {
		userFactor := svd.userFactor[innerUserId]
		itemFactor := svd.itemFactor[innerItemId]
		ret += floats.Dot(userFactor, itemFactor)
	}
	return ret
}

func (svd *SVD) Fit(trainSet TrainSet, options Options) {
	// Setup options
	nFactors := options.GetInt("nFactors", 100)
	nEpochs := options.GetInt("nEpochs", 20)
	lr := options.GetFloat64("lr", 0.005)
	reg := options.GetFloat64("reg", 0.02)
	//biased := options.GetBool("biased", true)
	initMean := options.GetFloat64("initMean", 0)
	initStdDev := options.GetFloat64("initStdDev", 0.1)
	// Initialize parameters
	svd.trainSet = trainSet
	svd.userBias = make([]float64, trainSet.UserCount())
	svd.itemBias = make([]float64, trainSet.ItemCount())
	svd.userFactor = make([][]float64, trainSet.UserCount())
	svd.itemFactor = make([][]float64, trainSet.ItemCount())
	for innerUserId := range svd.userFactor {
		svd.userFactor[innerUserId] = newNormalVector(nFactors, initMean, initStdDev)
	}
	for innerItemId := range svd.itemFactor {
		svd.itemFactor[innerItemId] = newNormalVector(nFactors, initMean, initStdDev)
	}
	// Create buffers
	a := make([]float64, nFactors)
	b := make([]float64, nFactors)
	// Stochastic Gradient Descent
	users, items, ratings := trainSet.Interactions()
	for epoch := 0; epoch < nEpochs; epoch++ {
		for i := 0; i < trainSet.Length(); i++ {
			userId, itemId, rating := users[i], items[i], ratings[i]
			innerUserId := trainSet.ConvertUserId(userId)
			innerItemId := trainSet.ConvertItemId(itemId)
			userBias := svd.userBias[innerUserId]
			itemBias := svd.itemBias[innerItemId]
			userFactor := svd.userFactor[innerUserId]
			itemFactor := svd.itemFactor[innerItemId]
			// Compute error
			diff := svd.Predict(userId, itemId) - rating
			// Update global bias
			gradGlobalBias := diff
			svd.globalBias -= lr * gradGlobalBias
			// Update user bias
			gradUserBias := diff + reg*userBias
			svd.userBias[innerUserId] -= lr * gradUserBias
			// Update item bias
			gradItemBias := diff + reg*itemBias
			svd.itemBias[innerItemId] -= lr * gradItemBias
			// Update user latent factor
			copy(a, itemFactor)
			mulConst(diff, a)
			copy(b, userFactor)
			mulConst(reg, b)
			floats.Add(a, b)
			mulConst(lr, a)
			floats.Sub(svd.userFactor[innerUserId], a)
			// Update item latent factor
			copy(a, userFactor)
			mulConst(diff, a)
			copy(b, itemFactor)
			mulConst(reg, b)
			floats.Add(a, b)
			mulConst(lr, a)
			floats.Sub(svd.itemFactor[innerItemId], a)
		}
	}
}
