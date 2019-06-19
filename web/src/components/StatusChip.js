import React from 'react'
import Chip from '@material-ui/core/Chip'
import DoneIcon from '@material-ui/icons/Done'
import RollBackIcon from '@material-ui/icons/Cached'
import FailIcon from '@material-ui/icons/Clear'
import { MuiThemeProvider, createMuiTheme } from '@material-ui/core/styles'

const tagTheme = createMuiTheme({
  palette: {
    primary: {
      main: '#4caf50'
    },
    secondary: {
      main: '#f44336'
    },
    default: {
      main: '#9e9e9e'
    }
  }
})

class StatusChip extends React.Component {

  getChip(state) {
    switch (state) {
      case 'deployed':
        return (
          <MuiThemeProvider theme={tagTheme}>
            <Chip
              icon={<DoneIcon />}
              label="Deployed"
              color="primary"
              variant="outlined"
              clickable
            />
          </MuiThemeProvider>
        )
      case 'deploying':
        return (
          <MuiThemeProvider theme={tagTheme}>
            <Chip
              icon={<RollBackIcon />}
              label="Deploying"
              color="default"
              variant="outlined"
              clickable
            />
          </MuiThemeProvider>
        )
      case 'failed':
        return (
          <MuiThemeProvider theme={tagTheme}>
            <Chip
              icon={<FailIcon />}
              label="Failed"
              color="secondary"
              variant="outlined"
              clickable
            />
          </MuiThemeProvider>
        )
      case 'rollback':
        return (
          <MuiThemeProvider theme={tagTheme}>
            <Chip
              icon={<RollBackIcon />}
              label="Rolled Back"
              color="default"
              variant="outlined"
              clickable
            />
          </MuiThemeProvider>
        )
      default:
        return null
    }
  }

  render() {
    return (
      <div>
        {this.getChip(this.props.status)}
      </div>
    )
  }
}

export default StatusChip