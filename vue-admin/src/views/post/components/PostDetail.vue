<template>
  <div class="createPost-container">
    <el-form ref="postForm" :model="postForm" :rules="rules" class="form-container">
      <sticky :z-index="10" :class-name="'sub-navbar '+postForm.state">
        <el-button v-loading="loading" style="margin-left: 10px;" type="success" @click="submitForm">
          Publish
        </el-button>
        <el-button v-loading="loading" type="warning" @click="draftForm">
          Draft
        </el-button>
      </sticky>

      <div class="createPost-main-container">
        <el-row>

          <el-col :span="24">
            <el-form-item style="margin-bottom: 40px;" prop="title">
              <MDinput v-model="postForm.title" :maxlength="100" name="name" required>
                Title
              </MDinput>
            </el-form-item>

          </el-col>
        </el-row>

        <el-form-item prop="markdown_content" style="margin-bottom: 30px;">
          <markdown-editor
            ref="markdownEditor"
            v-model="postForm.markdown_content"
            :language="language"
            height="300px"
            :options="{hideModeSwitch:true,previewStyle:'tab'}"
          />
        </el-form-item>
      </div>
    </el-form>
  </div>
</template>

<script>
import MDinput from '@/components/MDinput'
import Sticky from '@/components/Sticky'
import MarkdownEditor from '@/components/MarkdownEditor'
import { fetchPost, createPost, updatePost } from '@/api/post'

const defaultForm = {
  id: undefined,
  title: '',
  content: '',
  markdown_content: '',
  state: 0
}

export default {
  name: 'PostDetail',
  components: { MDinput, Sticky, MarkdownEditor },
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
      postForm: Object.assign({}, defaultForm),
      loading: false,
      userListOptions: [],
      rules: {
        title: [{ validator: validateRequire }],
        markdown_content: [{ validator: validateRequire }]
      },
      tempRoute: {},
      languageTypeList: {
        'en': 'en_US',
        'zh': 'zh_CN'
      }
    }
  },
  computed: {
    lang() {
      return this.$store.getters.language
    },
    language() {
      return this.languageTypeList[this.$store.getters.language]
    }
  },
  created() {
    if (this.isEdit) {
      this.postForm.id = this.$route.params && this.$route.params.id
      this.fetchData()
    }

    // Why need to make a copy of this.$route here?
    // Because if you enter this page and quickly switch tag, may be in the execution of the setTagsViewTitle function, this.$route is no longer pointing to the current page
    // https://github.com/PanJiaChen/vue-element-admin/issues/1221
    this.tempRoute = Object.assign({}, this.$route)
  },
  methods: {
    fetchData() {
      fetchPost(this.postForm.id).then(response => {
        this.postForm = response.data

        // set tagsview title
        this.setTagsViewTitle()

        // set page title
        this.setPageTitle()
      }).catch(err => {
        console.log(err)
      })
    },
    setTagsViewTitle() {
      const title = this.lang === 'zh' ? '编辑文章' : 'Edit Article'
      const route = Object.assign({}, this.tempRoute, { title: `${title}-${this.postForm.id}` })
      this.$store.dispatch('tagsView/updateVisitedView', route)
    },
    setPageTitle() {
      const title = 'Edit Article'
      document.title = `${title} - ${this.postForm.id}`
    },
    createPost() {
      this.postForm.content = this.$refs.markdownEditor.getHtml()
      this.loading = true
      if (this.postForm.id) {
        updatePost(this.postForm.id, this.postForm).then(response => {
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
        createPost(this.postForm).then(response => {
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
      console.log(this.postForm)
      this.$refs.postForm.validate(valid => {
        if (valid) {
          this.postForm.state = 2
          this.createPost()
        } else {
          console.log('error submit!!')
          return false
        }
      })
    },
    draftForm() {
      if (this.postForm.content.length === 0 || this.postForm.title.length === 0) {
        this.$message({
          message: '请填写必要的标题和内容',
          type: 'warning'
        })
        return
      }
      this.postForm.state = 1
      this.createPost()
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
