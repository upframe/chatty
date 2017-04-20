function pingPong (text) {
  let answer = 'pong'

  if (text.toLowerCase() === 'pong') {
    answer = 'ping'
  }

  for (let i = 0; i < 4; i++) {
    if (text.charAt(i) === text.charAt(i).toUpperCase()) {
      answer = answer.replaceAt(i, answer.charAt(i).toUpperCase())
    }
  }

  return answer
}

module.exports = (message, users) => {
  let answer = ''

  switch (message.text.toLowerCase()) {
    case 'ping':
    case 'pong':
      answer = pingPong(message.text)
      break
    case 'tell me a joke':
      answer = "I'm not fully grown up yet! Sorry :anguished:"
      break
    default:
      answer = `Sorry ${users[message.user].first_name}, I didn't quite understand what you just said :disappointed:`
  }

  return answer
}
