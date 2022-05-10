import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)


const store = new Vuex.Store({
  // 前端存储的 uid 全部是 string ，后端存储的全部是 int
  state: {
    queue: []
  },
  mutations: {
    addUser(state, user) {
      state.queue.push(user)
    },
    addUsers(state, users) {
      for(let user of users) {
        state.queue.push(user)
      }
    },
    removeUser(state, user) {
      for(let i=0; i<state.queue.length; i++) {
        if(state.queue[i].uid === user.uid) {
          state.queue.splice(i,1)
        }
      }
    },
    removeAllUsers(state) {
      state.queue.splice(0, state.queue.length)
    },
    reSort(state, oldIndex, newIndex) {
      state.queue.splice(newIndex, 0, state.queue.splice(oldIndex, 1)[0]);
    }
  }
})

export default store
