import React from 'react'
import Chip from '@material-ui/core/Chip'
import DoneIcon from '@material-ui/icons/Done'
import PendingIcon from '@material-ui/icons/Cached'
import FailIcon from '@material-ui/icons/Clear'
import NotInstalledIcon from '@material-ui/icons/Info'
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
      case 'DEPLOYED':
        return (
          <MuiThemeProvider theme={tagTheme}>
            <Chip
              icon={<DoneIcon />}
              label="Deployed"
              color="primary"
              variant="outlined"
            />
          </MuiThemeProvider>
        )
      case 'DELETED':
        return (
          <MuiThemeProvider theme={tagTheme}>
            <Chip
              icon={<FailIcon />}
              label="Deleted"
              color="secondary"
              variant="outlined"
            />
          </MuiThemeProvider>
        )
      case 'SUPERSEDED':
        return (
          <MuiThemeProvider theme={tagTheme}>
            <Chip
              icon={<DoneIcon />}
              label="Superseded"
              color="primary"
              variant="outlined"
            />
          </MuiThemeProvider>
        )
      case 'FAILED':
        return (
          <MuiThemeProvider theme={tagTheme}>
            <Chip
              icon={<FailIcon />}
              label="Failed"
              color="secondary"
              variant="outlined"
            />
          </MuiThemeProvider>
        )
      case 'DELETING':
        return (
          <MuiThemeProvider theme={tagTheme}>
            <Chip
              icon={<PendingIcon />}
              label="Deleting"
              color="secondary"
              variant="outlined"
            />
          </MuiThemeProvider>
        )
      case 'PENDING_INSTALL':
        return (
          <MuiThemeProvider theme={tagTheme}>
            <Chip
              icon={<PendingIcon />}
              label="Deploying"
              color="primary"
              variant="outlined"
            />
          </MuiThemeProvider>
        )
      case 'PENDING_UPGRADE':
        return (
          <MuiThemeProvider theme={tagTheme}>
            <Chip
              icon={<PendingIcon />}
              label="Upgrading"
              color="primary"
              variant="outlined"
            />
          </MuiThemeProvider>
        )
      case 'PENDING_ROLLBACK':
        return (
          <MuiThemeProvider theme={tagTheme}>
            <Chip
              icon={<PendingIcon />}
              label="Rolling Back"
              color="secondary"
              variant="outlined"
            />
          </MuiThemeProvider>
        )
      default: // UNKNOWN default status case
        return (
          <MuiThemeProvider theme={tagTheme}>
            <Chip
              icon={<NotInstalledIcon />}
              label="Unknown"
              color="default"
              variant="outlined"
            />
          </MuiThemeProvider>
        )
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