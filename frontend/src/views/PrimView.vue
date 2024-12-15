<template>
  <div class="prim-container">
    <h2>Prim 最小生成树算法</h2>
    
    <!-- 图形展示区域 -->
    <div ref="graphContainer" class="graph-container"></div>
    
    <!-- 控制面板 -->
    <div class="control-panel">
      <el-button type="primary" @click="runAlgorithm">运行算法</el-button>
      <el-button @click="resetGraph">重置图形</el-button>
    </div>
    
    <!-- 说明面板 -->
    <div class="info-panel">
      <h3>算法说明</h3>
      <p>Prim算法是一种用于寻找最小生成树的贪心算法。它的基本思想是从一个起始顶点开始，逐步选择代价最小的边来扩展树。</p>
    </div>
  </div>
</template>

<script>
import * as d3 from 'd3'
import axios from 'axios'

export default {
  name: 'PrimView',
  data() {
    return {
      nodes: [],
      edges: [],
      svg: null
    }
  },
  mounted() {
    this.initGraph()
  },
  methods: {
    initGraph() {
      // 初始化D3图形
      const width = 600
      const height = 400
      
      this.svg = d3.select(this.$refs.graphContainer)
        .append('svg')
        .attr('width', width)
        .attr('height', height)
      
      // 添加一些示例节点和边
      this.nodes = [
        {id: 0, x: 100, y: 100},
        {id: 1, x: 200, y: 150},
        {id: 2, x: 300, y: 100},
        {id: 3, x: 200, y: 250}
      ]
      
      this.edges = [
        {source: 0, target: 1, weight: 4},
        {source: 1, target: 2, weight: 3},
        {source: 2, target: 3, weight: 5},
        {source: 3, target: 0, weight: 2}
      ]
      
      this.drawGraph()
    },
    
    drawGraph() {
      // 绘制边
      this.svg.selectAll('.edge')
        .data(this.edges)
        .join('line')
        .attr('class', 'edge')
        .attr('x1', d => this.nodes[d.source].x)
        .attr('y1', d => this.nodes[d.source].y)
        .attr('x2', d => this.nodes[d.target].x)
        .attr('y2', d => this.nodes[d.target].y)
        .style('stroke', '#999')
        .style('stroke-width', 2)
      
      // 绘制节点
      this.svg.selectAll('.node')
        .data(this.nodes)
        .join('circle')
        .attr('class', 'node')
        .attr('cx', d => d.x)
        .attr('cy', d => d.y)
        .attr('r', 10)
        .style('fill', '#69c')
    },
    
    async runAlgorithm() {
      try {
        const response = await axios.post('http://localhost:8080/api/prim', {
          graph: this.convertToAdjMatrix()
        })
        // 处理返回结果，更新图形显示
        console.log(response.data)
      } catch (error) {
        console.error('Error running Prim algorithm:', error)
      }
    },
    
    convertToAdjMatrix() {
      // 将当前图形转换为邻接矩阵
      const n = this.nodes.length
      const matrix = Array(n).fill().map(() => Array(n).fill(Infinity))
      
      this.edges.forEach(edge => {
        matrix[edge.source][edge.target] = edge.weight
        matrix[edge.target][edge.source] = edge.weight
      })
      
      return matrix
    },
    
    resetGraph() {
      this.svg.selectAll('*').remove()
      this.initGraph()
    }
  }
}
</script>

<style scoped>
.prim-container {
  padding: 20px;
}

.graph-container {
  border: 1px solid #ddd;
  margin: 20px 0;
  min-height: 400px;
}

.control-panel {
  margin: 20px 0;
}

.info-panel {
  background-color: #f5f7fa;
  padding: 15px;
  border-radius: 4px;
}

.node {
  cursor: pointer;
}

.edge {
  stroke-opacity: 0.6;
}
</style>
