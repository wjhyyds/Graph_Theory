<template>
  <div class="tsp-container">
    <h2>旅行商问题 (TSP)</h2>
    
    <!-- 图形展示区域 -->
    <div ref="graphContainer" class="graph-container"></div>
    
    <!-- 控制面板 -->
    <div class="control-panel">
      <el-button type="primary" @click="runAlgorithm">运行算法</el-button>
      <el-button @click="resetGraph">重置图形</el-button>
      <el-button @click="addNode">添加城市</el-button>
    </div>
    
    <!-- 结果展示 -->
    <div v-if="result" class="result-panel">
      <h3>最优路径：</h3>
      <p>路径顺序：{{ result.path.join(' -> ') }}</p>
      <p>总距离：{{ result.distance }}</p>
    </div>
    
    <!-- 说明面板 -->
    <div class="info-panel">
      <h3>算法说明</h3>
      <p>旅行商问题(TSP)是寻找访问所有城市一次并返回起点的最短路径。这是一个NP难问题，我们使用动态规划方法求解。</p>
    </div>
  </div>
</template>

<script>
import * as d3 from 'd3'
import axios from 'axios'

export default {
  name: 'TSPView',
  data() {
    return {
      nodes: [],
      edges: [],
      svg: null,
      result: null,
      simulation: null
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
      
      // 初始化力导向图
      this.simulation = d3.forceSimulation()
        .force('link', d3.forceLink().id(d => d.id))
        .force('charge', d3.forceManyBody().strength(-300))
        .force('center', d3.forceCenter(width / 2, height / 2))
      
      // 示例数据
      this.nodes = [
        {id: 0, x: 100, y: 100},
        {id: 1, x: 200, y: 150},
        {id: 2, x: 300, y: 100},
        {id: 3, x: 200, y: 250}
      ]
      
      this.updateGraph()
    },
    
    updateGraph() {
      // 更新力导向图
      this.simulation
        .nodes(this.nodes)
        .on('tick', () => this.drawGraph())
      
      this.simulation.force('link')
        .links(this.createEdges())
      
      this.drawGraph()
    },
    
    createEdges() {
      // 创建完全图的边
      const edges = []
      for (let i = 0; i < this.nodes.length; i++) {
        for (let j = i + 1; j < this.nodes.length; j++) {
          edges.push({
            source: this.nodes[i],
            target: this.nodes[j],
            weight: Math.round(Math.sqrt(
              Math.pow(this.nodes[i].x - this.nodes[j].x, 2) +
              Math.pow(this.nodes[i].y - this.nodes[j].y, 2)
            ))
          })
        }
      }
      return edges
    },
    
    drawGraph() {
      // 清除旧的内容
      this.svg.selectAll('*').remove()
      
      const edges = this.createEdges()
      
      // 绘制边
      this.svg.selectAll('.edge')
        .data(edges)
        .join('line')
        .attr('class', 'edge')
        .attr('x1', d => d.source.x)
        .attr('y1', d => d.source.y)
        .attr('x2', d => d.target.x)
        .attr('y2', d => d.target.y)
        .style('stroke', '#999')
        .style('stroke-width', 1)
      
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
        .call(d3.drag()
          .on('start', this.dragstarted)
          .on('drag', this.dragged)
          .on('end', this.dragended))
      
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
        // 创建距离矩阵
        const n = this.nodes.length
        const graph = Array(n).fill().map(() => Array(n).fill(0))
        
        for (let i = 0; i < n; i++) {
          for (let j = 0; j < n; j++) {
            if (i !== j) {
              graph[i][j] = Math.round(Math.sqrt(
                Math.pow(this.nodes[i].x - this.nodes[j].x, 2) +
                Math.pow(this.nodes[i].y - this.nodes[j].y, 2)
              ))
            }
          }
        }
        
        const response = await axios.post('http://localhost:8080/api/tsp', {
          graph: graph
        })
        
        this.result = response.data.result
        this.highlightPath()
      } catch (error) {
        console.error('Error running TSP algorithm:', error)
        this.$message.error('运行算法时出错')
      }
    },
    
    highlightPath() {
      if (!this.result) return
      
      // 重置所有边的样式
      this.svg.selectAll('.edge')
        .style('stroke', '#999')
        .style('stroke-width', 1)
      
      // 高亮路径
      for (let i = 0; i < this.result.path.length - 1; i++) {
        const source = this.result.path[i]
        const target = this.result.path[i + 1]
        
        this.svg.selectAll('.edge')
          .filter(d => 
            (d.source.id === source && d.target.id === target) ||
            (d.source.id === target && d.target.id === source)
          )
          .style('stroke', '#f00')
          .style('stroke-width', 3)
      }
    },
    
    addNode() {
      const id = this.nodes.length
      const x = Math.random() * 500 + 50
      const y = Math.random() * 300 + 50
      
      this.nodes.push({id, x, y})
      this.updateGraph()
    },
    
    resetGraph() {
      this.result = null
      this.updateGraph()
    },
    
    // 拖拽相关函数
    dragstarted(event, d) {
      if (!event.active) this.simulation.alphaTarget(0.3).restart()
      d.fx = d.x
      d.fy = d.y
    },
    
    dragged(event, d) {
      d.fx = event.x
      d.fy = event.y
    },
    
    dragended(event, d) {
      if (!event.active) this.simulation.alphaTarget(0)
      d.fx = null
      d.fy = null
    }
  }
}
</script>

<style scoped>
.tsp-container {
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

.info-panel {
  background-color: #f5f7fa;
  padding: 15px;
  border-radius: 4px;
}

.node {
  cursor: move;
}

.edge {
  stroke-opacity: 0.6;
}
</style>
