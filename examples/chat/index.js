;(() => {
  // expectingMessage is set to true
  // if the user has just submitted a message
  // and so we should scroll the next message into view when received.
  let expectingMessage = false
  function dial() {
    const protocol = location.protocol === 'https' ? 'wss' : ws
    const conn = new WebSocket(`${protocol}://${location.host}/subscribe`)

    conn.addEventListener('close', ev => {
      appendLog(`WebSocket Disconnected code: ${ev.code}, reason: ${ev.reason}`, true)
      if (ev.code !== 1001) {
        appendLog('Reconnecting in 1s', true)
        setTimeout(dial, 1000)
      }
    })
    conn.addEventListener('open', ev => {
      console.info('websocket connected')
    })

    // This is where we handle messages received.
    conn.addEventListener('message', ev => {
      if (typeof ev.data !== 'string') {
        console.error('unexpected message type', typeof ev.data)
        return
      }
      let data;
      try {
        data = JSON.parse(ev.data)
      } catch (error) {
        console.log(ev.data)
        return
      }
      console.log(data)
      if (data.msg) {
        const p = appendLog(data.msg)
        if (expectingMessage) {
          p.scrollIntoView()
          expectingMessage = false
        }
      }
    })
  }
  dial()

  const messageLog = document.getElementById('message-log')
  const publishForm = document.getElementById('publish-form')
  const messageInput = document.getElementById('message-input')

  // appendLog appends the passed text to messageLog.
  function appendLog(text, error) {
    const p = document.createElement('p')
    // Adding a timestamp to each message makes the log easier to read.
    p.innerText = `${new Date().toLocaleTimeString()}: ${text}`
    if (error) {
      p.style.color = 'red'
      p.style.fontStyle = 'bold'
    }
    messageLog.append(p)
    return p
  }
  appendLog('Submit a message to get started!')

  // onsubmit publishes the message from the user when the form is submitted.
  publishForm.onsubmit = async ev => {
    ev.preventDefault()

    let msg = messageInput.value
    if (msg === '') {
      return
    }
    messageInput.value = ''

    expectingMessage = true
    try {
      msg = JSON.stringify({ msg })
      const resp = await fetch('/publish', {
        method: 'POST',
        body: msg,
      })
      if (resp.status !== 202) {
        throw new Error(`Unexpected HTTP Status ${resp.status} ${resp.statusText}`)
      }
    } catch (err) {
      appendLog(`Publish failed: ${err.message}`, true)
    }
  }
})()