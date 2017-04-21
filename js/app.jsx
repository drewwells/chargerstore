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
        <div style={chargeDone}>
          <p>CHARGING DONE {dDone.toLocaleString()}</p>
        </div>
      )
    }
		let now = new Date()
		let minutes = 60*1000
		let est120v = new Date(now.getTime() + c.v120_max.Minutes * minutes)
		let est240v = new Date(now.getTime() + c.v240.Minutes * minutes)

    return (
      <div>
        <span class="info">{this.state.charge.estimate?'Estimate':''}</span>
        <p class="">
          {pChrg}
        </p>
				<p>
					Estimate at 120V: {est120v.toLocaleString()}
				</p>
				<p>
					Estimate at 240V: {est240v.toLocaleString()}
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
        <p style={socStyle}>
          SOC: {soc}
        </p>
        <p>
			    Power: {this.state.power.Data}
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

const socStyle = {
	fontSize: 18,
	fontWeight: 'bold',
	color: 'purple',
}

const chargeDone = {
	borderStyle: 'solid',
	borderWidth: 2,
	borderColor: 'red',
	color: 'green',
	padding: 5,
	textAlign: 'center',
}

function App() {
  return (
    <div>
			<Status />
    </div>
  );
}

ReactDOM.render(<App />, document.getElementById('root'));
