const CLIENT_EVENTS = require('@slack/client').CLIENT_EVENTS
const RTM_EVENTS = require('@slack/client').RTM_EVENTS

var WebSocket = require('ws')
var RtmClient = require('@slack/client').RtmClient
var rtm = new RtmClient(process.env.SLACK_BOT_TOKEN)
var botID = ''
var users = {}
var firstMessage = false

RtmClient.prototype.sendThread = function(txt, channel, thread, callback) {
  this.send({
    type: 'message',
    channel: channel,
    text: txt,
    thread_ts: thread
  }, callback)
};

var Chat = {
  channel: 'C50SNCVSR',
  threads: {},
  slackHandler: message => {
    if (message.subtype === 'message_replied') return
    if (message.user === botID) return

    if (this.threads.hasOwnProperty(message.thread_ts)) {
      this.threads[message.thread_ts].send(message.text)
    }
  },
  socketHandler: conn => {
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

      rtm.sendThread(message, Chat.channel, thread, (err, res) => {
        if (err != null) console.log('There was an error!')
      })
    })

    conn.on('close', function open () {
      rtm.sendThread('*My fan just went offline.*', Chat.channel, thread, (err, res) => {
        if (err != null) console.log('There was an error!')
      })

      delete this.threads[thread]
    })
  } 
}

rtm.on(CLIENT_EVENTS.RTM.AUTHENTICATED, function (rtmStartData) {
  for (let user of rtmStartData.users) {
    if (user.deleted || user.is_bot) continue
    users[user.id] = user.profile
  }

  console.log(`Logged in as ${rtmStartData.self.name} of team ${rtmStartData.team.name}.`)
  botID = rtmStartData.self.ID

  var socket = new WebSocket.Server({ port: 80 })
  socket.on('connection', handleSocket)
})

rtm.on(RTM_EVENTS.MESSAGE, (message) => {
  if (message.channel === Chat.channel) Chat.slackHandler(message)

  if (!firstMessage) {
    // INVESTIGAR ISTO
    firstMessage = true
    return
  }

  // Direct messages! Fuck yeah!
  if (message.channel.startsWith('D')) {
    let answer = ''

    switch (message.text.toLowerCase()) {
      case 'ping':
        answer = 'pong'
        break
      case 'tell me a joke':
        answer = "I'm not fully grown up yet! Sorry :sad:"
        break
      default:
        answer = `Sorry ${users[message.user].first_name}, I didn't quite understand what you just said :disappointed:`
    }

    rtm.sendMessage(answer, message.channel)
  }
})

rtm.start()
