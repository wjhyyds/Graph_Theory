package algorithms

type BayesNode struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Parents      []int     `json:"parents"`
	Children     []int     `json:"children"`
	Probabilities []float64 `json:"probabilities"`
}

type BayesNetwork struct {
	Nodes []BayesNode `json:"nodes"`
}

type BayesResult struct {
	Network     BayesNetwork `json:"network"`
	Probability float64      `json:"probability"`
}

func Bayes(network BayesNetwork, evidence map[int]bool) BayesResult {
	// 计算给定证据下的概率
	prob := calculateProbability(network, evidence)

	return BayesResult{
		Network:     network,
		Probability: prob,
	}
}

func calculateProbability(network BayesNetwork, evidence map[int]bool) float64 {
	prob := 1.0

	// 对每个节点计算条件概率
	for _, node := range network.Nodes {
		if value, exists := evidence[node.ID]; exists {
			parentValues := make([]bool, len(node.Parents))
			for i, parentID := range node.Parents {
				parentValues[i] = evidence[parentID]
			}
			
			index := getIndexFromParentValues(parentValues)
			if value {
				prob *= node.Probabilities[index]
			} else {
				prob *= (1 - node.Probabilities[index])
			}
		}
	}

	return prob
}

func getIndexFromParentValues(parentValues []bool) int {
	index := 0
	for i, value := range parentValues {
		if value {
			index |= 1 << i
		}
	}
	return index
}
