function checkError(res) {
  if (!res.ok) {
    throw res
  }
  return res
}

function alertMessage(res) {
  res.json().then(json => {
    const { message } = json
    if (message.toLowerCase() === 'token is expired') {
      window.location.reload(true)
      return
    }
    alert(message[0].toUpperCase() + message.substr(1))
  })
}

export { checkError, alertMessage }
