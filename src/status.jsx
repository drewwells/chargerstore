import React, { Component } from 'react'
import Gauge from 'react-svg-gauge'

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

class Status extends Component {

  constructor(props) {
    super(props);
    this.state = {
      now: new Date(),
			currentMiles: 0,
			estMiles: 0,
    };
  }

  componentDidMount() {
    fetch("/api/v1/status", {})
      .then((r) => {
				return r.json()
			})
      .then((r) => {
        // this.setState(r)

				let charge = r.charge;
				let regained = charge.state.regained_kwh

				regained = Math.floor(regained*1000)/1000
				let soc = r.soc.Data

				let regained_pct = (soc + regained/16.5)

				let current_miles = this.prettyRound(this.milesFromPct(soc), 1)
				let regained_to_miles = this.prettyRound(this.milesFromPct(regained_pct), 1)

				this.setState(Object.assign({}, r, {
					currentMiles: current_miles,
					estMiles: regained_to_miles,
					estPct: regained_pct,
				}))

      })
  }

  componentWillUnmount() {
  }

	milesFromPct(pct) {
		let factor = 35/(0.84313728-0.2)
		return factor * (pct - .2)
	}

	renderGauge() {
		return (
			<div className="center">
				<Gauge value={this.state.currentMiles} max={35} width={160} height={130} label="Current Miles" />
				<Gauge value={this.state.estMiles} max={35} width={160} height={130} label="Regen Est. Miles" />
			</div>
		)
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

		let soc = this.state.soc.Data
		let pct = soc > 83 ? this.prettyRound(soc*100) + '%' : 'Full'
		let renderSOC = []

		// if we're charging, show estimate power level
		if (isCharging) {
			renderSOC.push(
        <p key="2">
			    EST: <span style={socStyle}>{this.state.estMiles} miles</span> ({this.prettyRound(this.state.estPct*100)}%)
			  </p>
			)
			renderSOC.push(
				<div key="3">
					<hr/>
					Charging times based on voltage input
					<p>
						120V: {est120v.toLocaleString()}
					</p>
					<p>
						240V: {est240v.toLocaleString()}
					</p>
				</div>
			)
		}

    return (
      <div>
        <span className="info">{this.state.charge.estimate?'SOC hasn\'t been reported in a while':''}</span>
        {pChrg}
				{renderSOC}
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
		let now = new Date()
		let charge = this.state.charge
		let last_reported_soc = this.prettyRound(-1*(new Date(charge.state.last_reported_soc)-now)/(60*60*1000))

    return (
      <div>
				<hr/>
				<h4>Raw Stats</h4>
				<div>
					<p>
						SOC: {this.prettyRound(this.state.soc.Data*100, 2)}% ({this.state.currentMiles}miles)
					</p>
					<p>
						SOC Last Reported: {last_reported_soc}hrs ago
					</p>
				</div>
				<p>
					Regained: {this.prettyRound(charge.state.regained_kwh)} kwh
				</p>
        <p>
			    Power: {this.state.power.Data} kw
        </p>
        <p>
          Volts: {this.state.volts.Data} volts
        </p>
        <p>
          Amps: {this.state.amps.Data} amps
        </p>
      </div>
    )
  }

  render() {
    return (
      <div>
				{this.renderGauge()}
        {this.renderChargeStatus()}
        {this.renderRawStats()}
      </div>
    )
  }
}

export default Status
