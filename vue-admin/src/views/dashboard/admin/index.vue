<template>
  <div class="dashboard-editor-container">
    <github-corner class="github-corner" />

    <h2>Hello World!</h2>

  </div>
</template>

<script>
import GithubCorner from '@/components/GithubCorner'

export default {
  name: 'DashboardAdmin',
  components: {
    GithubCorner
  },
  data() {
  },
  methods: {
  }
}
</script>

<style lang="scss" scoped>
.dashboard-editor-container {
  padding: 32px;
  background-color: rgb(240, 242, 245);
  position: relative;

  .github-corner {
    position: absolute;
    top: 0px;
    border: 0;
    right: 0;
  }

  .chart-wrapper {
    background: #fff;
    padding: 16px 16px 0;
    margin-bottom: 32px;
  }
}

@media (max-width:1024px) {
  .chart-wrapper {
    padding: 8px;
  }
}
</style>
