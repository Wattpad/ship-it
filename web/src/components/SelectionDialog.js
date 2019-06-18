import React from 'react'
import Dialog from '@material-ui/core/Dialog'
import DialogTitle from '@material-ui/core/DialogTitle'
import List from '@material-ui/core/List'
import { ListItem, IconButton, DialogContent, DialogActions, Button, Link } from '@material-ui/core'
import GitIcon from '../assets/octocat.png'

const urljoin = require('url-join')
const imgAlt = "not found"

class SelectionDialog extends React.Component {
  constructor(props) {
    super(props)
    this.state = { open: false }
  }

  handleClose = () => {
    this.setState({ open: false })
  }

  handleOpen = () => {
    this.setState({ open: true })
  }

  getMirandaURL() {
    var url = urljoin('https://github.com/Wattpad/highlander/tree/master/k8s/charts/services', this.props.serviceName)
    return url
  }

  getHighlanderURL() {
    var url = urljoin('https://github.com/Wattpad/highlander/tree/', this.props.gitref, this.props.highlanderPath)
    return url
  }

  render() {
    var highlanderURL = this.getHighlanderURL()
    var mirandaURL = this.getMirandaURL()
    return (
      <div>
        {
          this.state.open ?
            <Dialog onClose={this.handleClose} open={this.state.open}>
              <DialogTitle id="simple-dialog-title">Select Repository</DialogTitle>
              <DialogContent>
                <List>
                  <ListItem>
                    <IconButton href={mirandaURL}>
                      <img src={GitIcon} width="32" height="32" alt={imgAlt} />
                    </IconButton>
                    <Link href={mirandaURL}>Miranda</Link>
                  </ListItem>
                  <ListItem>
                    <IconButton href={highlanderURL}>
                      <img src={GitIcon} width="32" height="32" alt={imgAlt} />
                    </IconButton>
                    <Link href={highlanderURL}>Highlander</Link>
                  </ListItem>
                </List>
              </DialogContent>
              <DialogActions>
                <Button onClick={this.handleClose}>Close</Button>
              </DialogActions>
            </Dialog>
            :
            null
        }
        <IconButton onClick={this.handleOpen} width="32" height="32">
          <img src={GitIcon} width="32" height="32" alt={imgAlt} />
        </IconButton>
      </div>
    )
  }
}

export default SelectionDialog