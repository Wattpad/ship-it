import React from 'react'
import { AppBar, Toolbar, Typography } from '@material-ui/core'
import { createMuiTheme, MuiThemeProvider } from '@material-ui/core/styles'
import IconButton from '@material-ui/core/IconButton'
import ShipIcon from '../assets/passenger_ship.png'
import * as constants from '../Constants'
import * as themes from '../Themes'

class TopBar extends React.Component {
  constructor(props) {
    super(props)
    this.state = {}
  }

  render() {
    return (
      <MuiThemeProvider theme={themes.standard}>
        <AppBar position="static" color="primary">
          <Toolbar>
            <IconButton>
              <img src={ShipIcon} width="32" height="32" alt={constants.imgAlt} />
            </IconButton>
            <Typography variant="h6">
              Ship-it!
            </Typography>
          </Toolbar>
        </AppBar>
      </MuiThemeProvider>
    )
  }
}

export default TopBar