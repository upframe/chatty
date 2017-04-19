var WebSocket = require('ws')
var channel = 'C50SNCVSR'
var slack = require('@slack/client')

var rtm = null
var botID = ''

var threads = {}

function handleSocket (conn) {
  var thread = ''

  conn.on('message', function incoming (message) {
    if (thread === '') {
      rtm.sendMessage(message, channel, (err, res) => {
        if (err != null) console.log('There was an error!')
        thread = res.ts
        threads[thread] = conn
      })

      return
    }

    rtm.send({
      type: 'message',
      channel: channel,
      text: message,
      thread_ts: thread
    }, (err, res) => {
      if (err != null) console.log('There was an error!')
    })
  })

  conn.on('close', function open () {
    rtm.send({
      type: 'message',
      channel: channel,
      text: '*My fan just went offline.*',
      thread_ts: thread
    }, (err, res) => {
      if (err != null) console.log('There was an error!')
    })
    delete threads[thread]
  })
}

function handleSlack (message) {
  if (message.subtype === 'message_replied') return
  if (message.user === botID) return

  if (threads.hasOwnProperty(message.thread_ts)) {
    threads[message.thread_ts].send(message.text)
  }
}

module.exports = (r, i) => {
  rtm = r
  botID = i

  var socket = new WebSocket.Server({ port: 80 })
  socket.on('connection', handleSocket)
  rtm.on(slack.RTM_EVENTS.MESSAGE, handleSlack)
}
