<template>
  <div class="app-container">

    <div class="filter-container">
      <el-input v-model="listQuery.title" :placeholder="$t('table.title')" style="width: 200px;" class="filter-item" @keyup.enter.native="handleFilter" />
      <el-select v-model="listQuery.state" :placeholder="$t('table.state')" clearable class="filter-item" style="width: 130px">
        <el-option v-for="(value, key) in states" :key="key" :label="value" :value="key" />
      </el-select>
      <el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="handleFilter">
        {{ $t('table.search') }}
      </el-button>
      <el-button class="filter-item" style="margin-left: 10px;" type="primary" icon="el-icon-edit" @click="handleCreate">
        {{ $t('table.add') }}
      </el-button>
    </div>

    <el-table v-loading="listLoading" :data="list" border fit highlight-current-row style="width: 100%">
      <el-table-column align="center" label="ID" width="80">
        <template slot-scope="scope">
          <span>{{ scope.row.id }}</span>
        </template>
      </el-table-column>

      <el-table-column align="center" label="Title">
        <template slot-scope="scope">
          <span>{{ scope.row.title }}</span>
        </template>
      </el-table-column>

      <el-table-column class-name="state-col" label="State" width="110">
        <template slot-scope="{row}">
          <el-tag>
            {{ stateText(row.state) }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column width="180px" align="center" label="Created At">
        <template slot-scope="scope">
          <span>{{ scope.row.created_at }}</span>
        </template>
      </el-table-column>

      <el-table-column align="center" label="Actions" width="120">
        <template slot-scope="scope">
          <router-link :to="'/post/edit/'+scope.row.id">
            <el-button type="primary" size="small" icon="el-icon-edit" />
          </router-link>
          <el-popconfirm title="Are you sure you want to delete this item?" @onConfirm="handleDelete(scope.row.id)">
            <el-button slot="reference" type="danger" size="small" icon="el-icon-delete" />
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />
  </div>
</template>

<script>
import { fetchList, deletePost } from '@/api/post'
import Pagination from '@/components/Pagination' // Secondary package based on el-pagination

export default {
  name: 'PostList',
  components: { Pagination },
  data() {
    return {
      list: null,
      total: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 10
      },
      states: {
        'draft': 'Draft',
        'published': 'Published'
      }
    }
  },
  created() {
    this.getList()
  },
  methods: {
    getList() {
      this.listLoading = true
      fetchList(this.listQuery).then(response => {
        this.list = response.data.items
        this.total = response.data.total
        this.listLoading = false
      })
    },
    stateText(state) {
      const stateMap = {
        0: 'Draft',
        1: 'Published'
      }
      return stateMap[state]
    },
    handleDelete(id) {
      deletePost(id).then(response => {
        this.getList()
      })
    },
    handleFilter() {
      this.listQuery.page = 1
      this.getList()
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
