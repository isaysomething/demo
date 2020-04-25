<template>
  <div class="app-container documentation-container">
    <a class="document-btn" target="_blank" href="https://github.com/clevergo/clevergo">CleverGo</a>
    <a class="document-btn" target="_blank" href="https://github.com/clevergo/demo">Github</a>
    <a class="document-btn" target="_blank" href="https://clevergo.tech/docs">Documentation</a>
    <a class="document-btn" target="_blank" href="https://clevergo.tech/zh/docs">中文文档</a>
    <a class="document-btn" target="_blank" href="https://clevergo.tech/zh-hant/docs">中文文檔</a>
  </div>
</template>

<script>

export default {
  name: 'Documentation',
  data() {
    return {
    }
  }
}
</script>

<style lang="scss" scoped>
.documentation-container {
  margin: 50px;
  display: flex;
  flex-wrap: wrap;
  justify-content: space-evenly;

  .document-btn {
    flex-shrink: 0;
    display: block;
    cursor: pointer;
    background: black;
    color: white;
    height: 60px;
    width: 200px;
    margin-bottom: 16px;
    line-height: 60px;
    font-size: 20px;
    text-align: center;
  }
}
</style>
