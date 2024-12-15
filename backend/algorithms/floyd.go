package algorithms

type FloydResult struct {
	Distances [][]int `json:"distances"`
	Next      [][]int `json:"next"`
}

func Floyd(graph [][]int) FloydResult {
	n := len(graph)
	if n == 0 {
		return FloydResult{}
	}

	// 初始化距离矩阵和下一跳矩阵
	dist := make([][]int, n)
	next := make([][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int, n)
		next[i] = make([]int, n)
		for j := 0; j < n; j++ {
			dist[i][j] = graph[i][j]
			if graph[i][j] != 0 {
				next[i][j] = j
			} else {
				next[i][j] = -1
			}
		}
	}

	// Floyd-Warshall算法
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if dist[i][k] != 0 && dist[k][j] != 0 {
					if dist[i][j] == 0 || dist[i][j] > dist[i][k]+dist[k][j] {
						dist[i][j] = dist[i][k] + dist[k][j]
						next[i][j] = next[i][k]
					}
				}
			}
		}
	}

	return FloydResult{
		Distances: dist,
		Next:      next,
	}
}

// 获取从start到end的路径
func GetPath(next [][]int, start, end int) []int {
	if next[start][end] == -1 {
		return nil
	}

	path := []int{start}
	for start != end {
		start = next[start][end]
		path = append(path, start)
	}
	return path
}
