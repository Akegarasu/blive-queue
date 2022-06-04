<template>
  <div>
    <el-card style="margin: 5px">
      <el-alert
          title="温馨提示: 按住一行可以上下拖动排序~"
          type="success"
          show-icon
          style="margin-bottom: 10px">
      </el-alert>
      <el-button type="danger" icon="el-icon-delete" @click="removeAll">清除全部</el-button>
      <el-button type="primary" icon="el-icon-refresh" @click="syncData">同步</el-button>
      <el-divider direction="vertical"></el-divider>
      <el-button type="warning" icon="el-icon-video-pause" @click="pauseQueue">暂停排队</el-button>
      <el-button type="success" icon="el-icon-video-play" @click="continueQueue">继续排队</el-button>
    </el-card>
    <el-table-draggable>
      <el-table
          :data="tableData"
          border
          style="width: 100%; margin: 5px; border-radius: 5px;-webkit-box-shadow: 0 2px 12px 0 rgba(0,0,0,.1); box-shadow: 0 2px 12px 0 rgba(0,0,0,.1);">
        <el-table-column
            prop="nickname"
            label="序号"
            type="index"
            width="100">
        </el-table-column>
        <el-table-column
            prop="nickname"
            label="昵称">
        </el-table-column>
        <el-table-column
            prop="uid"
            label="UID">
        </el-table-column>
        <el-table-column label="详细信息">
          <template slot-scope="scope">
            <el-tag v-if="scope.row.level === '0'" type=""> 观众</el-tag>
            <el-tag v-else-if="scope.row.level === '3'" type="success"> 舰长</el-tag>
            <el-tag v-else-if="scope.row.level === '2'" type="warning"> 提督</el-tag>
            <el-tag v-else-if="scope.row.level === '1'" type="danger"> 总督</el-tag>
          </template>
        </el-table-column>
        <el-table-column fixed="right" label="操作">
          <template slot-scope="scope">
            <el-button @click="removeUser(scope.row)" type="danger" size="small">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-table-draggable>
  </div>
</template>

<script>
import client from '@/api/client'
import ElTableDraggable from '@/components/Draggable/SortableElTable'

export default {
  name: 'Control',
  components: {
    ElTableDraggable,
  },
  data() {
    return {
      tableData: this.$store.state.queue,
      testNow: 1
    }
  },
  methods: {
    removeUser(row) {
      client.emit('REMOVE_USER', row.uid)
    },
    removeAll() {
      client.emit('REMOVE_ALL')
      this.$message({
        message: '清除全部排队',
        duration: '1000'
      })
    },
    syncData() {
      client.syncData()
      this.$message({
        message: '同步中~',
        duration: '1000'
      })
    },
    pauseQueue() {
      client.emit('PAUSE')
      this.$message({
        message: '已暂停',
        duration: '1000'
      })
    },
    continueQueue() {
      client.emit('CONTINUE')
      this.$message({
        message: '已继续排队~',
        duration: '1000'
      })
    }
  },
  computed: {},
  mounted() {
    window.setTimeout(() => {
      client.syncData()
    }, 100)
  }
}
</script>
