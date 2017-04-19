const CLIENT_EVENTS = require('@slack/client').CLIENT_EVENTS;
const RTM_EVENTS = require('@slack/client').RTM_EVENTS;

var RtmClient = require('@slack/client').RtmClient;
var rtm = new RtmClient(process.env.SLACK_BOT_TOKEN);

rtm.on(CLIENT_EVENTS.RTM.AUTHENTICATED, function (rtmStartData) {
  console.log(`Logged in as ${rtmStartData.self.name} of team ${rtmStartData.team.name}.`);
  require('./chat')(rtm)
});

rtm.on(RTM_EVENTS.MESSAGE, function handleRtmMessage(message) {
  /*  console.log('Message:', message);


  if (message.channel != livechatChannel) return;

  if (message.text == "") return

  rtm.sendTyping(message.channel)

  if (message.subtype == "message_replied") {
    rtm.send({
      type: 'message',
      channel: message.channel,
      text: "I'm a fucking thread!",
      thread_ts: message.message.ts,
    })
  }
 */

 //this is no doubt the lamest possible message handler, but you get the idea
  //rtm.sendMessage(message.text, "#testes");

/*   rtm.sendMessage(message.text, message.channel, (err, res) => {
    console.log(err)
    console.log(res) // da para obter a TS aqui para começar a criar threads... OMFG
  }); */

});

rtm.start();
