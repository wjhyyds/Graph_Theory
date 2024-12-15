package algorithms

import "math"

type Edge struct {
	From   int `json:"from"`
	To     int `json:"to"`
	Weight int `json:"weight"`
}

func Prim(graph [][]int) []Edge {
	n := len(graph)
	if n == 0 {
		return nil
	}

	visited := make([]bool, n)
	minEdge := make([]int, n)
	parent := make([]int, n)
	result := make([]Edge, 0)

	// 初始化
	for i := 0; i < n; i++ {
		minEdge[i] = math.MaxInt32
		parent[i] = -1
	}

	// 从节点0开始
	minEdge[0] = 0

	for i := 0; i < n; i++ {
		// 找到最小边的顶点
		u := -1
		minWeight := math.MaxInt32
		for v := 0; v < n; v++ {
			if !visited[v] && minEdge[v] < minWeight {
				u = v
				minWeight = minEdge[v]
			}
		}

		if u == -1 {
			break
		}

		visited[u] = true

		// 如果不是第一个顶点，添加边到结果中
		if parent[u] != -1 {
			result = append(result, Edge{
				From:   parent[u],
				To:     u,
				Weight: graph[u][parent[u]],
			})
		}

		// 更新相邻顶点的最小边权值
		for v := 0; v < n; v++ {
			if !visited[v] && graph[u][v] != 0 && graph[u][v] < minEdge[v] {
				minEdge[v] = graph[u][v]
				parent[v] = u
			}
		}
	}

	return result
}
