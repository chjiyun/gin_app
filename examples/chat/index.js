;(() => {
  // expectingMessage is set to true
  // if the user has just submitted a message
  // and so we should scroll the next message into view when received.
  let expectingMessage = false
  function dial() {
    const protocol = location.protocol.includes('https') ? 'wss' : 'ws';
    const conn = new WebSocket(`${protocol}://${location.host}/subscribe`)
    const heartBeat = heartBeatCheck()

    conn.addEventListener('close', ev => {
      heartBeat.reset();
      appendLog('Connection closed', true)
      console.error(`WebSocket Disconnected code: ${ev.code}, reason: ${ev.reason}`)
      if (ev.code !== 1001) {
        appendLog('try to reconnect...', true)
        setTimeout(dial, 1000)
      }
    })
    conn.addEventListener('open', ev => {
      console.info('websocket connected')
      // heartBeat.start();
      setTimeout(() => {
        conn.send('ping')
      }, 3000)
    })

    // This is where we handle messages received.
    conn.addEventListener('message', ev => {
      // heartBeat.start();
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
    conn.addEventListener('error', ev => {
      heartBeat.reset();
    })
    
    function heartBeatCheck() {
      return {
        timeout: 10 * 1000, // 每10s向服务端发送一次消息
        serverTimeout: 30 * 1000, // 30s收不到服务端消息算超时
        timer: null,
        serverTimer: null,
        reset() { // 心跳检测重置
          clearTimeout(this.timer);
          clearTimeout(this.serverTimer);
          this.timer = null;
          this.serverTimer = null;
          return this;
        },
        start() { // 心跳检测启动
          this.reset();
          this.timer = setTimeout(() => { 
            conn.send('ping'); // 定时向服务端发送消息
            if (!this.serverTimer) {
              this.serverTimer = setTimeout(() => {
                // 关闭连接触发重连
                console.log(new Date().toLocaleString(), "not received pong, close the websocket");
                conn.close();
              }, this.serverTimeout);
            }
          }, this.timeout);
        },
      }
    }
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