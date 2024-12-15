<template>
  <div class="floyd-container">
    <h2>Floyd 多源最短路径算法</h2>
    
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
        <el-form-item label="终点">
          <el-select v-model="endNode" placeholder="选择终点">
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
    <div v-if="result" class="result-panel">
      <h3>所有节点间最短路径：</h3>
      <el-table :data="distanceTableData" style="width: 100%">
        <el-table-column
          v-for="node in nodes"
          :key="node.id"
          :prop="'node' + node.id"
          :label="'节点 ' + node.id"
        />
      </el-table>
    </div>
    
    <!-- 说明面板 -->
    <div class="info-panel">
      <h3>算法说明</h3>
      <p>Floyd-Warshall算法是一个用于寻找加权图中所有顶点对之间最短路径的算法。它通过动态规划的方法，逐步考虑所有可能的中间节点来更新最短路径。</p>
    </div>
  </div>
</template>

<script>
import * as d3 from 'd3'
import axios from 'axios'

export default {
  name: 'FloydView',
  data() {
    return {
      nodes: [],
      edges: [],
      svg: null,
      startNode: null,
      endNode: null,
      result: null
    }
  },
  computed: {
    distanceTableData() {
      if (!this.result) return []
      
      return this.nodes.map(node => {
        const row = {}
        this.nodes.forEach(otherNode => {
          row['node' + otherNode.id] = this.result.distances[node.id][otherNode.id] || '∞'
        })
        return row
      })
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
      try {
        // 转换为邻接矩阵
        const n = this.nodes.length
        const graph = Array(n).fill().map(() => Array(n).fill(0))
        this.edges.forEach(edge => {
          graph[edge.source][edge.target] = edge.weight
          graph[edge.target][edge.source] = edge.weight
        })
        
        const response = await axios.post('http://localhost:8080/api/floyd', {
          graph: graph
        })
        
        this.result = response.data.result
        
        if (this.startNode !== null && this.endNode !== null) {
          this.highlightPath()
        }
      } catch (error) {
        console.error('Error running Floyd algorithm:', error)
        this.$message.error('运行算法时出错')
      }
    },
    
    highlightPath() {
      if (!this.result || this.startNode === null || this.endNode === null) return
      
      // 重置所有边的样式
      this.svg.selectAll('.edge line')
        .style('stroke', '#999')
        .style('stroke-width', 2)
      
      // 获取路径
      let current = this.startNode
      while (current !== this.endNode) {
        const next = this.result.next[current][this.endNode]
        if (next === -1) break
        
        // 高亮路径
        this.svg.selectAll('.edge line')
          .filter(d => 
            (d.source === current && d.target === next) ||
            (d.source === next && d.target === current)
          )
          .style('stroke', '#f00')
          .style('stroke-width', 3)
        
        current = next
      }
    },
    
    resetGraph() {
      this.result = null
      this.startNode = null
      this.endNode = null
      this.drawGraph()
    }
  },
  watch: {
    startNode() {
      if (this.result) this.highlightPath()
    },
    endNode() {
      if (this.result) this.highlightPath()
    }
  }
}
</script>

<style scoped>
.floyd-container {
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
  cursor: pointer;
}

.edge {
  stroke-opacity: 0.6;
}

.edge text {
  font-size: 12px;
}
</style>
