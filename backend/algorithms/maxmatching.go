package algorithms

type MatchingResult struct {
	Matching []Edge `json:"matching"`
	Size     int    `json:"size"`
}

func MaxMatching(graph [][]int) MatchingResult {
	n := len(graph)
	if n == 0 {
		return MatchingResult{}
	}

	match := make([]int, n)
	for i := range match {
		match[i] = -1
	}

	result := MatchingResult{
		Matching: make([]Edge, 0),
	}

	// 对每个顶点尝试匹配
	for u := 0; u < n; u++ {
		if match[u] == -1 {
			visited := make([]bool, n)
			if augmentPath(graph, u, visited, match) {
				result.Size++
			}
		}
	}

	// 构建匹配边的列表
	for i := 0; i < n; i++ {
		if match[i] != -1 && i < match[i] {
			result.Matching = append(result.Matching, Edge{
				From: i,
				To:   match[i],
			})
		}
	}

	return result
}

// 寻找增广路径
func augmentPath(graph [][]int, u int, visited []bool, match []int) bool {
	for v := 0; v < len(graph); v++ {
		if graph[u][v] != 0 && !visited[v] {
			visited[v] = true

			if match[v] == -1 || augmentPath(graph, match[v], visited, match) {
				match[v] = u
				match[u] = v
				return true
			}
		}
	}
	return false
}
