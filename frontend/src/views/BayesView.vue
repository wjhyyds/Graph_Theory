<template>
  <div class="bayes-container">
    <h2>贝叶斯网络</h2>
    
    <!-- 图形展示区域 -->
    <div ref="graphContainer" class="graph-container"></div>
    
    <!-- 控制面板 -->
    <div class="control-panel">
      <el-form :inline="true">
        <el-form-item v-for="node in network.nodes" :key="node.id" :label="node.name">
          <el-switch v-model="evidence[node.id]" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="runAlgorithm">计算概率</el-button>
          <el-button @click="resetGraph">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    
    <!-- 结果展示 -->
    <div v-if="result" class="result-panel">
      <h3>计算结果：</h3>
      <p>条件概率：{{ result.probability.toFixed(4) }}</p>
    </div>
    
    <!-- 说明面板 -->
    <div class="info-panel">
      <h3>算法说明</h3>
      <p>贝叶斯网络是一种概率图模型，用于表示随机变量之间的条件依赖关系。通过贝叶斯定理，我们可以计算给定观测证据下的条件概率。</p>
    </div>
  </div>
</template>

<script>
import * as d3 from 'd3'
import axios from 'axios'

export default {
  name: 'BayesView',
  data() {
    return {
      network: {
        nodes: [
          {
            id: 0,
            name: '下雨',
            parents: [],
            children: [1, 2],
            probabilities: [0.2]
          },
          {
            id: 1,
            name: '地面湿',
            parents: [0],
            children: [],
            probabilities: [0.9, 0.1]  // P(湿|雨), P(湿|无雨)
          },
          {
            id: 2,
            name: '交通堵塞',
            parents: [0],
            children: [],
            probabilities: [0.7, 0.3]  // P(堵|雨), P(堵|无雨)
          }
        ]
      },
      evidence: {},
      svg: null,
      result: null
    }
  },
  mounted() {
    this.initGraph()
    this.resetEvidence()
  },
  methods: {
    initGraph() {
      const width = 600
      const height = 400
      
      this.svg = d3.select(this.$refs.graphContainer)
        .append('svg')
        .attr('width', width)
        .attr('height', height)
      
      // 设置节点位置
      this.network.nodes.forEach((node, index) => {
        node.x = 300
        node.y = 100 + index * 100
      })
      
      this.drawGraph()
    },
    
    drawGraph() {
      // 清除旧的内容
      this.svg.selectAll('*').remove()
      
      // 创建箭头标记
      this.svg.append('defs').append('marker')
        .attr('id', 'arrowhead')
        .attr('viewBox', '-0 -5 10 10')
        .attr('refX', 25)
        .attr('refY', 0)
        .attr('orient', 'auto')
        .attr('markerWidth', 6)
        .attr('markerHeight', 6)
        .attr('xoverflow', 'visible')
        .append('svg:path')
        .attr('d', 'M 0,-5 L 10 ,0 L 0,5')
        .attr('fill', '#999')
        .style('stroke', 'none')
      
      // 绘制边
      this.network.nodes.forEach(node => {
        node.children.forEach(childId => {
          const child = this.network.nodes.find(n => n.id === childId)
          this.svg.append('line')
            .attr('class', 'edge')
            .attr('x1', node.x)
            .attr('y1', node.y)
            .attr('x2', child.x)
            .attr('y2', child.y)
            .attr('marker-end', 'url(#arrowhead)')
            .style('stroke', '#999')
            .style('stroke-width', 2)
        })
      })
      
      // 绘制节点
      const nodes = this.svg.selectAll('.node')
        .data(this.network.nodes)
        .join('g')
        .attr('class', 'node')
        .attr('transform', d => `translate(${d.x},${d.y})`)
      
      nodes.append('circle')
        .attr('r', 20)
        .style('fill', d => this.evidence[d.id] ? '#69c' : '#fff')
        .style('stroke', '#69c')
        .style('stroke-width', 2)
      
      nodes.append('text')
        .attr('dy', '0.35em')
        .attr('text-anchor', 'middle')
        .style('fill', d => this.evidence[d.id] ? '#fff' : '#000')
        .text(d => d.name)
    },
    
    async runAlgorithm() {
      try {
        const response = await axios.post('http://localhost:8080/api/bayes', {
          network: this.network,
          evidence: this.evidence
        })
        
        this.result = response.data.result
      } catch (error) {
        console.error('Error running Bayes algorithm:', error)
        this.$message.error('运行算法时出错')
      }
    },
    
    resetEvidence() {
      this.evidence = {}
      this.network.nodes.forEach(node => {
        this.evidence[node.id] = false
      })
    },
    
    resetGraph() {
      this.resetEvidence()
      this.result = null
      this.drawGraph()
    }
  },
  watch: {
    evidence: {
      deep: true,
      handler() {
        this.drawGraph()
      }
    }
  }
}
</script>

<style scoped>
.bayes-container {
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

.el-form-item {
  margin-right: 20px;
}
</style>
