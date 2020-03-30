<template>
  <div class="captcha">
    <input v-model="id" type="hidden" name="captcha_id">
    <el-row :gutter="20">
      <el-col :span="12">
        <img :src="data" @click="reload">
      </el-col>
      <el-col :span="12">
        <el-form-item prop="captcha">
          <el-input v-model="captcha" type="text" name="captcha" @input="updateCaptcha($event)" />
        </el-form-item>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { captcha } from '@/api/captcha'
export default {
  name: 'Captcha',
  data() {
    return {
      id: '',
      captcha: '',
      data: ''
    }
  },
  mounted() {
    this.reload()
  },
  methods: {
    reload: function() {
      captcha().then(response => {
        this.updateCaptcha('')
        this.$emit('updateCaptchaID', response.data.id)
        this.data = response.data.data
      })
    },
    updateCaptcha: function(value) {
      this.captcha = value
      this.$emit('updateCaptcha', value)
    }
  }
}
</script>

<style scoped>
.captcha img {
  background-color: white;
  width: auto;
  max-height: 48px;
}
</style>
