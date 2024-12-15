package algorithms

import "math"

type Path struct {
	Distance int   `json:"distance"`
	Path     []int `json:"path"`
}

func Dijkstra(graph [][]int, start int) []Path {
	n := len(graph)
	if n == 0 {
		return nil
	}

	dist := make([]int, n)
	visited := make([]bool, n)
	prev := make([]int, n)

	// 初始化
	for i := 0; i < n; i++ {
		dist[i] = math.MaxInt32
		prev[i] = -1
	}
	dist[start] = 0

	for i := 0; i < n; i++ {
		// 找到最短距离的顶点
		u := -1
		minDist := math.MaxInt32
		for v := 0; v < n; v++ {
			if !visited[v] && dist[v] < minDist {
				u = v
				minDist = dist[v]
			}
		}

		if u == -1 {
			break
		}

		visited[u] = true

		// 更新相邻顶点的距离
		for v := 0; v < n; v++ {
			if !visited[v] && graph[u][v] != 0 {
				newDist := dist[u] + graph[u][v]
				if newDist < dist[v] {
					dist[v] = newDist
					prev[v] = u
				}
			}
		}
	}

	// 构建结果
	result := make([]Path, n)
	for i := 0; i < n; i++ {
		path := []int{i}
		for curr := i; curr != start && prev[curr] != -1; curr = prev[curr] {
			path = append([]int{prev[curr]}, path...)
		}
		result[i] = Path{
			Distance: dist[i],
			Path:     path,
		}
	}

	return result
}
