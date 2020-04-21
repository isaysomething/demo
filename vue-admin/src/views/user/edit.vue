<template>
  <role-detail :is-edit="true" />
</template>

<script>
import RoleDetail from './components/RoleDetail'

export default {
  name: 'EditRole',
  components: { RoleDetail }
}
</script>

