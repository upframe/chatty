const WebSocket = require('ws');
const channel = "C50SNCVSR";
const Slack = require('@slack/client');

var rtm = null;

function handleSocket(conn) {
  var randId = 'rand' + new Date().getTime();
  conn.on('message', function incoming(message) {
    rtm.sendMessage(message, channel);
  });
}

function handleSlack(message) {
  if (message.channel != channel) return;
  console.log(message)
}

module.exports = (r) => {
  rtm = r;

  var socket = new WebSocket.Server({ port: 80 });
  socket.on('connection', handleSocket);
  rtm.on(Slack.RTM_EVENTS.MESSAGE, handleSlack);
};
