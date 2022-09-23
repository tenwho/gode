package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Parameters
const (
	F         float64 = 0.5   // Scaling factor
	N         int     = 15    // Number of agents
	CR        float64 = 0.5   // Crossover rate
	Dimension int     = 2     // Dimension of the data
	IterMax   int     = 40    // Number of iterations
	Logging   bool    = false // If true, print data of each iteration
	Plotting  bool    = true  // If true, plot the trajectory to png and gif files
)

var (
	Xmin = [Dimension]float64{-15.0, -15.0} // Minimum value of each dimension
	Xmax = [Dimension]float64{15.0, 15.0}   // Maximum value of each dimention
)

// Target function to be minimized
func target(xi []float64) float64 {
	return xi[0]*xi[0] + xi[0]*xi[1] + xi[1]*xi[1] - 5*xi[0] - 5*xi[1] + 25
}

func main() {
	seed := time.Now().UnixNano()

	// Run
	best, bestScore, traj := optimize(seed)

	fmt.Println("Seed:", seed)
	fmt.Println("Result:", best)
	fmt.Println("Best Score:", bestScore)

	// Plot
	if Plotting {
		if Dimension != 2 {
			fmt.Println("Plotting is only supporting 2-Dimensional.")
		} else {
			if err := plot2D(traj); err != nil {
				fmt.Println("Failed to plot the trajectory:", err)
			}
		}
	}
}

// Optimize target function and return best values
// and the score and the trajectory.
func optimize(seed int64) ([]float64, float64, [IterMax][N][]float64) {
	rand.Seed(seed)

	var xs [N][]float64
	var scoreTable [N]float64

	// Generate initial values
	for i := 0; i < N; i++ {
		for j := 0; j < Dimension; j++ {
			xs[i] = append(xs[i], rand.Float64()*(Xmax[j]-Xmin[j])+Xmin[j])
		}
		scoreTable[i] = math.Inf(0)
	}

	// Variables to keep best values
	bestIdx := 0
	bestScore := math.Inf(0)

	// Trajectories of all iteration
	var traj [IterMax][N][]float64

	// Optimize
	for iter := 0; iter < IterMax; iter++ {
		// Generate a new agent candidate
		var newXs [N][]float64
		for i, xi := range xs {
			ia, ib, ic := pickThreeWithout(N, i)
			a := xs[ia]
			b := xs[ib]
			c := xs[ic]
			jr := rand.Intn(len(xi))
			var newXi []float64
			for j := 0; j < Dimension; j++ {
				if j == jr || rand.Float64() < CR {
					newXi = append(newXi, a[j]+F*(b[j]-c[j]))
				} else {
					newXi = append(newXi, xi[j])
				}
			}
			newXs[i] = newXi
		}
		// Compare score and update agents
		for i, xi := range xs {
			if math.IsInf(scoreTable[i], 0) {
				scoreTable[i] = target(xi)
			}
			newXi := newXs[i]
			newScore := target(newXi)
			if newScore < scoreTable[i] {
				xs[i] = newXi
				scoreTable[i] = newScore
				if newScore < bestScore {
					bestIdx = i
					bestScore = newScore
				}
			}
			// Save trajectory
			t := make([]float64, Dimension)
			copy(t, xs[i])
			traj[iter][i] = t
		}
		if Logging {
			fmt.Println(xs)
		}
	}

	return xs[bestIdx], bestScore, traj
}

// Return three number between 0 ~ n-1 randomly, except `wo`
func pickThreeWithout(n int, wo int) (int, int, int) {
	var ret []int
	for _, v := range rand.Perm(n) {
		if v != wo {
			ret = append(ret, v)
		}
		if len(ret) == 3 {
			return ret[0], ret[1], ret[2]
		}
	}
	// Unreachable, unless size of n is less than 4
	panic("Failed to pick two element")
}
