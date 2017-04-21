var http = require('http')
var rtm = null
var users = null

function pingPong (message) {
  let answer = message.text
  let replacement = null

  switch (answer.charAt(1)) {
    case 'I':
      replacement = 'O'
      break
    case 'i':
      replacement = 'o'
      break
    case 'O':
      replacement = 'I'
      break
    case 'o':
      replacement = 'i'
      break
    default:
      // Wut, what? This must not happen!
  }

  answer = answer.substr(0, 1) + replacement + answer.substr(2)
  rtm.sendMessage(answer, message.channel)
}

function makeFunOfUser (message) {
  let fname = encodeURIComponent(users[message.user].first_name)
  let lname = encodeURIComponent(users[message.user].last_name)

  let options = {
    host: 'api.icndb.com',
    path: `/jokes/random?escape=javascript&firstName=${fname}&lastName=${lname}`
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

function answer (message) {
  switch (message.text.toLowerCase()) {
    case 'hey':
    case 'hi':
    case 'Hello':
      rtm.sendMessage('Hi there!', message.channel)
      break
    case 'ping':
    case 'pong':
      pingPong(message)
      break
    case 'tell me a joke':
      makeFunOfUser(message)
      break
    case 'bye':
    case 'goodbye':
    case 'cya':
      rtm.sendMessage('Bye! Gonna miss you :kissing:', message.channel)
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
  answer: answer
}
