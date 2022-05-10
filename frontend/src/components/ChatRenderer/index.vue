<template>
  <yt-live-chat-renderer class="style-scope yt-live-chat-app" style="--scrollbar-width:11px;" hide-timestamps>
    <yt-live-chat-item-list-renderer class="style-scope yt-live-chat-renderer" allow-scroll>
      <div ref="scroller" id="item-scroller" class="style-scope yt-live-chat-item-list-renderer animated">
        <div ref="itemOffset" id="item-offset" class="style-scope yt-live-chat-item-list-renderer" style="height: 0px;">
          <div ref="items" id="items" class="style-scope yt-live-chat-item-list-renderer" style="overflow: hidden">
            <template v-for="(w, pos) in waitingList">
              <text-message :key="w.uid"
                            class="style-scope yt-live-chat-item-list-renderer"
                            :authorName="w.nickname"
                            :content="`${pos+1}. ${w.nickname}`"
              ></text-message>
            </template>
          </div>
        </div>
      </div>
    </yt-live-chat-item-list-renderer>
  </yt-live-chat-renderer>
</template>

<script>
import TextMessage from './TextMessage.vue'

export default {
  name: 'ChatRenderer',
  components: {
    TextMessage,
  },
  data () {
    return {
      waitingList: this.$store.state.queue,
    }
  },
  computed: {},
  watch: {
    // eslint-disable-next-line no-unused-vars
    waitingList (val) {
      (async () => {
        // 需要等待高度变化
        await this.$nextTick()
        console.log('clientHeight' + this.$refs.items.clientHeight)
        this.$refs.itemOffset.style.height = `${this.$refs.items.clientHeight}px`
        this.scrollToBottom()
      })()
    }
  },
  mounted () {
    this.scrollToBottom()
  },
  beforeDestroy () {},
  methods: {
    scrollToBottom () {
      this.$refs.scroller.scrollTop = Math.pow(2, 24)
    },
  }
}
</script>

<style src="@/assets/css/youtube/yt-html.css"></style>
<style src="@/assets/css/youtube/yt-live-chat-renderer.css"></style>
<style src="@/assets/css/youtube/yt-live-chat-item-list-renderer.css"></style>
