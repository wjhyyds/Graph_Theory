package algorithms

import (
    "math"
)

type FloydResult struct {
    Distances [][]int `json:"distances"`
    Next      [][]int `json:"next"`
    Path      []int   `json:"path"`
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
            if i == j {
                dist[i][j] = 0
                next[i][j] = j
            } else if graph[i][j] != math.MaxInt32 && graph[i][j] > 0 {
                dist[i][j] = graph[i][j]
                next[i][j] = j
            } else {
                dist[i][j] = math.MaxInt32
                next[i][j] = -1
            }
        }
    }

    // Floyd-Warshall算法
    for k := 0; k < n; k++ {
        for i := 0; i < n; i++ {
            for j := 0; j < n; j++ {
                if dist[i][k] != math.MaxInt32 && dist[k][j] != math.MaxInt32 {
                    newDist := dist[i][k] + dist[k][j]
                    if newDist < dist[i][j] && newDist >= 0 { // 添加溢出检查
                        dist[i][j] = newDist
                        next[i][j] = next[i][k]
                    }
                }
            }
        }
    }

    // 将无穷大的距离转换为-1（表示不可达）
    for i := 0; i < n; i++ {
        for j := 0; j < n; j++ {
            if dist[i][j] == math.MaxInt32 {
                dist[i][j] = -1
                next[i][j] = -1
            }
        }
    }

    return FloydResult{
        Distances: dist,
        Next:      next,
    }
}

// GetPath 返回从start到end的路径
func GetPath(next [][]int, start, end int) []int {
    if next[start][end] == -1 {
        return []int{}
    }

    path := []int{start}
    current := start
    for current != end {
        current = next[current][end]
        if current == -1 {
            return []int{}
        }
        path = append(path, current)
    }
    return path
}