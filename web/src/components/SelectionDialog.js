import React from 'react'
import Dialog from '@material-ui/core/Dialog'
import DialogTitle from '@material-ui/core/DialogTitle'
import List from '@material-ui/core/List'
import { ListItem, IconButton, DialogContent, DialogActions, Button, Link } from '@material-ui/core'
import GitIcon from '../assets/octocat.png'
import HelmIcon from '../assets/helm_icon.png'

import * as constants from '../Constants'

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

  render() {
    return (
      <div>
        {
          this.state.open ?
            <Dialog onClose={this.handleClose} open={this.state.open}>
              <DialogTitle id="simple-dialog-title">Select Repository</DialogTitle>
              <DialogContent>
                <List>
                  <ListItem>
                    <IconButton href={this.props.chart}>
                      <img src={HelmIcon} width="32" height="32" alt={constants.imgAlt} />
                    </IconButton>
                    <Link href={this.props.chart}>Chart Download</Link>
                  </ListItem>
                  <ListItem>
                    <IconButton href={this.props.source}>
                      <img src={GitIcon} width="32" height="32" alt={constants.imgAlt} />
                    </IconButton>
                    <Link href={this.props.source}>Source</Link>
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
          <img src={GitIcon} width="32" height="32" alt={constants.imgAlt} />
        </IconButton>
      </div>
    )
  }
}

export default SelectionDialog
