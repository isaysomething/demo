<template>
  <user-detail :is-edit="true" />
</template>

<script>
import UserDetail from './components/UserDetail'

export default {
  name: 'EditUser',
  components: { UserDetail }
}
</script>

