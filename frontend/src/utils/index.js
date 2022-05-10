export function mergeConfig (config, defaultConfig) {
  let res = {}
  for (let i in defaultConfig) {
    res[i] = i in config ? config[i] : defaultConfig[i]
  }
  return res
}

export function toBool (val) {
  if (typeof val === 'string') {
    return ['false', 'no', 'off', '0', ''].indexOf(val.toLowerCase()) === -1
  }
  return !!val
}

export function toInt (val, _default) {
  let res = parseInt(val)
  if (isNaN(res)) {
    res = _default
  }
  return res
}

export function formatCurrency (price) {
  return new Intl.NumberFormat('zh-CN', {
    minimumFractionDigits: price < 100 ? 2 : 0
  }).format(price)
}

export function getTimeTextHourMin (date) {
  let hour = date.getHours()
  let min = ('00' + date.getMinutes()).slice(-2)
  return `${hour}:${min}`
}

export function getUuid4Hex () {
  let chars = []
  for (let i = 0; i < 32; i++) {
    let char = Math.floor(Math.random() * 16).toString(16)
    chars.push(char)
  }
  return chars.join('')
}

export function parseDanmaku (command) {
  let info = command.info

  let roomId, medalLevel
  if (info[3]) {
    roomId = info[3][3]
    medalLevel = info[3][0]
  } else {
    roomId = medalLevel = 0
  }

  let uid = info[2][0]
  let isAdmin = info[2][2]
  let privilegeType = info[7]
  let authorType
  if (uid === this.roomOwnerUid) {
    authorType = 3
  } else if (isAdmin) {
    authorType = 2
  } else if (privilegeType !== 0) {
    authorType = 1
  } else {
    authorType = 0
  }

  let urank = info[2][5]
  return {
    uid: uid,
    timestamp: info[0][4] / 1000,
    authorName: info[2][1],
    authorType: authorType,
    content: info[1],
    privilegeType: privilegeType,
    isGiftDanmaku: !!info[0][9],
    authorLevel: info[4][0],
    isNewbie: urank < 10000,
    isMobileVerified: !!info[2][6],
    medalLevel: roomId === this.roomId ? medalLevel : 0,
    id: getUuid4Hex(),
    translation: ''
  }
}
