<template>
  <div>
    <p>
      <el-form
          :model="form"
          ref="form"
          label-width="150px"
          :rules="{
          roomId: [
            {
              required: true,
              message: $t('home.roomIdEmpty'),
              trigger: 'blur'
            },
            {
              type: 'integer',
              min: 1,
              message: $t('home.roomIdInteger'),
              trigger: 'blur'
            }
          ]
        }"
      >
        <el-tabs type="border-card">
          <el-tab-pane :label="$t('home.general')">
            <el-form-item :label="$t('home.roomId')" required prop="roomId">
              <el-input
                  v-model.number="form.roomId"
                  type="number"
                  min="1"
              ></el-input>
            </el-form-item>
            <el-row :gutter="20">
              <el-col :xs="24" :sm="8">
                <el-form-item label="仅允许舰长">
                  <el-switch v-model="form.guardOnly"></el-switch>
                </el-form-item>
              </el-col>
              <el-col :xs="24" :sm="8">
                <el-form-item label="排队关键词模糊匹配">
                  <el-popover
                      placement="top-start"
                      title="模糊匹配"
                      width="200"
                      trigger="hover"
                      content="如果弹幕中包含“排队”的关键词就会触发排队">
                    <el-switch v-model="form.fuzzyMatch" slot="reference"></el-switch>
                  </el-popover>
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :xs="24" :sm="8">
                <el-form-item label="最低牌子等级">
                  <el-popover
                      placement="top-start"
                      title="最低牌子等级"
                      width="200"
                      trigger="hover"
                      content="注意：如果这个参数非0，则未佩戴当前房间粉丝牌的人都无法排队">
                    <el-input
                        v-model.number="form.minMedalLevel"
                        type="number"
                        min="0"
                        slot="reference"
                    ></el-input>
                  </el-popover>
                </el-form-item>
                <el-form-item label="最大排队人数">
                  <el-input
                      v-model.number="form.maxQueueLength"
                      type="number"
                      min="1"
                  ></el-input>
                </el-form-item>
              </el-col>
            </el-row>
          </el-tab-pane>

          <el-tab-pane :label="$t('home.advanced')">
            <el-row :gutter="20">
              <!--              <el-form-item label="管理员UID">-->
              <!--                <el-input-->
              <!--                    v-model="form.admins"-->
              <!--                    type="textarea"-->
              <!--                    :rows="5"-->
              <!--                    :placeholder="$t('home.onePerLine')"-->
              <!--                ></el-input>-->
              <!--              </el-form-item>-->
              <el-form-item label="屏蔽UID">
                <el-input
                    v-model="form.blockUsers"
                    type="textarea"
                    :rows="5"
                    :placeholder="$t('home.onePerLine')"
                ></el-input>
              </el-form-item>
            </el-row>
          </el-tab-pane>
        </el-tabs>
      </el-form>
    </p>

    <p>
      <el-card>
        <el-form :model="form" label-width="150px">
          <el-form-item :label="$t('home.roomUrl')">
            <el-input
                ref="roomUrlInput"
                readonly
                :value="roomUrl"
                style="width: calc(100% - 8em); margin-right: 1em;"
            ></el-input>
            <el-button type="primary" @click="copyUrl"
            >{{ $t('home.copy') }}
            </el-button>
          </el-form-item>
          <el-form-item>
            <el-button
                type="primary"
                :disabled="!roomUrl"
                @click="connectRoom()"
            >
              连接到房间
            </el-button>
            <el-button type="info" :disabled="!roomUrl" @click="saveConfig()">
              保存设置
            </el-button>
            <el-button :disabled="!roomUrl" @click="enterRoom()">
              进入房间
            </el-button>
            <el-button @click="exportConfig"
            >{{ $t('home.exportConfig') }}
            </el-button>
            <el-button @click="importConfig"
            >{{ $t('home.importConfig') }}
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </p>
  </div>
</template>

<script>
import _ from 'lodash'
import download from 'downloadjs'

import {mergeConfig} from '@/utils'
import * as chatConfig from '@/api/chatConfig'
import client from '@/api/client'

export default {
  name: 'Home',
  data() {
    return {
      form: {
        roomId: parseInt(window.localStorage.roomId || '1'),
        ...chatConfig.getLocalConfig(),
      },
    }
  },
  computed: {
    roomUrl() {
      return this.getRoomUrl()
    },
  },
  watch: {
    roomUrl: _.debounce(function () {
      window.localStorage.roomId = this.form.roomId
      chatConfig.setLocalConfig(this.form)
    }, 500),
  },
  created() {
    if (window.localStorage.version !== chatConfig.VERSION) {
      this.$messagebox.alert('检测到您更新了排队姬~ 请记得在 obs 内选中弹幕姬的浏览器源，点击 “刷新” 来更新 obs 内的缓存！', '更新提示', {
        confirmButtonText: '确定',
        callback: () => {
          window.localStorage.version = chatConfig.VERSION
        }
      });
    }

    if (this.form.roomId !== 1) {
      setTimeout(() => {
        this.connectRoom()
      }, 1000)
    }
  },
  methods: {
    connectRoom() {
      chatConfig.setLocalConfig(this.form)
      client.emit('APPLY_RULE', this.form)
      client.emit('CONNECT_DANMAKU', this.form.roomId.toString())
      this.$message({
        message: '连接到房间',
        duration: '1000',
        type: 'success'
      })
    },
    saveConfig() {
      chatConfig.setLocalConfig(this.form)
      client.emit('APPLY_RULE', this.form)
      this.$message({
        message: '保存设置成功',
        duration: '1000',
        type: 'success'
      })
    },
    enterRoom() {
      window.open(
          this.roomUrl,
          `room ${this.form.roomId}`,
          'menubar=0,location=0,scrollbars=0,toolbar=0,width=600,height=600'
      )
    },
    getRoomUrl() {
      let resolved = this.$router.resolve({
        name: 'room',
        params: {roomId: this.form.roomId},
      })
      return `${window.location.protocol}//${window.location.host}${resolved.href}`
    },
    copyUrl() {
      this.$refs.roomUrlInput.select()
      document.execCommand('Copy')
    },
    exportConfig() {
      let cfg = mergeConfig(this.form, chatConfig.DEFAULT_CONFIG)
      download(
          JSON.stringify(cfg, null, 2),
          'blivequeue.json',
          'application/json'
      )
    },
    importConfig() {
      let input = document.createElement('input')
      input.type = 'file'
      input.accept = 'application/json'
      input.onchange = () => {
        let reader = new window.FileReader()
        reader.onload = () => {
          let cfg
          try {
            cfg = JSON.parse(reader.result)
          } catch (e) {
            this.$message.error(this.$t('home.failedToParseConfig') + e)
            return
          }
          cfg = mergeConfig(cfg, chatConfig.DEFAULT_CONFIG)
          this.form = {roomId: this.form.roomId, ...cfg}
        }
        reader.readAsText(input.files[0])
      }
      input.click()
    },
  }
}
</script>
