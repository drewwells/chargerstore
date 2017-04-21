class Status extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      now: new Date()
    };
  }

  componentDidMount() {
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
