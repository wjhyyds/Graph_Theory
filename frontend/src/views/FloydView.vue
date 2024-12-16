<template>
  <div class="floyd-container">
    <h2>Floyd 多源最短路径算法</h2>
    
    <!-- 图形展示区域 -->
    <div ref="graphContainer" class="graph-container"></div>
    
    <!-- 控制面板 -->
    <div class="control-panel">
      <el-form :inline="true">
        <el-form-item label="起点">
          <el-select v-model="startNode" placeholder="选择起点" @change="onStartNodeChange">
            <el-option
              v-for="node in nodeOptions"
              :key="node.id"
              :label="`节点 ${node.id}`"
              :value="node.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="终点">
          <el-select v-model="endNode" placeholder="选择终点" @change="onEndNodeChange">
            <el-option
              v-for="node in nodeOptions"
              :key="node.id"
              :label="`节点 ${node.id}`"
              :value="node.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="runAlgorithm">运行算法</el-button>
          <el-button @click="resetGraph">重置图形</el-button>
          <el-button @click="addNode">添加节点</el-button>
          <el-button @click="toggleDragMode" :type="isDragMode ? 'success' : ''">
            {{ isDragMode ? '拖拽模式开启' : '拖拽模式关闭' }}
          </el-button>
          <el-button @click="toggleDeleteMode" :type="isDeleteMode ? 'danger' : ''">
            {{ isDeleteMode ? '删除模式开启' : '删除模式关闭' }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 提示信息 -->
    <div class="instruction-panel" v-if="isDragMode || isDeleteMode">
      <p v-if="isDragMode">拖拽模式已开启：点击并拖动节点可以改变它们的位置</p>
      <p v-if="isDeleteMode">删除模式已开启：点击边可以删除它们</p>
    </div>
    
    <!-- 结果展示 -->
    <div v-if="result" class="result-panel">
      <h3>所有节点间最短路径：</h3>
      <el-table :data="distanceTableData" style="width: 100%">
        <el-table-column
          v-for="node in nodeOptions"
          :key="node.id"
          :prop="`node${node.id}`"
          :label="`节点 ${node.id}`"
        />
      </el-table>
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
      simulation: null,
      startNode: null,
      endNode: null,
      result: null,
      isDragMode: false,
      isDeleteMode: false,
      nextNodeId: 1,
      selectedStartNode: null,
      selectedEndNode: null
    }
  },

  computed: {
    // 添加计算属性用于节点选择框
    nodeOptions() {
      return this.nodes.map(node => ({
        id: node.id,
        label: `节点 ${node.id}`
      }))
    },
    // 添加计算属性用于距离表格数据
    distanceTableData() {
      if (!this.result || !this.result.distances) return []
      
      return this.nodes.map((node, i) => {
        const row = {}
        this.nodes.forEach((otherNode, j) => {
          const distance = this.result.distances[i][j]
          row[`node${otherNode.id}`] = distance === -1 ? '∞' : distance
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
        .on('wheel', (event) => {
          event.preventDefault() // 阻止滚轮事件
        })
      
      this.simulation = d3.forceSimulation(this.nodes)
        .force('link', d3.forceLink(this.edges).distance(100))
        .force('charge', d3.forceManyBody().strength(-300))
        .force('center', d3.forceCenter(width / 2, height / 2))
        .on('tick', this.updateGraph)
    },

    addNode() {
      const width = 600
      const height = 400
      const newNode = {
        id: this.nextNodeId++,
        x: width / 2 + (Math.random() - 0.5) * 100,
        y: height / 2 + (Math.random() - 0.5) * 100
      }
      
      this.nodes.push(newNode)
      
      // 如果已经有其他节点，创建到所有其他节点的边
      if (this.nodes.length > 1) {
        const otherNodes = this.nodes.filter(n => n.id !== newNode.id)
        otherNodes.forEach(otherNode => {
          const distance = Math.round(Math.sqrt(
            Math.pow(newNode.x - otherNode.x, 2) + 
            Math.pow(newNode.y - otherNode.y, 2)
          ))
          this.edges.push({
            source: newNode,
            target: otherNode,
            distance: distance
          })
        })
      }
      
      this.simulation.nodes(this.nodes)
      this.simulation.force('link').links(this.edges)
      this.simulation.alpha(1).restart()
      
      this.drawGraph()
    },

    drawGraph() {
      // 绘制边
      const links = this.svg.selectAll('.link')
        .data(this.edges)
        .join('line')
        .attr('class', 'link')
        .style('stroke', '#999')
        .style('stroke-width', 2)
        .on('click', (event, d) => {
          if (this.isDeleteMode) {
            this.deleteEdge(d)
          }
        })
        .on('mouseover', (event, d) => {
          if (this.isDeleteMode) {
            d3.select(event.target)
              .style('stroke', '#f00')
              .style('stroke-width', 3)
          }
        })
        .on('mouseout', (event, d) => {
          if (this.isDeleteMode) {
            d3.select(event.target)
              .style('stroke', '#999')
              .style('stroke-width', 2)
          }
        })

      // 绘制边的权重标签
      const edgeLabels = this.svg.selectAll('.edge-label')
        .data(this.edges)
        .join('text')
        .attr('class', 'edge-label')
        .style('font-size', '12px')
        .style('background', 'white')
        .style('padding', '2px')
        .text(d => d.distance)

      // 绘制节点
      const nodes = this.svg.selectAll('.node')
        .data(this.nodes)
        .join('g')
        .attr('class', 'node')
        .call(d3.drag()
          .on('start', this.dragstarted)
          .on('drag', this.dragged)
          .on('end', this.dragended))

      // 绘制节点的圆圈
      nodes.selectAll('circle')
        .data(d => [d])
        .join('circle')
        .attr('r', 15)
        .style('fill', d => {
          if (d.id === this.startNode) return '#4CAF50'
          if (d.id === this.endNode) return '#2196F3'
          return '#fff'
        })
        .style('stroke', '#000')
        .style('stroke-width', 2)

      // 绘制节点的文本标签
      nodes.selectAll('text')
        .data(d => [d])
        .join('text')
        .text(d => d.id)
        .attr('text-anchor', 'middle')
        .attr('dy', '.3em')
        .style('pointer-events', 'none')

      // 更新力导向图的节点和边
      this.simulation.nodes(this.nodes)
      this.simulation.force('link').links(this.edges)
      this.simulation.alpha(1).restart()
    },

    updateGraph() {
      // 更新边的位置
      this.svg.selectAll('.link')
        .attr('x1', d => d.source.x)
        .attr('y1', d => d.source.y)
        .attr('x2', d => d.target.x)
        .attr('y2', d => d.target.y)

      // 更新边权重标签的位置
      this.svg.selectAll('.edge-label')
        .attr('x', d => (d.source.x + d.target.x) / 2)
        .attr('y', d => (d.source.y + d.target.y) / 2)

      // 更新节点的位置
      this.svg.selectAll('.node')
        .attr('transform', d => `translate(${d.x},${d.y})`)
    },

    // 删除边
    deleteEdge(edge) {
      const index = this.edges.indexOf(edge)
      if (index > -1) {
        this.edges.splice(index, 1)
        this.drawGraph()
      }
    },

    // 切换删除模式
    toggleDeleteMode() {
      this.isDeleteMode = !this.isDeleteMode
      if (this.isDeleteMode) {
        this.isDragMode = false
      }
    },

    // 起点变化处理
    onStartNodeChange(value) {
      this.startNode = value
      this.drawGraph() // 重新绘制以更新节点颜色
    },

    // 终点变化处理
    onEndNodeChange(value) {
      this.endNode = value
      this.drawGraph() // 重新绘制以更新节点颜色
    },

    dragstarted(event, d) {
      if (!this.isDragMode) return
      if (!event.active) this.simulation.alphaTarget(0.3).restart()
      d.fx = d.x
      d.fy = d.y
    },

    dragged(event, d) {
      if (!this.isDragMode) return
      d.fx = event.x
      d.fy = event.y
      
      // 更新与该节点相连的边的距离
      this.edges.forEach(edge => {
        if (edge.source === d || edge.target === d) {
          const distance = Math.round(Math.sqrt(
            Math.pow(edge.source.x - edge.target.x, 2) + 
            Math.pow(edge.source.y - edge.target.y, 2)
          ))
          edge.distance = distance
        }
      })

      // 更新边的权重标签
      this.svg.selectAll(".edge-label")
        .text(d => d.distance)
        .attr("x", d => {
          const midX = (d.source.x + d.target.x) / 2
          const dx = d.target.x - d.source.x
          const dy = d.target.y - d.source.y
          const angle = Math.atan2(dy, dx)
          return midX + Math.sin(angle) * 10
        })
        .attr("y", d => {
          const midY = (d.source.y + d.target.y) / 2
          const dx = d.target.x - d.source.x
          const dy = d.target.y - d.source.y
          const angle = Math.atan2(dy, dx)
          return midY - Math.cos(angle) * 10
        })
    },

    dragended(event, d) {
      if (!this.isDragMode) return
      if (!event.active) this.simulation.alphaTarget(0)
      d.fx = null
      d.fy = null
    },

    toggleDragMode() {
      this.isDragMode = !this.isDragMode
    },

    resetGraph() {
      this.nodes = []
      this.edges = []
      this.nextNodeId = 1
      this.startNode = null
      this.endNode = null
      this.result = null
      this.initGraph()
    },

    async runAlgorithm() {
      if (!this.startNode || !this.endNode) {
        this.$message.warning('请选择起点和终点')
        return
      }

      // 构建邻接矩阵
      const n = this.nodes.length
      const graph = Array(n).fill().map(() => Array(n).fill(Number.MAX_SAFE_INTEGER))
      
      // 设置对角线为0
      for (let i = 0; i < n; i++) {
        graph[i][i] = 0
      }
      
      // 使用节点ID到索引的映射
      const nodeIdToIndex = new Map(this.nodes.map((node, index) => [node.id, index]))
      
      this.edges.forEach(edge => {
        const sourceIndex = nodeIdToIndex.get(edge.source.id)
        const targetIndex = nodeIdToIndex.get(edge.target.id)
        const distance = Math.round(Math.sqrt(
          Math.pow(edge.source.x - edge.target.x, 2) + 
          Math.pow(edge.source.y - edge.target.y, 2)
        ))
        
        if (sourceIndex !== undefined && targetIndex !== undefined && distance > 0) {
          graph[sourceIndex][targetIndex] = distance
          graph[targetIndex][sourceIndex] = distance // 无向图需要双向设置
        }
      })

      try {
        console.log('Sending graph to backend:', graph) // 调试日志
        const response = await axios.post('http://localhost:8080/api/floyd', { graph })
        console.log('Received response:', response.data) // 调试日志
        
        if (response.data && response.data.result) {
          // 从 result 字段中获取数据
          this.result = response.data.result
          
          // 使用相同的映射获取起点和终点的索引
          const startIndex = nodeIdToIndex.get(this.startNode)
          const endIndex = nodeIdToIndex.get(this.endNode)
          
          if (startIndex !== undefined && endIndex !== undefined) {
            this.highlightPath(startIndex, endIndex)
          }
        } else {
          throw new Error('Invalid response format from server')
        }
      } catch (error) {
        console.error('Error:', error)
        this.$message.error('算法执行失败：' + (error.response?.data?.error || error.message))
      }
    },

    highlightPath(startIndex, endIndex) {
      if (!this.result || !this.result.next) {
        console.error('No result or next matrix available')
        return
      }

      console.log('Highlighting path from', startIndex, 'to', endIndex)
      console.log('Next matrix:', this.result.next)

      // 重置所有边的样式
      this.svg.selectAll('.link')
        .style('stroke', '#999')
        .style('stroke-width', 2)

      // 获取路径
      const path = []
      let current = startIndex
      while (current !== endIndex) {
        const next = this.result.next[current][endIndex]
        if (next === -1) {
          console.log('No path found between nodes')
          break
        }
        
        path.push([current, next])
        current = next
      }

      console.log('Found path:', path)

      // 高亮路径
      this.svg.selectAll('.link')
        .style('stroke', (d) => {
          const sourceIndex = this.nodes.findIndex(n => n.id === d.source.id)
          const targetIndex = this.nodes.findIndex(n => n.id === d.target.id)
          
          return path.some(([a, b]) => 
            (a === sourceIndex && b === targetIndex) || 
            (a === targetIndex && b === sourceIndex)
          ) ? '#f00' : '#999'
        })
        .style('stroke-width', (d) => {
          const sourceIndex = this.nodes.findIndex(n => n.id === d.source.id)
          const targetIndex = this.nodes.findIndex(n => n.id === d.target.id)
          
          return path.some(([a, b]) => 
            (a === sourceIndex && b === targetIndex) || 
            (a === targetIndex && b === sourceIndex)
          ) ? 4 : 2
        })
    }
  }
}
</script>

<style scoped>
.floyd-container {
  padding: 20px;
}

.graph-container {
  width: 600px;
  height: 400px;
  border: 1px solid #ddd;
  margin: 20px 0;
  overflow: hidden;
}

.control-panel {
  margin-bottom: 20px;
}

.instruction-panel {
  margin: 10px 0;
  padding: 10px;
  background-color: #f5f5f5;
  border-radius: 4px;
}

.instruction-panel p {
  margin: 5px 0;
  color: #666;
}

.result-panel {
  margin-top: 20px;
}

.node circle {
  cursor: pointer;
}

.node text {
  font-size: 12px;
}

.link {
  cursor: pointer;
}

.edge-label {
  pointer-events: none;
  background: white;
  padding: 2px;
}

/* 起点和终点的样式 */
.node.start circle {
  fill: #4CAF50;
}

.node.end circle {
  fill: #2196F3;
}

/* 删除模式下的边样式 */
.link:hover {
  stroke-width: 3px;
}
</style>