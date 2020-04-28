<template>
  <div class="app-container">
    <el-table v-loading="listLoading" :data="list" border fit highlight-current-row style="width: 100%">

      <el-table-column align="center" label="ID" width="180">
        <template slot-scope="scope">
          <span>{{ scope.row.id }}</span>
        </template>
      </el-table-column>

      <el-table-column align="center" :label="$t('table.name')">
        <template slot-scope="scope">
          <span>{{ scope.row.name }}</span>
        </template>
      </el-table-column>

      <el-table-column align="center" :label="$t('table.actions')" width="120">
        <template slot-scope="scope">
          <router-link :to="'/role/edit/'+scope.row.id">
            <el-button type="primary" size="small" icon="el-icon-edit" />
          </router-link>
          <el-popconfirm title="Are you sure you want to delete this item?" @onConfirm="handleDelete(scope.row.id)">
            <el-button slot="reference" type="danger" size="small" icon="el-icon-delete" />
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

  </div>
</template>

<script>
import { queryRoles, deleteRole } from '@/api/role'

export default {
  name: 'PostList',
  data() {
    return {
      list: null,
      total: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 10
      }
    }
  },
  created() {
    this.getList()
  },
  methods: {
    getList() {
      this.listLoading = true
      queryRoles(this.listQuery).then(response => {
        this.list = response.data.items
        this.total = response.data.total
        this.listLoading = false
      })
    },
    statusText(status) {
      const statusMap = {
        0: 'Draft',
        1: 'Published'
      }
      return statusMap[status]
    },
    handleDelete(name) {
      deleteRole(name).then(response => {
        this.getList()
      })
    }
  }
}
</script>

<style scoped>
.edit-input {
  padding-right: 100px;
}
.cancel-btn {
  position: absolute;
  right: 15px;
  top: 10px;
}
</style>
