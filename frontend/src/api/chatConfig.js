import {mergeConfig} from '@/utils'

export const VERSION = "0.4.1"

export const DEFAULT_CONFIG = {
  guardOnly: false,
  minMedalLevel: 0,
  maxQueueLength: 20,
  admins: '',
  blockUsers: '',
  fuzzyMatch: false,
}

export function setLocalConfig (config) {
  config = mergeConfig(config, DEFAULT_CONFIG)
  window.localStorage.config = JSON.stringify(config)
}

export function getLocalConfig () {
  try {
    return mergeConfig(JSON.parse(window.localStorage.config), DEFAULT_CONFIG)
  } catch {
    return {...DEFAULT_CONFIG}
  }
}
