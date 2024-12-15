<template>
  <div id="app">
    <el-container>
      <el-aside width="200px">
        <el-menu router>
          <el-menu-item index="/prim">Prim算法</el-menu-item>
          <el-menu-item index="/dijkstra">Dijkstra算法</el-menu-item>
          <el-menu-item index="/floyd">Floyd算法</el-menu-item>
          <el-menu-item index="/bayes">Bayes算法</el-menu-item>
          <el-menu-item index="/max-matching">最大匹配</el-menu-item>
          <el-menu-item index="/tsp">旅行商问题</el-menu-item>
        </el-menu>
      </el-aside>
      <el-main>
        <router-view></router-view>
      </el-main>
    </el-container>
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style>
#app {
  font-family: Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  height: 100vh;
}

.el-container {
  height: 100%;
}

.el-aside {
  background-color: #f5f7fa;
}
</style>
