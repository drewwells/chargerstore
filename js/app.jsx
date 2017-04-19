class Status extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      now: new Date()
    };
  }

  componentDidMount() {
    this.setState({
  "amps": {
    "Data": 0,
    "PublishTime": "2017-04-19T13:25:07.858Z"
  },
  "done": "2017-04-18T22:38:25.906Z",
  "charge": {
    "state": {
      "last_reported_soc": "0001-01-01T00:00:00Z",
      "percent": 0,
      "deficit_kwh": 0,
      "regained_kwh": 0
    },
    "deficit": 0,
    "estimate": true,
    "published_at": "0001-01-01T00:00:00Z",
    "current": {
      "Duration": 1000*25*60*1000000,
      "Minutes": 25
    },
    "v120_standard": {
      "Duration": 0,
      "Minutes": 0
    },
    "v120_max": {
      "Duration": 0,
      "Minutes": 0
    },
    "v240": {
      "Duration": 0,
      "Minutes": 0
    }
  },
  "power": {
    "Data": 0.019,
    "PublishTime": "2017-04-18T22:38:25.906Z"
  },
  "soc": {
    "Data": 0.83529411,
    "PublishTime": "2017-04-19T10:44:29.439Z"
  },
  "volts": {
    "Data": 0,
    "PublishTime": "2017-04-19T13:25:07.858Z"
  }
    })
    return
    fetch("/api/v1/status", {})
      .then((r) => r.json())
      .then((r) => {
        console.log(r)
        this.setState(r)
      })

  }

  componentWillUnmount() {
  }

  renderChargeStatus() {
    let c = this.state.charge
    if (!c) {
      return (<div/>)
    }
    let isCharging = c.current.Duration > 0

    let pChrg = (<div />)
    if (isCharging) {
      let dDone = new Date(this.state.done)
      pChrg = (
        <div>
          <p>Finish at {dDone.toLocaleString()}</p>
        </div>
      )
    }

    return (
      <div>
        <span class="info">{this.state.charge.estimate?'Estimate':''}</span>
        <p class="">
          {pChrg}
        </p>
      </div>
    )
  }

  renderRawStats() {
    if (!this.state.soc) {
      return (<div/>)
    }
    let soc = this.state.soc.Data
    if (soc > 0.83) {
      soc = 'Full'
    }
    return (
      <div>
        <p>
          SOC: {soc}
        </p>
        <p>
          Volts: {this.state.volts.Data}
        </p>
        <p>
          Amps: {this.state.amps.Data}
        </p>
      </div>
    )
  }

  render() {
    return (
      <div>
        {this.renderChargeStatus()}
        {this.renderRawStats()}
      </div>
    )
  }
}

function App() {
  return (
    <div>
      <Status />
    </div>
  );
}

ReactDOM.render(<App />, document.getElementById('root'));
