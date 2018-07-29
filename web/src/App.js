import React, { Component } from 'react';

class App extends Component {
  constructor() {
    super()


    this.ws = new WebSocket('ws://localhost:1323/ws')
    this.state = {
      playlists: []
    }
  }

  componentDidMount() {
    const { ws } = this

    ws.onmessage = e => {
      this.setState({
        playlists: e.data
      })
    }
  }

  render() {
    const { ws } = this
    return (
      <div className="App">
        <button className="btn btn-primary"
                onClick={() => {
                  ws.send('Set playlists!')
                }}>Click</button>
        {this.state.playlists}
      </div>
    )
  }
}

export default App;
