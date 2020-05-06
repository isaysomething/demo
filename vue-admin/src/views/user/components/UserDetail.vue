<template>
  <div class="createPost-container">
    <el-form ref="form" :model="form" :rules="rules" class="form-container">
      <sticky :z-index="10" class-name="sub-navbar">
        <el-button v-loading="loading" style="margin-left: 10px;" type="success" @click="submitForm">
          Submit
        </el-button>
      </sticky>

      <div class="createPost-main-container">

        <el-form-item prop="username" :label="$t('user.username')">
          <el-input v-model="form.username" />
        </el-form-item>

        <el-form-item prop="email" :label="$t('user.email')">
          <el-input v-model="form.email" type="email" />
        </el-form-item>

        <el-form-item prop="name" :label="$t('user.password')">
          <el-input v-model="form.password" type="password" />
        </el-form-item>

        <el-form-item prop="name" :label="$t('table.state')">
          <el-select
            v-model="form.state"
            reserve-keyword
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
            v-model="form.roles"
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

      </div>
    </el-form>
  </div>
</template>

<script>
import Sticky from '@/components/Sticky'
import { getUser, createUser, updateUser } from '@/api/user'
import { queryRoles } from '@/api/role'
import UserStates from '@/constants/user-states'

const defaultForm = {
  id: undefined,
  username: '',
  email: '',
  password: '',
  state: 1,
  roles: []
}

export default {
  name: 'UserDetail',
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
      form: Object.assign({}, defaultForm),
      loading: false,
      userListOptions: [],
      rules: {
        username: [{ validator: validateRequire }],
        email: [{ validator: validateRequire }],
        roles: [{ validator: validateRequire }]
      },
      tempRoute: {},
      roleOptions: [],
      stateOptions: [
        { key: UserStates.ACTIVE, label: this.$t('user.active') },
        { key: UserStates.INACTIVE, label: this.$t('user.inactive') }
      ]
    }
  },
  computed: {
  },
  created() {
    this.queryRoles()
    if (this.isEdit) {
      this.form.id = this.$route.params && this.$route.params.id
      this.fetchData()
    }

    // Why need to make a copy of this.$route here?
    // Because if you enter this page and quickly switch tag, may be in the execution of the setTagsViewTitle function, this.$route is no longer pointing to the current page
    // https://github.com/PanJiaChen/vue-element-admin/issues/1221
    this.tempRoute = Object.assign({}, this.$route)
  },
  methods: {
    fetchData() {
      getUser(this.form.id).then(response => {
        this.form = response.data
        this.setTagsViewTitle()
        this.setPageTitle()
      }).catch(err => {
        console.log(err)
      })
    },
    setTagsViewTitle() {
      const title = this.$t('user.editUser')
      const route = Object.assign({}, this.tempRoute, { title: `${title}-${this.form.id}` })
      this.$store.dispatch('tagsView/updateVisitedView', route)
    },
    setPageTitle() {
      const title = this.$t('user.editUser')
      document.title = `${title} - ${this.form.id}`
    },
    createUser() {
      this.loading = true
      if (this.form.id) {
        updateUser(this.form.id, this.form).then(response => {
          this.$notify({
            title: this.$t('notify.success'),
            message: this.$t('notify.updated_successfully'),
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
        createUser(this.form).then(response => {
          this.$notify({
            title: this.$t('notify.success'),
            message: this.$t('notify.created_successfully'),
            type: 'success',
            duration: 2000
          })
          this.loading = false
          this.$router.push({
            name: 'EditUser',
            params: { id: response.data.id }
          })
        }).catch(err => {
          console.log(err)
          this.loading = false
        })
      }
    },
    submitForm() {
      this.$refs.form.validate(valid => {
        if (valid) {
          this.createUser()
        } else {
          console.log('error submit!!')
          return false
        }
      })
    },
    queryRoles(query) {
      queryRoles({ name: query }).then(response => {
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
