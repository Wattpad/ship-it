import React from 'react'
import { AppBar, Toolbar, Typography } from '@material-ui/core'
import MuiThemeProvider from '@material-ui/core/styles/MuiThemeProvider'
import IconButton from '@material-ui/core/IconButton'
import InputAdornment from '@material-ui/core/InputAdornment';
import TextField from '@material-ui/core/TextField'
import SearchIcon from '@material-ui/icons/Search';
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
            <div className="search-releases">
              <TextField
                fullWidth
                InputProps={{
                  type: "search",
                  placeholder: "Search by name, squad or status",
                  startAdornment: (
                    <InputAdornment position="start">
                      <SearchIcon />
                    </InputAdornment>
                  )
                }}
                onChange={e => this.props.onQueryChanged(e.target.value)}
              />
            </div>
          </Toolbar>
        </AppBar>
      </MuiThemeProvider>
    )
  }
}

export default TopBar
