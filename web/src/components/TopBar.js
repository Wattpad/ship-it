import React from 'react'
import { AppBar, Toolbar, Typography, InputBase } from '@material-ui/core'
import { withStyles, createMuiTheme, MuiThemeProvider } from '@material-ui/core/styles'
import { fade } from '@material-ui/core/styles/colorManipulator'
import PropTypes from 'prop-types'
import SearchIcon from '@material-ui/icons/Search'
import IconButton from '@material-ui/core/IconButton'
import ShipIcon from '../assets/passenger_ship.png'

const imgAlt = "not found"

const theme = createMuiTheme({
  palette: {
    primary: {
      main: '#FF6612'
    },
    secondary: {
      main: '#FEAF0A'
    }
  }
})

class TopBar extends React.Component {
  constructor(props) {
    super(props)
    this.state = {}
  }

  handleChange = (event) => {
    this.setState({
      searchText: event.target.value
    })
    console.log(this.state.searchText)
  }

  render() {
    const { classes } = this.props
    return (
      <MuiThemeProvider theme={theme}>
        <AppBar position="static" color="primary">
          <Toolbar>
            <IconButton>
              <img src={ShipIcon} width="32" height="32" alt={imgAlt} />
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

TopBar.propTypes = {
  classes: PropTypes.object.isRequired,
}

export default TopBar