import React from 'react'
import { AppBar, Toolbar, Typography, InputBase } from '@material-ui/core'
import { withStyles, createMuiTheme, MuiThemeProvider } from '@material-ui/core/styles'
import { fade } from '@material-ui/core/styles/colorManipulator'
import PropTypes from 'prop-types'
import SearchIcon from '@material-ui/icons/Search'
import IconButton from '@material-ui/core/IconButton'
import ShipIcon from '../assets/passenger_ship.png'

const styles = theme => ({
  root: {
    width: '100%',
  },
  grow: {
    flexGrow: 1,
  },
  menuButton: {
    marginLeft: -12,
    marginRight: 20,
  },
  title: {
    display: 'none',
    [theme.breakpoints.up('sm')]: {
      display: 'block',
    },
  },
  search: {
    position: 'relative',
    borderRadius: theme.shape.borderRadius,
    backgroundColor: fade(theme.palette.common.white, 0.15),
    '&:hover': {
      backgroundColor: fade(theme.palette.common.white, 0.25),
    },
    marginRight: theme.spacing.unit * 2,
    marginLeft: 0,
    width: '100%',
    [theme.breakpoints.up('sm')]: {
      marginLeft: theme.spacing.unit * 3,
      width: 'auto',
    },
  },
  searchIcon: {
    width: theme.spacing.unit * 9,
    height: '100%',
    position: 'absolute',
    pointerEvents: 'none',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  },
  inputRoot: {
    color: 'inherit',
    width: '100%',
  },
  inputInput: {
    paddingTop: theme.spacing.unit,
    paddingRight: theme.spacing.unit,
    paddingBottom: theme.spacing.unit,
    paddingLeft: theme.spacing.unit * 10,
    transition: theme.transitions.create('width'),
    width: '100%',
    [theme.breakpoints.up('md')]: {
      width: 200,
    },
  },
})

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
            {/* <div className={classes.search}>
              <div className={classes.searchIcon}>
                <SearchIcon />
              </div>
              <InputBase
                placeholder="Search…"
                onChange={this.handleChange}
                classes={{
                  root: classes.inputRoot,
                  input: classes.inputInput,
                }}
                color="secondary"
              />
            </div> */}
          </Toolbar>
        </AppBar>
      </MuiThemeProvider>
    )
  }
}

TopBar.propTypes = {
  classes: PropTypes.object.isRequired,
}

export default withStyles(styles)(TopBar)