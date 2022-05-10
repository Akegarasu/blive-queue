<template>
  <chat-renderer ref="renderer"></chat-renderer>
</template>

<script>
import client from '@/api/client'
import * as chatConfig from '@/api/chatConfig'
import ChatRenderer from '@/components/ChatRenderer'

export default {
  name: 'Room',
  components: {
    ChatRenderer
  },
  props: {
    roomId: {
      type: Number,
      default: null
    },
    strConfig: {
      type: Object,
      default: () => ({})
    }
  },
  data () {
    return {
      config: { ...chatConfig.DEFAULT_CONFIG },
    }
  },
  computed: {},
  mounted () {
    // 提示用户已加载
    this.$message({
      message: 'Loaded',
      duration: '500'
    })
    window.setTimeout(() => {
      client.syncData()
    }, 1500)
  },
  methods: {}
}
</script>
