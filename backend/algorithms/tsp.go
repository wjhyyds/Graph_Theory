package algorithms

import "math"

type TSPResult struct {
	Path     []int `json:"path"`
	Distance int   `json:"distance"`
}

func TSP(graph [][]int) TSPResult {
	n := len(graph)
	if n == 0 {
		return TSPResult{}
	}

	// 初始化dp数组
	dp := make([][]int, 1<<n)
	for i := range dp {
		dp[i] = make([]int, n)
		for j := range dp[i] {
			dp[i][j] = -1
		}
	}

	// 初始化路径数组
	path := make([][]int, 1<<n)
	for i := range path {
		path[i] = make([]int, n)
	}

	// 计算最短路径
	minDist := tspDP(graph, 1, 0, dp, path)

	// 重建路径
	result := TSPResult{
		Path:     reconstructPath(path, n),
		Distance: minDist,
	}

	return result
}

func tspDP(graph [][]int, mask int, pos int, dp [][]int, path [][]int) int {
	n := len(graph)
	if mask == (1<<n)-1 {
		return graph[pos][0]
	}

	if dp[mask][pos] != -1 {
		return dp[mask][pos]
	}

	ans := math.MaxInt32
	nextCity := -1

	for city := 0; city < n; city++ {
		if (mask&(1<<city)) == 0 {
			newDist := graph[pos][city] + tspDP(graph, mask|(1<<city), city, dp, path)
			if newDist < ans {
				ans = newDist
				nextCity = city
			}
		}
	}

	dp[mask][pos] = ans
	path[mask][pos] = nextCity
	return ans
}

func reconstructPath(path [][]int, n int) []int {
	result := make([]int, n+1)
	result[0] = 0
	mask := 1
	pos := 0

	for i := 1; i < n; i++ {
		pos = path[mask][pos]
		result[i] = pos
		mask |= 1 << pos
	}

	result[n] = 0
	return result
}
