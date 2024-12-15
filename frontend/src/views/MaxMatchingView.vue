<template>
  <div class="max-matching-container">
    <h2>最大匹配算法</h2>
    
    <!-- 图形展示区域 -->
    <div ref="graphContainer" class="graph-container"></div>
    
    <!-- 控制面板 -->
    <div class="control-panel">
      <el-button type="primary" @click="runAlgorithm">运行算法</el-button>
      <el-button @click="resetGraph">重置图形</el-button>
    </div>
    
    <!-- 结果展示 -->
    <div v-if="result" class="result-panel">
      <h3>最大匹配结果：</h3>
      <p>匹配数量：{{ result.size }}</p>
      <div v-for="(edge, index) in result.matching" :key="index" class="matching-item">
        节点 {{ edge.from }} - 节点 {{ edge.to }}
      </div>
    </div>
    
    <!-- 说明面板 -->
    <div class="info-panel">
      <h3>算法说明</h3>
      <p>最大匹配算法用于在无向图中寻找最大数量的边，使得这些边没有公共顶点。该算法常用于任务分配等问题。</p>
    </div>
  </div>
</template>

<script>
import * as d3 from 'd3'
import axios from 'axios'

export default {
  name: 'MaxMatchingView',
  data() {
    return {
      nodes: [],
      edges: [],
      svg: null,
      result: null
    }
  },
  mounted() {
    this.initGraph()
  },
  methods: {
    initGraph() {
      const width = 600
      const height = 400
      
      this.svg = d3.select(this.$refs.graphContainer)
        .append('svg')
        .attr('width', width)
        .attr('height', height)
      
      // 示例数据
      this.nodes = [
        {id: 0, x: 100, y: 100},
        {id: 1, x: 300, y: 100},
        {id: 2, x: 100, y: 300},
        {id: 3, x: 300, y: 300}
      ]
      
      this.edges = [
        {source: 0, target: 1},
        {source: 0, target: 2},
        {source: 1, target: 3},
        {source: 2, target: 3}
      ]
      
      this.drawGraph()
    },
    
    drawGraph() {
      // 清除旧的内容
      this.svg.selectAll('*').remove()
      
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
      const nodes = this.svg.selectAll('.node')
        .data(this.nodes)
        .join('g')
        .attr('class', 'node')
      
      nodes.append('circle')
        .attr('cx', d => d.x)
        .attr('cy', d => d.y)
        .attr('r', 15)
        .style('fill', '#69c')
      
      nodes.append('text')
        .attr('x', d => d.x)
        .attr('y', d => d.y)
        .attr('dy', '0.35em')
        .attr('text-anchor', 'middle')
        .style('fill', 'white')
        .text(d => d.id)
    },
    
    async runAlgorithm() {
      try {
        // 转换为邻接矩阵
        const n = this.nodes.length
        const graph = Array(n).fill().map(() => Array(n).fill(0))
        this.edges.forEach(edge => {
          graph[edge.source][edge.target] = 1
          graph[edge.target][edge.source] = 1
        })
        
        const response = await axios.post('http://localhost:8080/api/max-matching', {
          graph: graph
        })
        
        this.result = response.data.result
        this.highlightMatching()
      } catch (error) {
        console.error('Error running Maximum Matching algorithm:', error)
        this.$message.error('运行算法时出错')
      }
    },
    
    highlightMatching() {
      if (!this.result) return
      
      // 重置所有边的样式
      this.svg.selectAll('.edge')
        .style('stroke', '#999')
        .style('stroke-width', 2)
      
      // 高亮匹配边
      this.result.matching.forEach(edge => {
        this.svg.selectAll('.edge')
          .filter(d => 
            (d.source === edge.from && d.target === edge.to) ||
            (d.source === edge.to && d.target === edge.from)
          )
          .style('stroke', '#f00')
          .style('stroke-width', 3)
      })
    },
    
    resetGraph() {
      this.result = null
      this.drawGraph()
    }
  }
}
</script>

<style scoped>
.max-matching-container {
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

.result-panel {
  margin: 20px 0;
  padding: 15px;
  background-color: #f8f9fa;
  border-radius: 4px;
}

.matching-item {
  margin: 5px 0;
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
