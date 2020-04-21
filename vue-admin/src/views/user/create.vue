<template>
  <role-detail :is-edit="false" />
</template>

<script>
import RoleDetail from './components/RoleDetail'

export default {
  name: 'CreateRole',
  components: { RoleDetail }
}
</script>

