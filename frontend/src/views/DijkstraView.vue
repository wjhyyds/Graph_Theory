<template>
  <div class="dijkstra-container">
    <h2>Dijkstra 最短路径算法</h2>
    
    <!-- 图形展示区域 -->
    <div ref="graphContainer" class="graph-container"></div>
    
    <!-- 控制面板 -->
    <div class="control-panel">
      <el-form :inline="true">
        <el-form-item label="起点">
          <el-select v-model="startNode" placeholder="选择起点">
            <el-option
              v-for="node in nodes"
              :key="node.id"
              :label="'节点 ' + node.id"
              :value="node.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="runAlgorithm">运行算法</el-button>
          <el-button @click="resetGraph">重置图形</el-button>
        </el-form-item>
      </el-form>
    </div>
    
    <!-- 结果展示 -->
    <div v-if="paths.length" class="result-panel">
      <h3>最短路径结果：</h3>
      <div v-for="(path, index) in paths" :key="index" class="path-item">
        <span>到节点 {{ index }} 的最短路径：</span>
        <span>{{ path.path.join(' -> ') }}</span>
        <span>距离：{{ path.distance }}</span>
      </div>
    </div>
    
    <!-- 说明面板 -->
    <div class="info-panel">
      <h3>算法说明</h3>
      <p>Dijkstra算法用于计算从一个节点到其他所有节点的最短路径。它使用贪心策略，每次选择当前最短的路径进行扩展。</p>
    </div>
  </div>
</template>

<script>
import * as d3 from 'd3'
import axios from 'axios'

export default {
  name: 'DijkstraView',
  data() {
    return {
      nodes: [],
      edges: [],
      svg: null,
      startNode: null,
      paths: []
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
      // 清除旧的内容
      this.svg.selectAll('*').remove()
      
      // 绘制边
      const edges = this.svg.selectAll('.edge')
        .data(this.edges)
        .join('g')
        .attr('class', 'edge')
      
      edges.append('line')
        .attr('x1', d => this.nodes[d.source].x)
        .attr('y1', d => this.nodes[d.source].y)
        .attr('x2', d => this.nodes[d.target].x)
        .attr('y2', d => this.nodes[d.target].y)
        .style('stroke', '#999')
        .style('stroke-width', 2)
      
      edges.append('text')
        .attr('x', d => (this.nodes[d.source].x + this.nodes[d.target].x) / 2)
        .attr('y', d => (this.nodes[d.source].y + this.nodes[d.target].y) / 2)
        .text(d => d.weight)
      
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
      if (this.startNode === null) {
        this.$message.warning('请选择起点')
        return
      }
      
      try {
        // 转换为邻接矩阵
        const n = this.nodes.length
        const graph = Array(n).fill().map(() => Array(n).fill(0))
        this.edges.forEach(edge => {
          graph[edge.source][edge.target] = edge.weight
          graph[edge.target][edge.source] = edge.weight
        })
        
        const response = await axios.post('http://localhost:8080/api/dijkstra', {
          graph: graph,
          start: this.startNode
        })
        
        this.paths = response.data.result
        this.highlightPaths()
      } catch (error) {
        console.error('Error running Dijkstra algorithm:', error)
        this.$message.error('运行算法时出错')
      }
    },
    
    highlightPaths() {
      // 重置所有边的样式
      this.svg.selectAll('.edge line')
        .style('stroke', '#999')
        .style('stroke-width', 2)
      
      // 高亮最短路径
      this.paths.forEach(path => {
        for (let i = 0; i < path.path.length - 1; i++) {
          const source = path.path[i]
          const target = path.path[i + 1]
          
          this.svg.selectAll('.edge line')
            .filter(d => 
              (d.source === source && d.target === target) ||
              (d.source === target && d.target === source)
            )
            .style('stroke', '#f00')
            .style('stroke-width', 3)
        }
      })
    },
    
    resetGraph() {
      this.paths = []
      this.startNode = null
      this.drawGraph()
    }
  }
}
</script>

<style scoped>
.dijkstra-container {
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

.path-item {
  margin: 10px 0;
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

.edge text {
  font-size: 12px;
}
</style>
