var WebSocket = require('ws')
var channel = 'C50SNCVSR'
var slack = require('@slack/client')

var rtm = null
var bot = null
var threads = {}

function errorHandler (err, res) {
  if (err != null) console.log('Error on chat.js: ', err)
}

function handleSocket (conn) {
  let thread = null

  conn.on('message', function incoming (message) {
    if (thread != null) {
      rtm.sendThread(message, channel, thread, errorHandler)
      return
    }

    rtm.sendMessage(message, channel, (err, res) => {
      if (err != null) console.log('There was an error!')
      thread = res.ts
      threads[thread] = conn
    })

    return
  })

  conn.on('close', function open () {
    rtm.sendThread('*My fan just went offline.*', channel, thread, errorHandler)
    delete threads[thread]
  })
}

function handleSlack (message) {
  if (message.subtype === 'message_replied') return
  if (message.user === bot) return

  if (threads.hasOwnProperty(message.thread_ts)) {
    threads[message.thread_ts].send(message.text)
  }
}

module.exports = (r, i) => {
  rtm = r
  bot = i

  var socket = new WebSocket.Server({ port: 80 })
  socket.on('connection', handleSocket)
  rtm.on(slack.RTM_EVENTS.MESSAGE, handleSlack)
}
