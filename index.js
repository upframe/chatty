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

String.prototype.replaceAt = (index, replacement) => {
  return this.substr(0, index) + replacement + this.substr(index + replacement.length)
}

var rtm = new RtmClient(process.env.SLACK_BOT_TOKEN)
var bot = null
var users = {}
var firstMessage = false

rtm.on(CLIENT_EVENTS.RTM.AUTHENTICATED, function (rtmStartData) {
  for (let user of rtmStartData.users) {
    if (user.deleted || user.is_bot) continue
    users[user.id] = user.profile
  }

  console.log(`Logged in as ${rtmStartData.self.name} of team ${rtmStartData.team.name}.`)
  bot = rtmStartData.self.ID
  require('./chat')(rtm, bot)
  dm.init(rtm, users)
})

rtm.on(RTM_EVENTS.MESSAGE, (message) => {
  if (message.user === bot) return

  if (!firstMessage) {
    // TODO: INVESTIGAR ISTO
    firstMessage = true
    return
  }

  // Direct messages! Fuck yeah!
  if (message.channel.startsWith('D')) {
    dm.answer(message)
  }
})

rtm.start()
