const CLIENT_EVENTS = require('@slack/client').CLIENT_EVENTS
const RTM_EVENTS = require('@slack/client').RTM_EVENTS

var RtmClient = require('@slack/client').RtmClient
var rtm = new RtmClient(process.env.SLACK_BOT_TOKEN)

var botID = ''

var users = {}
var firstMessage = false

rtm.on(CLIENT_EVENTS.RTM.AUTHENTICATED, function (rtmStartData) {
  for (let user of rtmStartData.users) {
    if (user.deleted || user.is_bot) continue
    users[user.id] = user.profile
  }

  console.log(`Logged in as ${rtmStartData.self.name} of team ${rtmStartData.team.name}.`)
  botID = rtmStartData.self.ID
  require('./chat')(rtm, botID)
})

rtm.on(RTM_EVENTS.MESSAGE, (message) => {
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
