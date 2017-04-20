var http = require('http')
var rtm = null
var users = null

function pingPong (message) {
  let answer = 'pong'

  if (message.text.toLowerCase() === 'pong') {
    answer = 'ping'
  }

  for (let i = 0; i < 4; i++) {
    if (message.text.charAt(i) === message.text.charAt(i).toUpperCase()) {
      answer = answer.replaceAt(i, answer.charAt(i).toUpperCase())
    }
  }

  rtm.sendMessage(answer, message.channel)
}

function makeFunOfUser (message) {
  let fname = users[message.user].first_name
  let lname = users[message.user].last_name

  let options = {
    host: 'api.icndb.com',
    path: '/jokes/random?FirstName=' + fname + '&lastName=' + lname
  }

  let callback = function (response) {
    let str = ''

    response.on('data', function (obj) {
      str += obj
    })

    response.on('end', function () {
      let obj = JSON.parse(str)
      rtm.sendMessage(obj.value.joke, message.channel)
    })
  }

  http.request(options, callback).end()
}

function handler (message) {
  switch (message.text.toLowerCase()) {
    case 'ping':
    case 'pong':
      pingPong(message)
      break
    case 'tell me a joke':
      makeFunOfUser(message)
      break
    default:
      rtm.sendMessage(`Sorry ${users[message.user].first_name}, I didn't quite understand what you just said :disappointed:`, message.channel)
  }
}

module.exports = {
  init: (r, u) => {
    rtm = r
    users = u
  },
  answer: message => {
    handler(message)
  }
}
