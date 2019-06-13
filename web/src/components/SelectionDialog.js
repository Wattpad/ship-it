import React from 'react'
import Dialog from '@material-ui/core/Dialog'
import DialogTitle from '@material-ui/core/DialogTitle';
import List from '@material-ui/core/List'
import { ListItem, IconButton, DialogContent, DialogActions, Button } from '@material-ui/core';
import GitIcon from '../assets/octocat.png'

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
                    <IconButton>
                      <img src={GitIcon} width="32" height="32" alt={imgAlt} />
                    </IconButton>
                    Miranda
                            </ListItem>
                  <ListItem>
                    <IconButton>
                      <img src={GitIcon} width="32" height="32" alt={imgAlt} />
                    </IconButton>
                    Highlander
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