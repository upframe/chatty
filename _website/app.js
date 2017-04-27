'use strict'

document.addEventListener('DOMContentLoaded', () => {
  let conn = new window.WebSocket('ws://chatty.upframe.co/api/websocket/2de90faeba5091a93466cd5f5850085e')
  let messages = document.getElementById('messages')

  conn.onopen = function () {
    window.alert('You may start by clicking on the button in the page.')
  }

  conn.onmessage = function (event) {
    let newParaOAgrafo = '<p><strong>Them:</strong><br> ' + event.data + '</p>'
    messages.insertAdjacentHTML('beforeend', newParaOAgrafo)
  }

  conn.onclose = function (event) {
    window.alert('Connection was closed!')
  }

  document.getElementById('sender').addEventListener('click', () => {
    let message = window.prompt('What do you wanna send?')
    if (message === '') {
      window.alert('You are being too vague. Please fucking write something!')
      return
    }

    let newParaOAgrafo = '<p><strong>Me:</strong><br> ' + message + '</p>'
    messages.insertAdjacentHTML('beforeend', newParaOAgrafo)
    conn.send(message)
  })
})
