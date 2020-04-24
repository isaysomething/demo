<template>
  <div class="app-container">

    <div class="filter-container">
      <el-input v-model="listQuery.username" :placeholder="$t('user.username')" style="width: 200px;" class="filter-item" @keyup.enter.native="handleFilter" />
      <el-input v-model="listQuery.email" :placeholder="$t('user.email')" style="width: 200px;" class="filter-item" @keyup.enter.native="handleFilter" />
      <el-select v-model="listQuery.state" :placeholder="$t('table.state')" clearable class="filter-item" style="width: 130px">
        <el-option v-for="item in stateOptions" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
      <el-select v-model="listQuery.sort" style="width: 140px" class="filter-item" @change="handleFilter">
        <el-option v-for="item in sortOptions" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
      <el-select v-model="listQuery.direction" style="width: 140px" class="filter-item" @change="handleFilter">
        <el-option v-for="item in directionOptions" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
      <el-button class="filter-item" type="primary" icon="el-icon-search" @click="handleFilter">
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

      <el-table-column align="center" :label="$t('user.username')">
        <template slot-scope="scope">
          <span>{{ scope.row.username }}</span>
        </template>
      </el-table-column>

      <el-table-column width="280px" align="center" :label="$t('user.email')">
        <template slot-scope="scope">
          <span>{{ scope.row.email }}</span>
        </template>
      </el-table-column>

      <el-table-column class-name="state-col" :label="$t('table.state')" width="110">
        <template slot-scope="{row}">
          <el-tag>
            {{ stateText(row.state) }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column width="180px" align="center" :label="$t('table.created_at')">
        <template slot-scope="scope">
          <span>{{ scope.row.created_at | parseTime }}</span>
        </template>
      </el-table-column>

      <el-table-column width="180px" align="center" :label="$t('table.updated_at')">
        <template slot-scope="scope">
          <span>{{ scope.row.created_at | parseTime }}</span>
        </template>
      </el-table-column>

      <el-table-column align="center" :label="$t('table.actions')" width="120">
        <template slot-scope="scope">
          <router-link :to="'/user/edit/'+scope.row.id">
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
import { queryUsers, deleteUser } from '@/api/user'
import Pagination from '@/components/Pagination' // Secondary package based on el-pagination
import UserStates from '@/constants/user-states'

export default {
  name: 'UserList',
  components: { Pagination },
  data() {
    return {
      list: null,
      total: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 10,
        sort: 'created_at',
        direction: 'asc'
      },
      stateOptions: [
        { value: UserStates.ACTIVE, label: this.$t('user.active') },
        { value: UserStates.INACTIVE, label: this.$t('user.inactive') }
      ],
      sortOptions: [
        { value: 'created_at', label: this.$t('table.created_at') },
        { value: 'updated_at', label: this.$t('table.updated_at') }
      ],
      directionOptions: [
        { value: 'asc', label: this.$t('table.asc') },
        { value: 'desc', label: this.$t('table.desc') }
      ]
    }
  },
  created() {
    this.getList()
  },
  methods: {
    getList() {
      this.listLoading = true
      queryUsers(this.listQuery).then(response => {
        this.list = response.data.items
        this.total = response.data.total
        this.listLoading = false
      })
    },
    handleDelete(id) {
      deleteUser(id).then(response => {
        this.getList()
      })
    },
    handleCreate() {
      this.$router.push({
        name: 'CreateUser'
      })
    },
    handleFilter() {
      this.listQuery.page = 1
      this.getList()
    },
    stateText(state) {
      if (state === UserStates.INACTIVE) {
        return this.$t('user.inactive')
      }
      return this.$t('user.active')
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
