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
        console.log(r.charge.state)
        this.setState(r)
      })
  }

  componentWillUnmount() {
  }

	milesFromPct(pct) {
		pct = pct - .2
		let factor = 35/(0.84313728-0.2)
		return factor * pct
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

	prettyRound(num, places) {
		if (!places) {
			places = 3
		}
		let mult = Math.pow(10, places)
		return Math.floor(num*mult)/mult
	}

  renderRawStats() {
    if (!this.state.soc) {
      return (<div/>)
    }
		let charge = this.state.charge;
		let regained = charge.state.regained_kwh
		let now = new Date()
		let last_reported_soc = this.prettyRound(-1*(new Date(charge.state.last_reported_soc)-now)/(60*60*1000))

		regained = Math.floor(regained*1000)/1000
    let soc = this.state.soc.Data

		let regained_pct = (soc + regained/16.5)
		let current_miles = this.prettyRound(this.milesFromPct(soc), 2)
		let regained_to_miles = this.prettyRound(this.milesFromPct(regained_pct), 2)
		console.log(charge)
    return (
      <div>
        <p>
          SOC: <span style={socStyle}>{current_miles}miles</span> ({this.prettyRound(soc*100)}%)
        </p>
        <p>
			    EST: <span style={socStyle}>{regained_to_miles} miles</span> ({this.prettyRound(regained_pct*100)}%)
			  </p>
				<p>
					<h4>Raw Stats</h4>
				</p>
				<p>
					Last SOC: {last_reported_soc}hrs ago
				</p>
				<p>
					Regained: {this.prettyRound(regained)}kwh
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
