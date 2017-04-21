const CLIENT_EVENTS = require('@slack/client').CLIENT_EVENTS
const RTM_EVENTS = require('@slack/client').RTM_EVENTS
var RtmClient = require('@slack/client').RtmClient
var dm = require('./dm')

RtmClient.prototype.sendThread = function (txt, channel, thread, callback) {
  this.send({
    type: 'message',
    channel: channel,
    text: txt,
    thread_ts: thread
  }, callback)
}

var rtm = new RtmClient(process.env.SLACK_BOT_TOKEN)
var bot = null
var users = {}

rtm.on(CLIENT_EVENTS.RTM.AUTHENTICATED, function (rtmStartData) {
  for (let user of rtmStartData.users) {
    if (user.deleted || user.is_bot) continue
    users[user.id] = user.profile
  }

  bot = rtmStartData.self.id
  require('./chat')(rtm, bot)
  dm.init(rtm, users)
  console.log(`Logged in as ${rtmStartData.self.name} of team ${rtmStartData.team.name}.`)
})

rtm.on(RTM_EVENTS.MESSAGE, (message) => {
  if (message.user === bot) return
  if (message.hasOwnProperty('subtype')) return

  // Direct messages! Fuck yeah!
  if (message.channel.startsWith('D')) {
    dm.answer(message)
  }

  if (message.text.includes(`<@${bot}>`)) {
    rtm.sendMessage(
        `If you mention me again, I will break your nose <@${message.user}>! DM me if you wanna chat.`,
        message.channel
    )

    return
  }
})

rtm.start()
