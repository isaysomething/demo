<template>
  <div class="createPost-container">
    <el-form ref="form" :model="form" :rules="rules" class="form-container">

      <sticky :z-index="10" class="sub-navbar">
        <el-button v-loading="loading" style="margin-left: 10px;" type="success" @click="submitForm">
          Submit
        </el-button>
      </sticky>

      <div class="createPost-main-container">
        <el-row>

          <el-col>
            <el-form-item style="margin-bottom: 40px;" prop="id" label="ID">
              <el-input v-model="form.id" name="id" required />
            </el-form-item>
          </el-col>

          <el-col>
            <el-form-item style="margin-bottom: 40px;" prop="name" :label="$t('table.name')">
              <el-input v-model="form.name" name="name" required />
            </el-form-item>
          </el-col>

          <el-col>
            <el-select
              v-model="form.roles"
              multiple
              filterable
              remote
              reserve-keyword
              placeholder="Roles"
              :loading="loading"
            >
              <el-option
                v-for="item in roleOptions"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              />
            </el-select>
          </el-col>

          <el-col>

            <el-checkbox-group v-model="form.permissions">
              <el-form-item v-for="group in permissionGroups" :key="group.id" :label="group.name">
                <el-checkbox v-for="permission in group.permissions" :key="permission.id" :label="permission.id">
                  {{ permission.name }}
                </el-checkbox>
              </el-form-item>
            </el-checkbox-group>

          </el-col>

        </el-row>

      </div>
    </el-form>
  </div>
</template>

<script>
import Sticky from '@/components/Sticky'
import { queryRoles, queryRole, createRole, updateRole, queryPermissionGroups } from '@/api/role'

export default {
  name: 'RoleDetail',
  components: { Sticky },
  props: {
    isEdit: {
      type: Boolean,
      default: false
    }
  },
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
      form: {
        id: undefined,
        name: '',
        permissions: ['user:info']
      },
      loading: false,
      userListOptions: [],
      rules: {
        title: [{ validator: validateRequire }],
        content: [{ validator: validateRequire }]
      },
      tempRoute: {},
      permissionGroups: [],
      roleOptions: []
    }
  },
  computed: {
    lang() {
      return this.$store.getters.language
    }
  },
  created() {
    if (this.isEdit) {
      this.form.id = this.$route.params && this.$route.params.id
      this.fetchData()
    }
    this.getPermissions()
    this.queryRoles()

    // Why need to make a copy of this.$route here?
    // Because if you enter this page and quickly switch tag, may be in the execution of the setTagsViewTitle function, this.$route is no longer pointing to the current page
    // https://github.com/PanJiaChen/vue-element-admin/issues/1221
    this.tempRoute = Object.assign({}, this.$route)
  },
  methods: {
    checkPermission(arg1, arg2) {
      console.log(arg1, arg2)
    },
    fetchData() {
      queryRole(this.form.id).then(response => {
        this.form = response.data
        console.log(this.form.permissions)

        // set tagsview title
        this.setTagsViewTitle()

        // set page title
        this.setPageTitle()
      }).catch(err => {
        console.log(err)
      })
    },
    getPermissions() {
      queryPermissionGroups().then(response => {
        this.permissionGroups = response.data
      }).catch(err => {
        console.log(err)
      })
    },
    setTagsViewTitle() {
      const title = this.$t('role.editRole')
      const route = Object.assign({}, this.tempRoute, { title: `${title}-${this.form.id}` })
      this.$store.dispatch('tagsView/updateVisitedView', route)
    },
    setPageTitle() {
      const title = 'Edit Article'
      document.title = `${title} - ${this.form.id}`
    },
    createRole() {
      this.loading = true
      if (this.isEdit) {
        updateRole(this.form.id, this.form).then(response => {
          this.$notify({
            title: '成功',
            message: '更新文章成功',
            type: 'success',
            duration: 2000
          })
          this.loading = false
          this.fetchData()
        }).catch(err => {
          console.log(err)
          this.loading = false
        })
      } else {
        createRole(this.form).then(response => {
          this.$notify({
            title: '成功',
            message: '发布文章成功',
            type: 'success',
            duration: 2000
          })
          this.loading = false
          this.$router.push({
            name: 'EditPost',
            params: { id: response.data.id }
          })
        }).catch(err => {
          console.log(err)
          this.loading = false
        })
      }
    },
    submitForm() {
      console.log(this.form)
      this.$refs.form.validate(valid => {
        if (valid) {
          this.createRole()
        } else {
          console.log('error submit!!')
          return false
        }
      })
    },
    queryRoles() {
      queryRoles({ exclude: [this.form.id] }).then(response => {
        this.roleOptions = response.data.items
      })
    }
  }
}
</script>

<style lang="scss" scoped>
@import "~@/styles/mixin.scss";

.createPost-container {
  position: relative;

  .createPost-main-container {
    padding: 40px 45px 20px 50px;

    .postInfo-container {
      position: relative;
      @include clearfix;
      margin-bottom: 10px;

      .postInfo-container-item {
        float: left;
      }
    }
  }

  .word-counter {
    width: 40px;
    position: absolute;
    right: 10px;
    top: 0px;
  }
}

.article-textarea /deep/ {
  textarea {
    padding-right: 40px;
    resize: none;
    border: none;
    border-radius: 0px;
    border-bottom: 1px solid #bfcbd9;
  }
}
</style>
