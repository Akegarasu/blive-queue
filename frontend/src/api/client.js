import store from '@/store/index'
import axios from 'axios'

const HEARTBEAT_INTERVAL = 10 * 1000
const RECEIVE_TIMEOUT = HEARTBEAT_INTERVAL + 5 * 1000

class Client {
  constructor () {
    this.websocket = null
    this.connected = false
    this.waitSending = []
    this.retryCount = 0
    this.isDestroying = false
    this.heartbeatTimerId = null
    this.receiveTimeoutTimerId = null
  }

  start () {
    this.wsConnect()
  }

  stop () {
    this.isDestroying = true
    if (this.websocket) {
      this.websocket.close()
    }
  }

  syncData () {
    axios.get('/api/sync').then(
      (ret) => {
        if (ret.status === 200 && ret.data.cmd === 'SYNC') {
          console.log(ret.data.data)
          store.commit('removeAllUsers')
          store.commit('addUsers', ret.data.data)
        }
      }
    )
  }

  wsConnect () {
    if (this.isDestroying) {
      return
    }
    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
    const host = window.location.host
    const url = `${protocol}://${host}/eio`
    console.log(url)
    this.websocket = new WebSocket(url)
    this.websocket.onopen = this.onWsOpen.bind(this)
    this.websocket.onclose = this.onWsClose.bind(this)
    this.websocket.onmessage = this.onWsMessage.bind(this)
  }

  onWsOpen () {
    this.connected = true
    this.retryCount = 0
    this.heartbeatTimerId = window.setInterval(this.sendHeartbeat.bind(this), HEARTBEAT_INTERVAL)
    this.refreshReceiveTimeoutTimer()
  }

  emit (cmd, data) {
    this.send({
      'cmd': cmd,
      'data': data
    })
  }

  send (packet) {
    if (!this.connected) {
      this.waitSending.push(packet)
    } else {
      // 处理未发出的包
      if (this.waitSending.length !== 0) {
        for (let i = 0; i < this.waitSending.length; i++) {
          this.websocket.send(this.waitSending.pop())
        }
      }
      this.websocket.send(JSON.stringify(packet))
    }
  }

  sendHeartbeat () {
    this.emit('HEARTBEAT')
  }

  refreshReceiveTimeoutTimer () {
    if (this.receiveTimeoutTimerId) {
      window.clearTimeout(this.receiveTimeoutTimerId)
    }
    this.receiveTimeoutTimerId = window.setTimeout(this.onReceiveTimeout.bind(this), RECEIVE_TIMEOUT)
  }

  onReceiveTimeout () {
    window.console.warn('接收消息超时')
    this.receiveTimeoutTimerId = null

    // 直接丢弃阻塞的websocket，不等onclose回调了
    this.websocket.onopen = this.websocket.onclose = this.websocket.onmessage = null
    this.websocket.close()
    this.onWsClose()
  }

  onWsClose () {
    this.connected = false
    this.websocket = null
    if (this.heartbeatTimerId) {
      window.clearInterval(this.heartbeatTimerId)
      this.heartbeatTimerId = null
    }
    if (this.receiveTimeoutTimerId) {
      window.clearTimeout(this.receiveTimeoutTimerId)
      this.receiveTimeoutTimerId = null
    }

    if (this.isDestroying) {
      return
    }
    window.console.warn(`掉线重连中${++this.retryCount}`)
    window.setTimeout(this.wsConnect.bind(this), 1000)
  }

  onWsMessage (event) {
    this.refreshReceiveTimeoutTimer()
    let { cmd, data } = JSON.parse(event.data)
    if (cmd !== 'HEARTBEAT') {
      console.log(cmd, data)
    }
    switch (cmd) {
      case 'HEARTBEAT': {
        break
      }
      case 'REMOVE_USER': {
        store.commit('removeUser', {
          uid: data
        })
        break
      }
      case 'ADD_USER': {
        let d = JSON.parse(data)
        store.commit('addUser', {
          nickname: d.info[2][1],
          uid: d.info[2][0].toString(),
          level: d.info[7].toString()
        })
        break
      }
      case 'REMOVE_ALL': {
        store.commit('removeAllUsers')
        break
      }
      case 'RESORT': {
        this.syncData()
        break
      }
    }
  }
}

const client = new Client()
client.start()
export default client
