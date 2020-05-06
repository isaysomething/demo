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
          <el-button type="primary" size="mini" icon="el-icon-edit" @click="handleUpdate(scope.row)" />
          <el-popconfirm title="Are you sure you want to delete this item?" @onConfirm="handleDelete(scope.row.id)">
            <el-button slot="reference" type="danger" size="mini" icon="el-icon-delete" />
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />

    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
      <el-form ref="dataForm" :rules="rules" :model="temp" label-position="left" label-width="70px" style="width: 400px; margin-left:50px;">
        <el-form-item prop="username" :label="$t('user.username')">
          <el-input v-model="temp.username" />
        </el-form-item>
        <el-form-item prop="email" :label="$t('user.email')">
          <el-input v-model="temp.email" type="email" />
        </el-form-item>
        <el-form-item prop="name" :label="$t('user.password')">
          <el-input v-model="temp.password" type="password" />
        </el-form-item>
        <el-form-item prop="name" :label="$t('table.state')">
          <el-select
            v-model="temp.state"
            placeholder="$t('table.state')"
          >
            <el-option
              v-for="item in stateOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item prop="name" :label="$t('user.roles')">
          <el-select
            v-model="temp.roles"
            multiple
            filterable
            remote
            reserve-keyword
            placeholder="角色"
            :remote-method="queryRoles"
            :loading="loading"
          >
            <el-option
              v-for="item in roleOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">
          {{ $t('button.cancel') }}
        </el-button>
        <el-button type="primary" @click="dialogStatus==='create'?createData():updateData()">
          {{ $t('button.confirm') }}
        </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { queryUsers, deleteUser, createUser, updateUser } from '@/api/user'
import { queryRoles } from '@/api/role'
import Pagination from '@/components/Pagination' // Secondary package based on el-pagination
import UserStates from '@/constants/user-states'

export default {
  name: 'UserList',
  components: { Pagination },
  data() {
    const validateRequire = (rule, value, callback) => {
      if (value === '') {
        this.$message({
          message: rule.field + '为必传项',
          type: 'error'
        })
        callback(new Error(rule.field + '为必传项'))
      } else {
        callback()
      }
    }
    return {
      list: null,
      total: 0,
      listLoading: true,
      loading: false,
      rules: {
        username: [{ validator: validateRequire }],
        email: [{ validator: validateRequire }],
        state: [{ validator: validateRequire }],
        roles: [{ validator: validateRequire }]
      },
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
      ],
      roleOptions: [],
      temp: {
        id: undefined,
        username: '',
        email: '',
        password: '',
        state: UserStates.ACTIVE,
        roles: []
      },
      dialogFormVisible: false,
      dialogStatus: '',
      textMap: {
        update: 'Edit',
        create: 'Create'
      }
    }
  },
  created() {
    this.getList()
    this.queryRoles()
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
    queryRoles(query) {
      queryRoles({ name: query }).then(response => {
        this.roleOptions = response.data.items
      })
    },
    handleDelete(id) {
      deleteUser(id).then(response => {
        this.getList()
      })
    },
    resetTemp() {
      this.temp = {
        id: undefined,
        username: '',
        email: '',
        password: '',
        state: undefined,
        roles: []
      }
    },
    handleCreate() {
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    createData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          createUser(this.temp).then((response) => {
            this.temp = response.data
            this.list.unshift(this.temp)
            this.dialogFormVisible = false
            this.$notify({
              title: this.$t('notify.success'),
              message: this.$t('notify.created_successfully'),
              type: 'success',
              duration: 2000
            })
          })
        }
      })
    },
    handleUpdate(row) {
      this.temp = Object.assign({}, row) // copy obj
      this.temp.timestamp = new Date(this.temp.timestamp)
      this.dialogStatus = 'update'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    updateData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          const tempData = Object.assign({}, this.temp)
          tempData.timestamp = +new Date(tempData.timestamp) // change Thu Nov 30 2017 16:41:05 GMT+0800 (CST) to 1512031311464
          updateUser(tempData.id, tempData).then(() => {
            const index = this.list.findIndex(v => v.id === this.temp.id)
            this.list.splice(index, 1, this.temp)
            this.dialogFormVisible = false
            this.$notify({
              title: 'Success',
              message: 'Update Successfully',
              type: 'success',
              duration: 2000
            })
          })
        }
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
